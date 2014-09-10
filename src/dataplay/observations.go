package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

type Observations struct {
	ObservationId int       `json:"observation_id"`
	Comment       string    `json:"comment, omitempty"`
	X             string    `json:"x"`
	Y             string    `json:"y"`
	User          UserData  `json:"user"`
	Created       time.Time `json:"created, omitempty"`
	Valid         int       `json:"validations, omitempty"`
	Invalid       int       `json:"invalidations, omitempty"`
}

type UserData struct {
	Username   string `json:"name, omitempty"`
	Reputation int    `json:"reputation, omitempty"`
	Avatar     string `json:"avatar, omitempty"`
	Discoverer bool   `json:"discoverer, omitempty"`
}

// add observation
func AddObservation(vid int, uid int, comment string, x string, y string) (string, *appError) {
	observation := Observation{}
	addObs := false

	e := DB.Where("validated_id= ?", vid).Where("x =?", x).Where("y =?", y).First(&observation).Error
	if e == gorm.RecordNotFound {
		addObs = true
	}

	if addObs {
		observation.Comment = comment
		observation.ValidatedId = vid
		observation.Uid = uid
		observation.X = x
		observation.Y = y
		observation.Created = time.Now()

		validated := Validated{}
		err := DB.Where("validated_id= ?", vid).First(&validated).Error
		if err != nil {
			return "", &appError{err, "Database query failed (find validated)", http.StatusInternalServerError}
		}
		Reputation(validated.Uid, discObs) // add points to rep of user who discovered chart when their discovery receives an observation

		err1 := AddActivity(uid, "c", observation.Created) // add to activities
		if err1 != nil {
			return "", err1
		}

		err2 := DB.Save(&observation).Error
		if err2 != nil {
			return "", &appError{err2, "Database query failed (Save observation)", http.StatusInternalServerError}
		}

		err3 := DB.Where("validated_id= ?", vid).Where("x =?", x).Where("y =?", y).First(&observation).Error
		if err3 != nil {
			return "", &appError{err3, "Database query failed - add observation (find observation)", http.StatusInternalServerError}
		}
	}

	return strconv.Itoa(observation.ObservationId), nil
}

// get all observations for a particular chart
func GetObservations(vid int) ([]Observations, *appError) {
	validated := make([]Validated, 0)
	err := DB.Where("validated_id= ?", vid).Find(&validated).Error
	if err != nil {
		return nil, &appError{err, "Database query failed - get observation (find validation)", http.StatusInternalServerError}
	}

	observation := make([]Observation, 0)
	obsData := make([]Observations, 0)
	var tmpOD Observations

	err1 := DB.Where("validated_id= ?", vid).Find(&observation).Error
	if err1 != nil {
		return nil, &appError{err1, "Database query failed  - get observation (find observation)", http.StatusInternalServerError}
	}

	for _, o := range observation {
		tmpOD.Comment = o.Comment
		tmpOD.X = o.X
		tmpOD.Y = o.Y
		tmpOD.Valid = o.Valid
		tmpOD.Invalid = o.Invalid
		tmpOD.Created = o.Created
		tmpOD.ObservationId = o.ObservationId

		user := make([]User, 0)
		err2 := DB.Where("uid= ?", o.Uid).Find(&user).Error
		if err2 != nil {
			return nil, &appError{err2, "Database query failed - get observation - no such user", http.StatusInternalServerError}
		}

		tmpOD.User.Username = user[0].Username
		tmpOD.User.Avatar = user[0].Avatar
		tmpOD.User.Reputation = user[0].Reputation

		if validated[0].Uid == observation[0].Uid { // if commenter discovered the chart
			tmpOD.User.Discoverer = true
		} else {
			tmpOD.User.Discoverer = false
		}

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

	if params["vid"] == "" {
		http.Error(res, "no validated id.", http.StatusBadRequest)
		return ""
	}

	if params["x"] == "" {
		http.Error(res, "no x value.", http.StatusBadRequest)
		return ""
	}
	if params["y"] == "" {
		http.Error(res, "no y value.", http.StatusBadRequest)
		return ""
	}

	vid, err := strconv.Atoi(params["vid"])
	if err != nil {
		http.Error(res, "bad validated id", http.StatusBadRequest)
		return ""
	}

	uid, err1 := GetUserID(session)
	if err1 != nil {
		http.Error(res, err1.Message, err1.Code)
		return ""
	}

	result, err2 := AddObservation(vid, uid, params["comment"], params["x"], params["y"])
	if err2 != nil {
		http.Error(res, err2.Message, http.StatusBadRequest)
		return err2.Message
	}

	return result
}

func GetObservationsHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	if params["vid"] == "" {
		return "no validated id"
	}

	vid, err := strconv.Atoi(params["vid"])
	if err != nil {
		http.Error(res, "bad validated id", http.StatusBadRequest)
		return "bad validated id"
	}

	obs, err1 := GetObservations(vid)
	if err1 != nil {
		http.Error(res, err1.Message, http.StatusBadRequest)
		return err1.Message
	}

	r, err2 := json.Marshal(obs)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return "Unable to parse JSON"
	}

	return string(r)
}

func GetObservationsQ(params map[string]string) string {
	vid, err := strconv.Atoi(params["vid"])

	if err != nil || vid < 0 {
		return "Observations could not be retrieved"
	}

	result, err1 := GetObservations(vid)
	if err1 != nil {
		return "Observations could not be retrieved"
	}

	r, err2 := json.Marshal(result)
	if err2 != nil {
		return "json error"
	}

	return string(r)
}
