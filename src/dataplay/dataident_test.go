package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
)

func TestIdentifyTable(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"id": "",
	}

	result := IdentifyTable(response,request,prams)

	Convey("When no ID parameter is provided", t, func() {
		So(result, ShouldEqual, "")
		})

	// prams["id"] = "gold"
	// result = SetDefaults(response, request, prams)

	// Convey("When ID parameter is provided", t, func() {
	// 	So(result, ShouldNotBeBlank)
	// 	})
}


func TestAttemptToFindMatches(t *testing.T) {
}

func TestFindStringMatches(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"x": "",
		"name": "",
	}

	result := FindStringMatches(response,request,prams)

	Convey("When no ID parameter is provided", t, func() {
		So(result, ShouldEqual, "")
		})

	// prams["x"] = "test"
	// prams["word"] = "test"

	// result = FindStringMatches(response,request,prams)

	// Convey("When ID parameter is provided", t, func() {
	// 	So(response.Code, ShouldEqual, http.StatusInternalServerError)
	// })
}

func TestGetRelatedDatasetByStrings(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"guid": "",
	}

	result := FindStringMatches(response,request,prams)

	Convey("When no guid parameter is provided", t, func() {
		So(result, ShouldEqual, "")
		})
}

func TestSuggestColType(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"id": "",
	}

	result := SuggestColType(response,request,prams)

	Convey("When no ID parameter is provided", t, func() {
		So(result, ShouldEqual, "")
		})

		// prams["id"] = "gold"
	// result = SetDefaults(response, request, prams)

	// Convey("When ID parameter is provided", t, func() {
	// 	So(result, ShouldNotBeBlank)
	// 	})

}

func TestConvertIntoStructArrayAndSort (t *testing.T) {
	unsorted := map[string]int{
		"b": 2,
		"c": 3,
		"a": 1,
	}

	result := ConvertIntoStructArrayAndSort(unsorted)
	resultStr, _ := json.Marshal(result)

	Convey("Unsorted map bca should return sorted map abc", t, func() {
		So(string(resultStr), ShouldEqual, "[{\"Key\":\"a\"},{\"Key\":\"b\"},{\"Key\":\"c\"}]")
	})

}

func TestStringInSlice (t *testing.T) {
	test := []string{"a", "b", "c"}

	result:= StringInSlice("a", test)

	Convey("Should find string \"a\" in slice \"a\",\"b\",\"c\" ", t, func() {
		So(result, ShouldBeTrue)
	})

	result = StringInSlice("x", test)

	Convey("Should not find string \"x\" in slice \"a\",\"b\",\"c\" ", t, func() {
		So(result, ShouldBeFalse)
	})
}
