package main

import (
	"encoding/json"
	// "fmt"
	"github.com/jinzhu/gorm"
	"github.com/pmylund/sortutil"
	"net/http"
	"strconv"
	"time"
)

type HomeData struct {
	Label string
	Value int
}

type ProfileDiscovery struct {
	PatternId     int       `json:"patternid"`
	Title         string    `json:"title"`
	ApiString     string    `json:"apistring"`
	DiscoveryDate time.Time `json:"discoverydate"`
}

type ProfileObservation struct {
	ObservationId int       `json:"observationid"`
	Title         string    `json:"charttitle"`
	ApiString     string    `json:"apistring"`
	DiscoveryDate time.Time `json:"discoverydate"`
	Comment       string    `json:"comment"`
}

type DataExpert struct {
	Username   string `json:"username"`
	Avatar     string `json:"avatar"`
	Reputation int    `json:"reputation"`
	EmailMD5   []byte `json:"MD5email"`
}

type UserActivity struct {
	ActivityStr1 string    `json:"string"`
	PatternId    int       `json:"patternid"`
	Created      float64   `json:"-"`
	Time         time.Time `json:"time"`
}

func ActivityCheck(a string) string {
	switch a {
	case "c":
		return "Comment"
	case "ic":
		return "Invalidated Chart"
	case "vc":
		return "Validated Chart"
	case "io":
		return "Invalidated Observation"
	case "vo":
		return "Validated Observation"
	default:
		return "Unknown"
	}
}

func AddActivity(uid int, atype string, ts time.Time, disid int, obsid int) *appError {
	act := Activity{
		Uid:     uid,
		Type:    ActivityCheck(atype),
		Created: ts,
	}

	err := DB.Save(&act).Error
	if err != nil {
		return &appError{err, "Database query failed (Save)", http.StatusInternalServerError}
	}

	return nil
}

func GetProfileObservationsHttp(res http.ResponseWriter, req *http.Request) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	observation := []Observation{}
	err1 := DB.Where("uid = ?", uid).Find(&observation).Error
	if err1 != nil && err1 != gorm.RecordNotFound {
		http.Error(res, "Database query failed", http.StatusInternalServerError)
		return ""
	}

	profileObservations := make([]ProfileObservation, 0)

	for _, o := range observation {
		var tmp ProfileObservation
		tmp.ObservationId = o.ObservationId
		tmp.DiscoveryDate = o.Created
		tmp.Comment = o.Comment

		discTmp := Discovered{}
		err := DB.Where("discovered_id = ?", o.DiscoveredId).Find(&discTmp).Error
		if err != nil && err != gorm.RecordNotFound {
			http.Error(res, "Database query failed", http.StatusInternalServerError)
			return ""
		}

		if discTmp.CorrelationId == 0 {
			tmp.ApiString = "chart/related/" + discTmp.RelationId
			var td TableData
			json.Unmarshal(discTmp.Json, &td)
			tmp.Title = td.Title
		} else {
			cid := strconv.Itoa(discTmp.CorrelationId)
			tmp.ApiString = "chart/correlated/" + cid
			var cd CorrelationData
			json.Unmarshal(discTmp.Json, &cd)
			tmp.Title = cd.Table1.Title + " correlated with " + cd.Table2.Title
			if cd.Table3.Title != "" {
				tmp.Title += " correlated with " + cd.Table3.Title
			}
		}

		profileObservations = append(profileObservations, tmp)
	}

	r, err2 := json.Marshal(profileObservations)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func GetDiscoveriesHttp(res http.ResponseWriter, req *http.Request) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return "Could not validate user"
	}

	var discovered []Discovered
	err1 := DB.Where("uid = ?", uid).Find(&discovered).Error
	if err1 != nil {
		return "not found"
	}

	profileDiscoveries := make([]ProfileDiscovery, 0)

	for _, d := range discovered {
		var tmp ProfileDiscovery
		tmp.PatternId = d.DiscoveredId
		tmp.DiscoveryDate = d.Created

		if d.CorrelationId == 0 {
			tmp.ApiString = "chart/related/" + d.RelationId
			var td TableData
			json.Unmarshal(d.Json, &td)
			tmp.Title = td.Title
		} else {
			cid := strconv.Itoa(d.CorrelationId)
			tmp.ApiString = "chart/correlated/" + cid
			var cd CorrelationData
			json.Unmarshal(d.Json, &cd)
			tmp.Title = cd.Table1.Title + " correlated with " + cd.Table2.Title
			if cd.Table3.Title != "" {
				tmp.Title += " correlated with " + cd.Table3.Title
			}
		}

		profileDiscoveries = append(profileDiscoveries, tmp)
	}

	r, err2 := json.Marshal(profileDiscoveries)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return "Unable to parse JSON"
	}

	return string(r)
}

