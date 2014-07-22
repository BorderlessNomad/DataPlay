package main

import (
	"crypto/rand"
	"fmt"
	"github.com/fzzy/radix/redis"
	"net/http"
	"os"
	"time"
)

func IsUserLoggedIn(res http.ResponseWriter, req *http.Request) bool {
	cookie, _ := req.Cookie("DPSession")

	fmt.Println("IsUserLoggedIn Cookie:", cookie.Value)

	c, err := GetRedisConnection()
	if cookie != nil && err == nil {
		defer c.Close()
		r := c.Cmd("GET", cookie.Value)
		i, err := r.Int() // Get back from Redis the Int value of that cookie.
		if err != nil {
			return false
		}

		// There might be cases where redis could store 0 (meaning there is no logged in user)
		// for that session, Meaning that we need to check for when this happens.
		if i != 0 {
			return true
		} else {
			return false // there is no zero user.
		}
	} else {
		return false
	}
}

func GetUserID(res http.ResponseWriter, req *http.Request) int {
	cookie, _ := req.Cookie("DPSession")

	fmt.Println("GetUserID Cookie:", cookie.Value)

	c, err := GetRedisConnection()
	if cookie != nil && err == nil {
		defer c.Close()

		r := c.Cmd("GET", cookie.Value)
		i, err := r.Int() // Get back from Redis the Int value of that cookie.
		if err != nil {
			return 0
		}

		// There might be cases where redis could store 0 (meaning there is no logged in user)
		// for that session, Meaning that we need to check for when this happens.
		if i != 0 {
			return i
		} else {
			return 0 // there is no zero user.
		}
	} else {
		return 0
	}
}

func SetSession(res http.ResponseWriter, req *http.Request, userid int) (e error) {
	NewSessionID := randString(64)
	c, err := GetRedisConnection()
	if err != nil {
		return fmt.Errorf("Could not connect to redis server to make session")
	}
	defer c.Close()
	r := c.Cmd("SET", NewSessionID, userid)
	if r.Err != nil {
		return fmt.Errorf("Could not store session in Redis") // I'm not sure how this would ever happen (Plane crash in mid query?) but protecting against it.
	}

	NewCookie := &http.Cookie{
		Name:    "DPSession",
		Value:   NewSessionID,
		Path:    "/",
		Expires: time.Now().AddDate(1, 0, 0), // +1 Year
	}
	http.SetCookie(res, NewCookie)
	return e
}

func ClearSession(res http.ResponseWriter, req *http.Request) (e error) {
	cookie, _ := req.Cookie("DPSession")
	fmt.Println("ClearSession Cookie:", cookie.Value)

	c, errc := GetRedisConnection()
	if errc != nil {
		return fmt.Errorf("Could not connect to redis server to make session")
	}

	defer c.Close()

	if cookie == nil {
		return fmt.Errorf("No cookie found")
	}

	get := c.Cmd("GET", cookie.Value)
	_, errg := get.Int() // Get back from Redis the Int value of that cookie.
	if errg != nil {
		return fmt.Errorf("Could not find session in Redis")
	}

	set := c.Cmd("SET", cookie.Value, 0)
	if set.Err != nil {
		return fmt.Errorf("Could not update session in Redis")
	}

	NewCookie := &http.Cookie{
		Name:    "DPSession",
		Value:   "",
		Path:    "/",
		Expires: time.Now().AddDate(-1, 0, 0), // -1 Year = Expired
	}
	http.SetCookie(res, NewCookie)
	return e
}

// Gives a nice "clean" random string from the a-z, A-Z and 0-9
func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func GetRedisConnection() (c *redis.Client, err error) {
	redishost := "10.0.0.2:6379"
	if os.Getenv("redishost") != "" {
		redishost = os.Getenv("redishost")
	}

	c, err = redis.DialTimeout("tcp", redishost, time.Duration(10)*time.Second)

	if err != nil {
		Logger.Println("Could not connect to the redis server. Is it running? Sessions wont work otherwise !!1")
	}

	return c, err
}
