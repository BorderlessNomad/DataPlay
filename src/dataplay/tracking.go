package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type VisitedForm struct {
	Guid json.RawMessage `json:"guid"`
	Info json.RawMessage `json:"info"`
}

func GetLastVisitedHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
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

func GetLastVisited(uid int) ([]interface{}, *appError) {
	data := make([]interface{}, 0)
	var err error

	if uid != 0 {
		/* Anonymous struct for storing results */
		results := []struct {
			Tracking
			Title string
		}{}

		query := DB.Select("MAX (priv_tracking.id) id, priv_tracking.guid, (SELECT index.title FROM index WHERE index.guid = priv_tracking.guid LIMIT 1) as title, MAX (priv_tracking.created) created")
		query = query.Joins("LEFT JOIN index ON index.guid = priv_tracking.guid")
		query = query.Where("title != ?", "")
		query = query.Where("priv_tracking.user = ?", uid)
		query = query.Group("priv_tracking.guid")
		query = query.Order("created DESC")
		query = query.Order("priv_tracking.guid DESC")
		query = query.Limit(5)
		err = query.Find(&results).Error
		// err = DB.Select("MAX (priv_tracking.id) id, priv_tracking.guid, (SELECT index.title FROM index WHERE index.guid = priv_tracking.guid LIMIT 1) as title, MAX (priv_tracking.created) created").Joins("LEFT JOIN index ON index.guid = priv_tracking.guid").Where("title != ?", "").Where("priv_tracking.user = ?", uid).Group("priv_tracking.guid").Order("created DESC").Order("priv_tracking.guid DESC").Limit(5).Find(&results).Error

		if err != nil && err != gorm.RecordNotFound {
			return nil, &appError{err, "Database query failed (Select Tracking)", http.StatusInternalServerError}
		}

		trackingInfo := []TrackingInfo{}
		loadInfo := false
		if len(results) > 0 {
			loadInfo = true
			ids := make([]int, 0)
			for _, record := range results {
				ids = append(ids, record.Id)
			}

			err = DB.Select("id, info").Where(ids).Order("id DESC").Find(&trackingInfo).Error
			if err != nil && err != gorm.RecordNotFound {
				return nil, &appError{err, "Database query failed (Select Info)", http.StatusInternalServerError}
			}
		}

		for i, result := range results {
			var info map[string]interface{}
			if loadInfo {
				e := json.Unmarshal(trackingInfo[i].Info, &info)
				if e != nil {
					Logger.Println("Info Parse Error", e)
				}
			}

			r := HasTableGotLocationData(result.Guid)

			data = append(data, map[string]interface{}{
				"guid":  SanitizeString(result.Guid),
				"title": SanitizeString(result.Title),
				"info":  info,
				"map":   r,
			})
		}
	}

	return data, nil
}

func TrackVisitedHttp(res http.ResponseWriter, req *http.Request, visited VisitedForm) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	uid, err1 := GetUserID(session)
	if err1 != nil {
		http.Error(res, err1.Message, err1.Code)
		return ""
	}

	var guid string
	err_guid := json.Unmarshal(visited.Guid, &guid)
	if err_guid != nil {
		http.Error(res, "Unable to parse JSON (GUID)", http.StatusInternalServerError)
		return ""
	}

	var info_map map[string]interface{}
	err_info_map := json.Unmarshal(visited.Info, &info_map)
	if err_info_map != nil {
		http.Error(res, "Unable to parse JSON (INFO Map)", http.StatusInternalServerError)
		return ""
	}

	info, err_info := json.Marshal(info_map)
	if err_info != nil {
		http.Error(res, "Unable to parse JSON (INFO)", http.StatusInternalServerError)
		return ""
	}

	err2 := TrackVisited(uid, guid, string(info))
	if err2 != nil {
		http.Error(res, err2.Message, err2.Code)
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

func TrackVisited(user int, guid string, info string) *appError {
	tracking := Tracking{
		User:    user,
		Guid:    guid,
		Info:    info,
		Created: time.Now(),
	}

	err := DB.Save(&tracking).Error
	if err != nil {
		return &appError{err, "Database query failed (Save)", http.StatusInternalServerError}
	}

	Logger.Println("Tracking page hit to:", tracking.Guid, "by user:", tracking.User, "[ #", tracking.Id, "]")

	return nil
}
