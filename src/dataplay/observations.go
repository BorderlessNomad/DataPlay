package main

import (
	"github.com/codegangsta/martini"
	"net/http"
	"strconv"
	"time"
)

type Observations struct {
	Comment string
	X       string
	Y       string
}

// add observation
func AddObservation(id int, uid int, comment string, x string, y string) *appError {
	obs := Observation{}
	obs.Comment = comment
	obs.PatternId = id
	obs.Discoverer = uid
	obs.X = x
	obs.Y = y
	obs.Created = time.Now()

	err := DB.Save(&obs).Error
	if err != nil {
		return &appError{err, "Database query failed (Save)", http.StatusInternalServerError}
	}

	return nil
}

// get all observations for a particular chart
func GetObservations(id int) ([]Observations, *appError) {
	obs := make([]Observation, 0)
	obsData := make([]Observations, 0)
	var tmpOD Observations

	err := DB.Where("patternid= ?", id).Find(&obs).Error
	if err != nil {
		return obsData, &appError{err, "Database query failed (Save)", http.StatusInternalServerError}
	}

	for _, v := range obs {
		tmpOD.Comment = v.Comment
		tmpOD.X = v.X
		tmpOD.Y = v.Y
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

	uid, e := strconv.Atoi(params["uid"])
	if e != nil {
		http.Error(res, "bad uid", http.StatusBadRequest)
		return "bad uid"
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

	_, err := GetObservations(id)
	if err != nil {
		http.Error(res, "could not get observations", http.StatusBadRequest)
		return "could not get observations"
	}

	return "observations retrieved"
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

	return "Observations added"

}

func GetObservationsQ(params map[string]string) string {
	id, e := strconv.Atoi(params["id"])
	if e != nil || id < 0 {
		return "Observations could not be retrieved"
	}
	_, err := GetObservations(id)
	if err != nil {
		return "Observations could not be retrieved"
	}

	return "Observations returned"
}
