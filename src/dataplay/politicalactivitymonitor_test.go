package main

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetPoliticalActivityHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{
		"type": "d",
	}

	x := time.Now()
	Convey("Should return departments PoliticalActivity", t, func() {
		result := GetPoliticalActivityHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})
	y := time.Now()
	fmt.Println("ROBOCOP", y.Sub(x).Seconds())

	params["type"] = "e"

	x = time.Now()
	Convey("Should return events PoliticalActivity", t, func() {
		result := GetPoliticalActivityHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})
	y = time.Now()
	fmt.Println("ROBOCOP2", y.Sub(x).Seconds())

	// params["type"] = "r"
	// x = time.Now()
	// Convey("Should return regions PoliticalActivity", t, func() {
	// 	result := GetPoliticalActivityHttp(res, req, params)
	// 	So(result, ShouldNotBeNil)
	// })
	// y = time.Now()
	// fmt.Println("ROBOCOP3", y.Sub(x).Seconds())
	params["type"] = "p"

	Convey("Should return popular PoliticalActivity", t, func() {
		result := GetPoliticalActivityHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})

	// WriteCass()
}
