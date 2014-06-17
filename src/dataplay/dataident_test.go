package main

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIdentifyTable (t *testing.T) {
	Convey("With no parameter", t, func() {
		IdentifyTableNoParam(t)
	})
	// Convey("With a parameter", t, func() {
	// 	IdentifyTableWithParam(t)
	// })
}

func IdentifyTableNoParam(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{"id": ""}

	result := IdentifyTable(response, request, prams)

	Convey("When no ID parameter is provided", func() {
		So(result, ShouldEqual, "")
	})
}

// func IdentifyTableWithParam(t *testing.T) {
// 	request, _ := http.NewRequest("GET", "/", nil)
// 	response := httptest.NewRecorder()
// 	prams := map[string]string{"id": "gold"}

// 	result := IdentifyTable(response, request, prams)

// 	Convey("When ID parameter is provided", func() {
// 		So(result, ShouldEqual, "")
// 	})
// }

func TestCheckColExists(t *testing.T){
	Cols := []ColType{{"X", "0"}, {"Y", "0"}}

	result := CheckColExists(Cols,"X")
	Convey("When column exists", t, func() {
		So(result, ShouldBeTrue)
	})
	result = CheckColExists(Cols,"Z")
	Convey("When column does not exist", t, func() {
		So(result, ShouldBeFalse)
	})

}

func TestAttemptToFindMatches(t *testing.T) {
}

func TestFindStringMatches(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"word":    "",
	}

	result := FindStringMatches(response, request, prams)

	Convey("When no ID parameter is provided", t, func() {
		So(result, ShouldEqual, "")
	})

	// prams["x"] = "test"
	// prams["word"] = "test"

	// result = FindStringMatches(response,request,prams)

	// Convey("When ID parameter is provided", t, func() {
	// 	So(response.Code, ShouldEqual, http.StatusBadRequest)
	// })
}

func TestGetRelatedDatasetByStrings(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"guid": "",
	}

	result := FindStringMatches(response, request, prams)

	Convey("When no guid parameter is provided", t, func() {
		So(result, ShouldEqual, "")
	})
}

func TestSuggestColType(t *testing.T) {
	// request, _ := http.NewRequest("POST", "/", nil)
	// response := httptest.NewRecorder()
	// // prams := map[string]string{
	// // 	"id": "",
	// // }

	// // result := SuggestColType(response, request, prams)

	// // Convey("When no ID parameter is provided", t, func() {
	// // 	So(result, ShouldEqual, "")
	// // })

	// prams := map[string]string{
	// 	"table" : "gold",
	// 	"col" : "price",
	// }
	// result := SuggestColType(response, request, prams)

	// Convey("When ID parameter is provided", t, func() {
	// 	So(result, ShouldBeBlank)
	// })

}

func TestConvertIntoStructArrayAndSort(t *testing.T) {
	unsorted := map[string]int{
		"b": 2,
		"c": 3,
		"a": 1,
	}

	result := ConvertIntoStructArrayAndSort(unsorted)
	resultStr, _ := json.Marshal(result)

	Convey("Unsorted map bca should return sorted map abc", t, func() {
		So(string(resultStr), ShouldEqual, `[{"Key":"a","Value":1},{"Key":"b","Value":2},{"Key":"c","Value":3}]`)
	})

}

func TestStringInSlice(t *testing.T) {
	test := []string{"a", "b", "c"}

	result := StringInSlice("a", test)

	Convey("Should find string \"a\" in slice \"a\",\"b\",\"c\" ", t, func() {
		So(result, ShouldBeTrue)
	})

	result = StringInSlice("x", test)

	Convey("Should not find string \"x\" in slice \"a\",\"b\",\"c\" ", t, func() {
		So(result, ShouldBeFalse)
	})
}
