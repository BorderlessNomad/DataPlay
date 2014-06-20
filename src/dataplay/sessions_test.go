package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestIsUserLoggedIn(t *testing.T) {
	handleClearSession(t)
	handleLoggedIn(t)
	handleClearSession(t) // log out
	handleLoggedOut(t)
}

func TestClearSession(t *testing.T) {
	handleClearSession(t)
}

func TestRandString(t *testing.T) {
	result := randString(5)

	Convey("When Random String length is 5", t, func() {
		So(len(result), ShouldEqual, 5)
		So(len(result), ShouldNotEqual, 6)
		So(result, ShouldNotContainSubstring, "!")
	})
}

func handleLoggedIn(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=glyn@dataplay.com&password=123456"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()
	HandleLogin(response, request)

	NewSessionID := randString(64)
	c, _ := GetRedisConnection()
	defer c.Close()
	c.Cmd("SET", NewSessionID, 181)

	NewCookie := &http.Cookie{
		Name:    "DPSession",
		Value:   NewSessionID,
		Path:    "/",
		Expires: time.Now().AddDate(1, 0, 0),
	}

	http.SetCookie(response, NewCookie)
	request.Header.Set("Cookie", NewCookie.String())
	status := IsUserLoggedIn(response, request)

	Convey("When User is Logged In", t, func() {
		So(status, ShouldEqual, true)
	})
}

func handleLoggedOut(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=glyn@dataplay.com&password=123456"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()
	status := IsUserLoggedIn(response, request)

	Convey("When User is Logged Out", t, func() {
		So(status, ShouldEqual, false)
	})
}

func handleClearSession(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	NewSessionID := randString(64)
	c, _ := GetRedisConnection()
	defer c.Close()
	c.Cmd("SET", NewSessionID, 1)

	NewCookie := &http.Cookie{
		Name:    "DPSession",
		Value:   NewSessionID,
		Path:    "/",
		Expires: time.Now().AddDate(1, 0, 0),
	}
	http.SetCookie(response, NewCookie)
	request.Header.Set("Cookie", NewCookie.String())
	e_return := ClearSession(response, request)

	Convey("When session cookie is cleared", t, func() {
		So(e_return, ShouldBeNil)
	})
}
