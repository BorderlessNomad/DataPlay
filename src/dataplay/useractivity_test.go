package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddActivity(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{
		"uid":  "1",
		"type": "c",
	}

	Convey("Should add activity", t, func() {
		result := AddActivityHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})

	params["type"] = "X"

	Convey("Should fail to add activity", t, func() {
		result := AddActivityHttp(res, req, params)
		So(result, ShouldEqual, "Unknown activity type")
	})

}
