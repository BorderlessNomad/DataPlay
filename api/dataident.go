package api

import (
	msql "../databasefuncs"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

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

func IdentifyTable(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	// This function checks to see if the data has been imported yet or still is in need of importing
	database := msql.GetDB()
	defer database.Close()
	if prams["id"] == "" {
		http.Error(res, "There was no ID request", http.StatusBadRequest)
		return ""
	}
	results := FetchTableCols(string(prams["id"]), database)

	returnobj := IdentifyResponce{
		Cols:    results,
		Request: prams["id"],
	}
	b, _ := json.Marshal(returnobj)
	return string(b)
}

func FetchTableCols(guid string, database *sql.DB) (output []ColType) {
	if guid == "" {
		return output
	}

	var tablename string
	database.QueryRow("SELECT TableName FROM `priv_onlinedata` WHERE GUID = ? LIMIT 1", guid).Scan(&tablename)
	if tablename == "" {
		return output
	}

	var createcode string
	database.QueryRow("SHOW CREATE TABLE "+tablename).Scan(&tablename, &createcode)
	if createcode == "" {
		return output
	}
	results := ParseCreateTableSQL(createcode)
	return results
}

func BuildREArrayForCreateTable(input string) []string {
	re := ".*?(`.*?`).*?((?:[a-z][a-z]+))" // http://i.imgur.com/dkbyB.jpg
	// This regex looks for things that look like
	// `colname` INT,

	var sqlRE = regexp.MustCompile(re)
	results := sqlRE.FindStringSubmatch(input)
	return results
}

func ParseCreateTableSQL(input string) []ColType {
	returnerr := make([]ColType, 0) // Setup the array that I will be append()ing to.
	SQLLines := strings.Split(input, "\n")
	// The mysql server gives you the SQL create code formatted. So I exploit this by
	// using it to split the system up by \n

	for c, line := range SQLLines {
		if c != 0 && strings.HasPrefix(strings.TrimSpace(line), "`") { // Clipping off the create part since its useless for me.
			results := BuildREArrayForCreateTable(line)
			if len(results) == 3 {
				// We expect there to be 3 matches from the Regex, if not then we probs don't
				// have what we want
				DeQuoted := strings.Replace(results[1], "`", "", -1)
				NewCol := ColType{
					Name:    DeQuoted,
					Sqltype: results[2],
				}
				returnerr = append(returnerr, NewCol)
			}
		}
	}
	return returnerr
}

func SuggestColType(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	database := msql.GetDB()
	defer database.Close()
	if prams["table"] == "" || prams["col"] == "" {
		http.Error(res, "There was no ID request", http.StatusBadRequest)
		return ""
	}

	var tablename string
	database.QueryRow("SELECT TableName FROM `priv_onlinedata` WHERE GUID = ? LIMIT 1", prams["table"]).Scan(&tablename)
	if tablename == "" {
		http.Error(res, "Could not find that table", http.StatusNotFound)
		return ""
	}

	var createcode string
	database.QueryRow("SHOW CREATE TABLE "+tablename).Scan(&tablename, &createcode)
	if createcode == "" {
		http.Error(res, `Uhh, That table does not seem to acutally exist.
		this really should not happen. 
		Check if someone have been messing around in the database.`, http.StatusBadRequest)
		return ""
	}
	if CheckIfColExists(createcode, prams["col"]) {
		// Alrighty so I am now going to go though the whole table
		// and check what the data looks like
		// What that means for now is I am going to try and convert them all to ints and see if any of them breaks, If they do not, then I will suggest
		// that they be ints!
		rows, e := database.Query(fmt.Sprintf("SELECT `%s` FROM `%s`", prams["col"], tablename))
		if e == nil {
			for rows.Next() {
				var TestSubject string
				rows.Scan(&TestSubject)
				_, e := strconv.ParseInt(TestSubject, 10, 64)
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
	return "This isnt suppose to happen"
}

func CheckIfColExists(createcode string, targettable string) bool {

	SQLLines := strings.Split(createcode, "\n")

	for c, line := range SQLLines {
		if c != 0 { // Clipping off the create part since its useless for me.
			results := BuildREArrayForCreateTable(line)
			if len(results) == 3 { // 3 is the amount you would expect the regex to match. Even though we only use the first part
				if results[1] == "`"+targettable+"`" {
					return true
				}
			}
		}
	}
	return false
}

func AttemptToFindMatches(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	// m.Get("/api/findmatches/:id/:x/:y", api.AttemptToFindMatches)
	database := msql.GetDB()
	defer database.Close()
	RealTableName := getRealTableName(prams["id"], database, res)
	if RealTableName == "Error" {
		http.Error(res, "Could not find that table", http.StatusInternalServerError)
		return ""
	}

	CCode := ""
	database.QueryRow(fmt.Sprintf("SHOW CREATE TABLE `%s`;", RealTableName)).Scan(&RealTableName, &CCode)
	if !CheckIfColExists(CCode, prams["x"]) || !CheckIfColExists(CCode, prams["y"]) {
		http.Error(res, "Could not find the X or Y", http.StatusInternalServerError)
		return ""
	}

	// Now we need to check if it exists in the stats table. so we can compare its poly to other poly's
	HitCount := 0
	database.QueryRow("SELECT COUNT(*) FROM `priv_statcheck` WHERE `table` = ? AND `x` = ? AND `y` = ?", RealTableName, prams["x"], prams["y"]).Scan(&HitCount)

	if HitCount == 0 {
		http.Error(res, "Cannot find the poly code for that table x and y combo. It's probs not there because its not possible", http.StatusBadRequest)
		return ""
	}

	id := 0
	table := ""
	x := ""
	y := ""
	p1 := ""
	p2 := ""
	p3 := ""
	xstart := ""
	xend := ""

	database.QueryRow("SELECT * FROM `priv_statcheck` WHERE `table` = ? AND `x` = ? AND `y` = ? LIMIT 1", RealTableName, prams["x"], prams["y"]).Scan(&id, &table, &x, &y, &p1, &p2, &p3, &xstart, &xend)
	fmt.Println(id, table, x, y, p1, p2, p3, xstart, xend)

	return "wat"
}

func FindStringMatches(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	database := msql.GetDB()
	defer database.Close()

	if prams["word"] == "" {
		http.Error(res, "Please add a word", http.StatusBadRequest)
		return ""
	}

	Results := make([]StringMatchResult, 0)

	name := ""
	count := 0
	if prams["x"] != "" {
		rows, e := database.Query("SELECT tablename,count FROM priv_stringsearch WHERE x = ? AND value = ?", prams["x"], prams["word"])
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
		rows, e := database.Query("SELECT tablename,count FROM priv_stringsearch WHERE value = ?", prams["word"])
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
	database := msql.GetDB()
	defer database.Close()

	RealTableName := getRealTableName(prams["guid"], database, res)
	if RealTableName == "Error" {
		http.Error(res, "Could not find that table", http.StatusInternalServerError)
		return ""
	}

	jobs := make([]ScanJob, 0)
	var CreateSQL string
	database.QueryRow(fmt.Sprintf("SHOW CREATE TABLE `DataCon`.`%s`", RealTableName)).Scan(&RealTableName, &CreateSQL)

	Bits := ParseCreateTableSQL(CreateSQL)
	for _, bit := range Bits {
		if bit.Sqltype == "varchar" {
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
		q, e := database.Query(fmt.Sprintf("SELECT `%s` FROM `%s`", v.X, v.TableName))

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
	// Okay so at this point checkingdict has all of the vals of the strings that we should search for.
	for k, v := range checkingdict {
		if v < 5 || k == "" {
			// Lets be sensible here
			continue
		}
		var cnt int
		e := database.QueryRow("SELECT COUNT(*) FROM priv_stringsearch WHERE value = ? LIMIT 1", k).Scan(&cnt)
		if e == nil && cnt != 0 {
			tablelist := make([]string, 0)
			r, e := database.Query("SELECT `priv_onlinedata`.GUID FROM priv_stringsearch, priv_onlinedata, `index` WHERE (value = ?) AND `priv_stringsearch`.tablename = `priv_onlinedata`.TableName AND `priv_onlinedata`.GUID = `index`.GUID AND priv_stringsearch.count > 2", k)
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
				Match:  k,
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
