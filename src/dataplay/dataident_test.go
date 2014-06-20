package main

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDBSetUp(t *testing.T) {
	DBSetup() // database needs to be set up as the first test or it will not be initialised until main
}

func TestIdentifyTable(t *testing.T) {
	Convey("With no parameter", t, func() {
		IdentifyTableNoParam(t)
	})

	Convey("With a parameter", t, func() {
		IdentifyTableWithParam(t)
	})
}

func IdentifyTableNoParam(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{"": ""}

	Convey("When no ID parameter is provided", func() {
		result := IdentifyTable(response, request, prams)
		So(result, ShouldEqual, "")
	})
}

func IdentifyTableWithParam(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{"id": "gold"}

	Convey("When ID parameter is provided", func() {
		result := IdentifyTable(response, request, prams)
		So(result, ShouldNotBeBlank)
	})
}

func TestFetchTableCols(t *testing.T) {
	result := FetchTableCols("")

	Convey("When no guid passed no column names are returned", t, func() {
		So(result, ShouldBeNil)
	})
}

func TestGetSQLTableSchema(t *testing.T) {
	Convey("When dbname > 0", t, func() {
		result := GetSQLTableSchema("test_table", "test_db")
		So(result, ShouldNotBeNil)
	})
}
func TestCheckColExists(t *testing.T) {
	Cols := []ColType{{"X", "0"}, {"Y", "0"}}
	result := CheckColExists(Cols, "X")

	Convey("When column exists", t, func() {
		So(result, ShouldBeTrue)
	})

	Convey("When column does not exist", t, func() {
		result = CheckColExists(Cols, "Z")
		So(result, ShouldBeFalse)
	})
}

func TestAttemptToFindMatches(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"id": "gdp",
		"x":  "year",
		"y":  "gdp",
	}

	Convey("When attempting to find matches", t, func() {
		result := AttemptToFindMatches(response, request, prams)
		So(result, ShouldEqual, "wat")
	})

}

func TestFindStringMatches(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"x":    "",
		"word": "",
	}

	result := FindStringMatches(response, request, prams)

	Convey("When no ID parameter is provided", t, func() {
		So(result, ShouldEqual, "")
	})

	Convey("When ID parameter is provided", t, func() {
		prams["x"] = "postal_code"
		prams["word"] = "B37 7YE"
		result = FindStringMatches(response, request, prams)
		So(result, ShouldNotBeBlank)
	})
}

func TestGetRelatedDatasetByStrings(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"guid": "",
	}

	result := GetRelatedDatasetByStrings(response, request, prams)
	Convey("When no guid parameter is provided", t, func() {
		So(result, ShouldEqual, "")
	})

	Convey("When guid parameter is provided", t, func() {
		prams["guid"] = "hips"
		result := GetRelatedDatasetByStrings(response, request, prams)
		So(result, ShouldNotBeBlank)
	})
}

func TestSuggestColType(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"table": "",
		"col":   "",
	}

	result := SuggestColType(response, request, prams)

	Convey("When no ID parameter is provided", t, func() {
		So(result, ShouldEqual, "")
	})

	Convey("When ID parameter is provided", t, func() {
		prams["table"] = "gold"
		prams["col"] = "price"
		result = SuggestColType(response, request, prams)
		So(result, ShouldEqual, "true")
	})

}

func TestConvertIntoStructArrayAndSort(t *testing.T) {
	unsorted := map[string]int{
		"b": 2,
		"c": 3,
		"a": 1,
	}

	Convey("Unsorted map bca should return sorted map abc", t, func() {
		result := ConvertIntoStructArrayAndSort(unsorted)
		resultStr, _ := json.Marshal(result)
		So(string(resultStr), ShouldEqual, `[{"Key":"a","Value":1},{"Key":"b","Value":2},{"Key":"c","Value":3}]`)
	})

}

func TestStringInSlice(t *testing.T) {
	test := []string{"a", "b", "c"}
	result := StringInSlice("a", test)

	Convey("Should find string \"a\" in slice \"a\",\"b\",\"c\" ", t, func() {
		So(result, ShouldBeTrue)
	})

	Convey("Should not find string \"x\" in slice \"a\",\"b\",\"c\" ", t, func() {
		result = StringInSlice("x", test)
		So(result, ShouldBeFalse)
	})
}
