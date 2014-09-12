package main

import (
	"encoding/json"
	"fmt"
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
	Email      string `json:"email, omitempty"`
}

type ObservationComment struct {
	DiscoveryId string `json:"did" binding:"required"`
	X           string `json:"x" binding:"required"`
	Y           string `json:"y" binding:"required"`
	Comment     string `json:"comment" binding:"required"`
}

// add observation
func AddObservation(did int, uid int, comment string, x string, y string) (string, *appError) {
	observation := Observation{}
	observation.Comment = comment
	observation.DiscoveredId = did
	observation.Uid = uid
	observation.X = x
	observation.Y = y
	observation.Created = time.Now()

	discovered := Discovered{}
	err := DB.Where("discovered_id = ?", did).First(&discovered).Error
	if err != nil {
		return "Database query failed (find discovered)", &appError{err, "Database query failed (find discovered)", http.StatusInternalServerError}
	}

	Reputation(discovered.Uid, discObs) // add points to rep of user who discovered chart when their discovery receives an observation

	err1 := AddActivity(uid, "c", observation.Created) // add to activities
	if err1 != nil {
		return "", err1
	}

	err2 := DB.Save(&observation).Error
	if err2 != nil {
		return "Database query failed (Save observation)", &appError{err2, "Database query failed (Save observation)", http.StatusInternalServerError}
	}

	err3 := DB.Where("discovered_id= ?", did).Where("x =?", x).Where("y =?", y).First(&observation).Error
	if err3 != nil {
		return "Database query failed - add observation (find observation)", &appError{err3, "Database query failed - add observation (find observation)", http.StatusInternalServerError}
	}

	fmt.Println("MACARINA", observation)

	return strconv.Itoa(observation.ObservationId), nil
}

// get all observations for a particular chart
func GetObservations(did int) ([]Observations, *appError) {
	discovered := make([]Discovered, 0)
	err := DB.Where("discovered_id = ?", did).Find(&discovered).Error
	if err != nil {
		return nil, &appError{err, "Database query failed - get observation (find discovered)", http.StatusInternalServerError}
	}

	observation := make([]Observation, 0)
	obsData := make([]Observations, 0)
	var tmpOD Observations

	err1 := DB.Where("discovered_id = ?", did).Order("observation_id ASC").Find(&observation).Error
	if err1 != nil && err1 != gorm.RecordNotFound {
		return nil, &appError{err1, "Database query failed  - get observation (find observation)", http.StatusInternalServerError}
	} else if err1 == gorm.RecordNotFound {
		return obsData, nil //Empty map
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
		tmpOD.User.Email = GetMD5Hash(user[0].Email)

		if discovered[0].Uid == observation[0].Uid { // if commenter discovered the chart
			tmpOD.User.Discoverer = true
		} else {
			tmpOD.User.Discoverer = false
		}

		obsData = append(obsData, tmpOD)
	}

	return obsData, nil
}

func AddObservationHttp(res http.ResponseWriter, req *http.Request, observation ObservationComment) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	if observation.DiscoveryId == "" || observation.X == "" || observation.Y == "" || observation.Comment == "" {
		http.Error(res, "Invalid/missing request parameters.", http.StatusBadRequest)
		return ""
	}

	did, err := strconv.Atoi(observation.DiscoveryId)
	if err != nil {
		http.Error(res, "Bad Discovery id", http.StatusBadRequest)
		return ""
	}

	uid, err1 := GetUserID(session)
	if err1 != nil {
		http.Error(res, err1.Message, err1.Code)
		return err1.Message
	}

	result, err2 := AddObservation(did, uid, observation.Comment, observation.X, observation.Y)
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

	if params["did"] == "" {
		return "no discovered id"
	}

	did, err := strconv.Atoi(params["did"])
	if err != nil {
		http.Error(res, "bad discovered id", http.StatusBadRequest)
		return "bad discovered id"
	}

	obs, err1 := GetObservations(did)
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
	did, err := strconv.Atoi(params["did"])

	if err != nil || did < 0 {
		return "Observations could not be retrieved"
	}

	result, err1 := GetObservations(did)
	if err1 != nil {
		return "Observations could not be retrieved"
	}

	r, err2 := json.Marshal(result)
	if err2 != nil {
		return "json error"
	}

	return string(r)
}
