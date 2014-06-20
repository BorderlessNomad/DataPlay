package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckImportStatus(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	params := map[string]string{
		"id": "",
	}

	CheckImportStatus(response, request, params)

	Convey("When no ID parameter is provided", t, func() {
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})

	Convey("When ID parameter is provided", t, func() {
		params["id"] = "gold"
		result := CheckImportStatus(response, request, params)
		So(result, ShouldEqual, "{\"State\":\"online\",\"Request\":\"gold\"}")
	})
}
