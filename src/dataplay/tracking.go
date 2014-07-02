package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func GetLastVisitedHttp(res http.ResponseWriter, req *http.Request) string {
	uid := GetUserID(res, req)

	result, error := GetLastVisited(uid)
	if error != nil {
		http.Error(res, error.Message, error.Code)
		return ""
	}

	r, err := json.Marshal(result)
	if err != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	/* We ALWAYS return something [[], [], ...] or [] */
	return string(r)
}

func GetLastVisitedQ(params map[string]string) string {
	if params["user"] == "" {
		return ""
	}

	uid, error := strconv.Atoi(params["user"])
	if error != nil {
		return ""
	}

	result, error := GetLastVisited(uid)
	if error != nil {
		return ""
	}

	r, error := json.Marshal(result)
	if error != nil {
		return ""
	}

	return string(r)
}

func GetLastVisited(uid int) ([][]string, *appError) {
	data := make([][]string, 0)

	if uid != 0 {
		/* Anonymous struct for storing results */
		results := []struct {
			Tracking
			Title string
		}{}

		err := DB.Select("DISTINCT ON (priv_tracking.guid) guid, priv_tracking.id, (SELECT index.title FROM index WHERE index.guid = priv_tracking.guid LIMIT 1) as title").Where("priv_tracking.user = ?", uid).Order("guid desc").Order("priv_tracking.id desc").Limit(5).Find(&results).Error

		if err != nil && err != gorm.RecordNotFound {
			return nil, &appError{err, "Database query failed", http.StatusServiceUnavailable}
		}

		for _, result := range results {
			r := HasTableGotLocationData(result.Guid)

			data = append(data, []string{
				result.Guid,
				result.Title,
				r,
			})
		}
	}

	return data, nil
}

func HasTableGotLocationData(datasetGUID string) string {
	cols := FetchTableCols(datasetGUID)

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
