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
	cassandraTimeout := 1 * time.Minute
	cassandraMaxRetries := 5

	if os.Getenv("DP_CASSANDRA_HOST") != "" {
		cassandraHost = os.Getenv("DP_CASSANDRA_HOST")
	}

	if os.Getenv("DP_CASSANDRA_PORT") != "" {
		cassandraPort, _ = strconv.Atoi(os.Getenv("DP_CASSANDRA_PORT"))
	}

	if os.Getenv("DP_CASSANDRA_TIMEOUT") != "" {
		timeoutDuration, _ := strconv.Atoi(os.Getenv("DP_CASSANDRA_TIMEOUT"))
		cassandraTimeout = time.Duration(timeoutDuration) * time.Minute
	}

	if os.Getenv("DP_CASSANDRA_MAX_RETRIES") != "" {
		cassandraMaxRetries, _ = strconv.Atoi(os.Getenv("DP_CASSANDRA_MAX_RETRIES"))
	}

	cluster := gocql.NewCluster(cassandraHost)
	cluster.Port = cassandraPort
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Compressor = gocql.SnappyCompressor{}
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: cassandraMaxRetries}
	cluster.Timeout = cassandraTimeout

	Logger.Println("Connecting to Cassandara " + cassandraHost + ":" + strconv.Itoa(cassandraPort))
	session, err := cluster.CreateSession()

	if err != nil {
		Logger.Println("Could not connect to the Cassandara server.")
		return nil, err
	}

	Logger.Println("Connected!")

	return session, nil
}
