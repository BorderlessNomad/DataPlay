package main

import (
	"database/sql"
	"regexp"
	"strings"
)

type IdentifyResponse struct {
	Cols    []ColType
	Request string
}

type ColType struct {
	Name    string
	Sqltype string
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
