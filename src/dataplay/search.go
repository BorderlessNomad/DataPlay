package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/jinzhu/gorm"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type SearchResult struct {
	Title        string
	GUID         string
	LocationData bool
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

func SearchForDataQ(params map[string]string) string {
	if params["user"] == "" {
		return ""
	}

	uid, e := strconv.Atoi(params["user"])
	if e != nil {
		return ""
	}

	result, err := SearchForData(uid, params["keyword"], params)
	if err != nil {
		return ""
	}

	r, e := json.Marshal(result)
	if e != nil {
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

	var total int = 0
	var offset int = 0
	var count int = 9

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

	term := keyword + "%" // e.g. "nhs" => "nhs%" (What about "%nhs"?)

	Logger.Println("Searching with Suffix Wildcard", term)

	var err error
	err = DB.Where("LOWER(title) LIKE LOWER(?)", term).Where("(owner = 0 OR owner = ?)", uid).Limit(count).Offset(offset).Find(&indices).Limit(-1).Offset(-1).Count(&total).Error
	if err != nil && err != gorm.RecordNotFound {
		return response, &appError{err, "Database query failed (SUFFIX)", http.StatusServiceUnavailable}
	}

	Response := ProcessSearchResults(term, indices, total, err)
	if len(Response.Results) == 0 {
		term := "%" + keyword + "%" // e.g. "nhs" => "%nhs%"

		Logger.Println("Searching with Prefix + Suffix Wildcard", term)

		err = DB.Where("LOWER(title) LIKE LOWER(?)", term).Where("(owner = 0 OR owner = ?)", uid).Limit(count).Offset(offset).Find(&indices).Limit(-1).Offset(-1).Count(&total).Error
		if err != nil && err != gorm.RecordNotFound {
			return response, &appError{err, "Database query failed (PREFIX + SUFFIX)", http.StatusServiceUnavailable}
		}

		Response = ProcessSearchResults(term, indices, total, err)
		if len(Response.Results) == 0 {
			term := "%" + strings.Replace(keyword, " ", "%", -1) + "%" // e.g. "nh s" => "%nh%s%"

			Logger.Println("Searching with Prefix + Suffix + Trim Wildcard", term)

			err = DB.Where("LOWER(title) LIKE LOWER(?)", term).Where("(owner = 0 OR owner = ?)", uid).Limit(count).Offset(offset).Find(&indices).Limit(-1).Offset(-1).Count(&total).Error
			if err != nil && err != gorm.RecordNotFound {
				return response, &appError{err, "Database query failed (PREFIX + SUFFIX + TRIM)", http.StatusServiceUnavailable}
			}

			Response = ProcessSearchResults(term, indices, total, err)
			if len(Response.Results) == 0 && (len(keyword) >= 3 && len(keyword) < 20) {
				term := "%" + keyword + "%" // e.g. "nhs" => "%nhs%"

				Logger.Println("Searching with Prefix + Suffix Wildcard in String Table", term)

				indicesAll := []Index{}
				query := DB.Table("priv_stringsearch, priv_onlinedata, index")
				query = query.Select("DISTINCT ON (priv_onlinedata.guid) priv_onlinedata.guid, index.title")
				query = query.Where("(LOWER(value) LIKE LOWER(?) OR LOWER(x) LIKE LOWER(?))", term, term)
				query = query.Where("priv_stringsearch.tablename = priv_onlinedata.tablename")
				query = query.Where("priv_onlinedata.guid = index.guid")
				query = query.Where("(owner = ? OR owner = ?)", 0, uid)
				query = query.Order("priv_onlinedata.guid")
				query = query.Order("priv_stringsearch.count DESC")
				err = query.Limit(count).Offset(offset).Find(&indices).Limit(-1).Offset(-1).Find(&indicesAll).Error

				total = len(indicesAll)

				if err != nil && err != gorm.RecordNotFound {
					return response, &appError{err, "Database query failed (PREFIX + SUFFIX + STRING)", http.StatusInternalServerError}
				}

				Response = ProcessSearchResults(term, indices, total, err)
				if len(Response.Results) == 0 && (len(keyword) >= 3 && len(keyword) < 20) {
					term := "%" + strings.Replace(keyword, " ", "%", -1) + "%" // e.g. "nh s" => "%nh%s%"

					Logger.Println("Searching with Prefix + Suffix + Trim Wildcard in String Table", term)

					indicesAll := []Index{}
					query := DB.Table("priv_stringsearch, priv_onlinedata, index")
					query = query.Select("DISTINCT ON (priv_onlinedata.guid) priv_onlinedata.guid, index.title")
					query = query.Where("(LOWER(value) LIKE LOWER(?) OR LOWER(x) LIKE LOWER(?))", term, term)
					query = query.Where("priv_stringsearch.tablename = priv_onlinedata.tablename")
					query = query.Where("priv_onlinedata.guid = index.guid")
					query = query.Where("(owner = ? OR owner = ?)", 0, uid)
					query = query.Order("priv_onlinedata.guid")
					query = query.Order("priv_stringsearch.count DESC")
					err = query.Limit(count).Offset(offset).Find(&indices).Limit(-1).Offset(-1).Find(&indicesAll).Error

					total = len(indicesAll)

					if err != nil && err != gorm.RecordNotFound {
						return response, &appError{err, "Database query failed (PREFIX + SUFFIX + STRING)", http.StatusInternalServerError}
					}

					Response = ProcessSearchResults(term, indices, total, err)
				}
			}
		}
	}

	return Response, nil
}

func ProcessSearchResults(term string, rows []Index, total int, e error) SearchResponse {
	if e != nil && e != gorm.RecordNotFound {
		check(e)
	}

	Results := make([]SearchResult, 0)

	for _, row := range rows {
		Location := HasTableGotLocationData(row.Guid)

		result := SearchResult{
			Title:        SanitizeString(row.Title),
			GUID:         SanitizeString(row.Guid),
			LocationData: Location,
		}

		Results = append(Results, result)
	}

	Response := SearchResponse{
		Keyword: term,
		Results: Results,
		Total:   total,
	}

	return Response
}

func SanitizeString(str string) string {
	return strings.Replace(str, "Ã‚Â£", "£", -1)
}

func AddSearchTerm(str string) {
	searchterm := SearchTerm{}

	err := DB.Where("term = ?", str).Find(&searchterm).Error
	if err != nil && err != gorm.RecordNotFound {
		panic(err)
	} else if err == gorm.RecordNotFound {
		searchterm.Count = 0
		searchterm.Term = str
	}

	searchterm.Count++
	err = DB.Save(&searchterm).Error
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
