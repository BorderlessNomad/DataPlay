package main

import (
	"crypto/sha1"
	"fmt"
	goq "github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func ImportAllDatasets(url string) {
	// This call is desgined to fetch all of the links that are needed to be downloaded (on data.gov.uk)
	// And then fire off go routines to go and fetch them.
	ourl := strings.Replace(url, "\r", "", -1)
	var e error
	var doc *goq.Document
	// fmt.Println("Loading URL '" + url + "'")
	if doc, e = goq.NewDocument(ourl); e == nil {
		// fmt.Println(doc.Html())
		runtime.GC() // :( why is this a thing
		time.Sleep(time.Millisecond * 250)
		doc.Find(".dropdown-menu").Each(func(i int, s *goq.Selection) { // The dropdown menu on data.gov that gives you the download options
			s.Find("a").Each(func(i int, s *goq.Selection) {
				url, exists := s.Attr("href")
				html, _ := s.Html()
				if exists && strings.Contains(html, "icon-download-alt") { // Select the download button, by the CSS icon it has
					go DownloadDataset(url, ourl)
				}
			})
		})
	}
	doc = nil
}

func JustBloddyHashIt(input []byte) string {
	hash := sha1.New()
	hash.Write(input)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// This function will inhale a CSV file and if it is not formatted stupidly (ALA: most of the stuff on data.gov.uk)
// it will make a SQL table for it and then put the fact that the data exists in the online table allowing the client
// to move along and to ensure that it wont be downloaded twice.
func DownloadDataset(url string, guid string) {

	fmt.Println("Downloading dataset", url)
	response, e := http.Get(url)
	if e == nil {
		// fmt.Println(response.Header)
		full, _ := ioutil.ReadAll(response.Body)
		filename := "./data/" + JustBloddyHashIt(full) + "_" + JustBloddyHashIt([]byte(guid))
		fmt.Println(guid + " | " + url + " => " + filename)
		if response.Header.Get("Content-Type") == "text/csv" {
			ioutil.WriteFile(filename+".csv", full, 0667)
		} else if response.Header.Get("Content-Type") == "application/vnd.ms-excel" {
			ioutil.WriteFile(filename+".xlsx", full, 0667)
		} else {
			ioutil.WriteFile(filename+".dunno", full, 0667)
		}
		full = nil     // Oh my god please GC I beg you
		response = nil // I'm sorry that I did string(42) that one time
		// just PLEASE GC
	}
}
