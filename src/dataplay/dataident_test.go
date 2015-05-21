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
	params := map[string]string{"": ""}

	Convey("When no ID parameter is provided", func() {
		result := IdentifyTable(response, request, params)
		So(result, ShouldEqual, "")
	})
}

func IdentifyTableWithParam(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	params := map[string]string{"id": "gold"}

	Convey("When ID parameter is provided", func() {
		result := IdentifyTable(response, request, params)
		So(result, ShouldNotBeBlank)
	})
}

func TestFetchTableCols(t *testing.T) {
	result := FetchTableCols("")

	Convey("When no guid passed no column names are returned", t, func() {
		So(result, ShouldBeNil)
	})

	Convey("Should return column names and types", t, func() {
		result = FetchTableCols("births")
		So(result, ShouldNotBeNil)
	})
}

func TestHasTableGotLocatonData(t *testing.T) {
	result := HasTableGotLocationData("tweets")

	Convey("Should find Lattitude and Longitude columns in dataset", t, func() {
		So(result, ShouldEqual, true)
	})

	Convey("Should not find Lattitude and Longitude columns in dataset", t, func() {
		result = HasTableGotLocationData("houseprices")
		So(result, ShouldEqual, false)
	})
}

func TestContainsTableCol(t *testing.T) {
	Cols := []ColType{{"X", "0"}, {"Y", "0"}}
	result := ContainsTableCol(Cols, "y")

	Convey("Should find Column name", t, func() {
		So(result, ShouldBeTrue)
	})

	Convey("Should not find Column name", t, func() {
		result = ContainsTableCol(Cols, "z")
		So(result, ShouldBeFalse)
	})
}

func TestGetSQLTableSchema(t *testing.T) {
	Convey("When dbname > 0", t, func() {
		result := GetSQLTableSchema("test_table")
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
	params := map[string]string{"id": ""}

	Convey("When id parameter is incorrect", t, func() {
		params["id"] = "qwerty1"
		result := AttemptToFindMatches(response, request, params)
		So(result, ShouldEqual, "")
	})

	Convey("When col parameters are incorrect", t, func() {
		params["id"] = "gdp"
		params["x"] = "qwerty1"
		params["y"] = "qwerty1"
		result := AttemptToFindMatches(response, request, params)
		So(result, ShouldEqual, "")
	})

	// Convey("When parameters are correct", t, func() {
	// 	params["id"] = "gdp"
	// 	params["x"] = "date"
	// 	params["y"] = "gdp"
	// 	result := AttemptToFindMatches(response, request, params)
	// 	So(result, ShouldNotBeBlank)
	// })
}

func TestFindStringMatches(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	params := map[string]string{
		"x":    "",
		"word": "",
	}

	result := FindStringMatches(response, request, params)

	Convey("When no ID parameter is provided", t, func() {
		So(result, ShouldEqual, "")
	})

	Convey("When invalid ID parameter is provided with invalid string to match", t, func() {
		params["x"] = "qwerty1"
		params["word"] = ""
		result = FindStringMatches(response, request, params)
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When valid ID parameter is provided with valid string to match", t, func() {
		params["x"] = "postal_code"
		params["word"] = "B37 7YE"
		result = FindStringMatches(response, request, params)
		So(result, ShouldNotBeBlank)
	})
}

func TestGetRelatedDatasetByStrings(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	params := map[string]string{
		"guid": "",
	}

	result := GetRelatedDatasetByStrings(response, request, params)
	Convey("When no guid parameter is provided", t, func() {
		So(result, ShouldEqual, "")
	})

	Convey("When guid parameter is provided", t, func() {
		params["guid"] = "hips"
		result := GetRelatedDatasetByStrings(response, request, params)
		So(result, ShouldNotBeBlank)
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
