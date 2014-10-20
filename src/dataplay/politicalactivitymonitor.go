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

var Today = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC) // override today's date
var FromDate = Today.AddDate(0, 0, -numdays)

type TermKey struct {
	KeyTerm  string
	MainTerm string
}

type PoliticalActivity struct {
	Term     string               `json:"term"`
	Mentions [numdays]PoliticalXY `json:"graph"`
	Val      int                  `json:"-"`
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
func DepartmentsPoliticalActivity() []PoliticalActivity {
	var dept []Departments // get all departments from postgres sql table
	var terms []TermKey
	err := DB.Find(&dept).Error

	if err != nil && err != gorm.RecordNotFound {
		return nil
	}

	var tmp TermKey

	for _, d := range dept {
		tmp.KeyTerm = d.Key
		tmp.MainTerm = d.Dept
		terms = append(terms, tmp)
	}

	return TermFrequency(terms)
}

// gets names of all events, checks for mentions in specified time period and returns ranked array of 15 most popular terms and their 30 day frequencies
func EventsPoliticalActivity() []PoliticalActivity {
	var event []Events
	var terms []TermKey
	err := DB.Find(&event).Error

	if err != nil && err != gorm.RecordNotFound {
		return nil
	}

	var tmp TermKey
	for _, e := range event {
		tmp.KeyTerm = e.Key
		tmp.MainTerm = e.Event
		terms = append(terms, tmp)
	}

	return TermFrequency(terms)
}

// gets names of all regions, checks for mentions in specified time period and returns ranked array of 15 most popular terms and their 30 day frequencies
func RegionsPoliticalActivity() []PoliticalActivity {
	var region []Regions
	var terms []TermKey
	err := DB.Find(&region).Error

	if err != nil && err != gorm.RecordNotFound {
		return nil
	}

	var tmp TermKey
	for _, r := range region {
		tmp.KeyTerm = r.Key
		tmp.MainTerm = r.Region
		terms = append(terms, tmp)
	}

	return TermFrequency(terms)
}

func TermFrequency(terms []TermKey) []PoliticalActivity {
	var date time.Time
	var name string
	politicalActivity := make([]PoliticalActivity, 0)

	session, _ := GetCassandraConnection("dp") // create connection to cassandra
	defer session.Close()

	iter1 := session.Query(`SELECT date, name FROM keyword WHERE date >= ? AND date < ? ALLOW FILTERING`, FromDate, Today).Iter()
	for iter1.Scan(&date, &name) {
		for _, term := range terms {
			if name == term.KeyTerm { // for any key term matches
				i := PaPlace(&politicalActivity, term.MainTerm)                                       // either get place of main term or add to array if doesn't exist
				dayindex := int((Today.Round(time.Hour).Sub(date.Round(time.Hour)) / 24).Hours() - 1) // get day index
				politicalActivity[i].Mentions[dayindex].Y++
			}
		}
	}

	iter2 := session.Query(`SELECT date, name FROM entity WHERE date >= ? AND date < ? ALLOW FILTERING`, FromDate, Today).Iter()
	for iter2.Scan(&date, &name) {
		for _, term := range terms {
			if name == term.KeyTerm { // for any key term matches
				i := PaPlace(&politicalActivity, term.MainTerm)                                       // either get place of main term or add to array if doesn't exist
				dayindex := int((Today.Round(time.Hour).Sub(date.Round(time.Hour)) / 24).Hours() - 1) // get day index
				politicalActivity[i].Mentions[dayindex].Y++
			}
		}
	}
	return RankPA(politicalActivity)
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

func PopularPoliticalActivity() [3]Popular {
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
		return popular
	}

	err = DB.Select("priv_users.username, count(priv_discovered.uid) as counter").Joins("LEFT JOIN priv_users ON priv_discovered.uid = priv_users.uid").Group("priv_users.username, priv_discovered.uid").Order("counter DESC").Limit(5).Find(&results).Error
	if err != nil && err != gorm.RecordNotFound {
		return popular
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

	return popular
}

func GetCassandraConnection(keyspace string) (*gocql.Session, error) {
	cassandraHost := "109.231.121.129"
	cassandraPort := 9042

	if os.Getenv("DP_CASSANDRA_HOST") != "" {
		cassandraHost = os.Getenv("DP_CASSANDRA_HOST")
	}

	if os.Getenv("DP_CASSANDRA_PORT") != "" {
		cassandraPort, _ = strconv.Atoi(os.Getenv("DP_CASSANDRA_PORT"))
	}

	cluster := gocql.NewCluster(cassandraHost)
	cluster.Port = cassandraPort
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
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

	if params["type"] == "d" {
		result = DepartmentsPoliticalActivity()
	} else if params["type"] == "e" {
		result = EventsPoliticalActivity()
	} else if params["type"] == "r" {
		result = RegionsPoliticalActivity()
	} else if params["type"] == "p" {
		pResult := PopularPoliticalActivity()
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

	r, err := json.Marshal(result)
	if err != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func GetPoliticalActivityQ(params map[string]string) string {
	var result []PoliticalActivity

	if params["type"] == "d" {
		result = DepartmentsPoliticalActivity()
	} else if params["type"] == "e" {
		result = EventsPoliticalActivity()
	} else if params["type"] == "r" {
		result = RegionsPoliticalActivity()
	} else if params["type"] == "p" {
		pResult := PopularPoliticalActivity()
		r, e := json.Marshal(pResult)
		if e != nil {
			return e.Error()
		}
		return string(r)
	} else {
		return "Bad type param"
	}

	r, e := json.Marshal(result)
	if e != nil {
		return e.Error()
	}

	return string(r)
}
