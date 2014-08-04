package main

import "fmt"
import "encoding/json"
import "io/ioutil"
import "net/http"
import "net/url"
import "strings"

// Does a single query
func query(input map[string]interface{}, connectorGuid string, client *http.Client, userguid string, apikey string) {

	inputString, _ := json.Marshal(input)

	Url, _ := url.Parse("https://api.import.io/store/connector/" + connectorGuid + "/_query")
	parameters := url.Values{}
	parameters.Add("_user", userguid)
	parameters.Add("_apikey", apikey)
	Url.RawQuery = parameters.Encode()

	request, _ := http.NewRequest("POST", Url.String(), strings.NewReader(string(inputString)))
	request.Header.Add("Content-Type", "application/json")
	resp, _ := client.Do(request)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("SCRAPEIT", string(body[:]))

}

func mainB() {

	client := &http.Client{}

	userguid := "cf592fba-bd1f-4128-8e98-e729c2bb7dec"
	apikey := "aledxqRLOCLFo9O7cYeeC58aotifmZbL2C57Mg1zicz6ZLVSY94xttvI9AjeV1Fw9DpBg2y/cbrNZXM23yiWBg=="

	// Query for tile BBC Scraper
	query(map[string]interface{}{
		"input": map[string]interface{}{
			"start_day":   "01",
			"start_month": "01",
			"start_year":  "2010",
			"end_day":     "10",
			"end_month":   "01",
			"end_year":    "2010",
		},
	}, "054f1f45-4755-40f5-b9fc-48ab8f80214e", client, userguid, apikey)

}
