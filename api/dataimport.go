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
	"os"
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

func ImportAllDatasets(url string, guid string) {
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
				go DownloadDataset(url, guid)
			}
		})
	})
}

func DownloadDataset(url string, guid string) {
	fmt.Println("Downloading dataset", url)
	response, _ := http.Get(url)
	fmt.Println(response.Header)
	if response.Header.Get("Content-Type") == "text/csv" {

		full, _ := ioutil.ReadAll(response.Body)
		hash := sha1.New()
		hash.Write([]byte(url))
		ioutil.WriteFile("./temp"+fmt.Sprintf("%x", hash.Sum(nil)), full, 0667)
		database := msql.GetDB()
		defer database.Close()
		rows := strings.Split(string(full[:]), "\n")
		colcount := len(strings.Split(rows[0], ","))
		var tablebuilder string
		tablebuilder = "CREATE TABLE `" + fmt.Sprintf("%x", hash.Sum(nil)) + "` ("
		for i := 0; i < colcount; i++ {
			tablebuilder = tablebuilder + "`Column " + fmt.Sprintf("%d", i) + "` TEXT NULL,"
		}
		tablebuilder = tablebuilder + "`Column " + fmt.Sprintf("%d", colcount) + "` TEXT NULL"
		tablebuilder = tablebuilder + ") COLLATE='latin1_swedish_ci' ENGINE=InnoDB;"
		fmt.Println(tablebuilder)
		_, e := database.Exec(tablebuilder)
		// panic(e)
		if e != nil {
			panic(e)
		}
		var tableloader string

		tableloader = tableloader + "LOAD DATA LOCAL INFILE '" + "./temp" + fmt.Sprintf("%x", hash.Sum(nil)) + "'\n"
		tableloader = tableloader + "INTO TABLE " + fmt.Sprintf("%x", hash.Sum(nil)) + "\n"
		tableloader = tableloader + "FIELDS TERMINATED BY ','" + "\n"
		tableloader = tableloader + "OPTIONALLY ENCLOSED BY '\"'" + "\n"
		tableloader = tableloader + "LINES TERMINATED BY '\\r\\n'\n"
		tableloader = tableloader + "IGNORE 1 LINES\n"
		tableloader = tableloader + "("
		for i := 0; i < colcount; i++ {
			tableloader = tableloader + "`Column " + fmt.Sprintf("%d", i) + "`,"
		}
		tableloader = tableloader + "`Column " + fmt.Sprintf("%d", colcount) + "`)"

		_, e = database.Exec(tableloader)
		if e != nil {
			panic(e)
		}
		database.Exec("INSERT INTO `priv_onlinedata` (`GUID`, `DatasetGUID`, `TableName`) VALUES (?, ?, ?);", guid, "IForGotWhatIwantedToPutHere", fmt.Sprintf("%x", hash.Sum(nil)))
		os.Remove("./temp" + fmt.Sprintf("%x", hash.Sum(nil)))
		// LOAD DATA INFILE 'detection.csv'
		// INTO TABLE calldetections
		// FIELDS TERMINATED BY ','
		// OPTIONALLY ENCLOSED BY '"'
		// LINES TERMINATED BY ',,,\r\n'
		// IGNORE 1 LINES
		// (date, name, type, number, duration, addr, pin, city, state, country, lat, log)

		// CREATE TABLE `removeme` (
		// 	`Column 1` TEXT NULL,
		// 	`Column 2` TEXT NULL,
		// 	`Column 3` TEXT NULL,
		// 	`Column 4` TEXT NULL,
		// 	`Column 5` TEXT NULL,
		// 	`Column 6` TEXT NULL,
		// 	`Column 7` TEXT NULL,
		// 	`Column 8` TEXT NULL,
		// 	`Column 9` TEXT NULL
		// )
		// COLLATE='latin1_swedish_ci'
		// ENGINE=InnoDB;
		// fmt.Println(colcount)
	}

}
