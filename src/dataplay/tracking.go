package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

func GetLastVisitedHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := params["session"]
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	result, err := GetLastVisited(uid)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	r, e := json.Marshal(result)
	if e != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func GetLastVisitedQ(params map[string]string) string {
	if params["user"] == "" {
		return ""
	}

	uid, e := strconv.Atoi(params["user"])
	if e != nil {
		return ""
	}

	result, err := GetLastVisited(uid)
	if err != nil {
		return ""
	}

	r, e := json.Marshal(result)
	if e != nil {
		return ""
	}

	return string(r)
}

func GetLastVisited(uid int) ([]interface{}, *appError) {
	data := make([]interface{}, 0)

	if uid != 0 {
		/* Anonymous struct for storing results */
		results := []struct {
			Tracking
			Title string
		}{}

		err := DB.Select("MAX (priv_tracking.id) id, priv_tracking.guid, (SELECT index.title FROM index WHERE index.guid = priv_tracking.guid LIMIT 1) as title, MAX (priv_tracking.created) created").Joins("LEFT JOIN index ON index.guid = priv_tracking.guid").Where("title != ?", "").Where("priv_tracking.user = ?", uid).Group("priv_tracking.guid").Order("created DESC").Order("priv_tracking.guid DESC").Limit(5).Find(&results).Error

		if err != nil && err != gorm.RecordNotFound {
			return nil, &appError{err, "Database query failed", http.StatusInternalServerError}
		}

		for _, result := range results {
			r := HasTableGotLocationData(result.Guid)

			data = append(data, map[string]interface{}{
				"guid":  SanitizeString(result.Guid),
				"title": SanitizeString(result.Title),
				"map":   r,
			})
		}
	}

	return data, nil
}

func TrackVisited(guid string, user int) {
	tracking := Tracking{
		User:    user,
		Guid:    guid,
		Created: time.Now(),
	}

	err := DB.Save(&tracking).Error
	if err != nil {
		Logger.Println(err)
	}

	Logger.Println("Tracking page hit to:", tracking.Guid, "by user:", tracking.User, "[ #", tracking.Id, "]")
}
