package main

import (
	"crypto/rand"
	"net/http"
	"time"
)

func IsUserLoggedIn(res http.ResponseWriter, req *http.Request) bool {
	cookie, _ := req.Cookie("DPSession")

	c, err := GetRedisConnection()
	if cookie == nil || err != nil {
		return false
	}

	defer c.Close()

	user, err := c.Cmd("GET", cookie.Value).Int() // Get back from Redis the Int value of that cookie.
	if err != nil {
		return false
	}

	// There might be cases where redis could store 0 (meaning there is no logged in user)
	// for that SESSION, Meaning that we need to check for when this happens.
	if user != 0 {
		return true
	}

	return false // there is no zero user.
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

	user, err := c.Cmd("GET", cookie).Int() // Get back from Redis the Int value of that cookie.
	if err != nil {
		return 0, &appError{err, "Unable to parse GET value", http.StatusInternalServerError}
	}

	// There might be cases where redis could store 0 (meaning there is no logged in user)
	// for that SESSION, Meaning that we need to check for when this happens.
	if user <= 0 {
		return 0, &appError{err, "No such User found.", http.StatusUnauthorized}
	}

	return user, nil
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
	err = c.Cmd("SET", NewSessionID, userid).Err
	if err != nil {
		return nil, &appError{err, "Unable to store SESSION in Redis", http.StatusInternalServerError}
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
		return nil, &appError{errc, "No SESSION found", http.StatusBadRequest}
	}

	_, err := c.Cmd("GET", cookie).Int() // Get back from Redis the Int value of that cookie.
	if err != nil {
		return nil, &appError{err, "Unable to find SESSION in Redis", http.StatusInternalServerError}
	}

	err = c.Cmd("SET", cookie, 0).Err
	if err != nil {
		return nil, &appError{err, "Unable to update SESSION in Redis", http.StatusInternalServerError}
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
