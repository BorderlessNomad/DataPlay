package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetLastVisited(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", strings.NewReader("user=1&guid=gold"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	c := &http.Cookie{
		Name:  "DPSession",
		Value: "1",
	}

	response := httptest.NewRecorder()

	Convey("On HTTP Request", t, func() {
		request.Header.Set("Cookie", c.String())
		result := GetLastVisited(response, request)

		Convey("Should Logout", func() {
			So(result, ShouldEqual, "[]")
		})
	})

}

func TestHasTableGotLocatonData(t *testing.T) {
	result := HasTableGotLocationData("tweets")

	Convey("Should find Lattitude and Longitude columns in dataset", t, func() {
		So(result, ShouldEqual, "true")
	})

	result = HasTableGotLocationData("houseprices")

	Convey("Should not find Lattitude and Longitude columns in dataset", t, func() {
		So(result, ShouldEqual, "false")
	})
}

func TestContainsTableCol(t *testing.T) {
	Cols := []ColType{{"X", "0"}, {"Y", "0"}}

	result := ContainsTableCol(Cols, "y")

	Convey("Should find Column name", t, func() {
		So(result, ShouldBeTrue)
	})

	result = ContainsTableCol(Cols, "z")

	Convey("Should not find Column name", t, func() {
		So(result, ShouldBeFalse)
	})
}

func TestTrackVisited(t *testing.T) {
	Convey("Track visited", t, func() {
		TrackVisited("", "")
	})
}
