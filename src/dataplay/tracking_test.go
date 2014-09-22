package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetLastVisited(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	params := map[string]string{
		"s": "",
	}

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
	request.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")

	result := GetLastVisitedHttp(response, request, params)
	Convey("Should get last visited", t, func() {
		So(result, ShouldNotBeBlank)
	})

}

func TestGetLastVisitedQ(t *testing.T) {
	m := make(map[string]string)

	Convey("Should return empty when no user", t, func() {
		result := GetLastVisitedQ(m)
		So(result, ShouldBeEmpty)
	})

	Convey("Should not return empty when there is a user", t, func() {
		m["user"] = "11"
		result := GetLastVisitedQ(m)
		So(result, ShouldNotBeEmpty)
	})
}

func TestTrackVisitedHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	var g = []byte(`[{"Guid": "gold"},{"x": "x"}]`)
	var in = []byte(`[
		{"Info": "test"}
	]`)
	v := VisitedForm{
		Guid: g,
		Info: in,
	}
	Convey("Track visited", t, func() {
		TrackVisitedHttp(res, req, v)
	})
}
