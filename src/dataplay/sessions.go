package main

import (
	"crypto/rand"
	"github.com/fzzy/radix/redis"
	"net/http"
	"os"
	"time"
)

func IsUserLoggedIn(res http.ResponseWriter, req *http.Request) bool {
	cookie, _ := req.Cookie("DPSession")
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

func GetUserID(cookie string) (int, *appError) {
	if len(cookie) <= 0 {
		return 0, &appError{nil, "No Session specified.", http.StatusUnauthorized}
	}

	c, err := GetRedisConnection()
	if err != nil {
		return 0, &appError{err, "Unable to connect with Redis server", http.StatusInternalServerError}
	}

	defer c.Close()

	r := c.Cmd("GET", cookie)
	i, err := r.Int() // Get back from Redis the Int value of that cookie.
	if err != nil {
		return 0, &appError{err, "Unable to parse GET value", http.StatusInternalServerError}
	}

	// There might be cases where redis could store 0 (meaning there is no logged in user)
	// for that session, Meaning that we need to check for when this happens.
	if i <= 0 {
		return 0, &appError{err, "No such session found.", http.StatusUnauthorized}
	}

	return i, nil
}

func SetSession(userid int) (*http.Cookie, *appError) {
	if userid <= 0 {
		return nil, &appError{nil, "No UserID specified.", http.StatusInternalServerError}
	}

	c, err := GetRedisConnection()
	if err != nil {
		return nil, &appError{err, "Unable to connect with Redis server", http.StatusInternalServerError}
	}

	defer c.Close()

	NewSessionID := randString(64)
	r := c.Cmd("SET", NewSessionID, userid)
	if r.Err != nil {
		return nil, &appError{r.Err, "Unable to store session in Redis", http.StatusInternalServerError}
	}

	NewCookie := &http.Cookie{
		Name:    "DPSession",
		Value:   NewSessionID,
		Path:    "/",
		Expires: time.Now().AddDate(1, 0, 0), // +1 Year
	}

	return NewCookie, nil
}

func ClearSession(cookie string) (*http.Cookie, *appError) {
	c, errc := GetRedisConnection()
	if errc != nil {
		return nil, &appError{errc, "Unable to connect with Redis server", http.StatusInternalServerError}
	}

	defer c.Close()

	if len(cookie) <= 0 {
		return nil, &appError{errc, "No session found", http.StatusBadRequest}
	}

	get := c.Cmd("GET", cookie)
	_, errg := get.Int() // Get back from Redis the Int value of that cookie.
	if errg != nil {
		return nil, &appError{errg, "Unable to find session in Redis", http.StatusInternalServerError}
	}

	set := c.Cmd("SET", cookie, 0)
	if set.Err != nil {
		return nil, &appError{set.Err, "Unable to update session in Redis", http.StatusInternalServerError}
	}

	NewCookie := &http.Cookie{
		Name:    "DPSession",
		Value:   "",
		Path:    "/",
		Expires: time.Now().AddDate(-1, 0, 0), // -1 Year = Expired
	}

	return NewCookie, nil
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
