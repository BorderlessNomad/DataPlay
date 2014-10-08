package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/pmylund/sortutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type NewsArticle struct {
	Date     time.Time `json:"date"`
	Title    string    `json:"title"`
	Url      string    `json:"url"`
	ImageUrl string    `json:"image_url"`
	Score    int       `json:"SCORE"`
}

func SearchForNewsHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	result, error := SearchForNews(params["terms"])
	if error != nil {
		http.Error(res, error.Message, error.Code)
		return ""
	}

	r, err1 := json.Marshal(result)
	if err1 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func SearchForNewsQ(params map[string]string) string {
	if params["user"] == "" {
		return ""
	}

	result, err := SearchForNews(params["terms"])
	if err != nil {
		return ""
	}

	r, e := json.Marshal(result)
	if e != nil {
		return ""
	}

	return string(r)
}

// search Cassandra for news articles relating to the search terms
// also searches sql tables to find relevant dates to tie in with table search
// also checks if dates are entered in the search
func SearchForNews(searchstring string) ([]NewsArticle, *appError) {

	session, _ := GetCassandraConnection("dp") // create connection to cassandra
	defer session.Close()

	newsArticles := []NewsArticle{}
	searchTerms := strings.Split(searchstring, "_")
	earliestDate := Earliest(searchTerms)

	iter := session.Query(`SELECT id, title, original_url, date, description FROM response WHERE date >= ? AND date < ? ALLOW FILTERING`, earliestDate, Today).Iter()

	var id []byte
	var date time.Time
	var originalUrl, title, description string

	for iter.Scan(&id, &title, &originalUrl, &date, &description) {
		termcount := 0

		for _, term := range searchTerms {
			termcount += TermCheck(term, description)
			termcount += TermCheck(term, title)        //if term in title too it adds weight
			termcount += DateCheck(earliestDate, date) // add weight if the article is from around the right month or year
		}

		if termcount > 0 {
			var tmpNA NewsArticle
			imageUrl := ""
			picIndex := 1000000
			iter2 := session.Query(`SELECT url, pic_index FROM image WHERE id = ? ALLOW FILTERING`, id).Iter()
			for iter2.Scan(&imageUrl, &picIndex) {
				if picIndex == 0 {
					tmpNA.ImageUrl = imageUrl
				}
			}
			tmpNA.Date = date
			tmpNA.Title = title
			tmpNA.Url = originalUrl
			tmpNA.Score = termcount
			newsArticles = append(newsArticles, tmpNA)
		}
	}

	sortutil.DescByField(newsArticles, "Score")
	return newsArticles, nil
}

// return 1 if the term is found in the passage
func TermCheck(term string, passage string) int {
	descriptions := strings.Split(passage, " ")
	for _, d := range descriptions {
		if d == term {
			return 1
		}
	}
	return 0
}

//return 2, 1 or 0 depending on how accurate the date is
func DateCheck(date1 time.Time, date2 time.Time) int {
	y1 := date1.Year()
	y2 := date2.Year()
	m1 := date1.Month()
	m2 := date2.Month()

	if date1.Nanosecond() > 0 { // check nanosecond flag to see if year and month or just year need handling
		if y1 == y2 && m1 == m2 {
			return 2 // if month and year are correct
		} else if y1 == y2 {
			return 1 // if only correct year
		}
	} else if y1 == y2 { // when just year is set and correct
		return 1
	}

	return 0 // no match
}

// returns earliest date from related search
func Earliest(terms []string) time.Time {
	var dates []time.Time

	for _, term := range terms {

		result, _ := SearchForData(0, term, nil)

		for i, _ := range result.Results {
			d := strings.Split(result.Results[i].PrimaryDate, " ")

			if len(d) == 2 {
				year, _ := strconv.Atoi(d[1])
				month := MonthNum(d[0])
				date := time.Date(year, month, 1, 0, 0, 0, 1, time.UTC) // year and month with nanosecond flag set
				dates = append(dates, date)
			} else if len(d) == 1 {
				year, _ := strconv.Atoi(d[0])
				date := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC) //just year
				dates = append(dates, date)
			}
		}
	}

	earliest := time.Date(2100, 0, 0, 0, 0, 0, 0, time.UTC)
	for _, d := range dates {
		if d.Before(earliest) {
			earliest = d
		}
	}

	fmt.Println("NATIVES", earliest)
	return earliest
}

// convert month name to time month type
func MonthNum(month string) time.Month {
	mstr := "January_February_March_April_May_June_July_August_September_October_November_December"
	m := strings.Split(mstr, "_")
	for i, _ := range m {
		if month == m[i] {
			return time.Month(i + 1)
		}
	}
	return time.Month(1)
}
