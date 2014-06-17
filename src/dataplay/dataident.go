package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/jinzhu/gorm"
	"net/http"
	"sort"
	"strconv"
	"strings"
	// "fmt"
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

type SuggestionResponce struct {
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

// This function checks to see if the data has been imported yet or still is in need of importing
func IdentifyTable(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	if prams["id"] == "" {
		http.Error(res, "There was no ID request", http.StatusBadRequest)
		return ""
	}
	results := FetchTableCols(string(prams["id"]))

	returnobj := IdentifyResponse{
		Cols:    results,
		Request: prams["id"],
	}
	b, _ := json.Marshal(returnobj)

	return string(b)
}

// This fetches a array of all the col names and their types.
func FetchTableCols(guid string) (output []ColType) {
	if guid == "" {
		return output
	}

	var tablename string
	tablename, e := getRealTableName(guid)
	if e != nil {
		return output
	}

	results := GetSQLTableSchema(tablename)

	return results
}

/**
 * @brief Get the SQL Scheme for a Table
 * @details Almost all of the SQLs support 'information_schema' database which stores metadata about
 * other databases, tables etc.
 *
 * @todo Apply caching to queries which goes to 'information_schema'
 * MySQL has something like innodb_stats_on_metadata=0 which will prevent statistic update upon quering 'information_schema'.
 * Also it won't make 'information_schema' to be stale when changes are made on corresponding metadata.
 *
 * @param string <Table Name>
 * @return <Table Schema>
 */
func GetSQLTableSchema(table string, databaseName ...string) []ColType {
	database := "dataplay"
	if len(databaseName) > 0 {
		database = databaseName[0]
	}

	tableSchema := []TableSchema{}
	err := DB.Select("column_name, data_type").Where("table_catalog = ?", database).Where("table_name = ?", table).Find(&tableSchema).Error

	if err != gorm.RecordNotFound {
		check(err)
	}

	schema := make([]ColType, 0)

	for _, row := range tableSchema {
		NewCol := ColType{
			Name:    row.ColumnName,
			Sqltype: row.DataType,
		}

		if NewCol.Sqltype == "character varying" {
			NewCol.Sqltype = "varchar"
		} else if NewCol.Sqltype == "numeric" {
			NewCol.Sqltype = "float"
		}

		schema = append(schema, NewCol)
	}

	return schema
}

func CheckColExists(schema []ColType, column string) bool {
	for _, val := range schema {
		if val.Name == column {
			return true
		}
	}

	return false
}

func SuggestColType(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	if prams["table"] == "" || prams["col"] == "" {
		http.Error(res, "There was no ID request", http.StatusBadRequest)
		return ""
	}

	onlineData := OnlineData{}
	err := DB.Select("tablename").Where("guid = ?", prams["table"]).Find(&onlineData).Error

	if err == gorm.RecordNotFound {
		http.Error(res, "Could not find that Table", http.StatusNotFound)
		return ""
	} else if err != nil {
		check(err)
		return ""
	}

	schema := GetSQLTableSchema(onlineData.Tablename)

	if !CheckColExists(schema, prams["col"]) {
		http.Error(res, "You have requested a Columns that does not exist.", http.StatusBadRequest)
		return ""
	}

	var data []string
	err = DB.Table(onlineData.Tablename).Pluck(prams["col"], &data).Error
	check(err)

	for _, val := range data {
		_, e := strconv.ParseFloat(val, 10)
		if e != nil {
			return "false"
		}
	}

	return "true"
}

// Unfinished function
func AttemptToFindMatches(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	// m.Get("/api/findmatches/:id/:x/:y", api.AttemptToFindMatches)
	RealTableName, e := getRealTableName(prams["id"], res)
	if e != nil {
		http.Error(res, "Could not find that Table", http.StatusInternalServerError)
		return ""
	}

	schema := GetSQLTableSchema(RealTableName)

	if !CheckColExists(schema, prams["x"]) || !CheckColExists(schema, prams["y"]) {
		http.Error(res, "Could not find the X or Y", http.StatusInternalServerError)
		return ""
	}

	/* Check if data exists in the stats table. so we can compare its poly to other Polynomial */
	stats := StatsCheck{}
	count := 0
	err := DB.Model(&stats).Where("LOWER(\"table\") = ?", strings.ToLower(RealTableName)).Where("LOWER(x) = ?", strings.ToLower(prams["x"])).Where("LOWER(y) = ?", strings.ToLower(prams["y"])).Count(&count).Find(&stats).Error
	check(err)

	if count == 0 {
		http.Error(res, "Unable to find the Polynomial code using given X and Y co-ordinates.", http.StatusBadRequest)
		return ""
	}

	Logger.Println(&stats)

	return "wat"
}

func FindStringMatches(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	if prams["word"] == "" {
		http.Error(res, "Please add a word", http.StatusBadRequest)
		return ""
	}

	Results := make([]StringMatchResult, 0)

	search := []StringSearch{}

	query := DB.Select("tablename, count").Where("LOWER(value) = ?", strings.ToLower(prams["word"]))
	if prams["x"] != "" {
		query = query.Where("x = ?", prams["x"])
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

func GetRelatedDatasetByStrings(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	RealTableName, e := getRealTableName(prams["guid"], res)
	if e != nil {
		http.Error(res, "Could not find that table", http.StatusInternalServerError)
		return ""
	}

	jobs := make([]ScanJob, 0)

	Bits := GetSQLTableSchema(RealTableName)

	/* Prepare a job list */
	for _, bit := range Bits {
		if bit.Sqltype == "varchar" || bit.Sqltype == "character varying" {
			newJob := ScanJob{
				TableName: RealTableName,
				X:         bit.Name,
			}

			jobs = append(jobs, newJob)
		}
	}

	checkingdict := make(map[string]int)

	for _, job := range jobs {
		var data []string
		err := DB.Table(job.TableName).Pluck(job.X, &data).Error

		if err != nil {
			http.Error(res, "Could not read from target table", http.StatusInternalServerError)
			return ""
		}

		/* Map all vars of this table and store it's count */
		for _, vars := range data {
			checkingdict[vars]++
		}
	}

	Combos := make([]PossibleCombo, 0)

	/* Build a dictionary of all 'strings' to be searched */
	Dictionary := ConvertIntoStructArrayAndSort(checkingdict)
	Amt := 0
	SizeLimit := 5
	for _, dict := range Dictionary {
		/* Why follwing thing is there? */
		if dict.Value < SizeLimit || dict.Key == "" {
			// Lets be sensible here ???Why @mayur/@glyn
			continue
		}

		Amt++

		if Amt > SizeLimit { // this acts as a "LIMIT 5" in the whole thing else this thing can literally takes mins to run.
			continue
		}

		search := StringSearch{}
		count := 0
		/* Why? */
		err := DB.Model(&search).Where("value = ?", dict.Value).Count(&count).Error

		check(err)

		if count != 0 {
			tablelist := make([]string, 0)

			var data = []string{}
			query := DB.Table("priv_onlinedata, priv_stringsearch, index")
			query = query.Where("priv_stringsearch.value = ?", dict.Value)
			query = query.Where("priv_stringsearch.count > ?", SizeLimit) //Why?
			query = query.Where("priv_stringsearch.tablename = priv_onlinedata.tablename")
			query = query.Where("priv_onlinedata.guid = index.guid")
			err := query.Pluck("priv_onlinedata.guid", &data).Error

			if err == gorm.RecordNotFound {
				continue
			} else if err != nil {
				http.Error(res, "Could not read off data lookups", http.StatusInternalServerError)
				return ""
			}

			for _, id := range data {
				if !StringInSlice(id, tablelist) {
					tablelist = append(tablelist, id)
				}
			}

			Combo := PossibleCombo{
				Match:  dict.Key,
				Tables: tablelist,
			}

			Combos = append(Combos, Combo)
		}
	}

	b, e := json.Marshal(Combos)
	if e != nil {
		http.Error(res, "JSON failed", http.StatusInternalServerError)
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

	sort.Sort(ByVal(in)) //What does this do?

	return in
}
