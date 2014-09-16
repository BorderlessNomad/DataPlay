package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/pmylund/sortutil"
	"net/http"
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

func SearchForNews(searchterms string) ([]NewsArticle, *appError) {
	newsarticles := []NewsArticle{}
	searchterm := strings.Split(searchterms, "_")
	session, _ := GetCassandraConnection("dp") // create connection to cassandra
	defer session.Close()

	/////////// @TODO: Weight based on search date///////////////
	// searchdate := ""
	// params := map[string]string{
	// 	"offset": "0",
	// 	"count":  "1",
	// }
	// searchresponse, _ := SearchForData(1, "gold", params)
	// onlineData := OnlineData{}
	// DB.Where("guid = ?", searchresponse.Results[0].GUID).Find(&onlineData)
	// if onlineData.Primarydate != "" {
	// 	searchdate = onlineData.Primarydate
	// }

	var id []byte
	var date time.Time
	var originalUrl, title, description, imageUrl string
	iter := session.Query(`SELECT id, title, original_url, date, description FROM response WHERE date >= ? AND date < ? ALLOW FILTERING`, FromDate, Today).Iter()
	for iter.Scan(&id, &title, &originalUrl, &date, &description) {
		termcount := 0
		for _, st := range searchterm {
			termcount += TermCheck(st, description)
			termcount += TermCheck(st, title) //if term in title too it adds weight
		}
		if termcount > 0 {
			var tmpNA NewsArticle
			session.Query(`SELECT url FROM image WHERE id = ? LIMIT 1 ALLOW FILTERING`, id).Scan(&imageUrl)
			tmpNA.Date = date
			tmpNA.Title = title
			tmpNA.Url = originalUrl
			tmpNA.ImageUrl = imageUrl
			tmpNA.Score = termcount
			newsarticles = append(newsarticles, tmpNA)
		}
	}

	// if err := iter.Close(); err != nil {
	// 	return newsarticles, err
	// }

<<<<<<< HEAD
	return newsarticle, nil
=======
	sortutil.DescByField(newsarticles, "Score")
	return newsarticles, nil
}
>>>>>>> feature/visualisations

func TermCheck(term string, passage string) int {
	descriptions := strings.Split(passage, " ")
	for _, d := range descriptions {
		if d == term {
			return 1
		}
	}
	return 0
}

// sort.Sort(byScore(newsarticles))
// type byScore []NewsArticle
// func (v byScore) Len() int      { return len(v) }
// func (v byScore) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
// func (v byScore) Less(i, j int) bool { v[i].Score < v[j].Score }
