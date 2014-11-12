package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/gocql/gocql"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
	"strconv"
	"time"
)

const numdays = 30

type TermKey struct {
	KeyTerm  string
	MainTerm string
}

type PoliticalActivity struct {
	Term     string                   `json:"term"`
	Mentions [numdays + 1]PoliticalXY `json:"graph"`
	Val      int                      `json:"-"`
}

type PoliticalXY struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Popular struct {
	Id       string     `json:"id"`
	Category string     `json:"category"`
	TA       [5]TermAmt `json:"top5"`
}

type TermAmt struct {
	Term   string `json:"term"`
	Amount int    `json:"amount"`
}

// gets names of all departments, checks for mentions in specified time period and returns ranked array of 15 most popular terms and their 30 day frequencies
func DepartmentsPoliticalActivity() ([]PoliticalActivity, error) {
	var dept []Departments // get all departments from postgres sql table
	var terms []TermKey

	err := DB.Find(&dept).Error
	if err != nil && err != gorm.RecordNotFound {
		return nil, err
	}

	var tmp TermKey
	for _, d := range dept {
		tmp.KeyTerm = d.Key
		tmp.MainTerm = d.Dept
		terms = append(terms, tmp)
	}

	termFrequency, termError := TermFrequency(terms)
	if termError != nil {
		return nil, termError
	}

	return termFrequency, nil
}

// gets names of all events, checks for mentions in specified time period and returns ranked array of 15 most popular terms and their 30 day frequencies
func EventsPoliticalActivity() ([]PoliticalActivity, error) {
	var event []Events
	var terms []TermKey

	err := DB.Find(&event).Error
	if err != nil && err != gorm.RecordNotFound {
		return nil, err
	}

	var tmp TermKey
	for _, e := range event {
		tmp.KeyTerm = e.Key
		tmp.MainTerm = e.Event
		terms = append(terms, tmp)
	}

	termFrequency, termError := TermFrequency(terms)
	if termError != nil {
		return nil, termError
	}

	return termFrequency, nil
}

// gets names of all regions, checks for mentions in specified time period and returns ranked array of 15 most popular terms and their 30 day frequencies
func RegionsPoliticalActivity() ([]PoliticalActivity, error) {
	var region []Regions
	var terms []TermKey

	err := DB.Find(&region).Error
	if err != nil && err != gorm.RecordNotFound {
		return nil, err
	}

	var tmp TermKey
	for _, r := range region {
		tmp.KeyTerm = r.Key
		tmp.MainTerm = r.Region
		terms = append(terms, tmp)
	}

	termFrequency, termError := TermFrequency(terms)
	if termError != nil {
		return nil, termError
	}

	return termFrequency, nil
}

func TermFrequency(terms []TermKey) ([]PoliticalActivity, error) {
	var date time.Time
	var name string
	politicalActivity := make([]PoliticalActivity, 0)
	now := time.Now()
	var today = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC) // override today's date
	var from = today.AddDate(0, 0, -numdays)

	session, err := GetCassandraConnection("dataplay") // create connection to cassandra
	if err != nil {
		return nil, err
	}
	defer session.Close()

	iter1 := session.Query(`SELECT date, name FROM keyword LIMIT 60000`).Iter()

	for iter1.Scan(&date, &name) {
		for _, term := range terms {
			if name == term.KeyTerm && (date.Equal(from) || date.After(from) && date.Before(today)) { // for any key term matches in date
				i := PaPlace(&politicalActivity, term.MainTerm)                                   // either get place of main term or add to array if doesn't exist
				dayindex := int((today.Round(time.Hour).Sub(date.Round(time.Hour)) / 24).Hours()) // get day index
				politicalActivity[i].Mentions[dayindex].Y++
			}
		}
	}
	err1 := iter1.Close()
	if err1 != nil {
		return nil, err1
	}

	iter2 := session.Query(`SELECT date, name FROM entity LIMIT 60000`).Iter()

	for iter2.Scan(&date, &name) {
		for _, term := range terms {
			if name == term.KeyTerm && (date.Equal(from) || date.After(from) && date.Before(today)) { // for any key term matches in date
				i := PaPlace(&politicalActivity, term.MainTerm)                                   // either get place of main term or add to array if doesn't exist
				dayindex := int((today.Round(time.Hour).Sub(date.Round(time.Hour)) / 24).Hours()) // get day index
				politicalActivity[i].Mentions[dayindex].Y++
			}
		}
	}
	err2 := iter2.Close()
	if err2 != nil {
		return nil, err2
	}

	return RankPA(politicalActivity), nil
}

