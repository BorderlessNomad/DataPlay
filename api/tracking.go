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
	value := sess.Value.(string)
	rows, e := database.Query("SELECT DISTINCT(guid) FROM priv_tracking WHERE user = ? ORDER BY id DESC LIMIT 5", value)
	result := make([]string, 0)
	if e == nil {
		for rows.Next() {
			var guid string
			rows.Scan(&guid)
			result = append(result, guid)
		}
	}
	if e != nil {
		fmt.Println(e)
	}
	b, _ := json.Marshal(result)
	return (string(b))
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
