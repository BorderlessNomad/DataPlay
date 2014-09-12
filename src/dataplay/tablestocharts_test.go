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

	params["tablename"] = "gdp"
	params["type"] = "bubble"
	params["x"] = "date"
	params["y"] = "gdp"
	params["z"] = "change"

	Convey("Should return xyz chartlist", t, func() {
		result := GetChartHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})
}

func TestGetChartCorrelatedHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{
		"cid": "114743",
		"uid": "1",
	}
	Convey("Should return Correlated chart", t, func() {
		result := GetChartCorrelatedHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})

	// params["uid"] = "3"
	// params["discovered"] = "true"

	// Convey("Should return xyz chartlist", t, func() {
	// 	result := GetChartCorrelatedHttp(res, req, params)
	// 	So(result, ShouldNotBeNil)
	// })
}

// func TestGetRelatedChartsHttp(t *testing.T) {
// 	req, _ := http.NewRequest("GET", "/", nil)
// 	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
// 	res := httptest.NewRecorder()
// 	params := map[string]string{
// 		"tablename": "gdp",
// 		"offset":    "0",
// 		"count":     "60",
// 	}
// 	Convey("Should return chartlist", t, func() {
// 		result := GetRelatedChartsHttp(res, req, params)
// 		So(result, ShouldNotBeNil)
// 	})
// }

// func TestGetCorrelatedChartsHttp(t *testing.T) {
// 	req, _ := http.NewRequest("GET", "/", nil)
// 	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
// 	res := httptest.NewRecorder()
// 	params := map[string]string{
// 		"tablename": "gold",
// 		"offset":    "0",
// 		"count":     "20",
// 		"search":    "true",
// 	}
// 	Convey("Should return correlated chartlist", t, func() {
// 		result := GetCorrelatedChartsHttp(res, req, params)
// 		So(result, ShouldNotBeNil)
// 	})
// }

// func TestGetDiscoveredChartsHttp(t *testing.T) {
// 	req, _ := http.NewRequest("GET", "/", nil)
// 	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
// 	res := httptest.NewRecorder()
// 	params := map[string]string{
// 		"tablename":  "gold",
// 		"offset":     "0",
// 		"count":      "5",
// 		"correlated": "true",
// 	}
// 	Convey("Should return chartlist", t, func() {
// 		result := GetDiscoveredChartsHttp(res, req, params)
// 		So(result, ShouldNotBeNil)
// 	})
// }

// func TestGetChartQ(t *testing.T) {
// 	params := map[string]string{
// 		"uid":        "1",
// 		"tablename":  "gold",
// 		"tablenum":   "5",
// 		"type":       "line",
// 		"x":          "price",
// 		"y":          "date",
// 		"discovered": "false",
// 	}
// 	Convey("Should add chart discovery", t, func() {
// 		result := GetChartQ(params)
// 		So(result, ShouldNotBeNil)
// 	})

// 	params["discovered"] = "true"

// 	Convey("Should return chart", t, func() {
// 		result := GetChartQ(params)
// 		So(result, ShouldNotBeNil)
// 	})
// }

// func TestGetChartCorrelatedQ(t *testing.T) {
// 	params := map[string]string{
// 		"id":         "114264",
// 		"uid":        "6",
// 		"discovered": "false",
// 	}

// 	Convey("Should discover and return correlated chart", t, func() {
// 		result := GetChartCorrelatedQ(params)
// 		So(result, ShouldNotBeNil)
// 	})

// 	params["uid"] = "3"
// 	params["discovered"] = "true"

// 	Convey("Should return correlated chart", t, func() {
// 		result := GetChartCorrelatedQ(params)
// 		So(result, ShouldNotBeNil)
// 	})
// }

// func TestGetRelatedChartsQ(t *testing.T) {
// 	params := map[string]string{
// 		"user":      "1",
// 		"offset":    "0",
// 		"count":     "60",
// 		"tablename": "gdp",
// 	}
// 	Convey("Should return chartlist", t, func() {
// 		result := GetRelatedChartsQ(params)
// 		So(result, ShouldNotBeNil)
// 	})
// }

// func TestGetCorrelatedChartsQ(t *testing.T) {
// 	params := map[string]string{
// 		"tablename":   "gdp",
// 		"offset":      "0",
// 		"count":       "60",
// 		"searchdepth": "10",
// 	}
// 	Convey("Should return chartlist", t, func() {
// 		result := GetCorrelatedChartsQ(params)
// 		So(result, ShouldNotBeNil)
// 	})
// }

// func TestGetDiscoveredChartsQ(t *testing.T) {
// 	params := map[string]string{
// 		"tablename":  "gold",
// 		"offset":     "0",
// 		"count":      "5",
// 		"correlated": "true",
// 	}
// 	Convey("Should return chartlist", t, func() {
// 		result := GetDiscoveredChartsQ(params)
// 		So(result, ShouldNotBeNil)
// 	})
// }