func PaPlace(pa *[]PoliticalActivity, t string) int {
	for i, p := range *pa {
		if p.Term == t {
			return i
		}
	}

	var tmp PoliticalActivity
	tmp.Term = t
	*pa = append(*pa, tmp)

	return len(*pa) - 1
}

// sort PA array and returns slice of top 15
func RankPA(activities []PoliticalActivity) []PoliticalActivity {
	for i, _ := range activities {
		total := 0
		for j, _ := range activities[i].Mentions {
			total += activities[i].Mentions[j].Y
			activities[i].Mentions[j].X = j
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

func PopularPoliticalActivity() ([3]Popular, error) {
	var popular [3]Popular

	popular[0].Id = "most_popular"
	popular[0].Category = "Most Popular Keywords"

	popular[1].Id = "top_correlated"
	popular[1].Category = "Top Correlated Keywords"

	popular[2].Id = "top_discoverers"
	popular[2].Category = "Top Discoverers"

	results := []struct {
		Discovered
		Username string
		Counter  int
	}{}

	// SELECT term, count from priv_searchterms order by count DESC limit 5
	searchterm := []SearchTerm{}
	err := DB.Select("term, count").Order("count desc").Limit(5).Find(&searchterm).Error
	if err != nil && err != gorm.RecordNotFound {
		return popular, err
	}

	err = DB.Select("priv_users.username, count(priv_discovered.uid) as counter").Joins("LEFT JOIN priv_users ON priv_discovered.uid = priv_users.uid").Group("priv_users.username, priv_discovered.uid").Order("counter DESC").Limit(5).Find(&results).Error
	if err != nil && err != gorm.RecordNotFound {
		return popular, err
	}

	n := 5
	if len(searchterm) < 5 {
		n = len(searchterm)
	}

	for i := 0; i < n; i++ {
		popular[0].TA[i].Term = searchterm[i].Term
		popular[0].TA[i].Amount = searchterm[i].Count

		popular[1].TA[i].Term = searchterm[i].Term
		popular[1].TA[i].Amount = searchterm[i].Count
	}

	n = 5
	if len(results) < 5 {
		n = len(results)
	}

	for i := 0; i < n; i++ {
		popular[2].TA[i].Term = results[i].Username
		popular[2].TA[i].Amount = results[i].Counter
	}

	return popular, nil
}

func GetCassandraConnection(keyspace string) (*gocql.Session, error) {
	// cassandraHost := "109.231.121.96"
	// cassandraPort := 9042
	cassandraHost := "10.0.0.2"
	cassandraPort := 49211

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

/////methods used by APIs////////////////////
func GetPoliticalActivityHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	var result []PoliticalActivity
	var err error

	if params["type"] == "d" {
		result, err = DepartmentsPoliticalActivity()
	} else if params["type"] == "e" {
		result, err = EventsPoliticalActivity()
	} else if params["type"] == "r" {
		result, err = RegionsPoliticalActivity()
	} else if params["type"] == "p" {
		pResult, errP := PopularPoliticalActivity()
		if errP != nil {
			http.Error(res, errP.Error(), http.StatusInternalServerError)
			return ""
		}

		r, err := json.Marshal(pResult)
		if err != nil {
			http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
			return ""
		}

		return string(r)
	} else {
		http.Error(res, "Bad type param", http.StatusInternalServerError)
		return ""
	}

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return ""
	}

	r, err := json.Marshal(result)
	if err != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}
