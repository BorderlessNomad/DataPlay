package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCharts(t *testing.T) {
	handleClearSession(t)
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=glyn@dataplay.com&password=123456"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()
	u := UserForm{}
	u.Username = "glyn@dataplay.com"
	u.Password = "123456"
	HandleLogin(response, request, u)

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
	params := map[string]string{
		"id": "181",
	}
	Charts(response, request, params)
}

func TestSearchOverlay(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	SearchOverlay(response, request)
}

func TestOverlay(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	Overlay(response, request)
}

func TestOverview(t *testing.T) {
	handleClearSession(t)
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=glyn@dataplay.com&password=123456"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()
	u := UserForm{}
	u.Username = "glyn@dataplay.com"
	u.Password = "123456"
	HandleLogin(response, request, u)

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
	params := map[string]string{
		"id": "181",
	}

	Overview(response, request, params)
}

func TestSearch(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	Search(response, request)
}

func TestMapTest(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	MapTest(response, request)
}