func GetValidatedDiscoveriesHttp(res http.ResponseWriter, req *http.Request) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return "Could not validate user"
	}

	discovered := []Discovered{}
	err1 := DB.Where("uid = ?", uid).Where("valid > ?", 0).Find(&discovered).Error
	if err1 != nil && err1 != gorm.RecordNotFound {
		http.Error(res, "Database query failed", http.StatusInternalServerError)
		return ""
	}

	profileDiscoveries := make([]ProfileDiscovery, 0)

	for _, d := range discovered {
		var tmp ProfileDiscovery
		tmp.PatternId = d.DiscoveredId
		tmp.DiscoveryDate = d.Created

		if d.CorrelationId == 0 {
			tmp.ApiString = "chart/" + d.RelationId
			var td TableData
			json.Unmarshal(d.Json, &td)
			tmp.Title = td.Title
		} else {
			cid := strconv.Itoa(d.CorrelationId)
			tmp.ApiString = "chartcorrelated/" + cid
			var cd CorrelationData
			json.Unmarshal(d.Json, &cd)
			tmp.Title = cd.Table1.Title + " correlated with " + cd.Table2.Title
			if cd.Table3.Title != "" {
				tmp.Title += " correlated with " + cd.Table3.Title
			}
		}

		profileDiscoveries = append(profileDiscoveries, tmp)
	}

	r, err2 := json.Marshal(profileDiscoveries)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func GetHomePageDataHttp(res http.ResponseWriter, req *http.Request) string {
	// session := req.Header.Get("X-API-SESSION")
	// if len(session) <= 0 {
	// 	http.Error(res, "Missing session parameter", http.StatusBadRequest)
	// 	return "Missing session parameter"
	// }

	var result [3]HomeData
	result[0].Label = "players"
	result[1].Label = "discoveries"
	result[2].Label = "datasets"
	err := DB.Model(User{}).Count(result[0].Value).Error
	if err != nil {
		return "not found"
	}

	err = DB.Model(Discovered{}).Count(result[1].Value).Error
	if err != nil {
		return "not found"
	}

	err = DB.Model(OnlineData{}).Count(result[2].Value).Error
	if err != nil {
		return "not found"
	}
	r, err2 := json.Marshal(result)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return "Unable to parse JSON"
	}

	return string(r)
}

func GetReputationHttp(res http.ResponseWriter, req *http.Request) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return "Could not validate user"
	}

	var rep int
	err1 := DB.Model(User{}).Select("reputation").Where("uid = ?", uid).Find(&rep).Error
	if err1 != nil {
		return "not found"
	}

	return string(rep)
}

func GetAmountDiscoveriesHttp(res http.ResponseWriter, req *http.Request) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return "Could not validate user"
	}

	count := 0
	err1 := DB.Model(Discovered{}).Where("uid = ?", uid).Count(&count).Error
	if err1 != nil {
		return "not found"
	}

	return string(count)
}

func GetDataExpertsHttp(res http.ResponseWriter, req *http.Request) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	users := []User{}
	err1 := DB.Order("reputation desc").Limit(5).Find(&users).Error
	if err1 != nil {
		return "not found"
	}

	var tmpDE DataExpert
	var dataExperts []DataExpert
	for _, u := range users {
		tmpDE.Username = u.Username
		tmpDE.Avatar = u.Avatar
		tmpDE.Reputation = u.Reputation
		tmpDE.EmailMD5 = Hash(u.Email)
		dataExperts = append(dataExperts, tmpDE)
	}

	r, err2 := json.Marshal(dataExperts)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return "Unable to parse JSON"
	}

	return string(r)
}

func GetActivityStreamHttp(res http.ResponseWriter, req *http.Request) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return "Could not validate user"
	}

	var activities []UserActivity
	t := time.Now()
	activities = AddHappenedTo(uid, activities, t)
	activities = AddInstigated(uid, activities, t)
	sortutil.AscByField(activities, "Created")

	n := 5
	if len(activities) < 5 {
		n = len(activities)
	}

	r, err2 := json.Marshal(activities[:n])
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return "Unable to parse JSON"
	}

	return string(r)
}

func AddInstigated(uid int, activities []UserActivity, t time.Time) []UserActivity {
	activity := []Activity{}

	err := DB.Order("created desc").Where("uid = ?", uid).Find(&activity).Error
	if err != nil {
		return activities
	}

	for _, a := range activity {
		tmpA := UserActivity{}
		if a.Type == "Comment" {
			tmpA.ActivityStr1 = "You commented on pattern "
			tmpA.PatternId = a.DiscoveredId
		} else if a.Type == "Validated Observation" {
			obs := Observation{}
			err = DB.Where("observation_id = ?", a.ObservationId).Find(&obs).Error
			if err != nil {
				tmpA.ActivityStr1 = "Bad validated observation activity 1"
				tmpA.PatternId = 0
			}
			user := User{}
			err = DB.Where("uid = ?", obs.Uid).Find(&user).Error
			if err != nil {
				tmpA.ActivityStr1 = "Bad validated observation activity 2"
				tmpA.PatternId = 0
			}
			tmpA.ActivityStr1 = "You agreed with " + user.Username + "'s observation on pattern "
			tmpA.PatternId = obs.DiscoveredId
		} else if a.Type == "Invalidated Observation" {
			obs := Observation{}
			err := DB.Where("observation_id = ?", a.ObservationId).Find(&obs).Error
			if err != nil {
				tmpA.ActivityStr1 = "Bad invalidated observation activity 1"
			}
			user := User{}
			err = DB.Where("uid = ?", obs.Uid).Find(&user).Error
			if err != nil {
				tmpA.ActivityStr1 = "Bad invalidated observation activity 2"
			}
			tmpA.ActivityStr1 = "You disagreed with " + user.Username + "'s observation on pattern "
			tmpA.PatternId = obs.DiscoveredId
		} else if a.Type == "Validated Chart" {
			tmpA.ActivityStr1 = "You validated pattern "
			tmpA.PatternId = a.DiscoveredId

		} else if a.Type == "Invalidated Chart" {
			tmpA.ActivityStr1 = "You invalidated pattern "
			tmpA.PatternId = a.DiscoveredId

		} else {
			tmpA.ActivityStr1 = "No activity"
			tmpA.PatternId = 0
		}
		tmpA.Created = t.Sub(a.Created).Seconds()
		tmpA.Time = a.Created
		activities = append(activities, tmpA)
	}
	return activities

}

