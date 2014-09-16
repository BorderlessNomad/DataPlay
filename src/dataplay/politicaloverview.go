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

type DatedTerm struct {
	Term string
	Date time.Time
	ID   []byte
}

// gets names of all departments, checks for mentions in specified time period and returns ranked array of 15 most popular terms and their 30 day frequencies
func DepartmentsPoliticalActivity() []PoliticalActivity {
	var dept []Departments // get all departments from postgres sql table
	var terms []string
	err := DB.Find(&dept).Error

	if err != nil && err != gorm.RecordNotFound {
		return nil
	}

	for _, d := range dept {
		terms = append(terms, d.GovDept)
	}

	return CheckThese(terms)
}

// gets names of all events, checks for mentions in specified time period and returns ranked array of 15 most popular terms and their 30 day frequencies
func EventsPoliticalActivity() []PoliticalActivity {
	var event []Events
	var terms []string
	err := DB.Find(&event).Error

	if err != nil && err != gorm.RecordNotFound {
		return nil
	}

	for _, e := range event {
		terms = append(terms, e.Event)
	}

	return CheckThese(terms)
}

// gets names of all regions, checks for mentions in specified time period and returns ranked array of 15 most popular terms and their 30 day frequencies
func RegionsPoliticalActivity() []PoliticalActivity {
	var region []Regions
	var terms []string
	err := DB.Find(&region).Error

	if err != nil && err != gorm.RecordNotFound {
		return nil
	}

	for _, r := range region {
		terms = append(terms, r.Town)
	}

	return CheckThese(terms)
}

// takes slice of terms, checks for the total number of occurences and returns a top 15 ranked array
func CheckThese(terms []string) []PoliticalActivity {
	politicalActivity := make([]PoliticalActivity, len(terms))
	DatedTerm := keywords(terms) // returns array

	for i, term := range terms { // check each term
		politicalActivity[i].Term = term // copy to politicalActivity array
		for _, dt := range DatedTerm {   // check through all dated terms
			if term == dt.Term { //if there's a match
				dayindex := int((Today.Round(time.Hour).Sub(dt.Date.Round(time.Hour)) / 24).Hours() - 1) // get day index
				politicalActivity[i].Mentions[dayindex]++                                                // increase the count for that term on that day
			}
		}
	}

	return RankPA(politicalActivity)
}

func keywords(terms []string) []DatedTerm {
	var id []byte
	var dateID []string
	var queryDate time.Time
	var DatedTerms []DatedTerm
	var tmpDT DatedTerm

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

	for _, term := range terms {
		iter := session.Query(`SELECT id FROM keyword WHERE name = ?`, term).Iter()
		for iter.Scan(&id) {
			var date time.Time
			date = DateAndId(id, dateID)
			if date.Year() > 1 {
				tmpDT.Term = term
				tmpDT.ID = id
				tmpDT.Date = date
				DatedTerms = append(DatedTerms, tmpDT)
			}
		}
	}

	for _, term := range terms {
		iter := session.Query(`SELECT id FROM entity WHERE name = ?`, term).Iter()
		for iter.Scan(&id) {
			var date time.Time
			date = DateAndId(id, dateID)
			if date.Year() > 1 {
				tmpDT.Term = term
				tmpDT.ID = id
				tmpDT.Date = date
				DatedTerms = append(DatedTerms, tmpDT)
			}
		}
	}

	return DatedTerms
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

func GetPoliticalActivityHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	var result []PoliticalActivity

	if params["type"] == "d" {
		result = DepartmentsPoliticalActivity()
	} else if params["type"] == "e" {
		result = EventsPoliticalActivity()
	} else if params["type"] == "r" {
		result = RegionsPoliticalActivity()
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
