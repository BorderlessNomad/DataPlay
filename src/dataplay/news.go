package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	// "github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type NewsArticle struct {
	Date     time.Time `json:"date"`
	Title    string    `json:"title"`
	Url      string    `json:"url"`
	ImageUrl string    `json:"image_url"`
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
	newsarticle := []NewsArticle{}
	session, _ := GetCassandraConnection("dp") // create connection to cassandra
	defer session.Close()

	// searchresponse := SearchForData(uid int, keyword string, params map[string]string)

	// Select response where date
	// check description in response and date in response and if match search term get id and original url and use id to get image url

	// SELECT id, original_url, description FROM response WHERE date >= '2010-01-03' AND date < '2010-01-06' ALLOW FILTERING;

	//  select url from image where id = 0x7405b41f61f32f1f3992ba137cbebf82


	// // add all dated dateID between -n days and today to array
	// iter := session.Query(`SELECT id, date FROM response WHERE date >= ? AND date < ? ALLOW FILTERING`, FromDate, Today).Iter()
	// for iter.Scan(&id, &queryDate) {
	// 	dateID = append(dateID, string(id[:len(id)])+"!"+queryDate.Format(time.RFC3339))
	// }

	// if err := iter.Close(); err != nil {
	// 	///return err
	// }

	// for _, term := range terms {
	// 	iter := session.Query(`SELECT id FROM keyword WHERE name = ?`, term).Iter()
	// 	for iter.Scan(&id) {
	// 		var date time.Time
	// 		date = DateAndId(id, dateID)
	// 		if date.Year() > 1 {
	// 			tmpDT.Term = term
	// 			tmpDT.ID = id
	// 			tmpDT.Date = date
	// 			DatedTerms = append(DatedTerms, tmpDT)
	// 		}
	// 	}
	// }

	return newsarticle, nil

}
