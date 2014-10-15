package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserTableHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{
		"order":  "uid",
		"offset": "3",
		"count":  "3",
	}
	Convey("Should return users for admin", t, func() {
		result := GetUserTableHttp(res, req, params)
		So(result, ShouldEqual, "")
	})
}
