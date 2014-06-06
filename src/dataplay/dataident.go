package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"net/http"
	"regexp"
	"sort"
	"strconv"
)

type CheckDict struct {
	Key   string
	value int
}

type IdentifyResponce struct {
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
	results := FetchTableCols(string(prams["id"]), Database.DB)

	returnobj := IdentifyResponce{
		Cols:    results,
		Request: prams["id"],
	}
	b, _ := json.Marshal(returnobj)
	return string(b)
}

// This fetches a array of all the col names and their types.
func FetchTableCols(guid string, database *sql.DB) (output []ColType) {
	if guid == "" {
		return output
	}

	var tablename string
	tablename, e := getRealTableName(guid, Database.DB, nil)
	if e != nil {
		return output
	}

	results := GetSQLTableSchema(tablename)

	return results
}

func GetSQLTableSchema(table string) []ColType {
	schema := make([]ColType, 0) // Setup the array that I will be append()ing to.

	rows, e := Database.DB.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_catalog = 'dataplay' AND table_name = $1", table)

	if e == nil {
		for rows.Next() {
			var column_name, data_type string
			rows.Scan(&column_name, &data_type)

			NewCol := ColType{
				Name:    column_name,
				Sqltype: data_type,
			}

			schema = append(schema, NewCol)
		}
	}

	return schema
}

// This is a shortcut function to do a the common action that is parsing the SQL create code.
// Regex looks for things that look like
// `colname` INT,
func BuildREArrayForCreateTable(input string) []string {
	re := ".*?(`.*?`).*?((?:[a-z][a-z]+))" // http://i.imgur.com/dkbyB.jpg
	// This regex looks for things that look like
	// `colname` INT,

	var sqlRE = regexp.MustCompile(re)
	results := sqlRE.FindStringSubmatch(input)
	return results
}

func SuggestColType(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	if prams["table"] == "" || prams["col"] == "" {
		http.Error(res, "There was no ID request", http.StatusBadRequest)
		return ""
	}

	var tablename string
	Database.DB.QueryRow("SELECT TableName FROM priv_onlinedata WHERE GUID = $1 LIMIT 1", prams["table"]).Scan(&tablename)
	if tablename == "" {
		http.Error(res, "Could not find that table", http.StatusNotFound)
		return ""
	}

	rows, e := Database.DB.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_catalog = 'dataplay' AND table_name = $1", tablename)
	if e != nil {
		http.Error(res, "Uhh, That table does not seem to acutally exist. This really should not happen. Check if someone have been messing around in the Database.", http.StatusBadRequest)
		return ""
	}

	if CheckIfColExists(rows, prams["col"]) {
		// Alrighty so I am now going to go though the whole table
		// and check what the data looks like
		// What that means for now is I am going to try and convert them all to ints and see if any of them breaks, If they do not, then I will suggest
		// that they be ints!
		rows, e := Database.DB.Query("SELECT $1 FROM $2", prams["col"], tablename)
		if e == nil {
			for rows.Next() {
				var TestSubject string
				rows.Scan(&TestSubject)
				_, e := strconv.ParseFloat(TestSubject, 10)
				if e != nil {
					return "false"
				}
			}

			return "true"
		}

		http.Error(res, fmt.Sprintf("Well somthing went wrong during the reading of that col, go and grab ben and show him this. %s", e), http.StatusInternalServerError)
		return ""
	} else {
		http.Error(res, "You have requested a col that does not exist. Please avoid doing this in the future.", http.StatusBadRequest)
		return "" // Shut up go
	}
}

func CheckIfColExists(rows *sql.Rows, column string) bool {
	for rows.Next() {
		var column_name, data_type string
		rows.Scan(&column_name, &data_type)

		if column_name == column {
			return true
		}
	}

	return false
}

// Unfinished function
func AttemptToFindMatches(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	// m.Get("/api/findmatches/:id/:x/:y", api.AttemptToFindMatches)
	RealTableName, e := getRealTableName(prams["id"], Database.DB, res)
	if e != nil {
		http.Error(res, "Could not find that table", http.StatusInternalServerError)
		return ""
	}

	rows, e := Database.DB.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_catalog = 'dataplay' AND table_name = $1", RealTableName)
	if e != nil {
		http.Error(res, "Uhh, That table does not seem to acutally exist. This really should not happen. Check if someone have been messing around in the Database.", http.StatusBadRequest)
		return ""
	}

	if !CheckIfColExists(rows, prams["x"]) || !CheckIfColExists(rows, prams["y"]) {
		http.Error(res, "Could not find the X or Y", http.StatusInternalServerError)
		return ""
	}

	// Now we need to check if it exists in the stats table. so we can compare its poly to other poly's
	HitCount := 0
	Database.DB.QueryRow("SELECT COUNT(*) FROM priv_statcheck WHERE table = $1 AND x = $2 AND y = $3", RealTableName, prams["x"], prams["y"]).Scan(&HitCount)

	if HitCount == 0 {
		http.Error(res, "Cannot find the poly code for that table x and y combo. It's probs not there because its not possible", http.StatusBadRequest)
		return ""
	}

	var id int = 0
	var table, x, y, p1, p2, p3, xstart, xend string

	Database.DB.QueryRow("SELECT * FROM priv_statcheck WHERE table = $1 AND x = $2 AND y = $3 LIMIT 1", RealTableName, prams["x"], prams["y"]).Scan(&id, &table, &x, &y, &p1, &p2, &p3, &xstart, &xend)
	Logger.Println(id, table, x, y, p1, p2, p3, xstart, xend)

	return "wat"
}

