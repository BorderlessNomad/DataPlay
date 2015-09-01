package main

import (
	"crypto/rand"
	"net/http"
	"strconv"
	"time"
)

func IsUserLoggedIn(res http.ResponseWriter, req *http.Request) bool {
	cookie, _ := req.Cookie("DPSession")

	c, err := GetRedisConnection()
	if cookie == nil || len(cookie.Value) < 1 || err != nil {
		return false
	}

	defer c.Close()

	r := c.Cmd("GET", cookie.Value)
	user, err := r.Str()

	// There might be cases where redis could store 0 (meaning there is no logged in user)
	// for that session, Meaning that we need to check for when this happens.
	if err != nil || user == "" || user == "0" {
		return false
	}

	return true
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
	i, err := r.Str()
	if err != nil {
		return 0, &appError{err, "Unable to parse GET value (" + err.Error() + ")", http.StatusInternalServerError}
	}

	// There might be cases where redis could store 0 (meaning there is no logged in user)
	// for that session, Meaning that we need to check for when this happens.
	if i == "" || i == "0" {
		return 0, &appError{err, "No such session found.", http.StatusUnauthorized}
	}

	user, _ := strconv.Atoi(i)
	return user, nil
}

func SetSession(userId int) (*http.Cookie, *appError) {
	if userId <= 0 {
		return nil, &appError{nil, "No UserID specified.", http.StatusInternalServerError}
	}

	c, err := GetRedisConnection()
	if err != nil {
		return nil, &appError{err, "Unable to connect with Redis server", http.StatusInternalServerError}
	}

	defer c.Close()

	NewSessionID := randString(64)
	r := c.Cmd("SET", NewSessionID, userId)
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
		return nil, &appError{errc, "Invalid or Empty session", http.StatusBadRequest}
	}

	get := c.Cmd("GET", cookie)
	_, errg := get.Str()
	if errg != nil {
		return nil, &appError{errg, "No such session found", http.StatusNotFound}
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
