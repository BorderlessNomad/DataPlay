package main

import (
	"database/sql"
	"encoding/json"
	// "fmt"
	"net/http"
	"strings"
)

type Tracking struct {
	Id   int `primaryKey:"yes"`
	User string
	Guid string
}

func (t Tracking) TableName() string {
	return "priv_tracking"
}

type Index struct {
	Guid    int `primaryKey:"yes"`
	Name    string
	Title   string
	Notes   string
	CkanUrl string
	Owner   int
}

func (i Index) TableName() string {
	return "index"
}

func GetLastVisited(res http.ResponseWriter, req *http.Request) string {
	uid := GetUserID(res, req)
	data := make([][]string, 0)

	if uid != 0 {
		/* Anonymous struct for storing results */
		results := []struct {
			Tracking
			Title string
		}{}

		err := DB.Select("DISTINCT(priv_tracking.guid), priv_tracking.id, (SELECT index.title FROM index WHERE index.guid = priv_tracking.guid LIMIT 1) as title").Where("priv_tracking.user = ?", uid).Order("priv_tracking.id desc").Limit(5).Find(&results).Error
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
	tracking := Tracking{
		User: user,
		Guid: guid,
	}

	err := DB.Save(&tracking).Error
	if err != nil {
		Logger.Println(err)
	}

	Logger.Println("Tracking page hit to:", tracking.Guid, "by user:", tracking.User)
}
