package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

const OBSERVATION_ACTION_NONE = ""
const OBSERVATION_ACTION_CREDITED = "credited"
const OBSERVATION_ACTION_DISCREDITED = "discredited"

/**
 * @TODO Make activity table to consts - Glyn
 */
const OBSERVATION_TYPE_CREDITED = "Credited Observation"
const OBSERVATION_TYPE_DISCREDITED = "Discredited Observation"

type Observations struct {
	ObservationId int       `json:"observation_id"`
	Comment       string    `json:"comment, omitempty"`
	X             string    `json:"x"`
	Y             string    `json:"y"`
	Created       time.Time `json:"created, omitempty"`
	Credited      int       `json:"credits, omitempty"`
	Discredited   int       `json:"discredits, omitempty"`
	Flagged       bool      `json:"flagged, omitempty"`
	User          UserData  `json:"user"`
	Action        string    `json:"action, omitempty"`
}

type UserData struct {
	Username   string `json:"username, omitempty"`
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

type CommunityObservation struct {
	Username  string `json:"username"`
	Avatar    string `json:"avatar, omitempty"`
	Comment   string `json:"comment"`
	PatternID int    `json:"patternid"`
	Link      string `json:"linkstring"`
	EmailMD5  string `json:"MD5email"`
}

// add observation to chart
func AddObservation(did int, uid int, comment string, x string, y string) (Observations, *appError) {
	obs := Observations{}
	observation := Observation{}
	observation.Comment = comment
	observation.DiscoveredId = did
	observation.Uid = uid
	observation.X = x
	observation.Y = y
	observation.Created = time.Now()
	observation.Flagged = false

	discovered := Discovered{}
	err := DB.Where("discovered_id = ?", did).Find(&discovered).Error
	if err != nil {
		return obs, &appError{err, "Database query failed (find discovered)", http.StatusInternalServerError}
	}

	Reputation(discovered.Uid, discObs) // add points to rep of user who discovered chart when their discovery receives an observation

	err1 := AddActivity(uid, "c", observation.Created, did, 0) // add to activities
	if err1 != nil {
		return obs, err1
	}

	err2 := DB.Save(&observation).Error
	if err2 != nil {
		return obs, &appError{err2, "Database query failed (Save observation)", http.StatusInternalServerError}
	}

	err3 := DB.Where("discovered_id= ?", did).Where("x =?", x).Where("y =?", y).First(&observation).Error
	if err3 != nil {
		return obs, &appError{err3, "Database query failed - add observation (find observation)", http.StatusInternalServerError}
	}

	obs.Comment = observation.Comment
	obs.X = observation.X
	obs.Y = observation.Y
	obs.Credited = observation.Credited
	obs.Discredited = observation.Discredited
	obs.Created = observation.Created
	obs.ObservationId = observation.ObservationId
	obs.Flagged = observation.Flagged

	user := User{}
	err5 := DB.Where("uid= ?", observation.Uid).Find(&user).Error
	if err5 != nil {
		return obs, &appError{err2, "Database query failed - get observation - no such user", http.StatusInternalServerError}
	}

	obs.User.Username = user.Username
	obs.User.Avatar = user.Avatar
	obs.User.Reputation = user.Reputation
	obs.User.Email = GetMD5Hash(user.Email)

	if discovered.Uid == observation.Uid {
		obs.User.Discoverer = true
	} else {
		obs.User.Discoverer = false
	}

	activity := Activity{}
	err4 := DB.Where("uid = ?", uid).Where("observation_id = ?", observation.ObservationId).Find(&activity).Error

	if err4 != nil && err4 != gorm.RecordNotFound {
		return obs, &appError{err4, "Database query failed - get activity - no such activity", http.StatusInternalServerError}
	} else if err3 == gorm.RecordNotFound {
		obs.Action = OBSERVATION_ACTION_NONE
	} else {
		if activity.Type == OBSERVATION_TYPE_CREDITED {
			obs.Action = OBSERVATION_ACTION_CREDITED
		} else if activity.Type == OBSERVATION_TYPE_DISCREDITED {
			obs.Action = OBSERVATION_ACTION_DISCREDITED
		} else {
			obs.Action = OBSERVATION_ACTION_NONE
		}
	}

	return obs, nil
}

// get all observations for a particular chart
func GetObservations(did int, uid int) ([]Observations, *appError) {
	observation := make([]Observation, 0)
	obsData := make([]Observations, 0)
	discovered := Discovered{}

	err := DB.Where("discovered_id = ?", did).Find(&discovered).Error
	if err != nil && err != gorm.RecordNotFound {
		return nil, &appError{err, "Database query failed - get observation (find discovered)", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return obsData, nil //Empty map
	}

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
		tmpOD.Credited = o.Credited
		tmpOD.Discredited = o.Discredited
		tmpOD.Created = o.Created
		tmpOD.ObservationId = o.ObservationId
		tmpOD.Flagged = o.Flagged

		user := User{}
		err2 := DB.Where("uid = ?", o.Uid).Find(&user).Error
		if err2 != nil {
			return nil, &appError{err2, "Database query failed - get observation - no such user", http.StatusInternalServerError}
		}

		tmpOD.User.Username = user.Username
		tmpOD.User.Avatar = user.Avatar
		tmpOD.User.Reputation = user.Reputation
		tmpOD.User.Email = GetMD5Hash(user.Email)

		if discovered.Uid == observation[0].Uid { // if commenter discovered the chart
			tmpOD.User.Discoverer = true
		} else {
			tmpOD.User.Discoverer = false
		}

		activity := Activity{}
		err3 := DB.Where("uid = ?", uid).Where("observation_id = ?", o.ObservationId).Find(&activity).Error
		if err3 != nil && err3 != gorm.RecordNotFound {
			return nil, &appError{err3, "Database query failed - get activity - no such activity", http.StatusInternalServerError}
		} else if err3 == gorm.RecordNotFound {
			tmpOD.Action = OBSERVATION_ACTION_NONE
		} else {
			if activity.Type == OBSERVATION_TYPE_CREDITED {
				tmpOD.Action = OBSERVATION_ACTION_CREDITED
			} else if activity.Type == OBSERVATION_TYPE_DISCREDITED {
				tmpOD.Action = OBSERVATION_ACTION_DISCREDITED
			} else {
				tmpOD.Action = OBSERVATION_ACTION_NONE
			}
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
		return ""
	}

	result, err2 := AddObservation(did, uid, observation.Comment, observation.X, observation.Y)
	if err2 != nil {
		http.Error(res, err2.Message, http.StatusBadRequest)
		return ""
	}

	r, err3 := json.Marshal(result)
	if err3 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func GetObservationsHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	if params["did"] == "" {
		http.Error(res, "No discovered id.", http.StatusBadRequest)
		return ""
	}

	did, err1 := strconv.Atoi(params["did"])
	if err1 != nil {
		http.Error(res, "bad discovered id", http.StatusBadRequest)
		return ""
	}

	obs, err2 := GetObservations(did, uid)
	if err2 != nil {
		http.Error(res, err2.Message, http.StatusBadRequest)
		return ""
	}

	r, err3 := json.Marshal(obs)
	if err3 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func GetObservationsQ(params map[string]string) string {
	did, err := strconv.Atoi(params["did"])

	if err != nil || did < 0 {
		return "Observations could not be retrieved"
	}

	result, err1 := GetObservations(did, 11)
	if err1 != nil {
		return "Observations could not be retrieved"
	}

	r, err2 := json.Marshal(result)
	if err2 != nil {
		return "json error"
	}

	return string(r)
}

func GetRecentObservationsHttp(res http.ResponseWriter, req *http.Request) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	observations := []Observation{}
	err := DB.Order("created desc").Limit(5).Find(&observations).Error
	if err != nil {
		return "not found"
	}

	var tmpCO CommunityObservation
	var communityObservations []CommunityObservation
	for _, o := range observations {
		user := User{}
		err1 := DB.Where("uid= ?", o.Uid).Find(&user).Error
		if err1 != nil {
			return "not found"
		}
		tmpCO.Username = user.Username
		tmpCO.Avatar = user.Avatar
		tmpCO.Comment = o.Comment
		tmpCO.PatternID = o.DiscoveredId

		discovered := Discovered{}
		err = DB.Where("discovered_id = ?", o.DiscoveredId).Find(&discovered).Error
		if err != nil {
			return "can't find dicovered to generate link"
		}
		if discovered.CorrelationId == 0 {
			tmpCO.Link = "charts/related/" + discovered.RelationId
		} else {
			tmpCO.Link = "chartcorrelated/" + strconv.Itoa(discovered.CorrelationId)
		}

		tmpCO.EmailMD5 = GetMD5Hash(user.Email)
		communityObservations = append(communityObservations, tmpCO)
	}

	r, err2 := json.Marshal(communityObservations)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return "Unable to parse JSON"
	}

	if r == nil {
		return "No observations have been made yet"
	} else {
		return string(r)
	}
}

func FlagObservationHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	if params["id"] == "" {
		http.Error(res, "Missing id parameter.", http.StatusBadRequest)
		return ""
	}

	id, _ := strconv.Atoi(params["id"])

	err := DB.Model(Observation{}).Where("observation_id= ?", id).Update("flagged", true).Error
	if err != nil {
		http.Error(res, "Missing session parameter.", http.StatusInternalServerError)
		return ""
	}

	return "success"
}
