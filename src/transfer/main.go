package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"os"
	"strconv"
	"time"
)

func main() {
	DataTransfer()
}

func DataTransfer() {
	f, _ := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()

	session1, err1 := GetCassandraConnection("dp")
	session2, err2 := GetCassandraConnection("dataplay")

	if err1 != nil {
		fmt.Println("ERROR CASSANDRA 1!")
	}
	defer session1.Close()

	if err2 != nil {
		fmt.Println("ERROR CASSANDRA 2 !")
	}
	defer session2.Close()

	var ToDate = time.Date(2014, 28, 1, 0, 0, 0, 0, time.UTC)
	var FromDate = ToDate.AddDate(0, -1, 0) //+1 to LOOP's i value

	fmt.Println("START")

	for i := 0; i < 1; i++ {
		fmt.Println("RESPONSE LOOP ", i, " start ", time.Now())
		var id []byte
		var pic_index int
		var date time.Time
		var description, url, title, name, pic_url, related_url string

		iter1 := session1.Query(`SELECT id, date, description, url, title FROM response WHERE date >= ? AND date < ? ALLOW FILTERING`, FromDate, ToDate).Iter()

		for iter1.Scan(&id, &date, &description, &url, &title) {
			session2.Query(`INSERT INTO response (date, dummy, description, url, title) VALUES (?, ?, ?, ?, ?)`, date, 1, description, url, title).Exec()

			iter2 := session1.Query(`SELECT url, pic_index FROM image WHERE id = ? ALLOW FILTERING`, id).Iter()

			for iter2.Scan(&pic_url, &pic_index) {
				if pic_index == 0 {
					session2.Query(`INSERT INTO image (date, dummy, pic_url, url) VALUES (?, ?, ?, ?)`, date, 1, pic_url, url).Exec()
				}
			}

			err := iter2.Close()
			if err != nil {
				fmt.Println("ERROR LOOP image", err.Error())
				f.WriteString("ERROR LOOP image " + strconv.Itoa(i) + " " + err.Error() + "     ")
			}

			iter3 := session1.Query(`SELECT description, title, url FROM related WHERE id = ? ALLOW FILTERING`, id).Iter()

			for iter3.Scan(&description, &title, &related_url) {
				session2.Query(`INSERT INTO related (date, dummy, description, title, related_url, url) VALUES (?, ?, ?, ?, ?, ?)`, date, 1, description, title, related_url, url).Exec()
			}

			err = iter3.Close()
			if err != nil {
				fmt.Println("ERROR LOOP related", err.Error())
				f.WriteString("ERROR LOOP related " + strconv.Itoa(i) + " " + err.Error() + "     ")
			}

		}

		err := iter1.Close()
		if err != nil {
			fmt.Println("ERROR LOOP image, related, response", err.Error())
			f.WriteString("ERROR LOOP image, related, response " + strconv.Itoa(i) + " " + err.Error() + "     ")
		}

		fmt.Println("SUCCESS LOOP image, related, response ", i, " at ", time.Now())

		iter4 := session1.Query(`SELECT date, name, url FROM keyword WHERE date >= ? AND date < ? LIMIT 65000 ALLOW FILTERING`, FromDate, ToDate).Iter()

		for iter4.Scan(&date, &name, &url) {
			session2.Query(`INSERT INTO keyword (date, dummy, name, url) VALUES (?, ?, ?, ?)`, date, 1, name, url).Exec()
		}

		err = iter4.Close()

		if err != nil {
			fmt.Println("ERROR LOOP keyword", err.Error())
			f.WriteString("ERROR LOOP keyword " + strconv.Itoa(i) + " " + err.Error() + "     ")
		}

		fmt.Println("SUCCESS LOOP keyword", i, " at ", time.Now())

		iter5 := session1.Query(`SELECT date, name, url FROM entity WHERE date >= ? AND date < ? LIMIT 65000 ALLOW FILTERING`, FromDate, ToDate).Iter()

		for iter5.Scan(&date, &name, &url) {
			session2.Query(`INSERT INTO entity (date, dummy, name, url) VALUES (?, ?, ?, ?)`, date, 1, name, url).Exec()
		}

		err = iter5.Close()
		if err != nil {
			fmt.Println("ERROR LOOP entity", err.Error())
			f.WriteString("ERROR LOOP entity " + strconv.Itoa(i) + " " + err.Error() + "     ")
		}

		fmt.Println("SUCCESS LOOP entity", i, time.Now())

		ToDate = ToDate.AddDate(0, -1, 0)
		FromDate = FromDate.AddDate(0, -1, 0)
		fmt.Println("TOTAL LOOP ", i, " COMPLETE ", 58-i, " MORE TO GO ", time.Now())
	}
}

func GetCassandraConnection(keyspace string) (*gocql.Session, error) {
	cassandraHost := "109.231.121.96"
	cassandraPort := 9042

	if os.Getenv("DP_CASSANDRA_HOST") != "" {
		cassandraHost = os.Getenv("DP_CASSANDRA_HOST")
	}

	if os.Getenv("DP_CASSANDRA_PORT") != "" {
		cassandraPort, _ = strconv.Atoi(os.Getenv("DP_CASSANDRA_PORT"))
	}

	cluster := gocql.NewCluster(cassandraHost)
	cluster.Timeout = 2 * time.Minute
	cluster.Port = cassandraPort
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Compressor = gocql.SnappyCompressor{}
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 5}
	session, err := cluster.CreateSession()

	if err != nil {
		fmt.Println("Could not connect to the Cassandara server.")
		return nil, err
	}

	return session, nil
}
