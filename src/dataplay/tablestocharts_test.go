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

func TestGetCorrelatedChart(t *testing.T) {
	m := map[string]string{
		"id": "113799",
	}
	Convey("Should return chartlist", t, func() {
		result := GetCorrelatedChartQ(m)
		So(result, ShouldNotBeNil)
	})
}

func TestGetRelatedChartsHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{
		"user":      "1",
		"offset":    "0",
		"count":     "60",
		"tablename": "gdp",
	}
	Convey("Should return chartlist", t, func() {
		result := GetRelatedChartsHttp(res, req, params)
		So(result, ShouldNotBeNil)
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
	Convey("Should return xy chartlist", t, func() {
		result := GetChartHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})

	params["tablename"] = "gdp"
	params["type"] = "bubble"
	params["x"] = "year"
	params["y"] = "gdp"
	params["z"] = "change"

	Convey("Should return xyz chartlist", t, func() {
		result := GetChartHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})
}
