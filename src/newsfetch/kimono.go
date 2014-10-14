package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Container struct {
	Name           string `json:"name"`
	Frequency      string `json:"frequency"`
	Version        int    `json:"version"`
	Newdata        bool   `json:"newdata"`
	Lastrunstatus  string `json:"lastrunstatus"`
	Thisversionrun string `json:"thisversionrun"`
	Lastsuccess    string `json:"lastsuccess"`
	Results        Result `json:"results"`
	Count          int    `json:"count"`
}

type Result struct {
	Collections []Collection `json:"collection1"`
}

type Collection struct {
	Properties Property `json:"property1"`
}

type Property struct {
	Text string `json:"text"`
	Href string `json:"href"`
}

/// gets the day's BBC News articles URLs via google
func DailyKimono() error {

	t := time.Now()

	sd := fmt.Sprintf("%02d", t.AddDate(0, 0, -1).Day())
	sm := fmt.Sprintf("%02d", t.AddDate(0, 0, -1).Month())
	sy := fmt.Sprintf("%4d", t.AddDate(0, 0, -1).Year())
	ed := fmt.Sprintf("%02d", t.Day())
	em := fmt.Sprintf("%02d", t.Month())
	ey := fmt.Sprintf("%4d", t.Year())

	tbs := "cdr%3A1%2Ccd_min%3A" + sd + "%2F" + sm + "%2F" + sy + "%2Ccd_max%3A" + ed + "%2F" + em + "%2F" + ey
	addr := "https://www.kimonolabs.com/api/bvt2gs12?apikey=J6hUZ7JPZbnc2ASz6YoG6J4DeI6TuxvL" + "&tbs=" + tbs
	fmt.Println("KIMONO", addr)

	resp, err := http.Get(addr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Got non 200 status code: %s %q", resp.Status, body)
	}

	// Read the JSON message from the body.
	container := Container{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&container); err != nil {
		return err
	}

	var output string
	for _, x := range container.Results.Collections {
		output += x.Properties.Href + ",\n"
	}

	err = ioutil.WriteFile("dailyURLs.txt", []byte(output), 0644)
	if err != nil {
		return err
	}

	f, _ := os.OpenFile("URLs.txt", os.O_APPEND, 0666)
	f.Write([]byte(output))
	f.Close()

	return nil
}
