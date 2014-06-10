package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func GetLastVisited(rw http.ResponseWriter, req *http.Request) string {
	uid := GetUserID(rw, req)

	if uid != 0 {

		query := `
		SELECT
			DISTINCT(t.guid),
			t.id,
			(
				SELECT
					i.title
				FROM
					index as i
				WHERE
					i.guid = t.guid
				LIMIT 1
			) as title
		FROM
			priv_tracking as t
		WHERE
			t.user = $1
		ORDER BY
			t.id DESC
		LIMIT 5`

		rows, e := DB.SQL.Query(query, uid)

		result := make([][]string, 0)

		if e == nil {
			// Read out the rows now we know that nothing went wrong.
			for rows.Next() {
				var guid string
				var title string

				rows.Scan(&guid, &title)

				r := HasTableGotLocationData(guid, DB.SQL)
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
	_, e := DB.SQL.Exec("INSERT INTO priv_tracking (user, guid) VALUES ($1, $2)", user, guid)
	if e != nil {
		Logger.Println(e)
	}

	Logger.Println("Tracking page hit to ", guid, "by", user)
}
