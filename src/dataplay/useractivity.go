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

func AddCommentHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return "Missing session parameter."
	}

	u, err1 := GetUserID(session)
	if err1 != nil {
		http.Error(res, err1.Message, err1.Code)
		return "Could not validate user"
	}

	t := time.Now()

	err2 := AddActivity(u, "c", t)
	if err2 != nil {
		http.Error(res, err2.Message, err2.Code)
		return err2.Message
	}

	var id []int
	err = DB.Model(Activity{}).Where("created = ?", t).Where("uid = ?", u).Where("type = ?", "Comment").Pluck("activity_id", &id).Error
	if err != nil {
		http.Error(res, "No activity found", http.StatusBadRequest)
		return "No activity found"
	}

	c := Comment{}
	c.Comment = params["comment"]
	c.ActivityId = id[0]

	err = DB.Save(&c).Error
	if err != nil {
		return "Database query failed"
	}

	activityStr := strconv.Itoa(id[0])
	return "Comment " + activityStr + " added successfully"
}
