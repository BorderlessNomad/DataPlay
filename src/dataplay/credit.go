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

type CreditRequest struct {
	Cid string `json:"cid"`
	Rid string `json:"rid"`
}

// given a small fraction of ratings there is a strong (95%) chance that the "real", final positive rating will be this value
// eg: gives expected (not necessarily current as there may have only been a few votes so far) value of positive ratings / total ratings
func RankCredits(credit int, discredit int) float64 {
	pos := float64(credit)
	tot := float64(credit + discredit)

	if tot == 0 {
		return 0
	}

	z := 1.96
	phat := pos / tot
	result := (phat + z*z/(2*tot) - z*math.Sqrt((phat*(1-phat)+z*z/(4*tot))/tot)) / (1 + z*z/tot)
	return result
}

// increment user discovered total for chart and rerank, return discovered id
func CreditChart(rcid string, uid int, credflag bool) (string, *appError) {
	t := time.Now()
	discovered := Discovered{}
	credit := Credit{}

	if strings.ContainsAny(rcid, "/") { // if a relation id
		err := DB.Where("relation_id = ?", rcid).First(&discovered).Error
		if err != nil && err != gorm.RecordNotFound {
			return "", &appError{err, ", database query failed (relation_id)", http.StatusInternalServerError}
		}
	} else { // if a correlation id of type int
		cid, _ := strconv.Atoi(rcid)
		if err != nil {
			return "", &appError{err, ", could not convert id to int", http.StatusInternalServerError}
		}

		err := DB.Where("correlation_id = ?", cid).First(&discovered).Error
		if err != nil && err != gorm.RecordNotFound {
			return "", &appError{err, ", database query failed (correlation_id)", http.StatusInternalServerError}
		}
	}

	if credflag {
		discovered.Credited++
		Reputation(discovered.Uid, discCredit) // add points for discovery credit
		AddActivity(uid, "vc", t, discovered.DiscoveredId, 0)
	} else {
		discovered.Discredited++
		Reputation(discovered.Uid, discDiscredit) // remove points for discovery discredit
		AddActivity(uid, "ic", t, discovered.DiscoveredId, 0)
	}
	discovered.Rating = RankCredits(discovered.Credited, discovered.Discredited)
	err1 := DB.Save(&discovered).Error
	if err1 != nil {
		return "", &appError{err1, ", database query failed - credit chart (Save discovered)", http.StatusInternalServerError}
	}
	credit.DiscoveredId = discovered.DiscoveredId
	credit.Uid = uid
	credit.Created = t
	credit.ObservationId = 0 // not an observation
	err2 := DB.Save(&credit).Error
	if err2 != nil {
		return "", &appError{err2, ", database query failed (Save credit)", http.StatusInternalServerError}
	}

	return strconv.Itoa(discovered.DiscoveredId), nil
}

// increment user discovered total for observation and rerank
func CreditObservation(oid int, uid int, credflag bool) *appError {
	t := time.Now()
	observation := Observation{}
	credit := Credit{}

	err := DB.Where("observation_id = ?", oid).First(&observation).Error
	if err != nil && err != gorm.RecordNotFound {
		return &appError{err, " Database query failed - credit observation (get)", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return &appError{err, ", no such observation found!", http.StatusNotFound}
	}

	if observation.Uid == uid {
		return &appError{err, ", you cannot credit your own comment", http.StatusNotFound}
	}

	cred := Credit{}
	err2 := DB.Where("observation_id= ?", observation.ObservationId).Where("uid= ?", uid).Find(&cred).Error
	if err2 != nil && err2 != gorm.RecordNotFound {
		return &appError{err2, ", observation query failed.", http.StatusInternalServerError}
	} else if cred.CreditId != 0 {
		return &appError{err2, ", user has already credited this observation.", http.StatusInternalServerError}
	}

	if credflag {
		observation.Credited++
		Reputation(observation.Uid, obsCredit) // add points for observation credit
		AddActivity(uid, "vo", t, 0, observation.ObservationId)
	} else {
		observation.Credited++
		Reputation(observation.Uid, obsDiscredit) // remove points for observation incredit
		AddActivity(uid, "io", t, 0, observation.ObservationId)
	}

	observation.Rating = RankCredits(observation.Credited, observation.Discredited)
	err = DB.Save(&observation).Error
	if err != nil {
		return &appError{err, ", database query failed - Unable to save an observation.", http.StatusInternalServerError}
	}

	credit.DiscoveredId = 0 // not a chart
	credit.Uid = uid
	credit.Created = time.Now()
	credit.ObservationId = oid
	credit.Credflag = credflag

	err = DB.Save(&credit).Error
	if err != nil {
		return &appError{err, ", database query failed - credit observation (Save credit)", http.StatusInternalServerError}
	}

	return nil
}

//////////////////////////////////////////////
func CreditChartHttp(res http.ResponseWriter, req *http.Request, params martini.Params, credit CreditRequest) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	credflag := false

	if params["credflag"] == "" { // if no credflag then skip credit and just return discovered id
		http.Error(res, "Missing credflag", http.StatusBadRequest)
		return ""
	} else {
		credflag, _ = strconv.ParseBool(params["credflag"])
	}

	var rcid string
	if credit.Cid == "" && credit.Rid == "" {
		http.Error(res, "No Relation/Correlation ID provided.", http.StatusBadRequest)
		return ""
	} else if credit.Cid == "" {
		rcid = credit.Rid
	} else {
		rcid = credit.Cid
	}

	uid, err1 := GetUserID(session)
	if err1 != nil {
		http.Error(res, err1.Message, err1.Code)
		return ""
	}

	result, err2 := CreditChart(rcid, uid, credflag)
	if err2 != nil {
		msg := ""
		if credflag {
			msg = "Could not credit chart" + err2.Message
		} else {
			msg = "Could not discredit chart" + err2.Message
		}

		http.Error(res, err2.Message+msg, http.StatusBadRequest)
		return ""
	}

	if credflag {
		return result
	} else {
		return result
	}
}

func CreditObservationHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
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
		return "Could not credit user"
	}

	credflag, err2 := strconv.ParseBool(params["credflag"])
	if err2 != nil {
		http.Error(res, "bad credit flag", http.StatusBadRequest)
		return "bad credit flag"
	}

	err3 := CreditObservation(oid, uid, credflag)
	if err3 != nil {
		msg := ""
		if credflag {
			msg = "Could not credit observation" + err3.Message
		} else {
			msg = "Could not discredit observation" + err3.Message
		}
		return msg
	}

	if credflag {
		return "Observation credited"
	} else {
		return "Observation discredited"
	}
}
