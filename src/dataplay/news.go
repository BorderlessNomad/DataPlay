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
	response := []NewsArticle{}
	return response, nil

}
