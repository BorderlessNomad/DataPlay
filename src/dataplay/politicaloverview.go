package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/jinzhu/gorm"
	"os"
	"strconv"
	"time"
)

const numdays = 30

const layout = "2006-01-02"

type OV struct {
	Label  string
	DayAmt [numdays]int
}

type DateID struct {
	ID   []byte
	Date time.Time
}

func DepartmentsOverview() []OV {
	session, _ := GetCassandraConnection("dp")
	defer session.Close()

	var dept []Departments
	err := DB.Find(&dept).Error

	if err != nil && err != gorm.RecordNotFound {
		return nil
	}

	ov := make([]OV, len(dept))
	timeFrom := time.Now().AddDate(0, 0, -numdays) // older than n days
	var id []byte
	var date time.Time
	var ids []DateID
	iter := session.Query(`SELECT id, date FROM response WHERE date > ?`, timeFrom).Iter()
	for iter.Scan(&id, &date) {
		var tmp DateID
		tmp.ID = id
		tmp.Date = date.Format(layout)
		ids = append(ids, tmp)
	}

	for _, i := range ids {
		for index, d := range dept {
			ov[index].Label = d
			s, c := 0, 0
			err := session.Query(`SELECT score FROM keyword WHERE id = ? AND name =?`, i.ID, d).Scan(&s)
			check(err)
			err = session.Query(`SELECT count FROM entity WHERE id = ? AND name =?`, i.ID, d).Scan(&c)
			check(err)
			now := time.Now()
			dayindex := (now.Round(time.Hour).Sub(i.Date.Round(time.Hour)) / 24) - 1
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
	timeFrom := time.Now().AddDate(0, 0, -numdays) // older than n days
	var id []byte
	var date time.Time
	var ids []DateID
	iter := session.Query(`SELECT id, date FROM response WHERE date > ?`, timeFrom).Iter()
	for iter.Scan(&id, &date) {
		var tmp DateID
		tmp.ID = id
		tmp.Date = date.Format(layout)
		ids = append(ids, tmp)
	}

	for _, e := range event {
		for index, d := range dept {
			ov[index].Label = d
			s, c := 0, 0
			err := session.Query(`SELECT score FROM keyword WHERE id = ? AND name =?`, i.ID, d).Scan(&s)
			check(err)
			err = session.Query(`SELECT count FROM entity WHERE id = ? AND name =?`, i.ID, d).Scan(&c)
			check(err)
			now := time.Now()
			dayindex := (now.Round(time.Hour).Sub(i.Date.Round(time.Hour)) / 24) - 1
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
	timeFrom := time.Now().AddDate(0, 0, -numdays) // older than n days
	var id []byte
	var date time.Time
	var ids []DateID
	iter := session.Query(`SELECT id, date FROM response WHERE date > ?`, timeFrom).Iter()
	for iter.Scan(&id, &date) {
		var tmp DateID
		tmp.ID = id
		tmp.Date = date.Format(layout)
		ids = append(ids, tmp)
	}

	for _, r := range region {
		for index, d := range dept {
			ov[index].Label = d
			s, c := 0, 0
			err := session.Query(`SELECT score FROM keyword WHERE id = ? AND name =?`, i.ID, d).Scan(&s)
			check(err)
			err = session.Query(`SELECT count FROM entity WHERE id = ? AND name =?`, i.ID, d).Scan(&c)
			check(err)
			now := time.Now()
			dayindex := (now.Round(time.Hour).Sub(i.Date.Round(time.Hour)) / 24) - 1
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
