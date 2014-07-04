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

func TestSearchForDataHttp(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	params := map[string]string{
		"s": "",
	}
	result := SearchForDataHttp(response, request, params)

	Convey("When no search parameter is provided", t, func() {
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When search parameter is 'nhs'", t, func() {
		params["s"] = "nhs"
		result = SearchForDataHttp(response, request, params)
		So(result, ShouldNotBeBlank)
	})

	Convey("When search parameter is 'hs'", t, func() {
		params["s"] = "hs"
		result = SearchForDataHttp(response, request, params)
		So(result, ShouldNotBeBlank)
	})

	Convey("When search parameter is 'n h s'", t, func() {
		params["s"] = "n h s"
		result = SearchForDataHttp(response, request, params)
		So(result, ShouldNotBeBlank)
	})

	Convey("When search parameter is 'freakshine'", t, func() {
		params["s"] = "freakshine"
		result = SearchForDataHttp(response, request, params)
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

func TestDumpTableHttp(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	params := map[string]string{"id": ""}

	Convey("When no ID parameter is provided", t, func() {
		DumpTableHttp(response, request, params)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When table name is incorrect ", t, func() {
		params["id"] = "qwerty1"
		DumpTableHttp(response, request, params)
		So(response.Code, ShouldNotBeNil)
	})

	Convey("When limits are not provided", t, func() {
		params["id"] = "gdp"
		DumpTableHttp(response, request, params)
		So(response.Code, ShouldNotBeNil)
	})

	Convey("When incorrect limits are provided", t, func() {
		params["offset"] = "-3000"
		params["count"] = "10.5"
		DumpTableHttp(response, request, params)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When correct limits are provided", t, func() {
		params["offset"] = "5"
		params["count"] = "10"
		DumpTableHttp(response, request, params)
		So(response.Code, ShouldNotBeNil)
	})
}

func TestDumpTableRangeHttp(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	params := map[string]string{"id": ""}

	Convey("When no id parameter is provided", t, func() {
		DumpTableRangeHttp(response, request, params)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When no x, startx or endx parameters are provided", t, func() {
		params["id"] = "gdp"
		DumpTableRangeHttp(response, request, params)
	})

	Convey("When bad range parameters are provided", t, func() {
		params["x"] = "year"
		params["startx"] = "-400"
		params["endx"] = "700000000.23"
		DumpTableRangeHttp(response, request, params)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When good parameters are provided", t, func() {
		params["x"] = "year"
		params["startx"] = "1970"
		params["endx"] = "2000"
		DumpTableRangeHttp(response, request, params)
	})
}

func TestDumpTableGroupedHttp(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	params := map[string]string{"id": ""}

	Convey("When no ID, x or y parameters are provided", t, func() {
		DumpTableGroupedHttp(response, request, params)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When invalid X Col parameter is provided", t, func() {
		params["id"] = "gdp"
		params["x"] = "qwerty1"
		params["y"] = "gdpindex"
		DumpTableGroupedHttp(response, request, params)
	})

	Convey("When invalid Y Col parameter is provided", t, func() {
		params["id"] = "gdp"
		params["x"] = "change"
		params["y"] = "qwerty1"
		DumpTableGroupedHttp(response, request, params)
	})

	Convey("When Y Col parameter is a date", t, func() {
		params["id"] = "gold"
		params["x"] = "price"
		params["y"] = "date"
		DumpTableGroupedHttp(response, request, params)
	})

	Convey("When valid parameters are provided", t, func() {
		params["id"] = "gdp"
		params["x"] = "change"
		params["y"] = "gdpindex"
		DumpTableGroupedHttp(response, request, params)
	})
}

func TestDumpTablePredictionHttp(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	params := map[string]string{"id": ""}

	Convey("When no ID, x or y parameters are provided", t, func() {
		DumpTablePredictionHttp(response, request, params)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When invalid X Col parameter is provided", t, func() {
		params["id"] = "gdp"
		params["x"] = "qwerty1"
		params["y"] = "gdpindex"
		DumpTablePredictionHttp(response, request, params)
	})

	Convey("When invalid Y Col parameter is provided", t, func() {
		params["id"] = "gdp"
		params["x"] = "change"
		params["y"] = "qwerty1"
		DumpTablePredictionHttp(response, request, params)
	})

	Convey("When valid parameters are provided which point to incompatible values", t, func() {
		params["id"] = "hips"
		params["x"] = "hospital"
		params["y"] = "90p"
		DumpTablePredictionHttp(response, request, params)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When valid parameters are provided", t, func() {
		params["id"] = "gdp"
		params["x"] = "change"
		params["y"] = "gdpindex"
		DumpTablePredictionHttp(response, request, params)
	})
}

func TestDumpReducedTableHttp(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	params := map[string]string{"id": ""}

	Convey("When no ID parameter is provided", t, func() {
		DumpReducedTableHttp(response, request, params)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When invalid table is provided", t, func() {
		params["id"] = "qwerty1"
		DumpReducedTableHttp(response, request, params)
	})

	Convey("When valid table is provided without parameters", t, func() {
		params["id"] = "gdp"
		DumpReducedTableHttp(response, request, params)
	})

	Convey("When invalid percent parameter is provided", t, func() {
		params["id"] = "gdp"
		params["percent"] = "-101.1"
		DumpReducedTableHttp(response, request, params)
	})

	Convey("When invalid min parameter is provided", t, func() {
		params["id"] = "gdp"
		params["percent"] = "-10"
		params["min"] = "a"
		DumpReducedTableHttp(response, request, params)
	})

	Convey("When parameter Y is a varchar or date", t, func() {
		params["id"] = "gold"
		params["x"] = "price"
		params["y"] = "date"
		params["percent"] = "10"
		params["min"] = "1"
		DumpReducedTableHttp(response, request, params)
	})

	Convey("When valid table and parameters are provided", t, func() {
		params["id"] = "gdp"
		params["x"] = "change"
		params["y"] = "gdpindex"
		params["percent"] = "10"
		params["min"] = "1"
		DumpReducedTableHttp(response, request, params)
	})
}

func TestConvertToFloat(t *testing.T) {
	typ := make([]interface{}, 13)
	typ[0] = float64(1.0)
	typ[1] = float32(1.0)
	typ[2] = int64(1)
	typ[3] = int32(1)
	typ[4] = int16(1)
	typ[5] = int8(1)
	typ[6] = uint64(1)
	typ[7] = uint32(1)
	typ[8] = uint16(1)
	typ[9] = uint8(1)
	typ[10] = int(1)
	typ[11] = uint(1)
	typ[12] = string("1")

	var boolean bool = true
	var res, chk float64

	Convey("Should be returned as float values", t, func() {
		for i := 0; i < 12; i++ {
			res, _ = ConvertToFloat(typ[i])
			So(res, ShouldHaveSameTypeAs, chk)
		}
		_, err := ConvertToFloat(boolean)
		So(err.Error(), ShouldEqual, "getFloat: unknown value is of incompatible type")
	})
}
