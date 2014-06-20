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
	params := map[string]string{
		"id": "",
	}
	result := SetDefaults(response, request, params)

	Convey("When ID parameter is not provided", t, func() {
		So(result, ShouldEqual, "")
	})

	Convey("When ID parameter is provided", t, func() {
		params["id"] = "test"
		result = SetDefaults(response, request, params)
		So(result, ShouldNotBeBlank)
	})
}

func TestGetDefaults(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	params := map[string]string{
		"id": "",
	}
	result := GetDefaults(response, request, params)

	Convey("When ID parameter is not provided", t, func() {
		So(result, ShouldEqual, "")
	})

	Convey("When incorrect ID parameter is provided", t, func() {
		params["id"] = "wightreq"
		result = GetDefaults(response, request, params)
		So(result, ShouldEqual, "{}")
	})

	Convey("When correct ID parameter is provided", t, func() {
		params["id"] = "hips"
		result = GetDefaults(response, request, params)
		So(result, ShouldEqual, "{\"x\":\"Hospital\",\"y\":\"70t79\"}")
	})
}
