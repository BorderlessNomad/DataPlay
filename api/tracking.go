package api

import (
	msql "../databasefuncs"
	"encoding/json"
	"fmt"
	"github.com/mattn/go-session-manager"
	"net/http"
)

func GetLastVisited(rw http.ResponseWriter, req *http.Request, monager *session.SessionManager) string {
	database := msql.GetDB()
	defer database.Close()
	sess := monager.GetSession(rw, req)
	if sess.Value != nil {
		value := sess.Value.(string)
		rows, e := database.Query("SELECT DISTINCT(guid),(SELECT Title FROM `index` WHERE `index`.GUID = priv_tracking.guid LIMIT 1) as a FROM priv_tracking WHERE user = ? ORDER BY id DESC LIMIT 5", value)
		result := make([][]string, 0)
		if e == nil {
			for rows.Next() {
				var guid string
				var title string
				var status string
				rows.Scan(&guid, &title)
				result2 := make([]string, 0)
				status = HasTableGotLocationData(guid)
				result2 = append(result2, guid)
				result2 = append(result2, title)
				result2 = append(result2, status)

				result = append(result, result2)
			}
		}
		if e != nil {
			fmt.Println(e)
		}
		b, _ := json.Marshal(result)
		return (string(b))
	}
	return ""
}

func HasTableGotLocationData(datasetGUID string) string {
	cols := FetchTableCols(datasetGUID)

	for _, col1 := range cols {
		if col1.Name == "lat" {
			for _, col2 := range cols {
				if col2.Name == "lon" || col2.Name == "long" {
					return "true"
				}
			}
		}
	}
	return "false"
}

func TrackVisited(guid string, user string) {
	database := msql.GetDB()
	defer database.Close()
	_, e := database.Exec("INSERT INTO `DataCon`.`priv_tracking` (`user`, `guid`) VALUES (?, ?);", user, guid)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println("Tracking page hit to ", guid, "by", user)
}
