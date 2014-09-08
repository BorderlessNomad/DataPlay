package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOverviewHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{
		"type": "d",
	}

	Convey("Should return departments overview", t, func() {
		result := GetOverviewHttp(res, req, params)
		So(result, ShouldEqual, "")
	})

	params["type"] = "e"

	Convey("Should return events overview", t, func() {
		result := GetOverviewHttp(res, req, params)
		So(result, ShouldEqual, "")
	})

	params["type"] = "r"

	Convey("Should return regions overview", t, func() {
		result := GetOverviewHttp(res, req, params)
		So(result, ShouldEqual, "")
	})
}
