package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

// func TestRankValidations(t *testing.T) {
// 	Convey("Should return ranking", t, func() {
// 		result := RankValidations(23, 15)
// 		So(result, ShouldEqual, 0.44717586998695963)
// 	})
// }

func TestValidateChartHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	Convey("Should validate chart", t, func() {
		params := map[string]string{}
		params["rcid"] = "114789"
		params["valflag"] = "false"
		result := ValidateChartHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})
}

func TestValidateObservationHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{}
	params["oid"] = "702"
	params["valflag"] = "true"
	Convey("Should validate observation", t, func() {
		result := ValidateObservationHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})
	// 	Convey("Should invalidate observation", t, func() {
	// 		params["valflag"] = "false"
	// 		result := ValidateObservationHttp(res, req, params)
	// 		So(result, ShouldEqual, "Observation invalidated")
	// 	})
}
