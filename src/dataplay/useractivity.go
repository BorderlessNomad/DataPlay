package main

import (
	"github.com/codegangsta/martini"
	"net/http"
	"strconv"
	"time"
)

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

func AddActivityHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return "Missing session parameter."
	}

	u, err := strconv.Atoi(params["uid"])
	if err != nil {
		http.Error(res, "Invalid uid.", http.StatusBadRequest)
		return "Invalid uid"
	}

	a := ActivityCheck(params["type"])
	if a == "Unknown" {
		http.Error(res, "Unknown activity type.", http.StatusBadRequest)
		return "Unknown activity type"
	}

	t := time.Now()

	err2 := AddActivity(u, params["type"], t)
	if err2 != nil {
		http.Error(res, err2.Message, err2.Code)
		return err2.Message
	}

	var id []int
	err = DB.Table("priv_activity").Where("created = ?", t).Where("uid = ?", u).Where("type = ?", a).Pluck("id", &id).Error
	if err != nil {
		http.Error(res, "No activity found", http.StatusBadRequest)
		return "No activity found"
	}

	activityStr := strconv.Itoa(id[0])
	return "Activity " + activityStr + " added successfully"
}

func AddActivityQ(params map[string]string) string {
	u, err := strconv.Atoi(params["uid"])
	if err != nil {
		return "bad uid"
	}

	a := ActivityCheck(params["type"])
	if a == "Unknown" {
		return a
	}
	t := time.Now()

	err2 := AddActivity(u, params["type"], t)
	if err2 != nil {
		return err2.Message
	}

	var actid int
	err = DB.Table("priv_activities").Where("date = ?", t).Where("uid = ?", u).Where("type = ?", a).Pluck("activityid", &actid).Error
	if err != nil {
		return "No activity found"
	}

	activityStr := strconv.Itoa(actid)
	return "Activity " + activityStr + " added successfully"
}
