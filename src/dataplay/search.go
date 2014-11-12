package main

import (
	"encoding/json"
	"github.com/ahirmayur/gorm"
	"github.com/codegangsta/martini"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type SearchResult struct {
	Title        string
	GUID         string
	LocationData bool
	PrimaryDate  string
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

	AddSearchTerm(keyword) // add to search term count

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

	indices := []Index{}

	// keyword := strings.Trim(keyword, " ")      // Remove first and last spaces (if any)
	// r, _ := regexp.Compile("/[^A-Za-z]+/")     // RegEx compile
	// keyword = r.ReplaceAllString(keyword, "%") // Replace all non-alphabets with %
	// term := "%" + keyword + "%"                // e.g. "nh s" => "%nh%s%"

	keyword = strings.Trim(keyword, " ")
	term := "%" + strings.Replace(keyword, " ", "%", -1) + "%" // e.g. "gold" => "%gold%", nh s" => "%nh%s%", "  cri m e " => "%cri%m%e%"

	Logger.Println("Searching for term: '%s'", term)

	query := DB.Where("LOWER(title) LIKE LOWER(?)", term)
	query = query.Or("LOWER(notes) LIKE LOWER(?)", term)
	query = query.Or("LOWER(name) LIKE LOWER(?)", term)

	err := query.Order("random()").Limit(count).Offset(offset).Find(&indices).Error
	if err != nil && err != gorm.RecordNotFound {
		return response, &appError{err, "Database query failed", http.StatusServiceUnavailable}
	}

	Response := ProcessSearchResults(term, indices)

	// Randomise order
	for i := range Response.Results {
		j := rand.Intn(i + 1)
		Response.Results[i], Response.Results[j] = Response.Results[j], Response.Results[i]
	}

	return Response, nil
}

func ProcessSearchResults(term string, rows []Index) SearchResponse {
	Results := make([]SearchResult, 0)

	for _, row := range rows {
		Location := HasTableGotLocationData(row.Guid)

		result := SearchResult{
			Title:        SanitizeString(row.Title),
			GUID:         SanitizeString(row.Guid),
			LocationData: Location,
			PrimaryDate:  row.PrimaryDate,
		}

		Results = append(Results, result)
	}

	Response := SearchResponse{
		Keyword: term,
		Results: Results,
		Total:   len(rows),
	}

	return Response
}

func SanitizeString(str string) string {
	return strings.Replace(str, "Ã‚Â£", "£", -1)
}

func AddSearchTerm(str string) {
	searchterm := SearchTerm{}

	err := DB.Where("term = ?", str).Find(&searchterm).Error
	if err == nil && err != gorm.RecordNotFound {
		searchterm.Count++
		err = DB.Save(&searchterm).Error
	} else if err == gorm.RecordNotFound {
		searchterm.Count = 0
		searchterm.Term = str
		searchterm.Count++
		err = DB.Save(&searchterm).Error
	}
}

// Takes all the key terms from the title, name and description in the index table and writes them to the datadictionary along with their frequency
func DataDict() {
	indices := []Index{}
	DB.Find(&indices)
	var terms []string
	re, _ := regexp.Compile("\\W")

	for _, ind := range indices {
		title := strings.ToLower(ind.Title)
		title = strings.Replace(title, "_", " ", -1)
		title = re.ReplaceAllString(title, " ")
		term := strings.Split(title, " ")
		for i, _ := range term {
			terms = append(terms, term[i])
		}

		name := strings.ToLower(ind.Name)
		name = strings.Replace(name, "_", " ", -1)
		name = re.ReplaceAllString(name, " ")
		term = strings.Split(name, " ")
		for i, _ := range term {
			terms = append(terms, term[i])
		}

		notes := strings.ToLower(ind.Notes)
		notes = strings.Replace(notes, "_", " ", -1)
		notes = re.ReplaceAllString(notes, " ")
		term = strings.Split(notes, " ")
		for i, _ := range term {
			terms = append(terms, term[i])
		}
	}

	var dict []Dictionary
	for _, t := range terms {
		termNotPresent := true
		for i, _ := range dict {
			if t == dict[i].Term {
				dict[i].Frequency += 1
				termNotPresent = false
			}
		}
		if termNotPresent && len(t) > 2 {
			var tmp Dictionary
			tmp.Term = t
			tmp.Frequency = 1
			dict = append(dict, tmp)
		}
	}

	for _, d := range dict {
		DB.Create(&d)
	}
}
