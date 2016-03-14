package main

import (
	"encoding/json"
	"fmt"
	"github.com/ahirmayur/gorm"
	"github.com/codegangsta/martini"
	"net/http"
	"sort"
	"strings"
)

type CheckDict struct {
	Key   string
	Value int
}

type IdentifyResponse struct {
	Cols    []ColType
	Request string
}

type ColType struct {
	Name    string
	Sqltype string
}

type Suggestionresponse struct {
	Request string
}

type StringMatchResult struct {
	Count int
	Match string
}

type ScanJob struct {
	TableName string
	X         string
}

type PossibleCombo struct {
	Match  string
	Tables []string
}

/**
 * @brief This function checks to see if the data has been imported yet or still is in need of importing
 * @details
 *
 * @param http
 * @param http
 * @param martini
 * @return
 */
func IdentifyTable(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	if params["id"] == "" {
		http.Error(res, "There was no ID request", http.StatusBadRequest)
		return ""
	}
	results := FetchTableCols(string(params["id"]))

	returnobj := IdentifyResponse{
		Cols:    results,
		Request: params["id"],
	}
	b, _ := json.Marshal(returnobj)

	return string(b)
}

/**
 * @brief This fetches an array of all the col names and their types.
 * @details
 *
 * @param string
 * @return
 */
func FetchTableCols(guid string) (output []ColType) {
	if guid == "" {
		return output
	}

	onlineData, e := GetOnlineDataByGuid(guid)
	if e != nil {
		return output
	}

	results := GetSQLTableSchema(onlineData.Tablename)

	return results
}

func HasTableGotLocationData(datasetGUID string) bool {
	cols := FetchTableCols(datasetGUID)

	if ContainsTableCol(cols, "lat") && (ContainsTableCol(cols, "lon") || ContainsTableCol(cols, "long")) {
		return true
	}

	return false
}

func ContainsTableCol(cols []ColType, target string) bool {
	for _, v := range cols {
		if strings.ToLower(v.Name) == target {
			return true
		}
	}

	return false
}

func CheckColExists(schema []ColType, column string) bool {
	for _, val := range schema {
		if val.Name == column {
			return true
		}
	}

	return false
}

/**
 * @brief Find matching data for given ID & cordinates
 * m.Get("/api/findmatches/:id/:x/:y", api.AttemptToFindMatches)
 *
 * @param http.ResponseWriter
 * @param http.Request
 * @param martini.Params
 *
 * @return JSON containing Matched data
 */
func AttemptToFindMatches(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	onlineData, e := GetOnlineDataByGuid(params["id"])
	if e != nil {
		http.Error(res, "Could not find that Table", http.StatusInternalServerError)
		return ""
	}

	schema := GetSQLTableSchema(onlineData.Tablename)

	if !CheckColExists(schema, params["x"]) || !CheckColExists(schema, params["y"]) {
		http.Error(res, "Could not find the X or Y", http.StatusInternalServerError)
		return ""
	}

	/* Check if data exists in the stats table. so we can compare its poly to other Polynomial */
	stats := StatsCheck{}
	count := 0
	err := DB.Model(&stats).Where(fmt.Sprintf("LOWER(%q) = ?", "table"), strings.ToLower(onlineData.Tablename)).Where("LOWER(x) = ?", strings.ToLower(params["x"])).Where("LOWER(y) = ?", strings.ToLower(params["y"])).Count(&count).Find(&stats).Error
	if err != nil && err != gorm.RecordNotFound {
		http.Error(res, "Unexpected error (AttemptToFindMatches).", http.StatusInternalServerError)
		return ""
	}

	if count == 0 || err == gorm.RecordNotFound {
		http.Error(res, "Unable to find the Polynomial code using given X and Y co-ordinates.", http.StatusBadRequest)
		return ""
	}

	b, e := json.Marshal(stats)
	if e != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(b)
}

func FindStringMatches(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	if params["word"] == "" {
		http.Error(res, "Please add a word", http.StatusBadRequest)
		return ""
	}

	Results := make([]StringMatchResult, 0)

	search := []StringSearch{}

	query := DB.Select("tablename, count").Where("LOWER(value) = ?", strings.ToLower(params["word"]))
	if params["x"] != "" {
		query = query.Where("x = ?", params["x"])
	}
	err := query.Find(&search).Error
	if err != nil && err != gorm.RecordNotFound {
		http.Error(res, "SQL error", http.StatusInternalServerError)
		return ""
	}

	for _, s := range search {
		result := StringMatchResult{
			Count: s.Count,
			Match: s.Tablename,
		}

		Results = append(Results, result)
	}

	b, e := json.Marshal(Results)
	if e != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(b)
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

type ByVal []CheckDict

func (a ByVal) Len() int           { return len(a) }
func (a ByVal) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByVal) Less(i, j int) bool { return a[i].Value < a[j].Value }

func ConvertIntoStructArrayAndSort(input map[string]int) (in []CheckDict) {
	in = make([]CheckDict, 0)
	for k, v := range input {
		newd := CheckDict{
			Key:   k,
			Value: v,
		}

		in = append(in, newd)
	}

	sort.Sort(ByVal(in))

	return in
}
