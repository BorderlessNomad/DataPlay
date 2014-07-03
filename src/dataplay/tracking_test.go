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

	result := GetLastVisited(response, request)
	Convey("Should get last visited", t, func() {
		So(result, ShouldNotBeBlank)
	})
}

func TestHasTableGotLocatonData(t *testing.T) {
	result := HasTableGotLocationData("tweets")

	Convey("Should find Lattitude and Longitude columns in dataset", t, func() {
		So(result, ShouldEqual, "true")
	})

	Convey("Should not find Lattitude and Longitude columns in dataset", t, func() {
		result = HasTableGotLocationData("houseprices")
		So(result, ShouldEqual, "false")
	})
}

func TestContainsTableCol(t *testing.T) {
	Cols := []ColType{{"X", "0"}, {"Y", "0"}}
	result := ContainsTableCol(Cols, "y")

	Convey("Should find Column name", t, func() {
		So(result, ShouldBeTrue)
	})

	Convey("Should not find Column name", t, func() {
		result = ContainsTableCol(Cols, "z")
		So(result, ShouldBeFalse)
	})
}

func TestTrackVisited(t *testing.T) {
	Convey("Track visited", t, func() {
		TrackVisited("gold", 11)
	})
}
