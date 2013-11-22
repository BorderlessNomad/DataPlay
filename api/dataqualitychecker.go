package api

import (
	cache "../cache"
	msql "../databasefuncs"
	"encoding/json"
	// "fmt"
	goq "github.com/PuerkitoBio/goquery"
	"github.com/codegangsta/martini"
	"github.com/mattn/go-session-manager"
	"net/http"
	"strings"
)

type CheckDataQualityResponce struct {
	Amount  int
	Request string
}

func CheckDataQuality(res http.ResponseWriter, req *http.Request, prams martini.Params, manager *session.SessionManager) string {

	database := msql.GetDB()
	defer database.Close()
	if prams["id"] == "" {
		http.Error(res, "There was no ID request", http.StatusBadRequest)
	}

	test := cache.GetCache("QCheck::" + prams["id"])
	if test == "" {
		var ckan_url string
		e := database.QueryRow("SELECT ckan_url FROM `index` WHERE GUID LIKE ? LIMIT 10", prams["id"]+"%").Scan(&ckan_url)

		url := strings.Replace(strings.Replace(ckan_url, "//", "/", -1), "http:/", "http://", 1)

		var doc *goq.Document

		if doc, e = goq.NewDocument(url); e != nil {
			http.Error(res, "Unable to fetch page", http.StatusServiceUnavailable)
		}

		var count int

		doc.Find(".icon-star").Each(func(i int, s *goq.Selection) {
			count++
		})

		returnobj := CheckDataQualityResponce{
			Amount:  count,
			Request: prams["id"],
		}
		b, _ := json.Marshal(returnobj)
		cache.SetCache("QCheck::"+prams["id"], string(b[:]))
		return string(b[:])
	} else {
		return test
	}
}
