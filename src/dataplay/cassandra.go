package main

import (
	"github.com/gocql/gocql"
	"os"
	"strconv"
	"time"
)

func GetCassandraConnection(keyspace string) (*gocql.Session, error) {
	cassandraHost := "109.231.122.107"
	cassandraPort := 9042

	if os.Getenv("DP_CASSANDRA_HOST") != "" {
		cassandraHost = os.Getenv("DP_CASSANDRA_HOST")
	}

	if os.Getenv("DP_CASSANDRA_PORT") != "" {
		cassandraPort, _ = strconv.Atoi(os.Getenv("DP_CASSANDRA_PORT"))
	}

	cluster := gocql.NewCluster(cassandraHost)
	cluster.Timeout = 1 * time.Minute
	cluster.Port = cassandraPort
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Compressor = gocql.SnappyCompressor{}
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 5}
	session, err := cluster.CreateSession()

	if err != nil {
		Logger.Println("Could not connect to the Cassandara server.")
		return nil, err
	}

	return session, nil
}
