package main

import (
	"encoding/json"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/codegangsta/martini"
	"github.com/pmylund/sortutil"
	// "io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Tweet struct {
	Text     string    `json:"comment"`
	Name     string    `json:"name"`
	User     string    `json:"username"`
	Created  time.Time `json:"created"`
	Retweets int       `json:"retweets"`
	Source   string    `json:"source"`
	Hashtags []string  `json:"hashtags"`
	Urls     []string  `json:"urls"`
	Media    []string  `json:"mediaurls"`
}

type Sanitized struct {
	Result string `json:"result"`
}

func GetTweetsHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	if params["searchterms"] == "" {
		http.Error(res, "No search term", http.StatusBadRequest)
		return ""
	}

	anaconda.SetConsumerKey("37a5BBIeovRJ6eit5Bv2HFCJV")
	anaconda.SetConsumerSecret("IyclAamSNVeCZkrOypNpCwMZYpZX6RMbRlN0TdL1NjghhyKlSU")
	api := anaconda.NewTwitterApi("2834288205-7noj46EGdEDsXRu9wRou4hEC7lkM3ptNC3bktvo", "6LGn5IcZcWGEfSvpGU6rzfp4rqEPc8GVRM23qFHJoJsOg")
	terms := strings.Split(params["searchterms"], "_")
	tweets := []Tweet{}

	for _, term := range terms {
		searchResult, _ := api.GetSearch(term, nil)
		tmpTweet := Tweet{}

		for _, tweet := range searchResult {
			if tweet.User.Lang == "en" && !strings.Contains(tweet.Text, "RT @") && !tweet.PossiblySensitive {
				tmpTweet.Created, _ = tweet.CreatedAtTime()

				tmpTweet.Retweets = tweet.RetweetCount
				tmpTweet.Source = tweet.Source
				// tmpTweet.Text = Sanitize(tweet.Text)
				fmt.Println("JUNGLE", tweet.Text)
				fmt.Println("JUNGLE2", tmpTweet.Text)
				tmpTweet.Name = tweet.User.Name
				tmpTweet.User = tweet.User.ScreenName

				for _, h := range tweet.Entities.Hashtags {
					tmpTweet.Hashtags = append(tmpTweet.Hashtags, h.Text)
				}

				for _, m := range tweet.Entities.Media {
					tmpTweet.Media = append(tmpTweet.Media, m.Url)
				}

				for _, u := range tweet.Entities.Urls {
					tmpTweet.Urls = append(tmpTweet.Urls, u.Url)
				}

				tweets = append(tweets, tmpTweet)
			}
		}
	}

	n := 10
	if len(tweets) < 10 {
		n = len(tweets)
	}

	sortutil.DescByField(tweets, "Retweets")
	r, err := json.Marshal(tweets[0:n])
	if err != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

// func Sanitize(text string) string {
// 	strings.Replace(text, " ", "%20", -1)
// 	url := "http://www.purgomalum.com/service/json?text=" + text
// 	resp, _ := http.Get(url)
// 	defer resp.Body.Close()
// 	r, _ := ioutil.ReadAll(resp.Body)
// 	result := Sanitized{}
// 	json.Unmarshal(r, &result)
// 	return string(result.Result)
// }
