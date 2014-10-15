package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"net/http"
	"strconv"
)

type UserEdit struct {
	Id            int    `json:"id"`
	Avatar        string `json:"avatar"`
	Username      string `json:"username"`
	RepDifference int    `json:"repdifference"`
	Admin         int    `json:"admin"`
	Enabled       bool   `json:"enable"`
}

type UserReturn struct {
	Uid        int    `json:"uid"`
	Email      string `json:"md5email"`
	Avatar     string `json:"avatar"`
	Username   string `json:"username"`
	Reputation int    `json:"reputation"`
	Usertype   int    `json:"usertype"`
}

type UserReturnAndCount struct {
	Users []UserReturn `json:"users"`
	Count int          `json:"count"`
}

func GetUserTableHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	userReturn := []UserReturn{}

	order := params["order"] + " asc"

	e := DB.Model(User{}).Select("uid, email, avatar, username, reputation, usertype").Order(order).Scan(&userReturn).Error
	if e != nil {
		http.Error(res, "Unable to get users", http.StatusInternalServerError)
		return ""
	}

	offset, _ := strconv.Atoi(params["offset"])
	count, _ := strconv.Atoi(params["count"])
	if count > len(userReturn) || count == 0 {
		count = len(userReturn)
	}

	userReturn = userReturn[offset:count]

	for _, u := range userReturn {
		u.Email = GetMD5Hash(u.Email)
	}

	userReturnAndCount := UserReturnAndCount{userReturn, len(userReturn)}

	r, err := json.Marshal(userReturnAndCount)
	if err != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func EditUserHttp(res http.ResponseWriter, req *http.Request, userEdit UserEdit) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	if userEdit.Id <= 0 {
		http.Error(res, "No user id", http.StatusBadRequest)
		return ""
	}

	user := User{}
	rep := 0

	err := DB.Model(User{}).Where("uid = ?", userEdit.Id).Pluck("reputation", &rep).Error

	if err != nil {
		http.Error(res, "failed to get user's reputation", http.StatusBadRequest)
		return ""
	}

	user.Uid = userEdit.Id
	user.Avatar = userEdit.Avatar
	user.Username = userEdit.Username
	user.Reputation = rep + userEdit.RepDifference
	user.Usertype = userEdit.Admin
	user.Enabled = userEdit.Enabled

	err = DB.Save(&user).Error
	if err != nil {
		http.Error(res, "failed to update user", http.StatusBadRequest)
		return ""
	}

	return "success"
}
