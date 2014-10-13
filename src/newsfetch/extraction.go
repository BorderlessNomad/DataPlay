package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	Host = "http://api.embed.ly"
)

type Client struct {
	key string
}

func NewClient(key string) *Client {
	return &Client{key}
}

func (c *Client) Extract(urls []string, options Options, startpos int) error {

	for i := startpos; i < len(urls); i += 10 {
		fmt.Printf("Extracting %d out of %d URLS\n", i, len(urls))
		fmt.Sprintf("Extracting")
		to := len(urls)
		if to > i+10 {
			to = i + 10
		}
		res, err := c.extract(urls[i:to], options, i)

		if err != nil {
			f, _ := os.OpenFile("log.txt", os.O_RDWR|os.O_APPEND, 0666)
			errStr := "FAILED AND RESTARTED ON URL " + strconv.Itoa(i) + " for reason " + err.Error()
			e := []byte(errStr + "\n")
			f.Write(e)
			fmt.Println("RE-STARTING...")
			Start(i)
		}

		reslen := to - i
		if reslen > len(res) {
			reslen = len(res)

		}

		for j := 0; j < reslen; j++ {
			writeToCass(res[j])
		}
	}

	return nil
}

// extract will call extract 10 urls at max.
func (c *Client) extract(urls []string, options Options, place int) ([]string, error) {

	addr := Host + "/1/extract?"
	for i, u := range urls {
		urls[i] = url.QueryEscape(u)
	}
	if len(urls) == 0 {
		return nil, errors.New("At least one Url is required")
	} else if len(urls) == 1 {
		if len(urls[0]) == 0 {
			return nil, errors.New("Url cannot be empty")
		}
		addr += "url=" + urls[0]
	} else {
		for _, url := range urls {
			if len(url) == 0 {
				return nil, errors.New("A Url cannot be empty")
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
		fmt.Println("RE-STARTING...")
		Start(place)
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
		h := Hash(tmpResp.Url)
		tmpResp.Id = h

		for i, _ := range tmpResp.Authors {
			tmpResp.Authors[i].Id = h
		}
		for i, _ := range tmpResp.Keywords {
			tmpResp.Keywords[i].Id = h
		}
		for i, _ := range tmpResp.Entities {
			tmpResp.Entities[i].Id = h
		}
		for i, _ := range tmpResp.RelatedArticles {
			tmpResp.RelatedArticles[i].Id = h
		}
		for i, _ := range tmpResp.Images {
			tmpResp.Images[i].Id = h
		}
		var p int64
		p = tmpResp.Published
		date := publishedDate(p, place+i)
		tmpResp.Date = date
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

// return md5 hash of string
func Hash(str string) []byte {
	data := []byte(str)
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// write json string to cassandra
func writeToCass(resp string) {
	session, _ := GetCassandraConnection("dp")
	defer session.Close()

	var r Response
	err := json.Unmarshal([]byte(resp), &r)
	if err != nil {
		panic(err)
	}

	if err := session.Query(`INSERT INTO response (id, original_url, url, type, provider_name, provider_url, 
		provider_display, favicon_url, title, description, date, published, lead, content) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.Id, r.OriginalUrl, r.Url, r.Type, r.ProviderName, r.ProviderUrl,
		r.ProviderDisplay, r.FaviconUrl, r.Title, r.Description, r.Date, r.Published, r.Lead, r.Content).Exec(); err != nil {
		fmt.Println("HELP1!", err)
	}

	for _, a := range r.Authors {
		if err := session.Query(`INSERT INTO author (id, date, name, url) 
			VALUES (?, ?, ?, ?)`,
			a.Id, r.Date, a.Name, a.Url).Exec(); err != nil {
			fmt.Println("HELP2!", err)
		}
	}

	for _, k := range r.Keywords {
		if err := session.Query(`INSERT INTO keyword (id, url, date, score, name) 
			VALUES (?, ?, ?, ?, ?)`,
			k.Id, r.Url, r.Date, k.Score, k.Name).Exec(); err != nil {
			fmt.Println("HELP3!", err)
		}
	}

	for _, e := range r.Entities {
		if err := session.Query(`INSERT INTO entity (id, url, date, count, name) 
			VALUES (?, ?, ?, ?, ?)`,
			e.Id, r.Url, r.Date, e.Count, e.Name).Exec(); err != nil {
			fmt.Println("HELP4!", err)
		}
	}

	for index, i := range r.Images {
		if err := session.Query(`INSERT INTO image (pic_index, id, date, caption, url, width, height, entropy, size) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			index, i.Id, r.Date, i.Caption, i.Url, i.Width, i.Height, i.Entropy, i.Size).Exec(); err != nil {
			fmt.Println("HELP5!", err)
		}
	}

	for _, ra := range r.RelatedArticles {
		if err := session.Query(`INSERT INTO related (id, date, description, title, url, thumbnail_width, score, thumbnail_height, thumbnail_url) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			ra.Id, r.Date, ra.Description, ra.Title, ra.Url, ra.ThumbnailWidth, ra.Score, ra.ThumbnailHeight, ra.ThumbnailUrl).Exec(); err != nil {
			fmt.Println("HELP6!", err)
		}
	}

	f, _ := os.OpenFile("jsonout.txt", os.O_APPEND, 0666)
	f.WriteString(resp)
	f.Close()
}
