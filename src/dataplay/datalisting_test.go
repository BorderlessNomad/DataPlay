package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchForData(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"s": "",
	}

	SearchForData(response,request,prams)

	Convey("When no search parameter is provided", t, func() {
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})
}

func TestGetEntry(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"id": "",
	}

	GetEntry(response,request,prams)

	Convey("When no ID parameter is provided", t, func() {
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

}

func TestDumpTable(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"id": "",
	}

	DumpTable(response,request,prams)

	Convey("When no ID parameter is provided", t, func() {
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})
}

func TestDumpTableRange(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"id": "x",
		"x": "",
		"startx": "",
		"endx": "",
	}

	DumpTableRange(response,request,prams)
	Convey("When no x, startx or endx parameters are provided", t, func() {
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	prams["id"]= ""

	DumpTableRange(response,request,prams)
	Convey("When no x, startx or endx parameters are provided", t, func() {
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})
}

func TestDumpTableGrouped(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"id": "x",
		"x": "",
		"y": "",
	}

	DumpTableGrouped(response,request,prams)
	Convey("When no ID, x or y parameters are provided", t, func() {
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

}

func TestDumpTablePrediction(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"id": "x",
		"x": "",
		"y": "",
	}

	DumpTablePrediction(response,request,prams)
	Convey("When no ID, x or y parameters are provided", t, func() {
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})
}

func TestGetCSV(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"id": "x",
		"x": "",
		"y": "",
	}

	DumpTableRange(response,request,prams)
	Convey("When no x or y parameters are provided", t, func() {
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	prams["id"]= ""

	DumpTableRange(response,request,prams)
	Convey("When no ID parameter is provided", t, func() {
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})
}
