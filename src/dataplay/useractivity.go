package main

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

type HomeData struct {
	Label string
	Value int
}

type ProfileDiscovery struct {
	PatternId     int       `json:"patternid"`
	Title         string    `json:"title"`
	ApiString     string    `json:"apistring"`
	DiscoveryDate time.Time `json:"discoverydate"`
}

type ProfileObservation struct {
	ObservationId int       `json:"observationid"`
	Title         string    `json:"charttitle"`
	ApiString     string    `json:"apistring"`
	DiscoveryDate time.Time `json:"discoverydate"`
	Comment       string    `json:"comment"`
}

type DataExpert struct {
	Username   string `json:"username"`
	Avatar     string `json:"avatar"`
	Reputation int    `json:"reputation"`
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

func GetProfileObservationsHttp(res http.ResponseWriter, req *http.Request) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	observation := []Observation{}
	err1 := DB.Where("uid = ?", uid).Find(&observation).Error
	if err1 != nil && err1 != gorm.RecordNotFound {
		http.Error(res, "Database query failed", http.StatusInternalServerError)
		return ""
	}

	profileObservations := make([]ProfileObservation, 0)

	for _, o := range observation {
		var tmp ProfileObservation
		tmp.ObservationId = o.ObservationId
		tmp.DiscoveryDate = o.Created
		tmp.Comment = o.Comment

		discTmp := Discovered{}
		err := DB.Where("discovered_id = ?", o.DiscoveredId).Find(&discTmp).Error
		if err != nil && err != gorm.RecordNotFound {
			http.Error(res, "Database query failed", http.StatusInternalServerError)
			return ""
		}

		if discTmp.CorrelationId == 0 {
			tmp.ApiString = "chart/related/" + discTmp.RelationId
			var td TableData
			json.Unmarshal(discTmp.Json, &td)
			tmp.Title = td.Title
		} else {
			cid := strconv.Itoa(discTmp.CorrelationId)
			tmp.ApiString = "chart/correlated/" + cid
			var cd CorrelationData
			json.Unmarshal(discTmp.Json, &cd)
			tmp.Title = cd.Table1.Title + " correlated with " + cd.Table2.Title
			if cd.Table3.Title != "" {
				tmp.Title += " correlated with " + cd.Table3.Title
			}
		}

		profileObservations = append(profileObservations, tmp)
	}

	r, err2 := json.Marshal(profileObservations)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func GetDiscoveriesHttp(res http.ResponseWriter, req *http.Request) string {
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

	var discovered []Discovered
	err1 := DB.Where("uid = ?", uid).Find(&discovered).Error
	if err1 != nil {
		return "not found"
	}

	profileDiscoveries := make([]ProfileDiscovery, 0)

	for _, d := range discovered {
		var tmp ProfileDiscovery
		tmp.PatternId = d.DiscoveredId
		tmp.DiscoveryDate = d.Created

		if d.CorrelationId == 0 {
			tmp.ApiString = "chart/related/" + d.RelationId
			var td TableData
			json.Unmarshal(d.Json, &td)
			tmp.Title = td.Title
		} else {
			cid := strconv.Itoa(d.CorrelationId)
			tmp.ApiString = "chart/correlated/" + cid
			var cd CorrelationData
			json.Unmarshal(d.Json, &cd)
			tmp.Title = cd.Table1.Title + " correlated with " + cd.Table2.Title
			if cd.Table3.Title != "" {
				tmp.Title += " correlated with " + cd.Table3.Title
			}
		}

		profileDiscoveries = append(profileDiscoveries, tmp)
	}

	r, err2 := json.Marshal(profileDiscoveries)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return "Unable to parse JSON"
	}

	return string(r)
}

func GetValidatedDiscoveriesHttp(res http.ResponseWriter, req *http.Request) string {
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

	discovered := []Discovered{}
	err1 := DB.Where("uid = ?", uid).Where("valid > ?", 0).Find(&discovered).Error
	if err1 != nil && err1 != gorm.RecordNotFound {
		http.Error(res, "Database query failed", http.StatusInternalServerError)
		return ""
	}

	profileDiscoveries := make([]ProfileDiscovery, 0)

	for _, d := range discovered {
		var tmp ProfileDiscovery
		tmp.PatternId = d.DiscoveredId
		tmp.DiscoveryDate = d.Created

		if d.CorrelationId == 0 {
			tmp.ApiString = "chart/" + d.RelationId
			var td TableData
			json.Unmarshal(d.Json, &td)
			tmp.Title = td.Title
		} else {
			cid := strconv.Itoa(d.CorrelationId)
			tmp.ApiString = "chartcorrelated/" + cid
			var cd CorrelationData
			json.Unmarshal(d.Json, &cd)
			tmp.Title = cd.Table1.Title + " correlated with " + cd.Table2.Title
			if cd.Table3.Title != "" {
				tmp.Title += " correlated with " + cd.Table3.Title
			}
		}

		profileDiscoveries = append(profileDiscoveries, tmp)
	}

	r, err2 := json.Marshal(profileDiscoveries)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func GetHomePageDataHttp(res http.ResponseWriter, req *http.Request) string {
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

func GetReputationHttp(res http.ResponseWriter, req *http.Request) string {
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

func GetAmountDiscoveriesHttp(res http.ResponseWriter, req *http.Request) string {
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

func GetDataExpertsHttp(res http.ResponseWriter, req *http.Request) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	users := []User{}
	err1 := DB.Order("reputation desc").Limit(5).Find(&users).Error
	if err1 != nil {
		return "not found"
	}

	var tmpDE DataExpert
	var dataExperts []DataExpert
	for _, u := range users {
		tmpDE.Username = u.Username
		tmpDE.Avatar = u.Avatar
		tmpDE.Reputation = u.Reputation
		dataExperts = append(dataExperts, tmpDE)
	}

	r, err2 := json.Marshal(dataExperts)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return "Unable to parse JSON"
	}

	return string(r)
}
