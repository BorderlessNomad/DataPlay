package main

import (
	"crypto/sha1"
	"fmt"
	goq "github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strings"
)

func ImportAllDatasets(url string, guid string) {
	// This call is desgined to fetch all of the links that are needed to be downloaded (on data.gov.uk)
	// And then fire off go routines to go and fetch them.
	var e error
	var doc *goq.Document
	fmt.Println("Loading URL", url)
	if doc, e = goq.NewDocument(url); e != nil {
		panic("Unable to fetch page!")
	}
	doc.Find(".dropdown-menu").Each(func(i int, s *goq.Selection) {
		s.Find("a").Each(func(i int, s *goq.Selection) {
			url, exists := s.Attr("href")
			html, _ := s.Html()
			if exists && strings.Contains(html, "icon-download-alt") {
				fmt.Println(url)
				go DownloadDataset(url, guid)
			}
		})
	})
}

func DownloadDataset(url string, guid string) {
	// This function will inhale a CSV file and if it is not formatted stupidly (ALA: most of the stuff on data.gov.uk)
	// it will make a SQL table for it and then put the fact that the data exists in the online table allowing the client
	// to move along and to ensure that it wont be downloaded twice.
	fmt.Println("Downloading dataset", url)
	response, _ := http.Get(url)
	fmt.Println(response.Header)
	if response.Header.Get("Content-Type") == "text/csv" {

		full, _ := ioutil.ReadAll(response.Body)
		hash := sha1.New()
		hash.Write([]byte(url))
		ioutil.WriteFile("./"+guid+"_"+fmt.Sprintf("%x", hash.Sum(nil)), full, 0667)
	}

}
