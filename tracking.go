package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func GetLastVisited(rw http.ResponseWriter, req *http.Request) string {
	database := GetDB()
	defer database.Close()
	if GetUserID(rw, req) != 0 {
		value := string(GetUserID(rw, req))
		rows, e := database.Query("SELECT DISTINCT(guid),(SELECT Title FROM `index` WHERE `index`.GUID = priv_tracking.guid LIMIT 1) as a FROM priv_tracking WHERE user = ? ORDER BY id DESC LIMIT 5", value)
		result := make([][]string, 0)
		if e == nil {

			// Read out the rows now we know that nothing went wrong.
			for rows.Next() {
				var guid string
				var title string

				rows.Scan(&guid, &title)

				r := HasTableGotLocationData(guid, database)
				result2 := []string{
					guid,
					title,
					r,
				}

				result = append(result, result2)
			}

		} else {
			fmt.Println(e)
		}
		b, _ := json.Marshal(result)
		return (string(b))
	}
	return ""
}

func HasTableGotLocationData(datasetGUID string, database *sql.DB) string {
	cols := FetchTableCols(datasetGUID, database)

	if containsTableCol(cols, "lat") && (containsTableCol(cols, "lon") || containsTableCol(cols, "long")) {
		return "true"
	}
	return "false"
}

func containsTableCol(cols []ColType, target string) bool {
	for _, v := range cols {
		if strings.ToLower(v.Name) == target {
			return true
		}
	}
	return false
}

func TrackVisited(guid string, user string) {
	database := GetDB()
	defer database.Close()
	_, e := database.Exec("INSERT INTO `DataCon`.`priv_tracking` (`user`, `guid`) VALUES (?, ?);", user, guid)
	if e != nil {
		Logger.Println(e)
	}
	Logger.Println("Tracking page hit to ", guid, "by", user)
}
