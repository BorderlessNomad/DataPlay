package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"net/http"
	"strconv"
	"time"
)

type Observations struct {
	ObservationId int    `json:"observation_id"`
	Comment       string `json:"comment"`
	X             string `json:"y"`
	Y             string `json:"y"`
}

// add observation
func AddObservation(id int, uid int, comment string, x string, y string) *appError {
	obs := Observation{}
	obs.Comment = comment
	obs.PatternId = id
	obs.Uid = uid
	obs.X = x
	obs.Y = y
	obs.Created = time.Now()

	vtd := Validated{}
	e := DB.Where("pattern_id= ?", id).First(&vtd).Error
	check(e)
	Reputation(vtd.Uid, discObs) // add points to rep of user who discovered chart when their discovery receives an observation

	err := AddActivity(uid, "c", obs.Created) // add to activities
	if err != nil {
		return err
	}

	err2 := DB.Save(&obs).Error
	if err2 != nil {
		return &appError{err2, "Database query failed (Save)", http.StatusInternalServerError}
	}

	return nil
}

// get all observations for a particular chart
func GetObservations(id int) ([]Observations, *appError) {
	obs := make([]Observation, 0)
	obsData := make([]Observations, 0)
	var tmpOD Observations

	err := DB.Where("pattern_id= ?", id).Find(&obs).Error
	if err != nil {
		return nil, &appError{err, "Database query failed (Save)", http.StatusInternalServerError}
	}

	for _, v := range obs {
		tmpOD.Comment = v.Comment
		tmpOD.X = v.X
		tmpOD.Y = v.Y
		tmpOD.ObservationId = v.ObservationId
		obsData = append(obsData, tmpOD)
	}

	return obsData, nil
}

func AddObservationHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	if params["id"] == "" {
		return "no observations id"
	}

	if params["uid"] == "" {
		return "no user id"
	}

	if params["x"] == "" {
		return "no x value"
	}
	if params["y"] == "" {
		return "no y value"
	}

	id, e := strconv.Atoi(params["id"])
	if e != nil {
		http.Error(res, "bad id", http.StatusBadRequest)
		return "bad id"
	}

	uid, err1 := GetUserID(session)
	if err1 != nil {
		http.Error(res, err1.Message, err1.Code)
		return "Could not validate user"
	}

	err := AddObservation(id, uid, params["comment"], params["x"], params["y"])
	if err != nil {
		http.Error(res, "could not add observation", http.StatusBadRequest)
		return "could not add observation"
	}

	return "observation added"
}

func GetObservationsHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	if params["id"] == "" {
		return "no observations id"
	}

	id, e := strconv.Atoi(params["id"])
	if e != nil {
		http.Error(res, "bad id", http.StatusBadRequest)
		return "bad id"
	}

	obs, err := GetObservations(id)
	if err != nil {
		http.Error(res, "could not get observations", http.StatusBadRequest)
		return "could not get observations"
	}

	r, err1 := json.Marshal(obs)
	if err1 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func AddObservationQ(params map[string]string) string {
	if params["id"] == "" {
		return "no id"
	}

	if params["uid"] == "" {
		return "no uid"
	}

	if params["x"] == "" {
		return "no x coordinate"
	}

	if params["y"] == "" {
		return "no y coordinate"
	}

	id, e := strconv.Atoi(params["id"])
	if e != nil || id < 0 {
		id = 0
	}

	uid, e := strconv.Atoi(params["uid"])
	if e != nil {
		return "bad uid"
	}

	err := AddObservation(id, uid, params["comment"], params["x"], params["y"])
	if err != nil {
		return "could not add observation"
	}

	return "observation added"
}

func GetObservationsQ(params map[string]string) string {
	id, e := strconv.Atoi(params["id"])

	if e != nil || id < 0 {
		return "Observations could not be retrieved"
	}

	result, err := GetObservations(id)
	if err != nil {
		return "Observations could not be retrieved"
	}

	r, e := json.Marshal(result)
	if e != nil {
		return "json error"
	}

	return string(r)
}