func FindStringMatches(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	if prams["word"] == "" {
		http.Error(res, "Please add a word", http.StatusBadRequest)
		return ""
	}

	Results := make([]StringMatchResult, 0)

	var name string
	var count int = 0
	if prams["x"] != "" {
		rows, e := Database.DB.Query("SELECT tablename, count FROM priv_stringsearch WHERE x = $1 AND value = $2", prams["x"], prams["word"])
		if e != nil {
			http.Error(res, "SQL error", http.StatusInternalServerError)
			return ""
		}

		for rows.Next() {
			rows.Scan(&name, &count)
			temp := StringMatchResult{
				Count: count,
				Match: name,
			}

			Results = append(Results, temp)
		}
	} else {
		rows, e := Database.DB.Query("SELECT tablename, count FROM priv_stringsearch WHERE value = $1", prams["word"])
		if e != nil {
			http.Error(res, "SQL error", http.StatusInternalServerError)
			return ""
		}
		for rows.Next() {
			rows.Scan(&name, &count)
			temp := StringMatchResult{
				Count: count,
				Match: name,
			}
			Results = append(Results, temp)
		}
	}

	b, e := json.Marshal(Results)
	if e != nil {
		http.Error(res, "Could not marshal JSON", http.StatusInternalServerError)
		return ""
	}

	return string(b)
}

func GetRelatedDatasetByStrings(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	RealTableName, e := getRealTableName(prams["guid"], Database.DB, res)
	if e != nil {
		http.Error(res, "Could not find that table", http.StatusInternalServerError)
		return ""
	}

	jobs := make([]ScanJob, 0)

	Bits := GetSQLTableSchema(RealTableName)

	for _, bit := range Bits {
		if bit.Sqltype == "character varying" {
			newJob := ScanJob{
				TableName: RealTableName,
				X:         bit.Name,
			}

			jobs = append(jobs, newJob)
		}
	}
	// Okay now we have a "job list" that looks a great deal like
	// the one in "makesearch_index" (Hint: Thats because its basically doing the same thing)
	checkingdict := make(map[string]int)

	for _, v := range jobs {
		q, e := Database.DB.Query("SELECT $1 FROM $2", v.X, v.TableName)

		if e != nil {
			http.Error(res, "Could not read from target table", http.StatusInternalServerError)
			return ""
		}
		// Count up all the vars in this col
		for q.Next() {
			var strout string
			q.Scan(&strout)
			checkingdict[strout]++
		}
	}

	Combos := make([]PossibleCombo, 0)
	Cdict := ConvertIntoStructArrayAndSort(checkingdict)
	var Amt int = 0
	// Okay so at this point checkingdict has all of the vals of the strings that we should search for.
	for _, v := range Cdict {
		if v.value < 5 || v.Key == "" {
			// Lets be sensible here
			continue
		}

		Amt++

		if Amt > 5 { // this acts as a "LIMIT 5" in the whole thing else this thing can literally takes mins to run.
			continue
		}

		var cnt int
		e := Database.DB.QueryRow("SELECT COUNT(*) FROM priv_stringsearch WHERE value = $1 LIMIT 1", v.value).Scan(&cnt)
		if e == nil && cnt != 0 {
			tablelist := make([]string, 0)
			r, e := Database.DB.Query("SELECT priv_onlinedata.GUID FROM priv_stringsearch, priv_onlinedata, index WHERE (value = $1) AND priv_stringsearch.tablename = priv_onlinedata.TableName AND priv_onlinedata.GUID = index.GUID AND priv_stringsearch.count > 5", v.value)
			if e != nil {
				http.Error(res, "Could not read off data lookups", http.StatusInternalServerError)
				return ""
			}

			res := ""
			for r.Next() {
				r.Scan(&res)
				if !stringInSlice(res, tablelist) {
					tablelist = append(tablelist, res)
				}
			}

			Combo := PossibleCombo{
				Match:  v.Key,
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

func stringInSlice(a string, list []string) bool {
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
func (a ByVal) Less(i, j int) bool { return a[i].value < a[j].value }

func ConvertIntoStructArrayAndSort(input map[string]int) (in []CheckDict) {
	in = make([]CheckDict, 0)
	for k, v := range input {
		newd := CheckDict{
			Key:   k,
			value: v,
		}
		in = append(in, newd)
	}
	sort.Sort(ByVal(in))
	return in
}
