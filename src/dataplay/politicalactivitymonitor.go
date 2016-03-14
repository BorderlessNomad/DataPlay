package main

import (
	"encoding/json"
	"github.com/ahirmayur/gorm"
	"github.com/codegangsta/martini"
	"net/http"
	"time"
)

const numdays = 365 // @todo: change this to '30' after importing new data in Cassandra - Mayur

type TermKey struct {
	KeyTerm  string
	MainTerm string
	DataFreq int `json:"omitempty"`
}

type PoliticalActivity struct {
	Term     string                   `json:"term"`
	Mentions [numdays + 1]PoliticalXY `json:"graph"`
	Val      int                      `json:"val"`
}

type PoliticalXY struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Popular struct {
	Id       string      `json:"id"`
	Category string      `json:"category"`
	TA       [10]TermAmt `json:"top"`
}

type TermAmt struct {
	Term   string `json:"term"`
	Amount int    `json:"amount"`
}

func KeywordsPoliticalActivity() ([]PoliticalActivity, error) {
	var dictionary []Dictionary
	var terms []TermKey

	err := DB.Find(&dictionary).Error
	if err != nil && err != gorm.RecordNotFound {
		return nil, err
	}

	var tmp TermKey
	for _, d := range dictionary {
		tmp.KeyTerm = d.Term
		tmp.MainTerm = d.Term
		tmp.DataFreq = d.Frequency
		terms = append(terms, tmp)
	}

	termFrequency, termError := TermFrequency(terms)
	if termError != nil {
		return nil, termError
	}

	termFrequency = CalcPAValsDatasets(termFrequency, dictionary)

	return RankPA(termFrequency), nil
}

func MediaPulsePoliticalActivity() ([]PoliticalActivity, error) {
	var dictionary []Dictionary
	var terms []TermKey

	err := DB.Find(&dictionary).Error
	if err != nil && err != gorm.RecordNotFound {
		return nil, err
	}

	var tmp TermKey
	for _, d := range dictionary {
		tmp.KeyTerm = d.Term
		tmp.MainTerm = d.Term
		tmp.DataFreq = d.Frequency
		terms = append(terms, tmp)
	}

	termFrequency, termError := TermFrequency(terms)
	if termError != nil {
		return nil, termError
	}

	termFrequency = CalcPAValsMedia(termFrequency)
	return RankPA(termFrequency), nil
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

	termFrequency = CalcPAValsMedia(termFrequency)
	return RankPA(termFrequency), nil
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

	iter1 := session.Query(`SELECT date, name FROM keyword WHERE date > ? LIMIT ? ALLOW FILTERING`, from, 40000).Iter()

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

	iter2 := session.Query(`SELECT date, name FROM entity WHERE date > ? LIMIT ? ALLOW FILTERING`, from, 40000).Iter()

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

	return politicalActivity, nil
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

// attach Vals based upon our datasets (frequency)
func CalcPAValsDatasets(activities []PoliticalActivity, dictionary []Dictionary) []PoliticalActivity {
	for i, _ := range activities {

		for j, _ := range activities[i].Mentions {
			activities[i].Mentions[j].X = j
		}

		freq := 0
		for _, d := range dictionary {
			if d.Term == activities[i].Term {
				freq = d.Frequency
			}
		}
		activities[i].Val = freq
	}

	return activities
}

// attach Vals based upon the media data (bbc mentions)
func CalcPAValsMedia(activities []PoliticalActivity) []PoliticalActivity {
	for i, _ := range activities {
		total := 0
		for j, _ := range activities[i].Mentions {
			total += activities[i].Mentions[j].Y
			activities[i].Mentions[j].X = j
		}
		activities[i].Val = total
	}

	return activities
}

// sort PA array and returns slice of top 15
func RankPA(activities []PoliticalActivity) []PoliticalActivity {
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

	if len(activities) < 15 {
		return activities
	}

	return activities[0:15]
}

func PopularPoliticalActivity() ([2]Popular, error) {
	var popular [2]Popular

	popular[0].Id = "most_popular"
	popular[0].Category = "Most Popular Keywords"

	popular[1].Id = "top_discoverers"
	popular[1].Category = "Top Discoverers"

	results := []struct {
		Discovered
		Username string
		Counter  int
	}{}

	// SELECT term, count from priv_searchterms order by count DESC limit 10
	searchterm := []SearchTerm{}
	err := DB.Select("term, count").Order("count desc").Limit(10).Find(&searchterm).Error
	if err != nil && err != gorm.RecordNotFound {
		return popular, err
	}

	err = DB.Select("priv_users.username, count(priv_discovered.uid) as counter").Joins("LEFT JOIN priv_users ON priv_discovered.uid = priv_users.uid").Group("priv_users.username, priv_discovered.uid").Order("counter DESC").Limit(10).Find(&results).Error
	if err != nil && err != gorm.RecordNotFound {
		return popular, err
	}

	n := 10
	if len(searchterm) < 10 {
		n = len(searchterm)
	}

	for i := 0; i < n; i++ {
		popular[0].TA[i].Term = searchterm[i].Term
		popular[0].TA[i].Amount = searchterm[i].Count
	}

	n = 10
	if len(results) < 10 {
		n = len(results)
	}

	for i := 0; i < n; i++ {
		popular[1].TA[i].Term = results[i].Username
		popular[1].TA[i].Amount = results[i].Counter
	}

	return popular, nil
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

	if params["type"] == "keywords" {
		result, err = KeywordsPoliticalActivity()
	} else if params["type"] == "mediapulse" {
		result, err = MediaPulsePoliticalActivity()
	} else if params["type"] == "regions" {
		result, err = RegionsPoliticalActivity()
	} else if params["type"] == "popular" {
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
