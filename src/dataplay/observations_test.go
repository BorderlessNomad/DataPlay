package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

// func TestAddObservationHttp(t *testing.T) {
// 	req, _ := http.NewRequest("POST", "/", nil)
// 	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
// 	res := httptest.NewRecorder()
// 	Convey("Should add observation", t, func() {
// 		obs := ObservationComment{}
// 		obs.DiscoveryId = "663"
// 		obs.Comment = "test comment"
// 		obs.X = "0"
// 		obs.Y = "0"
// 		result := AddObservationHttp(res, req, obs)
// 		So(result, ShouldEqual, "")
// 	})
// }

// func TestGetObservationsQ(t *testing.T) {
// 	Convey("Should get observations", t, func() {
// 		params := map[string]string{}
// 		params["did"] = "0"
// 		result := GetObservationsQ(params)
// 		So(result, ShouldEqual, "Observations could not be retrieved")
// 	})
// }

func TestGetObservationsHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	Convey("Should get observations", t, func() {
		params := map[string]string{}
		params["did"] = "592"
		result := GetObservationsHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})
}

func TestGetRecentObservationsHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()

	Convey("Should get observations", t, func() {
		result := GetRecentObservationsHttp(res, req)
		So(result, ShouldNotBeNil)
	})
}
