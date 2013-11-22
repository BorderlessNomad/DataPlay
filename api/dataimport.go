package api

import (
	msql "../databasefuncs"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	goq "github.com/PuerkitoBio/goquery"
	"github.com/codegangsta/martini"
	"github.com/mattn/go-session-manager"
	"io/ioutil"
	"net/http"
	"strings"
)

type CheckImportResponce struct {
	State   string
	Request string
}

func CheckImportStatus(res http.ResponseWriter, req *http.Request, prams martini.Params, manager *session.SessionManager) string {

	database := msql.GetDB()
	defer database.Close()
	if prams["id"] == "" {
		http.Error(res, "There was no ID request", http.StatusBadRequest)
	}

	var count int
	database.QueryRow("SELECT COUNT(*) FROM `priv_onlinedata` WHERE GUID = ? LIMIT 10", prams["id"]).Scan(&count)
	var state string

	if count != 0 {
		state = "online"
	} else {
		state = "offline"
	}

	returnobj := CheckImportResponce{
		State:   state,
		Request: prams["id"],
	}
	b, _ := json.Marshal(returnobj)
	return string(b[:])
}

func ImportAllDatasets(url string) {
	var e error
	var doc *goq.Document
	fmt.Println("Loading URL", url)
	if doc, e = goq.NewDocument(url); e != nil {
		panic("Unable to fetch page!")
	}
	// urls := make([]string, 0)
	doc.Find(".dropdown-menu").Each(func(i int, s *goq.Selection) {
		s.Find("a").Each(func(i int, s *goq.Selection) {
			url, exists := s.Attr("href")
			html, _ := s.Html()
			if exists && strings.Contains(html, "icon-download-alt") {
				fmt.Println(url)
				go DownloadDataset(url)
			}
		})
	})
}

func DownloadDataset(url string) {
	fmt.Println("Downloading dataset", url)
	response, _ := http.Get(url)
	full, _ := ioutil.ReadAll(response.Body)
	hash := sha1.New()
	hash.Write(url)
	ioutil.WriteFile("./temp"+fmt.Sprintf("%x", url), full, 0667)
	database := msql.GetDB()
	defer database.Close()
	rowlength := len(strings.Split(response, ","))

}
