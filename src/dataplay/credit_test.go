package main

//// FILE NEEDS RENAMING TO SOMETHING > "M..." TO RUN AFTER "M"ain file

// import (
// 	. "github.com/smartystreets/goconvey/convey"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// )

// func TestRankCredits(t *testing.T) {
// 	Convey("Should return ranking", t, func() {
// 		result := RankCredits(23, 15)
// 		So(result, ShouldEqual, 0.44717586998695963)
// 	})
// }

// func TestCreditChartHttp(t *testing.T) {
// 	req, _ := http.NewRequest("GET", "/", nil)
// 	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
// 	res := httptest.NewRecorder()
// 	Convey("Should credit chart", t, func() {
// 		params := map[string]string{}
// 		params["rcid"] = "116144"
// 		params["credflag"] = "false"
// 		result := CreditChartHttp(res, req, params)
// 		So(result, ShouldEqual, "")
// 	})
// }

// func TestCreditObservationHttp(t *testing.T) {
// 	req, _ := http.NewRequest("GET", "/", nil)
// 	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
// 	res := httptest.NewRecorder()
// 	params := map[string]string{}
// 	params["oid"] = "702"
// 	params["credflag"] = "true"
// 	Convey("Should credit observation", t, func() {
// 		result := CreditObservationHttp(res, req, params)
// 		So(result, ShouldNotBeNil)
// 	})
// 	Convey("Should discredit observation", t, func() {
// 		params["credflag"] = "false"
// 		result := CreditObservationHttp(res, req, params)
// 		So(result, ShouldEqual, "Observation discredited")
// 	})
// }
