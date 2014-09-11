package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"net/http"
	"time"
)

type HomeData struct {
	Label string
	Value int
}

func ActivityCheck(a string) string {
	switch a {
	case "c":
		return "Comment"
	case "ic":
		return "Invalidated Chart"
	case "vc":
		return "Validated Chart"
	case "io":
		return "Invalidated Observation"
	case "vo":
		return "Validated Observation"
	default:
		return "Unknown"
	}
}

func AddActivity(uid int, atype string, ts time.Time) *appError {
	act := Activity{
		Uid:     uid,
		Type:    ActivityCheck(atype),
		Created: ts,
	}

	err := DB.Save(&act).Error
	if err != nil {
		return &appError{err, "Database query failed (Save)", http.StatusInternalServerError}
	}

	return nil
}

func GetProfileObservationsHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return "Could not validate user"
	}

	var comments []string
	err1 := DB.Model(Observations{}).Select("comment").Where("uid = ?", uid).Find(&comments).Error
	if err1 != nil {
		return "not found"
	}

	r, err2 := json.Marshal(comments)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return "Unable to parse JSON"
	}

	return string(r)
}

func GetDiscoveriesHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return "Could not validate user"
	}

	var patterns []int
	err1 := DB.Model(Discovered{}).Select("discovered_id").Where("uid = ?", uid).Find(&patterns).Error
	if err1 != nil {
		return "not found"
	}

	r, err2 := json.Marshal(patterns)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return "Unable to parse JSON"
	}

	return string(r)
}

func GetValidatedDiscoveriesHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return "Could not validate user"
	}

	var patterns []int
	err1 := DB.Model(Discovered{}).Select("discovered_id").Where("uid = ?", uid).Where("valid > ?", 0).Find(&patterns).Error
	if err1 != nil {
		return "not found"
	}

	r, err2 := json.Marshal(patterns)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return "Unable to parse JSON"
	}

	return string(r)
}

func GetHomePageDataHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	var result [3]HomeData
	result[0].Label = "players"
	result[1].Label = "discoveries"
	result[2].Label = "datasets"
	err := DB.Model(User{}).Count(result[0].Value).Error
	if err != nil {
		return "not found"
	}

	err = DB.Model(Discovered{}).Count(result[1].Value).Error
	if err != nil {
		return "not found"
	}

	err = DB.Model(OnlineData{}).Count(result[2].Value).Error
	if err != nil {
		return "not found"
	}
	r, err2 := json.Marshal(result)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return "Unable to parse JSON"
	}

	return string(r)
}

func GetReputationHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return "Could not validate user"
	}

	var rep int
	err1 := DB.Model(User{}).Select("reputation").Where("uid = ?", uid).Find(&rep).Error
	if err1 != nil {
		return "not found"
	}

	return string(rep)
}

func GetAmountDiscoveriesHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return "Could not validate user"
	}

	count := 0
	err1 := DB.Model(Discovered{}).Where("uid = ?", uid).Count(&count).Error
	if err1 != nil {
		return "not found"
	}

	return string(count)
}
