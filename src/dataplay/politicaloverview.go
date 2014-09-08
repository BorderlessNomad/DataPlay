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
	"time"
)

const numdays = 30

var Today = time.Date(2010, 3, 1, 0, 0, 0, 0, time.UTC)
var FromDate = Today.AddDate(0, 0, -numdays)

type OV struct {
	Label  string       `json:"label"`
	DayAmt [numdays]int `json:"dayamounts"`
}

type DateID struct {
	ID   []byte
	Date time.Time
}

func DepartmentsOverview() []OV {
	session, _ := GetCassandraConnection("dp") // create connection to cassandra
	defer session.Close()
	var dept []Departments // get all departments from postgres sql table
	err := DB.Find(&dept).Error
	if err != nil && err != gorm.RecordNotFound {
		return nil
	}

	ov := make([]OV, len(dept)) // create overview output array
	var id []byte               //individual md5 hashed url id within cassandra
	var date time.Time
	var urlId []DateID

	// add all dated urlId between -n days and today to array
	iter := session.Query(`SELECT id, date FROM response WHERE date >= ? AND date < ? ALLOW FILTERING`, FromDate, Today).Iter()
	for iter.Scan(&id, &date) {
		var tmp DateID
		tmp.ID = id
		tmp.Date = date
		urlId = append(urlId, tmp)
	}

	// if err := iter.Close(); err != nil {
	// 	check(err)
	// }

	for _, ui := range urlId { // for each dated url
		for index, d := range dept { // check if any keywords or entities from the articles match the department and if so increment the counter for that date
			ov[index].Label = d.GovDept // add the dpeartment name as the label
			s, c := 0, 0
			err := session.Query(`SELECT score FROM keyword WHERE id = ? AND name =?`, ui.ID, d.GovDept).Scan(&s)
			if err != nil {
				s = 0
			}
			err = session.Query(`SELECT count FROM entity WHERE id = ? AND name =?`, ui.ID, d.GovDept).Scan(&c)
			if err != nil {
				c = 0
			}
			dayindex := int((Today.Round(time.Hour).Sub(ui.Date.Round(time.Hour)) / 24).Hours() - 1) // get day index
			if s > 0 {
				ov[index].DayAmt[dayindex]++
			}
			if c > 0 {
				ov[index].DayAmt[dayindex]++
			}
		}
	}

	return ov
}

func EventsOverview() []OV {
	session, _ := GetCassandraConnection("dp")
	defer session.Close()

	var event []Events
	err := DB.Find(&event).Error

	if err != nil && err != gorm.RecordNotFound {
		return nil
	}

	ov := make([]OV, len(event))
	var id []byte
	var date time.Time
	var urlId []DateID
	iter := session.Query(`SELECT id, date FROM response WHERE date >= ? AND date < ? ALLOW FILTERING`, FromDate, Today).Iter()
	for iter.Scan(&id, &date) {
		var tmp DateID
		tmp.ID = id
		tmp.Date = date
		urlId = append(urlId, tmp)
	}

	// if err := iter.Close(); err != nil {
	// 	check(err)
	// }

	for _, ui := range urlId {
		for index, e := range event {
			ov[index].Label = e.Event
			s, c := 0, 0
			err := session.Query(`SELECT score FROM keyword WHERE id = ? AND name = ?`, ui.ID, e.Keyword).Scan(&s)
			if err != nil {
				s = 0
			}
			err = session.Query(`SELECT count FROM entity WHERE id = ? AND name = ?`, ui.ID, e.Keyword).Scan(&c)
			if err != nil {
				c = 0
			}
			dayindex := int((Today.Round(time.Hour).Sub(ui.Date.Round(time.Hour)) / 24).Hours() - 1)
			if s > 0 {
				ov[index].DayAmt[dayindex]++
			}
			if c > 0 {
				ov[index].DayAmt[dayindex]++
			}
		}
	}
	return ov
}

func RegionsOverview() []OV {
	session, _ := GetCassandraConnection("dp")
	defer session.Close()

	var region []Regions
	err := DB.Find(&region).Error

	if err != nil && err != gorm.RecordNotFound {
		return nil
	}

	ov := make([]OV, len(region))
	var id []byte
	var date time.Time
	var urlId []DateID

	iter := session.Query(`SELECT id, date FROM response WHERE date >= ? AND date < ? ALLOW FILTERING`, FromDate, Today).Iter()
	for iter.Scan(&id, &date) {
		var tmp DateID
		tmp.ID = id
		tmp.Date = date
		urlId = append(urlId, tmp)
	}

	// if err := iter.Close(); err != nil {
	// 	check(err)
	// }

	for _, ui := range urlId {
		for index, r := range region {
			ov[index].Label = r.County
			s, c := 0, 0
			err := session.Query(`SELECT score FROM keyword WHERE id = ? AND name =?`, ui.ID, r.County).Scan(&s)
			if err != nil {
				s = 0
			}
			err = session.Query(`SELECT count FROM entity WHERE id = ? AND name =?`, ui.ID, r.County).Scan(&c)
			if err != nil {
				c = 0
			}
			dayindex := int((Today.Round(time.Hour).Sub(ui.Date.Round(time.Hour)) / 24).Hours() - 1)
			if s > 0 {
				ov[index].DayAmt[dayindex]++
			}
			if c > 0 {
				ov[index].DayAmt[dayindex]++
			}
		}
	}

	return ov
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

func GetOverviewHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	var result []OV

	if params["type"] == "d" {
		result = DepartmentsOverview()
	} else if params["type"] == "e" {
		result = EventsOverview()
	} else if params["type"] == "r" {
		result = RegionsOverview()
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
