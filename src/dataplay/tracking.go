package main

import (
	"database/sql"
	"encoding/json"
	// "fmt"
	"net/http"
	"strings"
)

func GetLastVisited(res http.ResponseWriter, req *http.Request) string {
	uid := GetUserID(res, req)
	data := make([][]string, 0)

	if uid != 0 {
		/* Anonymous struct for storing results */
		results := []struct {
			Tracking
			Title string
		}{}

		err := DB.Select("DISTINCT ON (priv_tracking.guid) guid, priv_tracking.id, (SELECT index.title FROM index WHERE index.guid = priv_tracking.guid LIMIT 1) as title").Where("priv_tracking.user = ?", uid).Order("guid desc").Order("priv_tracking.id desc").Limit(5).Find(&results).Error
		if err != nil {
			panic(err)
		}

		for _, result := range results {
			r := HasTableGotLocationData(result.Guid, DB.SQL)

			data = append(data, []string{
				result.Guid,
				result.Title,
				r,
			})
		}
	}

	/* We ALWAYS return something [[], [], ...] or [] */
	d, _ := json.Marshal(data)

	return string(d)
}

func HasTableGotLocationData(datasetGUID string, database *sql.DB) string {
	cols := FetchTableCols(datasetGUID, DB.SQL)

	if ContainsTableCol(cols, "lat") && (ContainsTableCol(cols, "lon") || ContainsTableCol(cols, "long")) {
		return "true"
	}

	return "false"
}

func ContainsTableCol(cols []ColType, target string) bool {
	for _, v := range cols {
		if strings.ToLower(v.Name) == target {
			return true
		}
	}

	return false
}

func TrackVisited(guid string, user string) {
	tracking := Tracking{
		User: user,
		Guid: guid,
	}

	err := DB.Save(&tracking).Error
	if err != nil {
		Logger.Println(err)
	}

	Logger.Println("Tracking page hit to:", tracking.Guid, "by user:", tracking.User, "[ #", tracking.Id, "]")
}
