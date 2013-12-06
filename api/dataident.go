package api

import (
	msql "../databasefuncs"
	"encoding/json"
	// "fmt"
	"github.com/codegangsta/martini"
	"net/http"
	"regexp"
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
	}

	var tablename string
	database.QueryRow("SELECT TableName FROM `priv_onlinedata` WHERE GUID = ? LIMIT 1", prams["id"]).Scan(&tablename)
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
	results := ParseCreateTableSQL(createcode)

	returnobj := IdentifyResponce{
		Cols:    results,
		Request: prams["id"],
	}
	b, _ := json.Marshal(returnobj)
	return string(b[:])
}

func ParseCreateTableSQL(input string) []ColType {
	re := ".*?(`.*?`).*?((?:[a-z][a-z]+))" // http://i.imgur.com/dkbyB.jpg
	var sqlRE = regexp.MustCompile(re)
	returnerr := make([]ColType, 0) // Setup the array that I will be append()ing to.
	SQLLines := strings.Split(input, "\n")
	for c, line := range SQLLines {
		if c != 0 { // Clipping off the create part since its useless for me.
			results := sqlRE.FindStringSubmatch(line)
			if len(results) == 3 {
				NewCol := ColType{
					Name:    results[1],
					Sqltype: results[2],
				}
				returnerr = append(returnerr, NewCol)
			}
		}
	}
	return returnerr
}
