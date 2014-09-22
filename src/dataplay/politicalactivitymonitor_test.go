package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPoliticalActivityHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{
		"type": "d",
	}

	Convey("Should return departments PoliticalActivity", t, func() {
		result := GetPoliticalActivityHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})

	params["type"] = "e"

	Convey("Should return events PoliticalActivity", t, func() {
		result := GetPoliticalActivityHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})

	params["type"] = "r"

	Convey("Should return regions PoliticalActivity", t, func() {
		result := GetPoliticalActivityHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})

	params["type"] = "p"

	Convey("Should return popular PoliticalActivity", t, func() {
		result := GetPoliticalActivityHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})

	// WriteCass()
}
