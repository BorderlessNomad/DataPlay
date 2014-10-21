package main

import (
	// . "github.com/smartystreets/goconvey/convey"
	// "net/http"
	// "net/http/httptest"
	// "strings"
	"testing"
	// "time"
	// "encoding/json"
	// "fmt"
)

func TestSearchForData(t *testing.T) {
	// request, _ := http.NewRequest("GET", "/", nil)
	// request.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	// response := httptest.NewRecorder()
	// params := map[string]string{
	// 	"keyword": "",
	// }
	// result := SearchForDataHttp(response, request, params)

	// Convey("When no search parameter is provided", t, func() {
	// 	So(response.Code, ShouldEqual, http.StatusBadRequest)
	// })
	// Convey("When search parameter is 'nhs'", t, func() {
	// 	params["keyword"] = "nhs"
	// 	result = SearchForDataHttp(response, request, params)
	// 	So(result, ShouldEqual, "")
	// })
	// Convey("When search parameter is 'hs'", t, func() {
	// 	params["keyword"] = "hs"
	// 	result = SearchForDataHttp(response, request, params)
	// 	So(result, ShouldEqual, "")
	// })
	// Convey("When search parameter is 'n h s'", t, func() {
	// 	params["keyword"] = "n h s"
	// 	result = SearchForDataHttp(response, request, params)
	// 	So(result, ShouldEqual, "")
	// })
	// Convey("When search parameter is 'nhs'", t, func() {
	// 	params["keyword"] = "nhs"
	// 	result, _ := SearchForData(0, "nhs", nil)
	// 	r, _ := json.Marshal(result)

	// 	So(result, ShouldEqual, "")
	// })
}

// 	//////////////////Q TESTS///////////////////

// 	Convey("When no search parameter is provided", t, func() {
// 		params["user"] = ""
// 		result = SearchForDataQ(params)
// 		So(result, ShouldBeBlank)
// 	})
// 	Convey("When bad user parameter is provided", t, func() {
// 		params["user"] = "q23x467bdf82123ff2344"
// 		result = SearchForDataQ(params)
// 		So(result, ShouldBeBlank)
// 	})
// 	Convey("When bad data parameter is provided", t, func() {
// 		params["user"] = "-98"
// 		params["keyword"] = "derpaderp"
// 		result = SearchForDataQ(params)
// 		So(result, ShouldEqual, "[]")
// 	})
// 	Convey("When search parameter is 'nhs'", t, func() {
// 		params["user"] = "1"
// 		params["keyword"] = "nhs"
// 		result = SearchForDataQ(params)
// 		So(result, ShouldNotBeBlank)
// 	})
// }

// func TestGetEntry(t *testing.T) {
// 	request, _ := http.NewRequest("GET", "/", nil)
// 	response := httptest.NewRecorder()
// 	params := map[string]string{
// 		"id": "",
// 	}

// 	Convey("When no ID parameter is provided", t, func() {
// 		GetEntry(response, request, params)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 	})

// 	Convey("When ID parameter with incorrect value is provided", t, func() {
// 		params["id"] = "derp"
// 		result := GetEntry(response, request, params)
// 		So(result, ShouldNotBeBlank)
// 	})

// 	Convey("When ID parameter with correct value is provided", t, func() {
// 		params["id"] = "gold"
// 		result := GetEntry(response, request, params)
// 		So(result, ShouldNotBeBlank)
// 	})
// }

// func TestScanRow(t *testing.T) {
// 	cols := []string{"cmxval", "bval", "ival", "i64val", "fval", "sval", "btval", "tval"}
// 	var cmxval complex128 = -1 + 3i //triggers "unexpected type"
// 	var bval bool = true
// 	var ival int = 1
// 	var i64val int64 = 1
// 	var fval float64 = 1.0
// 	var sval string = "a"
// 	btval := []byte("a")
// 	tval := time.Now()
// 	vals := []interface{}{cmxval, bval, ival, i64val, fval, sval, btval, tval}
// 	record := ScanRow(vals, cols)

// 	Convey("Scanrow", t, func() {
// 		So(record, ShouldNotEqual, 0)
// 	})

// }

// func TestDumpTable(t *testing.T) {
// 	request, _ := http.NewRequest("GET", "/", nil)
// 	response := httptest.NewRecorder()
// 	params := map[string]string{"id": ""}

