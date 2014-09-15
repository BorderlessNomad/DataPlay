package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/gocql/gocql"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const numdays = 30

var Today = time.Date(2010, 2, 1, 0, 0, 0, 0, time.UTC)
var FromDate = Today.AddDate(0, 0, -numdays)

type PoliticalActivity struct {
	Term     string       `json:"term"`
	Mentions [numdays]int `json:"mentionsperday"`
	Val      int          `json:"-"`
}

type DateID struct {
	ID   []byte
	Date time.Time
}

type KeyEnt struct {
	Name string
	Date time.Time
	ID   []byte
}

func WriteCass() {
	session, _ := GetCassandraConnection("dp") // create connection to cassandra
	defer session.Close()
	url := ""

	iter := session.Query(`SELECT original_url FROM response`).Iter()
	f, _ := os.OpenFile("dat1.txt", os.O_RDWR|os.O_APPEND, 0666)

	for iter.Scan(&url) {
		u := []byte(url + "\n")
		f.Write(u)
	}

}

func DepartmentsPoliticalActivity() []PoliticalActivity {
	var dept []Departments // get all departments from postgres sql table
	var deptNames []string
	err := DB.Find(&dept).Error

	if err != nil && err != gorm.RecordNotFound {
		return nil
	}

	for _, d := range dept {
		deptNames = append(deptNames, d.GovDept)
	}

	pa := make([]PoliticalActivity, len(dept))
	keyEnt := keywords(deptNames)

	for i, d := range dept {
		pa[i].Term = d.GovDept
		for _, ke := range keyEnt {
			if d.GovDept == ke.Name {
				dayindex := int((Today.Round(time.Hour).Sub(ke.Date.Round(time.Hour)) / 24).Hours() - 1) // get day index
				pa[i].Mentions[dayindex]++
			}
		}
	}

	return RankPA(pa)
}

// func EventsPoliticalActivity() []PoliticalActivity {
// 	var event []Events
// 	err := DB.Find(&event).Error

// 	if err != nil && err != gorm.RecordNotFound {
// 		return nil
// 	}

// 	pa := make([]PoliticalActivity, len(event))

// 	keyEnt := keywords()

// 	for i, e := range event {
// 		pa[i].Term = e.Event
// 		for _, ke := range keyEnt {
// 			if e.Event == ke.Name {
// 				dayindex := int((Today.Round(time.Hour).Sub(ke.Date.Round(time.Hour)) / 24).Hours() - 1) // get day index
// 				pa[i].Mentions[dayindex]++
// 			}
// 		}
// 	}

// 	return RankPA(pa)
// }

// func RegionsPoliticalActivity() []PoliticalActivity {
// 	var region []Regions
// 	err := DB.Find(&region).Error

// 	if err != nil && err != gorm.RecordNotFound {
// 		return nil
// 	}
// 	pa := make([]PoliticalActivity, len(region))

// 	keyEnt := keywords()

// 	for i, r := range region {
// 		pa[i].Term = r.Town
// 		for _, ke := range keyEnt {
// 			if r.Town == ke.Name {
// 				dayindex := int((Today.Round(time.Hour).Sub(ke.Date.Round(time.Hour)) / 24).Hours() - 1) // get day index
// 				pa[i].Mentions[dayindex]++
// 			}
// 		}
// 	}

// 	return RankPA(pa)
// }

func GetPoliticalActivityHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	var result []PoliticalActivity

	if params["type"] == "d" {
		result = DepartmentsPoliticalActivity()
		// } else if params["type"] == "e" {
		// 	result = EventsPoliticalActivity()
		// } else if params["type"] == "r" {
		// 	result = RegionsPoliticalActivity()
	} else {
		http.Error(res, "Bad type param", http.StatusInternalServerError)
		return "Bad type param"
	}

	r, err := json.Marshal(result)
	if err != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return "Unable to parse JSON"
	}

	return string(r)
}

func GetCassandraConnection(keyspace string) (*gocql.Session, error) {
	cassandraHost := "10.0.0.2"
	cassandraPort := 49236
	if os.Getenv("cassandrahost") != "" {
		cassandraHost = os.Getenv("cassandrahost")
	}
	if os.Getenv("cassandraport") != "" {
		cassandraPort, _ = strconv.Atoi(os.Getenv("cassandraport"))
	}

	cluster := gocql.NewCluster(cassandraHost)
	cluster.Port = cassandraPort
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()

	if err != nil {
		fmt.Println("Could not connect to the Cassandara server.")
	}

	return session, err
}

func keywords(names []string) []KeyEnt {
	var id []byte
	var dateID []string
	var queryDate time.Time
	var keyEnt []KeyEnt
	var tmpKE KeyEnt

	session, _ := GetCassandraConnection("dp") // create connection to cassandra
	defer session.Close()

	// add all dated dateID between -n days and today to array
	iter := session.Query(`SELECT id, date FROM response WHERE date >= ? AND date < ? ALLOW FILTERING`, FromDate, Today).Iter()
	for iter.Scan(&id, &queryDate) {
		dateID = append(dateID, string(id[:len(id)])+"!"+queryDate.Format(time.RFC3339))
	}

	if err := iter.Close(); err != nil {
		///return err
	}

	for _, n := range names {
		iter := session.Query(`SELECT id FROM keyword WHERE name = ?`, n).Iter()
		for iter.Scan(&id) {
			var date time.Time
			date = DateAndId(id, dateID)
			if date.Year() < 2 {
				tmpKE.Name = n
				tmpKE.ID = id
				tmpKE.Date = date
				keyEnt = append(keyEnt, tmpKE)
			}
		}
	}

	for _, n := range names {
		iter := session.Query(`SELECT id FROM entity WHERE name = ?`, n).Iter()
		for iter.Scan(&id) {
			var date time.Time
			date = DateAndId(id, dateID)
			if date.Year() < 2 {
				tmpKE.Name = n
				tmpKE.ID = id
				tmpKE.Date = date
				keyEnt = append(keyEnt, tmpKE)
			}
		}
	}

	return keyEnt
}

func DateAndId(id []byte, dateID []string) time.Time {
	var t time.Time
	for _, d := range dateID {
		split := strings.Split(d, "!")
		if string(id[:len(id)]) == split[0] {
			t, _ = time.Parse(time.RFC3339, split[1])
			return t
		}
	}
	return t
}

// sort PA array and return top 15
func RankPA(activities []PoliticalActivity) []PoliticalActivity {

	for i, _ := range activities {
		total := 0
		for j, _ := range activities[i].Mentions {
			total += activities[i].Mentions[j]
		}
		activities[i].Val = total
	}

	n := len(activities)
	chk := true
	var tmp PoliticalActivity

	for chk == true {
		newn := 0

		for i := 1; i < n; i++ {
			if activities[i].Val > activities[i-1].Val {
				tmp = activities[i]
				activities[i] = activities[i-1]
				activities[i-1] = tmp
				newn = i
			}
		}
		n = newn

		if n == 0 {
			chk = false
		}
	}

	return activities[0:15]
}
