package main

import (
	"encoding/json"
	"fmt"
	"github.com/ahirmayur/gorm"
	"github.com/codegangsta/martini"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type SearchResult struct {
	Title       string
	GUID        string
	PrimaryDate string
}

type SearchResponse struct {
	Keyword string
	Results []SearchResult
	Total   int
}

func SearchForDataHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	result, error := SearchForData(uid, params["keyword"], params)
	if error != nil {
		http.Error(res, error.Message, error.Code)
		return ""
	}

	r, err1 := json.Marshal(result)
	if err1 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

/**
 * @brief Search a given keyword in database
 * @details This method searches for a matching title with following conditions,
 * 		Suffix wildcard
 * 		Prefix & suffix wildcard
 * 		Prefix, suffix & trimmed spaces with wildcard
 * 		Prefix & suffix on previously searched terms
 */
func SearchForData(uid int, keyword string, params map[string]string) (SearchResponse, *appError) {
	response := SearchResponse{}
	if keyword == "" {
		return response, &appError{nil, "There was no search request", http.StatusBadRequest}
	}

	offset, count := 0, 9

	if params["offset"] != "" {
		var oE error
		offset, oE = strconv.Atoi(params["offset"])
		if oE != nil {
			return response, &appError{oE, "Invalid offset value.", http.StatusBadRequest}
		}
	}

	if params["count"] != "" {
		var cE error
		count, cE = strconv.Atoi(params["count"])
		if params["count"] != "" && cE != nil {
			return response, &appError{cE, "Invalid count value.", http.StatusBadRequest}
		}
	}

	// Search index
	indices := []Index{}
	keyword = strings.ToLower(strings.Trim(keyword, " "))
	term := "%" + strings.Replace(keyword, " ", "%", -1) + "%" // e.g. "gold" => "%gold%", nh s" => "%nh%s%", "  cri m e " => "%cri%m%e%"

	query := DB.Where(fmt.Sprintf("LOWER(%q) LIKE ?", "title"), term)
	query = query.Or(fmt.Sprintf("LOWER(%q) LIKE ?", "notes"), term)
	query = query.Or(fmt.Sprintf("LOWER(%q) LIKE ?", "name"), term)
	query = query.Or(fmt.Sprintf("LOWER(%q) LIKE ?", "desc"), term)

	err := query.Order("random()").Limit(count).Offset(offset).Find(&indices).Error
	if err != nil && err != gorm.RecordNotFound {
		return response, &appError{err, "Database query failed (Index - random)", http.StatusServiceUnavailable}
	} else if err == gorm.RecordNotFound {
		indices = make([]Index, 0)
	}

	searchResults := indices

	// Search table columns
	err1 := DB.Find(&indices).Error
	if err1 != nil {
		return response, &appError{err, "Database query failed (Index - all)", http.StatusServiceUnavailable}
	}

	for _, table := range indices {
		schema := GetSQLTableSchema(table.Guid)
		for _, column := range schema {
			if strings.Contains(strings.ToLower(column.Name), keyword) {
				searchResults = append(searchResults, table)
			}
		}
	}

	Response := ProcessSearchResults(term, searchResults)

	// Randomise order
	for i := range Response.Results {
		j := rand.Intn(i + 1)
		Response.Results[i], Response.Results[j] = Response.Results[j], Response.Results[i]
	}

	totalCharts := len(Response.Results)
	if offset > totalCharts {
		return SearchResponse{}, nil
	}

	last := offset + count
	if last > totalCharts {
		last = totalCharts
	}

	Response.Results = Response.Results[offset:last] // return slice

	searchTermErr := AddSearchTerm(keyword, 1) // add to search term count

	return Response, searchTermErr
}

func ProcessSearchResults(keyword string, rows []Index) SearchResponse {
	Results := make([]SearchResult, 0)

	for _, row := range rows {
		result := SearchResult{
			Title:       SanitizeString(row.Title),
			GUID:        SanitizeString(row.Guid),
			PrimaryDate: row.PrimaryDate,
		}

		Results = append(Results, result)
	}

	Response := SearchResponse{
		Keyword: keyword,
		Results: Results,
		Total:   len(rows),
	}

	return Response
}

func SanitizeString(str string) string {
	return strings.Replace(str, "Ã‚Â£", "£", -1)
}

func AddSearchTerm(str string, attempts int) *appError {
	searchterm := SearchTerm{}

	err := DB.Where("term = ?", str).Find(&searchterm).Error
	if err != nil && err != gorm.RecordNotFound {
		if attempts > 5 {
			return &appError{err, "Database query failed (Add Search Term)", http.StatusInternalServerError}
		}

		return AddSearchTerm(str, attempts+1)
	} else if err == gorm.RecordNotFound {
		searchterm.Count = 0
		searchterm.Term = str
	}

	searchterm.Count++
	err = DB.Save(&searchterm).Error

	if err != nil {
		return &appError{err, "Database query failed (Save Search Term)", http.StatusInternalServerError}
	}

	return nil
}

/**
 * @brief Build Data Dictionary
 * @details Takes all the key terms from the title, name and description in the index table and writes them to the datadictionary along with their frequency
 */
func BuildDataDictionary() {
	indices := []Index{}
	err := DB.Find(&indices).Error
	if err != nil {
		Logger.Println("Error" + err.Error())
		return
	}

	terms := make(map[string]int)
	re, _ := regexp.Compile("[^A-Za-z]+") //\W

	for _, index := range indices {
		title := re.ReplaceAllString(index.Title, "-")
		title = strings.ToLower(strings.Trim(title, "-"))
		keywords_title := strings.Split(title, "-")
		for _, k := range keywords_title {
			if len(k) > 2 {
				if _, ok := terms[k]; !ok {
					terms[k] = 0
				}

				terms[k]++
			}
		}

		name := re.ReplaceAllString(index.Name, "-")
		name = strings.ToLower(strings.Trim(name, "-"))
		keywords_name := strings.Split(name, "-")
		for _, k := range keywords_name {
			if len(k) > 2 {
				if _, ok := terms[k]; !ok {
					terms[k] = 0
				}

				terms[k]++
			}
		}

		desc := re.ReplaceAllString(index.Desc, "-")
		desc = strings.ToLower(strings.Trim(desc, "-"))
		keywords_desc := strings.Split(desc, "-")
		for _, k := range keywords_desc {
			if len(k) > 2 {
				if _, ok := terms[k]; !ok {
					terms[k] = 0
				}

				terms[k]++
			}
		}
	}

	for term, frequency := range terms {
		dictionary := Dictionary{
			Term:      term,
			Frequency: frequency,
		}

		err1 := DB.Create(&dictionary).Error
		if err1 != nil {
			Logger.Println("Error: " + err.Error())
		}
	}
}