// 	Convey("When no ID parameter is provided", t, func() {
// 		DumpTableHttp(response, request, params)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 	})
// 	Convey("When table name is incorrect ", t, func() {
// 		params["id"] = "qwerty1"
// 		DumpTableHttp(response, request, params)
// 		So(response.Code, ShouldNotBeNil)
// 	})
// 	Convey("When limits are not provided", t, func() {
// 		params["id"] = "gdp"
// 		DumpTableHttp(response, request, params)
// 		So(response.Code, ShouldNotBeNil)
// 	})
// 	Convey("When incorrect limits are provided", t, func() {
// 		params["offset"] = "-3000"
// 		params["count"] = "10.5"
// 		DumpTableHttp(response, request, params)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 	})
// 	Convey("When correct limits are provided", t, func() {
// 		params["offset"] = "5"
// 		params["count"] = "10"
// 		DumpTableHttp(response, request, params)
// 		So(response.Code, ShouldNotBeNil)
// 	})

// 	//////////////////////Q TESTS////////////////////////
// 	result := ""
// 	Convey("When no ID parameter is provided", t, func() {
// 		params["id"] = ""
// 		DumpTableQ(params)
// 		So(result, ShouldEqual, "")
// 	})
// 	Convey("When bad paramaters provided", t, func() {
// 		params["id"] = "qwerty1"
// 		params["offset"] = "-3000"
// 		params["count"] = "10.5"
// 		DumpTableQ(params)
// 		So(result, ShouldEqual, "")
// 	})
// 	Convey("When correct parameters are provided", t, func() {
// 		params["id"] = "gdp"
// 		params["offset"] = "5"
// 		params["count"] = "10"
// 		DumpTableQ(params)
// 		So(result, ShouldNotBeNil)
// 	})
// }

// func TestDumpTableRange(t *testing.T) {
// 	request, _ := http.NewRequest("GET", "/", nil)
// 	response := httptest.NewRecorder()
// 	params := map[string]string{"id": ""}

// 	Convey("When no id parameter is provided", t, func() {
// 		DumpTableRangeHttp(response, request, params)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 	})
// 	Convey("When table name is incorrect ", t, func() {
// 		params["id"] = "derpaderp"
// 		DumpTableRangeHttp(response, request, params)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 	})
// 	Convey("When no x, startx or endx parameters are provided", t, func() {
// 		params["id"] = "gdp"
// 		DumpTableRangeHttp(response, request, params)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 	})
// 	Convey("When bad range parameters are provided", t, func() {
// 		params["x"] = "year"
// 		params["startx"] = "-400"
// 		params["endx"] = "700000000.23"
// 		DumpTableRangeHttp(response, request, params)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 	})
// 	Convey("When invalid column parameters are provided", t, func() {
// 		params["x"] = "derp"
// 		DumpTableRangeHttp(response, request, params)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 	})

// 	//////////////////////Q TESTS////////////////////////

// 	result := ""
// 	Convey("When no id parameter is provided", t, func() {
// 		params["id"] = ""
// 		result = DumpTableRangeQ(params)
// 		So(result, ShouldEqual, "")
// 	})
// 	Convey("When empty parameters are provided", t, func() {
// 		params["id"] = "gdp"
// 		params["x"] = ""
// 		result = DumpTableRangeQ(params)
// 		So(result, ShouldEqual, "")
// 	})
// 	Convey("When bad parameters are provided", t, func() {
// 		params["x"] = "year"
// 		params["startx"] = "-400"
// 		params["endx"] = "700000000.23"
// 		result = DumpTableRangeQ(params)
// 		So(result, ShouldEqual, "")
// 	})
// 	// Convey("When good parameters are provided", t, func() {
// 	// 	params["x"] = "year"
// 	// 	params["startx"] = "1970"
// 	// 	params["endx"] = "2000"
// 	// 	result = DumpTableRangeQ(params)
// 	// 	So(result, ShouldNotBeNil)
// 	// })
// }

// func TestDumpTableGrouped(t *testing.T) {
// 	request, _ := http.NewRequest("GET", "/", nil)
// 	response := httptest.NewRecorder()
// 	params := map[string]string{"id": ""}

// 	Convey("When no ID, x or y parameters are provided", t, func() {
// 		DumpTableGroupedHttp(response, request, params)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 	})
// 	Convey("When invalid X Col parameter is provided", t, func() {
// 		params["id"] = "gdp"
// 		params["x"] = "qwerty1"
// 		params["y"] = "gdpindex"
// 		DumpTableGroupedHttp(response, request, params)
// 	})
// 	Convey("When invalid Y Col parameter is provided", t, func() {
// 		params["id"] = "gdp"
// 		params["x"] = "change"
// 		params["y"] = "qwerty1"
// 		DumpTableGroupedHttp(response, request, params)
// 	})
// 	Convey("When Y Col parameter is a date", t, func() {
// 		params["id"] = "gold"
// 		params["x"] = "price"
// 		params["y"] = "date"
// 		DumpTableGroupedHttp(response, request, params)
// 	})
// 	Convey("When valid parameters are provided", t, func() {
// 		params["id"] = "gdp"
// 		params["x"] = "change"
// 		params["y"] = "gdpindex"
// 		DumpTableGroupedHttp(response, request, params)
// 	})

