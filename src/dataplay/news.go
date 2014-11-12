package main

import (
	"encoding/json"
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
	Score    float64   `json:"score"`
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

// search Cassandra for news articles relating to the search terms
// also searches sql tables to find relevant dates to tie in with table search
// also checks if dates are entered in the search
func SearchForNews(searchstring string) ([]NewsArticle, *appError) {
	// now := time.Now()
	// var Today = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC) // override today's date

	session, _ := GetCassandraConnection("dataplay") // create connection to cassandra
	defer session.Close()

	newsArticles := []NewsArticle{}
	searchTerms := strings.Split(strings.Replace(searchstring, "_", " ", -1), "!")
	earliestDate := Earliest(searchTerms) // links with SQL database

	var date time.Time
	var url, imageUrl, title, description string

	// iter1 := session.Query(`SELECT title, url, date, description FROM response WHERE date >= ? AND date <= ? ALLOW FILTERING`, earliestDate, earliestDate.AddDate(0, 1, 0)).Iter()
	// for iter1.Scan(&title, &url, &date, &description) {
	// 	termcount := 0.0

	// 	for i, term := range searchTerms {
	// 		termcount += float64(TermCheck(term, description+" "+title)) * 1 / float64(i+1)
	// 		termcount += float64(DateCheck(earliestDate, date)) * 1 / float64(i+1) // add weight if the article is from around the right month or year
	// 	}

	// 	if termcount > 0 {
	// 		var tmpNA NewsArticle
	// 		tmpNA.Date = date
	// 		tmpNA.Title = title
	// 		tmpNA.Url = url
	// 		tmpNA.Score = termcount
	// 		newsArticles = append(newsArticles, tmpNA)
	// 	}
	// }

	iter2 := session.Query(`SELECT title, url, date, description FROM response LIMIT 6000`).Iter()
	for iter2.Scan(&title, &url, &date, &description) {
		termcount := 0.0

		for i, term := range searchTerms {
			termcount += float64(TermCheck(term, title)) * 1 / float64(i+1)
			termcount += float64(TermCheck(term, description)) * 1 / float64(i+1)
			termcount += float64(DateCheck(earliestDate, date)) * 1 / float64(i+1) // add weight if the article is from around the right month or year
		}

		if termcount > 0 {
			var tmpNA NewsArticle
			tmpNA.Date = date
			tmpNA.Title = title
			tmpNA.Url = url
			tmpNA.Score = termcount
			newsArticles = append(newsArticles, tmpNA)
		}
	}

	sortutil.DescByField(newsArticles, "Score")

	n := len(newsArticles)
	if len(newsArticles) > 8 {
		n = 8
	}

	newsSlice := newsArticles[0:n]

	for i, _ := range newsSlice {
		iter3 := session.Query(`SELECT pic_url, url FROM image WHERE date = ? ALLOW FILTERING`, newsSlice[i].Date).Iter()
		for iter3.Scan(&imageUrl, &url) {
			if url == newsSlice[i].Url {
				newsSlice[i].ImageUrl = imageUrl
			}
		}
	}
	return newsSlice, nil
}

// return 1 if the term is found in the passage
func TermCheck(term string, passage string) int {
	check := strings.Contains(strings.ToLower(passage), strings.ToLower(term))
	if check {
		return 1
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
			} else if d[0] != "" {
				year, _ := strconv.Atoi(d[0])
				date := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC) //just year
				dates = append(dates, date)
			}
		}
	}

	earliest := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, d := range dates {
		if d.Before(earliest) {
			earliest = d
		}
	}

	origin := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	if earliest.Before(origin) {
		earliest = origin
	}

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
