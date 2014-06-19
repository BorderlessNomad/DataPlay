package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"strings"
)

func TestCheckAuth(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=glyn@dataplay.com&password=123456"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()
	NewSessionID := randString(64)
	c, _ := GetRedisConnection()
	defer c.Close()
	c.Cmd("SET", NewSessionID, 1)

	NewCookie := &http.Cookie{
		Name:    "DPSession",
		Value:   NewSessionID,
		Path:    "/",
		Expires: time.Now().AddDate(1, 0, 0),
	}
	http.SetCookie(response, NewCookie)

	request.Header.Set("Cookie", NewCookie.String())

	HandleLogin(response, request)
	prams := map[string]string{
		"id": "181",
	}

	CheckAuth(response,request,prams)

	Convey("When user is/isn't logged in", t, func() {
		So(response.Code, ShouldNotBeNil)
	})
}


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

	prams["s"] = "nhs"
	result := SearchForData(response,request,prams)
	Convey("When search parameter is provided", t, func() {
		So(result, ShouldNotBeBlank)
	})

	prams["s"] = "¬"
	result = SearchForData(response,request,prams)
	Convey("When search parameter is provided but deep search required", t, func() {
		So(result, ShouldNotBeBlank)
	})
}

func TestProcessSearchResults(t *testing.T) {

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

	prams["id"] = "gold"
	result := GetEntry(response,request,prams)

	Convey("When no ID parameter is provided", t, func() {
		So(result, ShouldNotBeBlank)
	})
}

func TestScanRow(t *testing.T) {
	cols := []string{"cmxval", "bval", "ival", "i64val","fval","sval", "btval"}
	var cmxval complex128 = -1 + 3i //triggers "unexpected type"
	var bval bool = true
	var ival int = 1
	var i64val int64 = 1
	var fval float64 = 1.0
	var sval string = "a"
	btval := []byte("a")
	tval := time.Now()

	vals := []interface{}{cmxval, bval, ival, i64val, fval, sval, btval, tval}
	record := ScanRow(vals, cols)
	Convey("Scanrow", t, func() {
		So(record, ShouldNotEqual, 0)
	})

}

func TestDumpTable(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"top": "",
		"bot": "",
	}

	DumpTable(response,request,prams)

	Convey("When no ID parameter is provided", t, func() {
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	prams["top"] = "5"
	prams["bot"] = "10"
	DumpTable(response,request,prams)

	Convey("When no ID parameter is provided", t, func() {
		So(response.Code, ShouldNotBeNil)
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

func TestDumpReducedTable(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{"id": ""}

	DumpReducedTable(response,request,prams)
	Convey("When no ID parameter is provided", t, func() {
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

	GetCSV(response,request,prams)
	Convey("When no x or y parameters are provided", t, func() {
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	prams["id"]= ""

	GetCSV(response,request,prams)
	Convey("When no ID parameter is provided", t, func() {
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})
}