// 	////////////////////Q TESTS/////////////////

// 	result := ""
// 	Convey("When no ID, x or y parameters are provided", t, func() {
// 		params["id"] = ""
// 		params["x"] = ""
// 		params["y"] = ""
// 		result = DumpTableGroupedQ(params)
// 		So(result, ShouldEqual, "")
// 	})
// 	Convey("When invalid parameters is provided", t, func() {
// 		params["id"] = "derp"
// 		params["x"] = "change"
// 		params["y"] = "qwerty1"
// 		result = DumpTableGroupedQ(params)
// 		So(result, ShouldEqual, "")
// 	})
// 	Convey("When valid parameters are provided", t, func() {
// 		params["id"] = "gdp"
// 		params["x"] = "change"
// 		params["y"] = "gdpindex"
// 		result = DumpTableGroupedQ(params)
// 		So(result, ShouldNotBeNil)
// 	})
// }

// func TestDumpTablePrediction(t *testing.T) {
// 	request, _ := http.NewRequest("GET", "/", nil)
// 	response := httptest.NewRecorder()
// 	params := map[string]string{"id": ""}

// 	Convey("When no ID, x or y parameters are provided", t, func() {
// 		DumpTablePredictionHttp(response, request, params)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 	})

// 	Convey("When invalid X Col parameter is provided", t, func() {
// 		params["id"] = "gdp"
// 		params["x"] = "qwerty1"
// 		params["y"] = "gdpindex"
// 		DumpTablePredictionHttp(response, request, params)
// 	})

// 	Convey("When invalid Y Col parameter is provided", t, func() {
// 		params["id"] = "gdp"
// 		params["x"] = "change"
// 		params["y"] = "qwerty1"
// 		DumpTablePredictionHttp(response, request, params)
// 	})

// 	Convey("When valid parameters are provided which point to incompatible values", t, func() {
// 		params["id"] = "hips"
// 		params["x"] = "hospital"
// 		params["y"] = "90p"
// 		DumpTablePredictionHttp(response, request, params)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 	})

// 	Convey("When valid parameters are provided", t, func() {
// 		params["id"] = "gdp"
// 		params["x"] = "change"
// 		params["y"] = "gdpindex"
// 		DumpTablePredictionHttp(response, request, params)
// 	})

// 	////////////////////Q TESTS/////////////////

// 	result := ""

// 	Convey("When no ID, x or y parameters are provided", t, func() {
// 		params["id"] = ""
// 		params["x"] = ""
// 		params["y"] = ""
// 		result = DumpTablePredictionQ(params)
// 		So(result, ShouldEqual, "")
// 	})

// 	Convey("When invalid parameters are provided", t, func() {
// 		params["id"] = "derp"
// 		params["x"] = "qwerty1"
// 		params["y"] = "derp"
// 		result = DumpTablePredictionQ(params)
// 		So(result, ShouldEqual, "")
// 	})
// 	Convey("When valid parameters are provided", t, func() {
// 		params["id"] = "gdp"
// 		params["x"] = "change"
// 		params["y"] = "gdpindex"
// 		result = DumpTablePredictionQ(params)
// 		So(result, ShouldNotBeEmpty)
// 	})
// }

// func TestDumpReducedTable(t *testing.T) {
// 	request, _ := http.NewRequest("GET", "/", nil)
// 	response := httptest.NewRecorder()
// 	params := map[string]string{"id": ""}

