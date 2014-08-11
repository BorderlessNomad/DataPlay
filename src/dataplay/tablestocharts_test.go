package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRelatedCharts(t *testing.T) {
	m := map[string]string{
		"user":      "1",
		"offset":    "0",
		"count":     "6",
		"tablename": "gold",
	}
	Convey("Should return chartlist", t, func() {
		result := GetRelatedChartsQ(m)
		So(result, ShouldNotBeNil)
	})
}

// func TestGetRelatedChartsHttp(t *testing.T) {
// 	req, _ := http.NewRequest("GET", "/", nil)
// 	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
// 	res := httptest.NewRecorder()
// 	params := map[string]string{
// 		"user":      "1",
// 		"offset":    "0",
// 		"count":     "10",
// 		"tablename": "fe5e88f1c898b2ea870c928a3b94d5a1bf219d057e68010a018a73634dd",
// 	}
// 	Convey("Should return chartlist", t, func() {
// 		result := GetRelatedChartsHttp(res, req, params)
// 		So(result, ShouldEqual, "")
// 	})
// }

func TestGetRelatedChartsHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{
		"user":      "1",
		"offset":    "0",
		"count":     "6",
		"tablename": "gold",
	}
	Convey("Should return chartlist", t, func() {
		result := GetRelatedChartsHttp(res, req, params)
		So(result, ShouldEqual, "")
	})
}

func TestGetChart(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{
		"user":      "1",
		"tablename": "gold",
		"type":      "line",
		"x":         "price",
		"y":         "date",
	}
	Convey("Should return chartlist", t, func() {
		result := GetChartHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})
}
