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
	params := map[string]string{
		"id": "181",
	}

	Convey("When user is/isn't logged in", t, func() {
		CheckAuth(response, request, params)
		So(response.Code, ShouldNotBeNil)
	})
}

func TestSearchForData(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	params := map[string]string{
		"s": "",
	}
	result := SearchForData(response, request, params)

	Convey("When no search parameter is provided", t, func() {
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When search parameter is 'nhs'", t, func() {
		params["s"] = "nhs"
		result = SearchForData(response, request, params)
		So(result, ShouldNotBeBlank)
	})

	Convey("When search parameter is 'hs'", t, func() {
		params["s"] = "hs"
		result = SearchForData(response, request, params)
		So(result, ShouldNotBeBlank)
	})

	Convey("When search parameter is 'n h s'", t, func() {
		params["s"] = "n h s"
		result = SearchForData(response, request, params)
		So(result, ShouldNotBeBlank)
	})

	Convey("When search parameter is 'freakshine'", t, func() {
		params["s"] = "freakshine"
		result = SearchForData(response, request, params)
		So(result, ShouldNotBeBlank)
	})
}

func TestGetEntry(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	params := map[string]string{
		"id": "",
	}

	Convey("When no ID parameter is provided", t, func() {
		GetEntry(response, request, params)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When ID parameter with incorrect value is provided", t, func() {
		params["id"] = "derp"
		result := GetEntry(response, request, params)
		So(result, ShouldNotBeBlank)
	})

	Convey("When ID parameter with correct value is provided", t, func() {
		params["id"] = "gold"
		result := GetEntry(response, request, params)
		So(result, ShouldNotBeBlank)
	})
}

func TestScanRow(t *testing.T) {
	cols := []string{"cmxval", "bval", "ival", "i64val", "fval", "sval", "btval", "tval"}
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
	params := map[string]string{"id": ""}

	Convey("When no ID parameter is provided", t, func() {
		DumpTable(response, request, params)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When table name is incorrect ", t, func() {
		params["id"] = "qwerty1"
		DumpTable(response, request, params)
		So(response.Code, ShouldNotBeNil)
	})

	Convey("When limits are not provided", t, func() {
		params["id"] = "gdp"
		DumpTable(response, request, params)
		So(response.Code, ShouldNotBeNil)
	})

	Convey("When incorrect limits are provided", t, func() {
		params["offset"] = "-3000"
		params["count"] = "10.5"
		DumpTable(response, request, params)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When correct limits are provided", t, func() {
		params["offset"] = "5"
		params["count"] = "10"
		DumpTable(response, request, params)
		So(response.Code, ShouldNotBeNil)
	})
}

func TestDumpTableRange(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	params := map[string]string{"id": ""}

	Convey("When no id parameter is provided", t, func() {
		DumpTableRange(response, request, params)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When no x, startx or endx parameters are provided", t, func() {
		params["id"] = "gdp"
		DumpTableRange(response, request, params)
	})

	Convey("When bad range parameters are provided", t, func() {
		params["x"] = "year"
		params["startx"] = "-400"
		params["endx"] = "700000000.23"
		DumpTableRange(response, request, params)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When good parameters are provided", t, func() {
		params["x"] = "year"
		params["startx"] = "1970"
		params["endx"] = "2000"
		DumpTableRange(response, request, params)
	})
}

func TestDumpTableGrouped(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	params := map[string]string{"id": ""}

	Convey("When no ID, x or y parameters are provided", t, func() {
		DumpTableGrouped(response, request, params)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When invalid X Col parameter is provided", t, func() {
		params["id"] = "gdp"
		params["x"] = "qwerty1"
		params["y"] = "gdpindex"
		DumpTableGrouped(response, request, params)
	})

	Convey("When invalid Y Col parameter is provided", t, func() {
		params["id"] = "gdp"
		params["x"] = "change"
		params["y"] = "qwerty1"
		DumpTableGrouped(response, request, params)
	})

	Convey("When valid parameters are provided", t, func() {
		params["id"] = "gdp"
		params["x"] = "change"
		params["y"] = "gdpindex"
		DumpTableGrouped(response, request, params)
	})
}

func TestDumpTablePrediction(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	params := map[string]string{"id": ""}

	Convey("When no ID, x or y parameters are provided", t, func() {
		DumpTablePrediction(response, request, params)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When invalid X Col parameter is provided", t, func() {
		params["id"] = "gdp"
		params["x"] = "qwerty1"
		params["y"] = "gdpindex"
		DumpTablePrediction(response, request, params)
	})

	Convey("When invalid Y Col parameter is provided", t, func() {
		params["id"] = "gdp"
		params["x"] = "change"
		params["y"] = "qwerty1"
		DumpTablePrediction(response, request, params)
	})

	Convey("When valid parameters are provided which point to incompatible values", t, func() {
		params["id"] = "hips"
		params["x"] = "hospital"
		params["y"] = "90p"
		DumpTablePrediction(response, request, params)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When valid parameters are provided", t, func() {
		params["id"] = "gdp"
		params["x"] = "change"
		params["y"] = "gdpindex"
		DumpTablePrediction(response, request, params)
	})
}

func TestDumpReducedTable(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	params := map[string]string{"id": ""}

	Convey("When no ID parameter is provided", t, func() {
		DumpReducedTable(response, request, params)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When invalid table is provided", t, func() {
		params["id"] = "qwerty1"
		DumpReducedTable(response, request, params)
	})

	Convey("When valid table is provided without parameters", t, func() {
		params["id"] = "gdp"
		DumpReducedTable(response, request, params)
	})

	Convey("When invalid percent parameter is provided", t, func() {
		params["id"] = "gdp"
		params["percent"] = "-101.1"
		DumpReducedTable(response, request, params)
	})

	Convey("When invalid min parameter is provided", t, func() {
		params["id"] = "gdp"
		params["percent"] = "-10"
		params["min"] = "a"
		DumpReducedTable(response, request, params)
	})

	Convey("When valid table and parameters are provided", t, func() {
		params["id"] = "gdp"
		params["percent"] = "10"
		params["min"] = "1"
		DumpReducedTable(response, request, params)
	})
}
