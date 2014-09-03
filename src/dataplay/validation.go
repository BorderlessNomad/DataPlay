package main

import (
	"github.com/codegangsta/martini"
	"math"
	"net/http"
	"strconv"
	"time"
)

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

// increment user validated total for chart and rerank
func ValidateChart(id int, uid int, valflag bool) *appError {
	t := time.Now()
	vtd := Validated{}
	vdn := Validation{}

	err := DB.Where("pattern_id= ?", id).First(&vtd).Error
	check(err)

	if valflag {
		vtd.Valid++
		Reputation(vtd.Uid, discVal) // add points for discovery validation
		AddActivity(uid, "vc", t)
	} else {
		vtd.Invalid++
		Reputation(vtd.Uid, discInval) // remove points for discovery invalidation
		AddActivity(uid, "ic", t)
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

// increment user validated total for observation and rerank
func ValidateObservation(id int, uid int, valflag bool) *appError {
	t := time.Now()
	obs := Observation{}
	vdn := Validation{}

	err := DB.Where("observation_id= ?", id).First(&obs).Error
	check(err)

	if valflag {
		obs.Valid++
		Reputation(obs.Uid, obsVal) // add points for observation validation
		AddActivity(uid, "vo", t)
	} else {
		obs.Invalid++
		Reputation(obs.Uid, obsInval) // remove points for observation invalidation
		AddActivity(uid, "io", t)
	}

	obs.Rating = RankValidations(obs.Valid, obs.Invalid)
	err = DB.Save(&obs).Error
	check(err)

	vdn.PatternId = 0 // not a chart
	vdn.Validator = uid
	vdn.Created = time.Now()
	vdn.ObservationId = id
	vdn.Valflag = valflag

	err = DB.Save(&vdn).Error
	if err != nil {
		return &appError{err, "Database query failed (Save)", http.StatusInternalServerError}
	}

	return nil
}

//////////////////////////////////////////////
func ValidateChartHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
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
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
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
