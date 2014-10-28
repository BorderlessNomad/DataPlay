package main

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetChartHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{
		"tablename": "gdp",
		"tablenum":  "3",
		"type":      "line",
		"x":         "date",
		"y":         "gdp",
	}
	Convey("Should return xy chartlist", t, func() {
		result := GetChartHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})
}

// params["tablename"] = "gdp"
// params["type"] = "bubble"
// params["x"] = "date"
// params["y"] = "gdp"
// params["z"] = "change"

// Convey("Should return xyz chartlist", t, func() {
// 	result := GetChartHttp(res, req, params)
// 	So(result, ShouldEqual, "")
// })
// }

func TestGetChartCorrelatedHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{
		"cid": "115925",
	}
	Convey("Should return Correlated chart", t, func() {
		result := GetChartCorrelatedHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})
}

func TestGetRelatedChartsHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{
		"tablename": "gdp",
		"offset":    "0",
		"count":     "10",
	}
	Convey("Should return chartlist", t, func() {
		result := GetRelatedChartsHttp(res, req, params)
		So(result, ShouldEqual, "")
	})
}

func TestGetCorrelatedChartsHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{
		"tablename": "life",
		"offset":    "0",
		"count":     "10",
		"search":    "true",
	}
	Convey("Should return correlated chartlist", t, func() {
		x := time.Now()
		result := GetCorrelatedChartsHttp(res, req, params)
		y := time.Now()
		fmt.Println("CORRELATED_CHARTS_TIME_TAKEN", y.Sub(x).Seconds())
		So(result, ShouldEqual, "")
	})
}

func TestGetDiscoveredChartsHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{
		"tablename": "gold",
		// "offset":     "0",
		// "count":      "5",
		"correlated": "true",
	}
	Convey("Should return chartlist", t, func() {
		result := GetDiscoveredChartsHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})
}

func TestGetTopRatedChartsHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()

	Convey("Should return chartlist", t, func() {
		result := GetTopRatedChartsHttp(res, req)
		So(result, ShouldNotBeNil)
	})
}

func TestGetAwaitingCreditHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()

	Convey("Should return chartlist", t, func() {
		result := GetAwaitingCreditHttp(res, req)
		So(result, ShouldNotBeNil)
	})
}
