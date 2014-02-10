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

type SuggestionResponce struct {
	Request string
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
			if len(results) == 3 {
				if results[1] == "`"+targettable+"`" {
					return true
				}
			}
		}
	}
	return false
}

func AttemptToFindMatches(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	database := msql.GetDB()
	defer database.Close()
	// m.Get("/api/findmatches/:id/:x/:y", api.AttemptToFindMatches)

	return "wat"
}
