package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetDefaults(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"id": "",
	}

	result := SetDefaults(response,request,prams)

	Convey("When ID parameter is not provided", t, func() {
		So(result, ShouldEqual, "")
		})

	// prams["id"] = "0"
	// result = SetDefaults(response, request, prams)

	// Convey("When ID parameter is provided", t, func() {
	// 	So(result, ShouldNotBeBlank)
	// 	})
}

func TestGetDefaults(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	prams := map[string]string{
		"id": "",
	}

	result := GetDefaults(response,request,prams)

	Convey("When ID parameter is not provided", t, func() {
		So(result, ShouldEqual, "")
		})

// 	prams["id"] = "0"
// 	result = GetDefaults(response, request, prams)

// 	Convey("When ID parameter is provided", t, func() {
// 		So(result, ShouldNotBeBlank)
// 		})
}
