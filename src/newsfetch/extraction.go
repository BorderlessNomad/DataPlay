package main

import (
	"crypto/md5"
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

var (
	Host = "http://api.embed.ly"
)

type Client struct {
	key string
}

func NewClient(key string) *Client {
	return &Client{key}
}

func (c *Client) Extract(urls []string, options Options, startpos int) (error, int) {

	for i := startpos; i < len(urls); i += 10 {
		fmt.Printf("Extracting %d out of %d URLS\n", i, len(urls))
		fmt.Sprintf("Extracting")
		to := len(urls)
		if to > i+10 {
			to = i + 10
		}
		res, err := c.extract(urls[i:to], options, i)

		if err != nil {
			return err, i
		}

		reslen := to - i
		if reslen > len(res) {
			reslen = len(res)

		}

		for j := 0; j < reslen; j++ {
			writeToCass(res[j])
		}
	}

	return nil, 0
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
		return nil, errors.New("bad request")
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

// return md5 hash of string
func Hash(str string) []byte {
	data := []byte(str)
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// write json string to cassandra

func writeToCass(resp string) {
	session, _ := GetCassandraConnection("dataplay")
	defer session.Close()

	var r Response
	err := json.Unmarshal([]byte(resp), &r)
	if err != nil {
		panic(err)
	}

	date := time.Unix(r.Published/1000, 0)

	if err := session.Query(`INSERT INTO response (date, dummy, description, url, title)
		VALUES (?, ?, ?, ?, ?)`,
		date, 1, r.Description, r.Url, r.Title).Exec(); err != nil {
		fmt.Println("HELP1!", err)
	}

	for _, k := range r.Keywords {
		if err := session.Query(`INSERT INTO keyword (date, dummy, name, url)
			VALUES (?, ?, ?, ?)`,
			date, 1, k.Name, r.Url).Exec(); err != nil {
			fmt.Println("HELP3!", err)
		}
	}

	for _, e := range r.Entities {
		if err := session.Query(`INSERT INTO entity (date, dummy, name, url)
			VALUES (?, ?, ?, ?)`,
			date, 1, e.Name, r.Url).Exec(); err != nil {
			fmt.Println("HELP4!", err)
		}
	}

	if len(r.Images) > 0 { // check if there are any images
		if err := session.Query(`INSERT INTO image (date, dummy, pic_url, url)
			VALUES (?, ?, ?, ?)`,
			date, 1, r.Images[0].Url, r.Url).Exec(); err != nil {
			fmt.Println("HELP5!", err)
		}
	}

	for _, ra := range r.Related {
		if err := session.Query(`INSERT INTO related (date, dummy, description, title, related_url, url)
			VALUES (?, ?, ?, ?, ?, ?)`,
			date, 1, ra.Description, ra.Title, ra.Url, r.Url).Exec(); err != nil {
			fmt.Println("HELP6!", err)
		}
	}
}