func AddHappenedTo(uid int, activities []UserActivity, t time.Time) []UserActivity {
	vDisc := []Validation{}

	err = DB.Select("priv_validations.discovered_id, priv_validations.created, priv_validations.uid, priv_validations.valflag").Joins("LEFT JOIN priv_discovered AS d ON priv_validations.discovered_id = d.discovered_id").Where("d.uid = ?", uid).Where("priv_validations.discovered_id > ?", 0).Order("priv_validations.created DESC").Find(&vDisc).Error
	if err != nil && err != gorm.RecordNotFound {
		return activities
	}

	vObs := []struct {
		Validation
		Comment string
		Did     int
	}{}

	err = DB.Select("o.discovered_id as did, o.comment as comment, priv_validations.created, priv_validations.uid, priv_validations.valflag").Joins("LEFT JOIN priv_observations AS o ON priv_validations.observation_id = o.observation_id").Where("o.uid = ?", uid).Where("priv_validations.observation_id > ?", 0).Order("priv_validations.created DESC").Find(&vObs).Error
	if err != nil && err != gorm.RecordNotFound {
		return activities
	}

	activity := []Activity{}

	err = DB.Select("priv_activity.discovered_id, priv_activity.created, priv_activity.uid").Joins("LEFT JOIN priv_discovered as d ON priv_activity.discovered_id = d.discovered_id").Where("d.uid = ?", uid).Where("priv_activity.type = ?", "Comment").Order("priv_activity.created DESC").Find(&activity).Error
	if err != nil && err != gorm.RecordNotFound {
		return activities
	}

	for _, d := range vDisc {
		tmpA := UserActivity{}
		user := User{}
		err = DB.Where("uid = ?", d.Uid).Find(&user).Error
		if err != nil {
			tmpA.ActivityStr1 = "Bad invalidated observation activity 2"
		}

		if d.Valflag == true {
			tmpA.ActivityStr1 = "You gained " + strconv.Itoa(discVal) + " reputation when " + user.Username + " validated your pattern "
			tmpA.PatternId = d.DiscoveredId
			tmpA.Created = t.Sub(d.Created).Seconds()
			tmpA.Time = d.Created
		} else {
			tmpA.ActivityStr1 = "You lost " + strconv.Itoa(discInval) + " reputation when " + user.Username + " invalidated your pattern "
			tmpA.PatternId = d.DiscoveredId
			tmpA.Created = t.Sub(d.Created).Seconds()
			tmpA.Time = d.Created
		}

		activities = append(activities, tmpA)
	}

	for _, o := range vObs {
		tmpA := UserActivity{}
		user := User{}
		err = DB.Where("uid = ?", o.Uid).Find(&user).Error
		if err != nil {
			tmpA.ActivityStr1 = "Bad invalidated observation activity 2"
		}

		if o.Valflag == true {
			tmpA.ActivityStr1 = "You gained " + strconv.Itoa(obsVal) + " reputation when " + user.Username + " validated your observation on pattern "
			tmpA.PatternId = o.Did
			tmpA.Created = t.Sub(o.Created).Seconds()
			tmpA.Time = o.Created
		} else {
			tmpA.ActivityStr1 = "You lost " + strconv.Itoa(obsInval) + " reputation when " + user.Username + " invalidated your observation on pattern "
			tmpA.PatternId = o.Did
			tmpA.Created = t.Sub(o.Created).Seconds()
			tmpA.Time = o.Created
		}

		activities = append(activities, tmpA)
	}

	for _, a := range activity {
		tmpA := UserActivity{}
		user := User{}
		err = DB.Where("uid = ?", a.Uid).Find(&user).Error
		if err != nil {
			tmpA.ActivityStr1 = "Bad invalidated observation activity 2"
		}

		tmpA.ActivityStr1 = "You gained " + strconv.Itoa(discObs) + " reputation when " + user.Username + " commented on your pattern "
		tmpA.PatternId = a.DiscoveredId
		tmpA.Created = t.Sub(a.Created).Seconds()
		tmpA.Time = a.Created
		activities = append(activities, tmpA)
	}

	return activities

}
