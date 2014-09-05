package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/jinzhu/gorm"
	"os"
	"strconv"
)

type OV struct {
	Label  string
	DayAmt [30]int
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

	for _, d := range dept {

	}

	return ov
}

func EventsOverview() []OV {
	var event []Events
	err := DB.Find(&event).Error

	if err != nil && err != gorm.RecordNotFound {
		return nil
	}

	ov := make([]OV, len(event))

	for _, e := range event {

	}

	return ov
}

func RegionsOverview() []OV {
	var region []Regions
	err := DB.Find(&region).Error
	if err != nil && err != gorm.RecordNotFound {
		return nil
	}

	ov := make([]OV, len(region))

	for _, r := range region {

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
