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

// increment user validated total for chart and rerank, id = 0 for new validation
func ValidateChart(valflag bool, patternid int, correlated bool, json []byte, originid int, uid int) *appError {
	val := Validated{}
	vld := Validation{}

	if patternid == 0 { // if new validation
		val.Discoverer = uid
		val.Created = time.Now()
		val.Correlated = correlated
		val.Rating = RankValidations(1, 0)
		if valflag {
			val.Rating = RankValidations(1, 0)
			val.Valid = 1
			val.Invalid = 0
		} else {
			val.Rating = 0
			val.Valid = 0
			val.Invalid = 1
		}
		val.Json = json
		val.OriginId = originid

		err := DB.Save(&val).Error
		if err != nil {
			return &appError{err, "Database query failed (Save)", http.StatusInternalServerError}
		}

	} else { // for pre existing charts
		err := DB.Where("patternid= ?", patternid).Find(&val).Error
		check(err)

		if valflag {
			val.Valid++
		} else {
			val.Invalid++
		}
		val.Rating = RankValidations(val.Valid, val.Invalid)

		err = DB.Save(&val).Error
		check(err)

		vld.PatternId = patternid
		vld.ObservationId = 0
		vld.Validator = uid
		vld.ValidationType = "chart"
		vld.Created = time.Now()

		err = DB.Save(&vld).Error
		if err != nil {
			return &appError{err, "Database query failed (Save)", http.StatusInternalServerError}
		}
	}

	return nil
}

// increment user validated total for observation and rerank, id = 0 to add new observation
func ValidateObservation(valflag bool, obsid int, text string, patternid int, uid int, coordinates string) *appError {
	obs := Observation{}
	vld := Validation{}

	if obsid == 0 {
		obs.Text = text
		obs.PatternId = patternid
		obs.Discoverer = uid
		obs.Coordinates = coordinates
		if valflag {
			obs.Rating = RankValidations(1, 0)
			obs.Valid = 1
			obs.Invalid = 0
		} else {
			obs.Rating = 0
			obs.Valid = 0
			obs.Invalid = 1
		}
		obs.Created = time.Now()

		err := DB.Save(&obs).Error
		if err != nil {
			return &appError{err, "Database query failed (Save)", http.StatusInternalServerError}
		}

	} else {
		err := DB.Where("patternid= ?", obsid).Find(&obs).Error
		check(err)

		if valflag {
			obs.Valid++
		} else {
			obs.Invalid++
		}
		obs.Rating = RankValidations(obs.Valid, obs.Invalid)

		err = DB.Save(&obs).Error
		check(err)

		vld.PatternId = 0
		vld.ObservationId = obsid
		vld.Validator = uid
		vld.ValidationType = "observation"
		vld.Created = time.Now()

		err = DB.Save(&vld).Error
		if err != nil {
			return &appError{err, "Database query failed (Save)", http.StatusInternalServerError}
		}
	}

	return nil
}

func ValidateChartHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	if params["uid"] == "" {
		return "no user id"
	}

	valflag, e := strconv.ParseBool(params["valflag"])
	if e != nil {
		http.Error(res, "bad validation flag", http.StatusBadRequest)
		return "bad validation flag"
	}

	patternid, e := strconv.Atoi(params["patternid"])
	if e != nil || patternid < 0 {
		patternid = 0
	}

	correlated, e := strconv.ParseBool(params["correlated"])
	if e != nil {
		http.Error(res, "bad correlation flag", http.StatusBadRequest)
		return "bad correlation flag"
	}

	originid, e := strconv.Atoi(params["originid"])
	if e != nil {
		http.Error(res, "bad originid", http.StatusBadRequest)
		return "bad originid"
	}

	uid, e := strconv.Atoi(params["uid"])
	if e != nil {
		http.Error(res, "bad uid", http.StatusBadRequest)
		return "bad uid"
	}

	json := []byte(params["json"])

	err := ValidateChart(valflag, patternid, correlated, json, originid, uid)
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

	valflag, e := strconv.ParseBool(params["valflag"])
	if e != nil {
		http.Error(res, "bad validation flag", http.StatusBadRequest)
		return "bad validation flag"
	}

	obsid, e := strconv.Atoi(params["obsid"])
	if e != nil || obsid < 0 {
		obsid = 0
	}

	patternid, e := strconv.Atoi(params["patternid"])
	if e != nil {
		http.Error(res, "bad patternid", http.StatusBadRequest)
		return "bad patternid"
	}

	if params["uid"] == "" {
		return "no user id"
	}

	uid, e := strconv.Atoi(params["uid"])
	if e != nil {
		http.Error(res, "bad uid", http.StatusBadRequest)
		return "bad uid"
	}

	if params["coordinates"] == "" {
		return "bad coordinates"
	}

	err := ValidateObservation(valflag, obsid, params["text"], patternid, uid, params["coordinates"])
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
	if params["uid"] == "" {
		return "no user id"
	}

	valflag, e := strconv.ParseBool(params["valflag"])
	if e != nil {
		return "bad validation flag"
	}

	patternid, e := strconv.Atoi(params["patternid"])
	if e != nil || patternid < 0 {
		patternid = 0
	}

	correlated, e := strconv.ParseBool(params["correlated"])
	if e != nil {
		return "bad correlation flag"
	}

	originid, e := strconv.Atoi(params["originid"])
	if e != nil {
		return "bad originid"
	}

	uid, e := strconv.Atoi(params["uid"])
	if e != nil {
		return "bad uid"
	}

	json := []byte(params["json"])

	err := ValidateChart(valflag, patternid, correlated, json, originid, uid)
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
	valflag, e := strconv.ParseBool(params["valflag"])
	if e != nil {
		return "bad validation flag"
	}

	obsid, e := strconv.Atoi(params["obsid"])
	if e != nil || obsid < 0 {
		obsid = 0
	}

	patternid, e := strconv.Atoi(params["patternid"])
	if e != nil {
		return "bad patternid"
	}

	if params["uid"] == "" {
		return "no user id"
	}

	uid, e := strconv.Atoi(params["uid"])
	if e != nil {
		return "bad uid"
	}

	if params["coordinates"] == "" {
		return "bad coordinates"
	}

	err := ValidateObservation(valflag, obsid, params["text"], patternid, uid, params["coordinates"])
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
