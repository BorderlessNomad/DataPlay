package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetChartHttp(t *testing.T) {
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

func TestGetCorrelatedChartHttp(t *testing.T) {
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
		result := GetCorrelatedChartHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})

	params["tablename"] = "gdp"
	params["type"] = "bubble"
	params["x"] = "year"
	params["y"] = "gdp"
	params["z"] = "change"

	Convey("Should return xyz chartlist", t, func() {
		result := GetCorrelatedChartHttp(res, req, params)
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

func TestGetCorrelatedChartsHttp(t *testing.T) {
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
		result := GetCorrelatedChartsHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})
}

func TestGetValidatedChartsHttp(t *testing.T) {
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
		result := GetValidatedChartsHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})
}

func TestGetChartQ(t *testing.T) {
	params := map[string]string{
		"user":      "1",
		"tablename": "gold",
		"type":      "line",
		"x":         "price",
		"y":         "date",
	}
	Convey("Should return xy chartlist", t, func() {
		result := GetChartQ(params)
		So(result, ShouldNotBeNil)
	})

	params["tablename"] = "gdp"
	params["type"] = "bubble"
	params["x"] = "year"
	params["y"] = "gdp"
	params["z"] = "change"

	Convey("Should return xyz chartlist", t, func() {
		result := GetChartQ(params)
		So(result, ShouldNotBeNil)
	})
}

func TestGetCorrelatedChartQ(t *testing.T) {
	params := map[string]string{
		"user":      "1",
		"tablename": "gold",
		"type":      "line",
		"x":         "price",
		"y":         "date",
	}

	Convey("Should return xy chartlist", t, func() {
		result := GetCorrelatedChartQ(params)
		So(result, ShouldNotBeNil)
	})

	params["tablename"] = "gdp"
	params["type"] = "bubble"
	params["x"] = "year"
	params["y"] = "gdp"
	params["z"] = "change"

	Convey("Should return xyz chartlist", t, func() {
		result := GetCorrelatedChartQ(params)
		So(result, ShouldNotBeNil)
	})
}

func TestGetRelatedChartsQ(t *testing.T) {
	params := map[string]string{
		"user":      "1",
		"offset":    "0",
		"count":     "60",
		"tablename": "gdp",
	}
	Convey("Should return chartlist", t, func() {
		result := GetRelatedChartsQ(params)
		So(result, ShouldNotBeNil)
	})
}

func TestGetCorrelatedChartsQ(t *testing.T) {
	params := map[string]string{
		"user":      "1",
		"offset":    "0",
		"count":     "60",
		"tablename": "gdp",
	}
	Convey("Should return chartlist", t, func() {
		result := GetCorrelatedChartsQ(params)
		So(result, ShouldNotBeNil)
	})
}

func TestGetValidatedChartsQ(t *testing.T) {
	params := map[string]string{
		"user":      "1",
		"offset":    "0",
		"count":     "60",
		"tablename": "gdp",
	}
	Convey("Should return chartlist", t, func() {
		result := GetValidatedChartsQ(params)
		So(result, ShouldNotBeNil)
	})
}
