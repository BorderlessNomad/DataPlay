package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"os"
	"strconv"
)

func GetCassandraConnection(keyspace string) (*gocql.Session, error) {
	cassandraHost := "10.0.0.2"
	cassandraPort := 49211
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
