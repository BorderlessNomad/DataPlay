package main

import (
	"github.com/codegangsta/martini"
	"github.com/jinzhu/gorm"
	"math"
	"net/http"
	"strconv"
	"strings"
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

// increment user validated total for chart and rerank, return validated id
func ValidateChart(rcid string, uid int, valflag bool, skipval bool) (string, *appError) {
	t := time.Now()
	validated := Validated{}
	validation := Validation{}

	if strings.ContainsAny(rcid, "_") { // if a relation id
		err := DB.Where("relation_id= ?", rcid).First(&validated).Error
		if err != nil {
			return "", &appError{err, "Database query failed", http.StatusInternalServerError}
		}
	} else { // if a correlation id of type int
		cid, _ := strconv.Atoi(rcid)
		if err != nil {
			return "", &appError{err, "Could not convert id to int", http.StatusInternalServerError}
		}
		err := DB.Where("correlation_id= ?", cid).First(&validated).Error
		if err != nil {
			return "", &appError{err, "Database query failed (cid)", http.StatusInternalServerError}
		}
	}

	if !skipval {
		if valflag {
			validated.Valid++
			Reputation(validated.Uid, discVal) // add points for discovery validation
			AddActivity(uid, "vc", t)
		} else {
			validated.Invalid++
			Reputation(validated.Uid, discInval) // remove points for discovery invalidation
			AddActivity(uid, "ic", t)
		}
		validated.Rating = RankValidations(validated.Valid, validated.Invalid)

		err1 := DB.Save(&validated).Error
		if err1 != nil {
			return "", &appError{err1, "Database query failed - validate chart (Save validated)", http.StatusInternalServerError}
		}

		validation.ValidatedId = validated.ValidatedId
		validation.Validator = uid
		validation.Created = t
		validation.ObservationId = 0 // not an observation

		err2 := DB.Save(&validation).Error
		if err2 != nil {
			return "", &appError{err2, "Database query failed (Save validaition)", http.StatusInternalServerError}
		}
	}

	return strconv.Itoa(validated.ValidatedId), nil
}

// increment user validated total for observation and rerank
func ValidateObservation(oid int, uid int, valflag bool) *appError {
	t := time.Now()
	obs := Observation{}
	validation := Validation{}

	err := DB.Where("observation_id = ?", oid).First(&obs).Error
	if err != nil && err != gorm.RecordNotFound {
		return &appError{err, "Database query failed - validate observation (get)", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return &appError{err, "No such observation found!", http.StatusNotFound}
	}

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
	if err != nil {
		return &appError{err, "Database query failed - Unable to save an observation.", http.StatusInternalServerError}
	}

	validation.ValidatedId = 0 // not a chart
	validation.Validator = uid
	validation.Created = time.Now()
	validation.ObservationId = oid
	validation.Valflag = valflag

	err = DB.Save(&validation).Error
	if err != nil {
		return &appError{err, "Database query failed - validate observation (Save validation)", http.StatusInternalServerError}
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

	skipval := false
	valflag := false

	if params["valflag"] == "" { // if no valflag then skip validation and just return validated id
		skipval = true
	} else {
		valflag, _ = strconv.ParseBool(params["valflag"])
	}

	if params["rcid"] == "" {
		return "no chart id"
	}

	uid, err1 := GetUserID(session)
	if err1 != nil {
		http.Error(res, err1.Message, err1.Code)
		return "Could not validate user"
	}

	result, err2 := ValidateChart(params["rcid"], uid, valflag, skipval)
	if err2 != nil {
		msg := ""
		if valflag {
			msg = " could not validate chart"
		} else {
			msg = " could not invalidate chart"
		}

		http.Error(res, msg, http.StatusBadRequest)
		return err2.Message + msg
	}

	if valflag {
		return result
	} else {
		return result
	}
}

func ValidateObservationHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	oid, err := strconv.Atoi(params["oid"])
	if err != nil || oid < 0 {
		oid = 0
	}

	uid, err1 := GetUserID(session)
	if err1 != nil {
		http.Error(res, err1.Message, err1.Code)
		return "Could not validate user"
	}

	valflag, err2 := strconv.ParseBool(params["valflag"])
	if err2 != nil {
		http.Error(res, "bad validation flag", http.StatusBadRequest)
		return "bad validation flag"
	}

	err3 := ValidateObservation(oid, uid, valflag)
	if err3 != nil {
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
