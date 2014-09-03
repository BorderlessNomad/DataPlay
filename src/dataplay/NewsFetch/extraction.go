package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const layout = "2006-01-02"

var (
	Host = "http://api.embed.ly"
)

type Client struct {
	key string
}

func NewClient(key string) *Client {
	return &Client{key}
}

// The main exported function that will extract urls.
func (c *Client) Extract(urls []string, options Options) ([]string, error) {

	responses := make([]string, len(urls))

	for i := 0; i < len(urls); i += 10 {
		to := len(urls)
		if to > i+10 {
			to = i + 10
		}
		res, err := c.extract(urls[i:to], options, i)

		if err != nil {
			return nil, err
		}

		reslen := to - i
		if reslen > len(res) {
			reslen = len(res)

		}
		for j := 0; j < reslen; j++ {
			responses[i+j] = res[j]
		}
	}

	return responses, nil
}

// extract will call extract 10 urls at max.
func (c *Client) extract(urls []string, options Options, place int) ([]string, error) {

	addr := Host + "/1/extract?"
	for i, u := range urls {
		urls[i] = url.QueryEscape(u)
	}
	if len(urls) == 0 {
		return nil, errors.New("At least one URL is required")
	} else if len(urls) == 1 {
		if len(urls[0]) == 0 {
			return nil, errors.New("URL cannot be empty")
		}
		addr += "url=" + urls[0]
	} else {
		for _, url := range urls {
			if len(url) == 0 {
				return nil, errors.New("A URL cannot be empty")
			}
		}
		addr += "urls=" + strings.Join(urls, ",")
	}

	v := url.Values{}
	v.Add("key", c.key)

	// Make the request.
	addr += "&" + v.Encode() + "&format=json"
	resp, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Got non 200 status code: %s %q", resp.Status, body)
	}

	// Read the JSON message from the body.
	response := []Response{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, err
	}

	responses := make([]string, len(response))

	for i, r := range response {
		var tmpResp Response
		tmp, _ := json.Marshal(r)
		json.Unmarshal(tmp, &tmpResp)
		var p int64
		p = tmpResp.Published
		date := publishedDate(p, place+i)
		tmpResp.Date = date.Format(layout)
		result, _ := json.Marshal(tmpResp)
		responses[i] = string(result)
	}

	return responses, nil
}

// addInt adds an int if non-zero.
func addInt(v *url.Values, name string, value int) {
	if value > 0 {
		v.Add(name, strconv.Itoa(value))
	}
}

// addBool adds a boolean value if set to true.
func addBool(v *url.Values, name string, value bool) {
	if value {
		v.Add(name, "true")
	}
}

// takes published date in raw format and returns UTC date format.
// If there is no published date, function will return the probable month and year of the article based on the array position of the url (approx 4358 urls per month starting from 2010-01-01)
func publishedDate(date int64, place int) time.Time {
	var i int64
	var x int
	var t time.Time
	published := date
	i = 0

	if published == 0 {
		x = place / 4358                                // get likely month number
		t = time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC) // starting from 2010-01-01
		t = t.AddDate(0, x, 0)                          // add number of months to base date

	} else {
		published = date / 1000
		t = time.Unix(published, i)
	}
	return t
}
