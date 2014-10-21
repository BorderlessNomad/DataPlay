package main

import (
	"encoding/json"
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
	EmailMD5   string `json:"MD5email"`
}

type UserActivity struct {
	Activity  string    `json:"activitystring"`
	Link      string    `json:"linkstring"`
	PatternId int       `json:"patternid"`
	Created   float64   `json:"-"`
	Time      time.Time `json:"time"`
}

func ActivityCheck(a string) string {
	switch a {
	case "c":
		return "Comment"
	case "ic":
		return "Discredited Chart"
	case "vc":
		return "Credited Chart"
	case "io":
		return "Discredited Observation"
	case "vo":
		return "Credited Observation"
	default:
		return "Unknown"
	}
}

func AddActivity(uid int, atype string, ts time.Time, disid int, obsid int) *appError {
	act := Activity{
		Uid:           uid,
		Type:          ActivityCheck(atype),
		Created:       ts,
		DiscoveredId:  disid,
		ObservationId: obsid,
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
		return ""
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	var discovered []Discovered
	err1 := DB.Where("uid = ?", uid).Find(&discovered).Error
	if err1 != nil && err1 != gorm.RecordNotFound {
		http.Error(res, "Database query failed! (Discovered)", http.StatusInternalServerError)
		return ""
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
		return ""
	}

	return string(r)
}

func GetCreditedDiscoveriesHttp(res http.ResponseWriter, req *http.Request) string {
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

	discovered := []Discovered{}
	err1 := DB.Where("uid = ?", uid).Where("credited > ?", 0).Find(&discovered).Error
	if err1 != nil && err1 == gorm.RecordNotFound {
		return "this user has yet to make any discoveries"
	} else if err1 != nil {
		http.Error(res, "Database query failed! (Discovered)", http.StatusInternalServerError)
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
	var result [3]HomeData
	result[0].Label = "players"
	result[0].Value = 0

	result[1].Label = "discoveries"
	result[1].Value = 0

	result[2].Label = "datasets"
	result[2].Value = 0

	err := DB.Model(User{}).Count(&result[0].Value).Error
	if err != nil && err != gorm.RecordNotFound {
		http.Error(res, "Database query failed! (User)", http.StatusInternalServerError)
		return ""
	}

	err = DB.Model(Discovered{}).Count(&result[1].Value).Error
	if err != nil && err != gorm.RecordNotFound {
		http.Error(res, "Database query failed! (Discovered)", http.StatusInternalServerError)
		return ""
	}

	err = DB.Model(OnlineData{}).Count(&result[2].Value).Error
	if err != nil && err != gorm.RecordNotFound {
		http.Error(res, "Database query failed! (OnlineData)", http.StatusInternalServerError)
		return ""
	}

	r, err2 := json.Marshal(result)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func GetReputationHttp(res http.ResponseWriter, req *http.Request) string {
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

	rep := 0
	err1 := DB.Model(User{}).Select("reputation").Where("uid = ?", uid).Find(&rep).Error
	if err1 != nil && err1 != gorm.RecordNotFound {
		http.Error(res, "Database query failed! (User)", http.StatusInternalServerError)
		return ""
	}

	return string(rep)
}

func GetAmountDiscoveriesHttp(res http.ResponseWriter, req *http.Request) string {
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

	count := 0
	err1 := DB.Model(Discovered{}).Where("uid = ?", uid).Count(&count).Error
	if err1 != nil && err1 != gorm.RecordNotFound {
		http.Error(res, "Database query failed! (Discovered)", http.StatusInternalServerError)
		return ""
	}

	return string(count)
}

func GetDataExpertsHttp(res http.ResponseWriter, req *http.Request) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	users := []User{}
	err1 := DB.Order("reputation DESC").Limit(5).Find(&users).Error
	if err1 != nil && err1 != gorm.RecordNotFound {
		http.Error(res, "Database query failed! (User)", http.StatusInternalServerError)
		return ""
	}

	var tmpDE DataExpert
	var dataExperts []DataExpert
	for _, u := range users {
		tmpDE.Username = u.Username
		tmpDE.Avatar = u.Avatar
		tmpDE.Reputation = u.Reputation
		tmpDE.EmailMD5 = GetMD5Hash(u.Email)
		dataExperts = append(dataExperts, tmpDE)
	}

	r, err2 := json.Marshal(dataExperts)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func GetActivityStreamHttp(res http.ResponseWriter, req *http.Request) string {
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

	var activities []UserActivity
	t := time.Now()
	activities = AddHappenedTo(uid, activities, t)
	activities = AddInstigated(uid, activities, t)
	sortutil.AscByField(activities, "Created")

	n := 5
	if len(activities) < n {
		n = len(activities)
	}

	r, err2 := json.Marshal(activities[:n])
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func GetActivityStreamQ(params map[string]string) string {
	if params["user"] == "" {
		return "no user id"
	}

	uid, e := strconv.Atoi(params["user"])
	if e != nil {
		return e.Error()
	}

	var activities []UserActivity
	t := time.Now()
	activities = AddHappenedTo(uid, activities, t)
	activities = AddInstigated(uid, activities, t)
	sortutil.AscByField(activities, "Created")

	n := 5
	if len(activities) < n {
		n = len(activities)
	}

	r, err := json.Marshal(activities[:n])
	if err != nil {
		return err.Error()
	}

	return string(r)
}

func AddInstigated(uid int, activities []UserActivity, t time.Time) []UserActivity {
	activity := []Activity{}

	err := DB.Where("uid = ?", uid).Order("created DESC").Find(&activity).Error
	if err != nil {
		return activities
	}

	for _, a := range activity {
		tmpA := UserActivity{}

		if a.Type == "Comment" {
			tmpA.Activity = "You commented on pattern "
			tmpA.PatternId = a.DiscoveredId

			discovered := Discovered{}
			err = DB.Where("discovered_id = ?", a.DiscoveredId).Find(&discovered).Error
			if err != nil {
				return activities
			}

			if discovered.CorrelationId == 0 {
				tmpA.Link = "charts/related/" + discovered.RelationId
			} else {
				tmpA.Link = "chartcorrelated/" + strconv.Itoa(discovered.CorrelationId)
			}
		} else if a.Type == "Credited Observation" {
			obs := Observation{}
			err = DB.Where("observation_id = ?", a.ObservationId).Find(&obs).Error
			if err != nil {
				tmpA.Activity = "Bad credited observation activity 1"
				tmpA.PatternId = 0
			}

			user := User{}
			err = DB.Where("uid = ?", obs.Uid).Find(&user).Error
			if err != nil {
				tmpA.Activity = "Bad credited observation activity 2"
				tmpA.PatternId = 0
			}

			tmpA.Activity = "You agreed with " + user.Username + "'s observation on pattern "

			tmpA.PatternId = obs.DiscoveredId
			discovered := Discovered{}
			err = DB.Where("discovered_id = ?", obs.DiscoveredId).Find(&discovered).Error
			if err != nil {
				return activities
			}

			if discovered.CorrelationId == 0 {
				tmpA.Link = "charts/related/" + discovered.RelationId
			} else {
				tmpA.Link = "chartcorrelated/" + strconv.Itoa(discovered.CorrelationId)
			}
		} else if a.Type == "Discredited Observation" {
			obs := Observation{}
			err := DB.Where("observation_id = ?", a.ObservationId).Find(&obs).Error
			if err != nil {
				tmpA.Activity = "Bad discredited observation activity 1"
			}

			user := User{}
			err = DB.Where("uid = ?", obs.Uid).Find(&user).Error
			if err != nil {
				tmpA.Activity = "Bad discredited observation activity 2"
			}

			tmpA.Activity = "You disagreed with " + user.Username + "'s observation on pattern "
			tmpA.PatternId = obs.DiscoveredId
			discovered := Discovered{}
			err = DB.Where("discovered_id = ?", obs.DiscoveredId).Find(&discovered).Error
			if err != nil {
				return activities
			}

			if discovered.CorrelationId == 0 {
				tmpA.Link = "charts/related/" + discovered.RelationId
			} else {
				tmpA.Link = "chartcorrelated/" + strconv.Itoa(discovered.CorrelationId)
			}
		} else if a.Type == "Credited Chart" {
			tmpA.Activity = "You credited pattern "
			tmpA.PatternId = a.DiscoveredId

			discovered := Discovered{}
			err = DB.Where("discovered_id = ?", a.DiscoveredId).Find(&discovered).Error
			if err != nil {
				return activities
			}

			if discovered.CorrelationId == 0 {
				tmpA.Link = "charts/related/" + discovered.RelationId
			} else {
				tmpA.Link = "chartcorrelated/" + strconv.Itoa(discovered.CorrelationId)
			}
		} else if a.Type == "Discredited Chart" {
			tmpA.Activity = "You discredited pattern "
			tmpA.PatternId = a.DiscoveredId
			discovered := Discovered{}
			err = DB.Where("discovered_id = ?", a.DiscoveredId).Find(&discovered).Error
			if err != nil {
				return activities
			}

			if discovered.CorrelationId == 0 {
				tmpA.Link = "charts/related/" + discovered.RelationId
			} else {
				tmpA.Link = "chartcorrelated/" + strconv.Itoa(discovered.CorrelationId)
			}
		} else {
			tmpA.Activity = "No activity"
			tmpA.PatternId = 0
		}

		tmpA.Created = t.Sub(a.Created).Seconds()
		tmpA.Time = a.Created
		activities = append(activities, tmpA)
	}

	return activities
}

func AddHappenedTo(uid int, activities []UserActivity, t time.Time) []UserActivity {
	vDisc := []Credit{}

	err = DB.Select("priv_credits.discovered_id, priv_credits.created, priv_credits.uid, priv_credits.credflag").Joins("LEFT JOIN priv_discovered AS d ON priv_credits.discovered_id = d.discovered_id").Where("d.uid = ?", uid).Where("priv_credits.discovered_id > ?", 0).Order("priv_credits.created DESC").Find(&vDisc).Error
	if err != nil && err != gorm.RecordNotFound {
		return activities
	}

	vObs := []struct {
		Credit
		Comment string
		Did     int
	}{}

	err = DB.Select("o.discovered_id as did, o.comment as comment, priv_credits.created, priv_credits.uid, priv_credits.credflag").Joins("LEFT JOIN priv_observations AS o ON priv_credits.observation_id = o.observation_id").Where("o.uid = ?", uid).Where("priv_credits.observation_id > ?", 0).Order("priv_credits.created DESC").Find(&vObs).Error
	if err != nil {
		return activities
	}

	activity := []Activity{}

	err = DB.Select("priv_activity.discovered_id, priv_activity.created, priv_activity.uid").Joins("LEFT JOIN priv_discovered as d ON priv_activity.discovered_id = d.discovered_id").Where("d.uid = ?", uid).Where("priv_activity.type = ?", "Comment").Order("priv_activity.created DESC").Find(&activity).Error
	if err != nil {
		return activities
	}

	for _, d := range vDisc {
		tmpA := UserActivity{}
		user := User{}
		err = DB.Where("uid = ?", d.Uid).Find(&user).Error
		if err != nil {
			tmpA.Activity = "Bad discredited observation activity 2"
		}

		if d.Credflag == true {
			tmpA.Activity = "You gained " + strconv.Itoa(discCredit) + " reputation when " + user.Username + " credited your pattern "
			tmpA.PatternId = d.DiscoveredId
			tmpA.Created = t.Sub(d.Created).Seconds()
			tmpA.Time = d.Created
		} else {
			tmpA.Activity = "You lost " + strconv.Itoa(discDiscredit) + " reputation when " + user.Username + " discredited your pattern "
			tmpA.PatternId = d.DiscoveredId
			tmpA.Created = t.Sub(d.Created).Seconds()
			tmpA.Time = d.Created
		}

		discovered := Discovered{}
		err = DB.Where("discovered_id = ?", d.DiscoveredId).Find(&discovered).Error
		if err != nil {
			return activities
		}

		if discovered.CorrelationId == 0 {
			tmpA.Link = "charts/related/" + discovered.RelationId
		} else {
			tmpA.Link = "chartcorrelated/" + strconv.Itoa(discovered.CorrelationId)
		}

		activities = append(activities, tmpA)
	}

	for _, o := range vObs {
		tmpA := UserActivity{}
		user := User{}
		err = DB.Where("uid = ?", o.Uid).Find(&user).Error
		if err != nil {
			tmpA.Activity = "Bad discredited observation activity 2"
		}

		if o.Credflag == true {
			tmpA.Activity = "You gained " + strconv.Itoa(obsCredit) + " reputation when " + user.Username + " credited your observation on pattern "
			tmpA.PatternId = o.Did
			tmpA.Created = t.Sub(o.Created).Seconds()
			tmpA.Time = o.Created
		} else {
			tmpA.Activity = "You lost " + strconv.Itoa(obsDiscredit) + " reputation when " + user.Username + " discredited your observation on pattern "
			tmpA.PatternId = o.Did
			tmpA.Created = t.Sub(o.Created).Seconds()
			tmpA.Time = o.Created
		}

		discovered := Discovered{}
		err = DB.Where("discovered_id = ?", o.Did).Find(&discovered).Error
		if err != nil {
			return activities
		}

		if discovered.CorrelationId == 0 {
			tmpA.Link = "charts/related/" + discovered.RelationId
		} else {
			tmpA.Link = "chartcorrelated/" + strconv.Itoa(discovered.CorrelationId)
		}

		activities = append(activities, tmpA)
	}

	for _, a := range activity {
		tmpA := UserActivity{}
		user := User{}
		err = DB.Where("uid = ?", a.Uid).Find(&user).Error
		if err != nil {
			tmpA.Activity = "Bad discredited observation activity 2"
		}

		tmpA.Activity = "You gained " + strconv.Itoa(discObs) + " reputation when " + user.Username + " commented on your pattern "
		tmpA.PatternId = a.DiscoveredId
		tmpA.Created = t.Sub(a.Created).Seconds()
		tmpA.Time = a.Created

		discovered := Discovered{}
		err = DB.Where("discovered_id = ?", a.DiscoveredId).Find(&discovered).Error
		if err != nil {
			return activities
		}

		if discovered.CorrelationId == 0 {
			tmpA.Link = "charts/related/" + discovered.RelationId
		} else {
			tmpA.Link = "chartcorrelated/" + strconv.Itoa(discovered.CorrelationId)
		}

		activities = append(activities, tmpA)
	}

	return activities
}
