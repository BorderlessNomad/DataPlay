package main

import (
	"encoding/json"
	"github.com/ChimeraCoder/anaconda"
	"github.com/codegangsta/martini"
	"github.com/kennygrant/sanitize"
	"github.com/pmylund/sortutil"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Tweet struct {
	Text       string    `json:"comment"`
	Name       string    `json:"name"`
	User       string    `json:"username"`
	UserAvatar string    `json:"useravatar"`
	Created    time.Time `json:"created"`
	Retweets   int       `json:"retweets"`
	Id         string    `json:"id"`
	Hashtags   []string  `json:"hashtags"`
	Urls       []string  `json:"urls"`
	Media      []string  `json:"mediaurls"`
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

	anaconda.SetConsumerKey("okJJRSM9oHTyuYY3ePoYbjtXA")
	anaconda.SetConsumerSecret("zj6iaxrJzwY2B9HxiZk7atdo2yQxGheJSoUiIq7Z8fQ7yut7Kk")
	api := anaconda.NewTwitterApi("73992277-Bf13OY4Nn0MmI9amLyDcfUe57MXc2tk6YdQ0eLiVA", "VZ3JUSdgG8NkZvmmnI7oiLHINlCLiTZUVBKJv1VDdqK4Q")
	terms := strings.Split(params["searchterms"], "_")
	tweets := []Tweet{}

	for _, term := range terms {
		v := url.Values{}
		v.Set("count", "40") // take sample 40 tweets but will only use first best 10
		searchResult, _ := api.GetSearch(term, v)

		tmpTweet := Tweet{}

		for _, tweet := range searchResult.Statuses {
			if tweet.User.Lang == "en" && !strings.Contains(tweet.Text, "RT @") && !tweet.PossiblySensitive && len(tweets) <= 10 {
				tmpTweet.Created, _ = tweet.CreatedAtTime()

				tmpTweet.Retweets = tweet.RetweetCount
				tmpTweet.Id = tweet.IdStr
				tmpTweet.Text = ProfanityCheck(tweet.Text)
				tmpTweet.Name = tweet.User.Name
				tmpTweet.User = tweet.User.ScreenName
				tmpTweet.UserAvatar = tweet.User.ProfileImageURL

				for _, h := range tweet.Entities.Hashtags {
					tmpTweet.Hashtags = append(tmpTweet.Hashtags, h.Text)
				}

				for _, m := range tweet.Entities.Media {
					tmpTweet.Media = append(tmpTweet.Media, m.Url)
				}

				for _, u := range tweet.Entities.Urls {
					tmpTweet.Urls = append(tmpTweet.Urls, u.Url)
				}

				if tmpTweet.Text != "" { // if the profanity check was passed
					tweets = append(tweets, tmpTweet)
				}
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

func ProfanityCheck(text string) string {
	strings.Replace(text, " ", "%20", -1)
	t := sanitize.Path(text)
	url := "http://www.purgomalum.com/service/json?text=" + t
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	r, _ := ioutil.ReadAll(resp.Body)
	sanitized := Sanitized{}
	json.Unmarshal(r, &sanitized)

	if !strings.Contains(sanitized.Result, "***") { // if there were no profane words returns the original string
		return text
	} else {
		return ""
	}
}
