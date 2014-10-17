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
		So(result, ShouldNotBeNil)
	})
}

func TestGetObservationsTableHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{
		"order":   "username",
		"offset":  "0",
		"count":   "100",
		"flagged": "",
	}
	Convey("Should return users for admin", t, func() {
		result := GetObservationsTableHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})
}

// func TestDeleteObservationHttp(t *testing.T) {
// 	req, _ := http.NewRequest("GET", "/", nil)
// 	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
// 	res := httptest.NewRecorder()
// 	params := map[string]string{
// 		"id": "754",
// 	}
// 	Convey("Should delete obs", t, func() {
// 		result := DeleteObservationHttp(res, req, params)
// 		So(result, ShouldEqual, "")
// 	})
// }
