package main

import (
	"github.com/codegangsta/martini"
	"math"
	"net/http"
	"strconv"
	"time"
)

type Observations struct {
	Comment string
	X       string
	Y       string
}

// given a small fraction of ratings there is a strong (95%) chance that the "real", final positive rating will be this value
// eg: gives expected (not necessarily current as there may have only been a few votes so far) value of positive ratings / total ratings
func RankValidations(valid int, invalid int) float64 {
	pos := float64(valid)
	tot := float64(valid + invalid)

	if tot == 0 {
		return 0
	}

	z := 1.96
	phat := pos / tot
	result := (phat + z*z/(2*tot) - z*math.Sqrt((phat*(1-phat)+z*z/(4*tot))/tot)) / (1 + z*z/tot)
	return result
}

// increment user validated total for chart and rerank, id = 0 for new validation
func ValidateChart(id int, uid int, valflag bool) *appError {
	vtd := Validated{}
	vdn := Validation{}

	err := DB.Where("patternid= ?", id).First(&vtd).Error
	check(err)

	if valflag {
		vtd.Valid++
	} else {
		vtd.Invalid++
	}
	vtd.Rating = RankValidations(vtd.Valid, vtd.Invalid)

	err = DB.Save(&vtd).Error
	check(err)

	vdn.PatternId = id
	vdn.Validator = uid
	vdn.Created = time.Now()
	vdn.ObservationId = 0 // not an observation

	err = DB.Save(&vdn).Error
	if err != nil {
		return &appError{err, "Database query failed (Save)", http.StatusInternalServerError}
	}

	return nil
}

// increment user validated total for observation and rerank, id = 0 to add new observation
func ValidateObservation(id int, uid int, valflag bool) *appError {
	obs := Observation{}
	vdn := Validation{}

	err := DB.Where("observationid= ?", id).First(&obs).Error
	check(err)

	if valflag {
		obs.Valid++
	} else {
		obs.Invalid++
	}

	obs.Rating = RankValidations(obs.Valid, obs.Invalid)
	err = DB.Save(&obs).Error
	check(err)

	vdn.PatternId = 0 // not a chart
	vdn.Validator = uid
	vdn.Created = time.Now()
	vdn.ObservationId = id

	err = DB.Save(&vdn).Error
	if err != nil {
		return &appError{err, "Database query failed (Save)", http.StatusInternalServerError}
	}

	return nil
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

//////////////////////////////////////////////
func ValidateChartHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	if params["id"] == "" {
		return "no pattern id"
	}

	if params["uid"] == "" {
		return "no user id"
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

	valflag, e := strconv.ParseBool(params["valflag"])
	if e != nil {
		http.Error(res, "bad validation flag", http.StatusBadRequest)
		return "bad validation flag"
	}

	err := ValidateChart(id, uid, valflag)
	if err != nil {
		msg := ""
		if valflag {
			msg = "Could not validate chart"
		} else {
			msg = "Could not invalidate chart"
		}
		http.Error(res, msg, http.StatusBadRequest)
		return msg
	}

	if valflag {
		return "Chart validated"
	} else {
		return "Chart invalidated"
	}
}

func ValidateObservationHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return "Missing session parameter."
	}

	id, e := strconv.Atoi(params["id"])
	if e != nil || id < 0 {
		id = 0
	}

	if params["uid"] == "" {
		return "no user id"
	}

	uid, e := strconv.Atoi(params["uid"])
	if e != nil {
		http.Error(res, "bad uid", http.StatusBadRequest)
		return "bad uid"
	}

	valflag, e := strconv.ParseBool(params["valflag"])
	if e != nil {
		http.Error(res, "bad validation flag", http.StatusBadRequest)
		return "bad validation flag"
	}

	err := ValidateObservation(id, uid, valflag)
	if err != nil {
		msg := ""
		if valflag {
			msg = "Could not validate observation"
		} else {
			msg = "Could not invalidate observation"
		}
		return msg
	}

	if valflag {
		return "Observation validated"
	} else {
		return "Observation invalidated"
	}
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

func ValidateChartQ(params map[string]string) string {
	if params["id"] == "" {
		return "no id"
	}

	if params["uid"] == "" {
		return "no user id"
	}

	id, e := strconv.Atoi(params["id"])
	if e != nil || id < 0 {
		id = 0
	}

	uid, e := strconv.Atoi(params["uid"])
	if e != nil {
		return "bad uid"
	}

	valflag, e := strconv.ParseBool(params["valflag"])
	if e != nil {
		return "bad validation flag"
	}

	err := ValidateChart(id, uid, valflag)
	if err != nil {
		msg := ""
		if valflag {
			msg = "Could not validate chart"
		} else {
			msg = "Could not invalidate chart"
		}
		return msg
	}

	if valflag {
		return "Chart validated"
	} else {
		return "Chart invalidated"
	}
}

func ValidateObservationQ(params map[string]string) string {
	if params["id"] == "" {
		return "no id"
	}

	id, e := strconv.Atoi(params["id"])
	if e != nil || id < 0 {
		id = 0
	}

	if params["uid"] == "" {
		return "no user id"
	}

	uid, e := strconv.Atoi(params["uid"])
	if e != nil {
		return "bad uid"
	}
	valflag, e := strconv.ParseBool(params["valflag"])
	if e != nil {
		return "bad validation flag"
	}

	err := ValidateObservation(id, uid, valflag)
	if err != nil {
		msg := ""
		if valflag {
			msg = "Could not validate observation"
		} else {
			msg = "Could not invalidate observation"
		}
		return msg
	}

	if valflag {
		return "Observation validated"
	} else {
		return "Observation invalidated"
	}
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
