package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

	Convey("When user is/isn't logged in", t, func() {
		CheckAuth(response, request, prams)
		So(response.Code, ShouldNotBeNil)
	})
}

func TestSearchForData(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"s": "",
	}
	result := SearchForData(response, request, prams)

	Convey("When no search parameter is provided", t, func() {
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When search parameter is 'nhs'", t, func() {
		prams["s"] = "nhs"
		result = SearchForData(response, request, prams)
		So(result, ShouldNotBeBlank)
	})

	Convey("When search parameter is 'hs'", t, func() {
		prams["s"] = "hs"
		result = SearchForData(response, request, prams)
		So(result, ShouldNotBeBlank)
	})

	Convey("When search parameter is 'n h s'", t, func() {
		prams["s"] = "n h s"
		result = SearchForData(response, request, prams)
		So(result, ShouldNotBeBlank)
	})

	Convey("When search parameter is 'freakshine'", t, func() {
		prams["s"] = "freakshine"
		result = SearchForData(response, request, prams)
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

	Convey("When no ID parameter is provided", t, func() {
		GetEntry(response, request, prams)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When no ID parameter is provided", t, func() {
		prams["id"] = "gold"
		result := GetEntry(response, request, prams)
		So(result, ShouldNotBeBlank)
	})
}

func TestScanRow(t *testing.T) {
	cols := []string{"cmxval", "bval", "ival", "i64val", "fval", "sval", "btval"}
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
	prams := map[string]string{"id": ""}

	Convey("When no ID parameter is provided", t, func() {
		DumpTable(response, request, prams)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When limits are not provided", t, func() {
		prams["id"] = "gdp"
		DumpTable(response, request, prams)
		So(response.Code, ShouldNotBeNil)
	})

	Convey("When limits are provided", t, func() {
		prams["offset"] = "5"
		prams["count"] = "10"
		DumpTable(response, request, prams)
		So(response.Code, ShouldNotBeNil)
	})
}

func TestDumpTableRange(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{"id": ""}

	Convey("When no id parameter is provided", t, func() {
		DumpTableRange(response, request, prams)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When no x, startx or endx parameters are provided", t, func() {
		prams["id"] = "gdp"
		DumpTableRange(response, request, prams)
	})

	// Convey("When no x, startx or endx parameters are provided", t, func() {
	// 	prams["x"] = "1"
	// 	prams["startx"] = "2"
	// 	prams["endx"] = "3"
	// 	DumpTableRange(response, request, prams)
	// })
}

func TestDumpTableGrouped(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{"id": ""}

	Convey("When no ID, x or y parameters are provided", t, func() {
		DumpTableGrouped(response, request, prams)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When valid parameters are provided", t, func() {
		prams["id"] = "gdp"
		prams["x"] = "change"
		prams["y"] = "gdpindex"
		DumpTableGrouped(response, request, prams)
	})
}

func TestDumpTablePrediction(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{"id": ""}

	Convey("When no ID, x or y parameters are provided", t, func() {
		DumpTablePrediction(response, request, prams)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When valid parameters are provided", t, func() {
		prams["id"] = "gdp"
		prams["x"] = "change"
		prams["y"] = "gdpindex"
		DumpTablePrediction(response, request, prams)
	})
}

func TestDumpReducedTable(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{"id": ""}

	Convey("When no ID parameter is provided", t, func() {
		DumpReducedTable(response, request, prams)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When valid parameters are provided", t, func() {
		prams["id"] = "gdp"
		DumpReducedTable(response, request, prams)
	})
}
