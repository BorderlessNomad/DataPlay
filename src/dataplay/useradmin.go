package main

import (
	bcrypt "code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	"github.com/codegangsta/martini"
	"net/http"
	"strconv"
)

type UserEdit struct {
	Uid              int    `json:"uid"`
	Avatar           string `json:"avatar"`
	Email            string `json:"email"`
	Username         string `json:"username"`
	ReputationPoints int    `json:"reputationpoints"`
	Admin            int    `json:"admin"`
	Enabled          bool   `json:"enabled"`
	Password         string `json:"password"`
}

type UserReturn struct {
	Uid        int    `json:"uid"`
	MD5Email   string `json:"md5email"`
	Email      string `json:"email"`
	Avatar     string `json:"avatar"`
	Username   string `json:"username"`
	Reputation int    `json:"reputation"`
	Usertype   int    `json:"usertype"`
	Enabled    bool   `json:"enabled"`
}

type UserReturnAndCount struct {
	Users []UserReturn `json:"users"`
	Count int          `json:"count"`
}

type ObservationReturn struct {
	Comment       string `json:"comment"`
	Uid           int    `json:"uid"`
	Username      string `json:"username"`
	Flagged       bool   `json:"flagged"`
	ObservationId int    `json:"observationid"`
}

type ObsReturnAndCount struct {
	Observations []ObservationReturn `json:"comments"`
	Count        int                 `json:"count"`
}

func GetUserTableHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	userReturn := []UserReturn{}

	order := params["order"] + " asc"

	e := DB.Model(User{}).Select("uid, email, email, avatar, username, reputation, usertype, enabled").Order(order).Scan(&userReturn).Error
	if e != nil {
		http.Error(res, "Unable to get users", http.StatusInternalServerError)
		return ""
	}

	offset, _ := strconv.Atoi(params["offset"])
	count, _ := strconv.Atoi(params["count"])
	if offset+count > len(userReturn) || count == 0 {
		userReturn = userReturn[offset:len(userReturn)]
	} else {
		userReturn = userReturn[offset : offset+count]
	}

	for i, _ := range userReturn {
		userReturn[i].MD5Email = GetMD5Hash(userReturn[i].Email)
	}

	uCount := 0
	DB.Model(User{}).Count(&uCount)

	userReturnAndCount := UserReturnAndCount{userReturn, uCount}

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

	if userEdit.Uid <= 0 {
		http.Error(res, "No user id", http.StatusBadRequest)
		return ""
	}

	user := User{}

	err := DB.Model(User{}).Where("uid = ?", userEdit.Uid).Find(&user).Error

	if err != nil {
		http.Error(res, "failed to get user's reputation", http.StatusBadRequest)
		return ""
	}

	// fields to update
	if userEdit.Email != "" {
		user.Email = userEdit.Email
	}
	if userEdit.Avatar != "" {
		user.Avatar = userEdit.Avatar
	}
	if userEdit.Username != "" {
		user.Username = userEdit.Username
	}
	if userEdit.ReputationPoints != 0 {
		user.Reputation = user.Reputation + userEdit.ReputationPoints
	}
	if userEdit.Admin != user.Usertype {
		user.Usertype = userEdit.Admin
	}
	if userEdit.Enabled != user.Enabled {
		user.Enabled = userEdit.Enabled
	}

	if userEdit.Password != "" { // generate whatever password has been passed

		hashedPassword, err1 := bcrypt.GenerateFromPassword([]byte(userEdit.Password), bcrypt.DefaultCost)
		if err1 != nil {
			http.Error(res, "Unable to generate password hash.", http.StatusInternalServerError)
			return ""
		}
		user.Password = string(hashedPassword)
		err = DB.Save(&user).Error // update or add record
		if err != nil {
			http.Error(res, "failed to update user", http.StatusBadRequest)
			return ""
		}
	} else { // do not update password field
		err = DB.Save(&user).Error
		if err != nil {
			http.Error(res, "failed to update user", http.StatusBadRequest)
			return ""
		}

	}

	return "success"
}

func GetObservationsTableHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	observationReturn := []ObservationReturn{}

	ob := Observation{}
	u := User{}
	uCount := 0
	joinStr := "JOIN " + u.TableName() + " ON " + u.TableName() + ".uid = " + ob.TableName() + ".uid"
	selectStr := "comment, " + ob.TableName() + ".uid, username, flagged, observation_id"
	order := params["order"] + " asc"

	if params["flagged"] == "true" {
		e := DB.Model(ob).Select(selectStr).Joins(joinStr).Order(order).Where("flagged = ?", true).Scan(&observationReturn).Error
		if e != nil {
			http.Error(res, "Unable to get observations", http.StatusInternalServerError)
			return ""
		}
		DB.Model(Observation{}).Where("flagged = ?", true).Count(&uCount)
	} else {
		e := DB.Model(ob).Select(selectStr).Joins(joinStr).Order(order).Scan(&observationReturn).Error
		if e != nil {
			http.Error(res, "Unable to get observations", http.StatusInternalServerError)
			return ""
		}
		DB.Model(Observation{}).Count(&uCount)
	}

	offset, _ := strconv.Atoi(params["offset"])
	count, _ := strconv.Atoi(params["count"])
	if offset+count > len(observationReturn) || count == 0 {
		observationReturn = observationReturn[offset:len(observationReturn)]
	} else {
		observationReturn = observationReturn[offset : offset+count]
	}

	obsReturnAndCount := ObsReturnAndCount{observationReturn, uCount}

	r, err := json.Marshal(obsReturnAndCount)
	if err != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func DeleteObservationHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	if params["id"] == "" {
		http.Error(res, "Missing observation id", http.StatusBadRequest)
		return ""
	}

	oid, _ := strconv.Atoi(params["id"])
	observation := Observation{}

	e := DB.Where("observation_id = ?", oid).Find(&observation).Error

	if e != nil {
		http.Error(res, "Unable to find observation", http.StatusInternalServerError)
		return ""
	}

	uid := observation.Uid

	e = DB.Where("observation_id = ?", oid).Delete(&observation).Error
	if e != nil {
		http.Error(res, "Unable to delete observation", http.StatusInternalServerError)
	}

	Reputation(uid, obsSpam) // knock off some points

	return "success"
}