// 	Convey("When no ID parameter is provided", t, func() {
// 		DumpReducedTableHttp(response, request, params)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 	})
// 	Convey("When invalid table is provided", t, func() {
// 		params["id"] = "qwerty1"
// 		DumpReducedTableHttp(response, request, params)
// 	})
// 	Convey("When valid table is provided without parameters", t, func() {
// 		params["id"] = "gdp"
// 		DumpReducedTableHttp(response, request, params)
// 	})
// 	Convey("When invalid percent parameter is provided", t, func() {
// 		params["id"] = "gdp"
// 		params["percent"] = "-101.1"
// 		DumpReducedTableHttp(response, request, params)
// 	})
// 	Convey("When invalid min parameter is provided", t, func() {
// 		params["id"] = "gdp"
// 		params["percent"] = "-10"
// 		params["min"] = "a"
// 		DumpReducedTableHttp(response, request, params)
// 	})
// 	Convey("When parameter Y is a varchar or date", t, func() {
// 		params["id"] = "gold"
// 		params["x"] = "price"
// 		params["y"] = "date"
// 		params["percent"] = "10"
// 		params["min"] = "1"
// 		DumpReducedTableHttp(response, request, params)
// 	})
// 	Convey("When parameter X is invalid", t, func() {
// 		params["id"] = "gold"
// 		params["x"] = "badcolX"
// 		params["y"] = "badcolY"
// 		params["percent"] = "10"
// 		params["min"] = "1"
// 		DumpReducedTableHttp(response, request, params)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 	})
// 	Convey("When parameter Y is invalid", t, func() {
// 		params["id"] = "gold"
// 		params["x"] = "price"
// 		params["y"] = "bacdcolY"
// 		params["percent"] = "10"
// 		params["min"] = "1"
// 		DumpReducedTableHttp(response, request, params)
// 		So(response.Code, ShouldEqual, http.StatusBadRequest)
// 	})
// 	Convey("When valid table and parameters are provided", t, func() {
// 		params["id"] = "gdp"
// 		params["x"] = "change"
// 		params["y"] = "gdpindex"
// 		params["percent"] = "10"
// 		params["min"] = "1"
// 		DumpReducedTableHttp(response, request, params)
// 	})

// 	//////////////////Q TESTS//////////////////

// 	result := ""
// 	Convey("When invalid parameters are passed", t, func() {
// 		params["id"] = ""
// 		params["x"] = "derp"
// 		params["y"] = "derp"
// 		params["percent"] = "-10"
// 		params["min"] = "1"
// 		result = DumpReducedTableQ(params)
// 		So(result, ShouldEqual, "")
// 	})
// 	Convey("When valid table and parameters are provided", t, func() {
// 		params["id"] = "gdp"
// 		params["x"] = "change"
// 		params["y"] = "gdpindex"
// 		params["percent"] = "10"
// 		params["min"] = "1"
// 		result = DumpReducedTableQ(params)
// 		So(result, ShouldNotBeEmpty)
// 	})
// }

// func TestConvertToFloat(t *testing.T) {
// 	typ := make([]interface{}, 13)
// 	typ[0] = float64(1.0)
// 	typ[1] = float32(1.0)
// 	typ[2] = int64(1)
// 	typ[3] = int32(1)
// 	typ[4] = int16(1)
// 	typ[5] = int8(1)
// 	typ[6] = uint64(1)
// 	typ[7] = uint32(1)
// 	typ[8] = uint16(1)
// 	typ[9] = uint8(1)
// 	typ[10] = int(1)
// 	typ[11] = uint(1)
// 	typ[12] = string("1")

// 	var boolean bool = true
// 	var res, chk float64

// 	Convey("Should be returned as float values", t, func() {
// 		for i := 0; i < 12; i++ {
// 			res, _ = ConvertToFloat(typ[i])
// 			So(res, ShouldHaveSameTypeAs, chk)
// 		}
// 		_, err := ConvertToFloat(boolean)
// 		So(err.Error(), ShouldEqual, "ConvertToFloat: Unknown value is of incompatible type")
// 	})
// }

// func TestAddSearchTerm(t *testing.T) {

// 	Convey("Add search term", t, func() {
// 		AddSearchTerm("hello")
// 	})
// }

// func TestMainDate(t *testing.T) {
// 	t1 := time.Date(2012, 3, 1, 0, 0, 0, 0, time.UTC)
// 	t2 := time.Date(2012, 5, 2, 0, 0, 0, 0, time.UTC)
// 	t3 := time.Date(2012, 3, 3, 0, 0, 0, 0, time.UTC)
// 	t4 := time.Date(2012, 2, 4, 0, 0, 0, 0, time.UTC)
// 	t5 := time.Date(2012, 3, 5, 0, 0, 0, 0, time.UTC)
// 	t6 := time.Date(2012, 2, 6, 0, 0, 0, 0, time.UTC)
// 	t7 := time.Date(2012, 5, 7, 0, 0, 0, 0, time.UTC)

// 	dv := make([]DateVal, 7)
// 	dv[0].Date = t1
// 	dv[1].Date = t2
// 	dv[2].Date = t3
// 	dv[3].Date = t4
// 	dv[4].Date = t5
// 	dv[5].Date = t6
// 	dv[6].Date = t7

// 	Convey("Get main date", t, func() {
// 		result := MainDate(dv)
// 		So(result, ShouldEqual, "March 2012")
// 	})
// }

// func TestPrimaryDate(t *testing.T) {
// 	Convey("Get main dates", t, func() {
// 		PrimaryDate()
// 	})
// }
