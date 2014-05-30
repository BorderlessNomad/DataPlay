package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsUserLoggedIn(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	Convey("On HTTP Request", t, func() {
		IsUserLoggedIn(response, request)

		Convey("When user is logged in", func() {
			So(response.Code, ShouldEqual, http.StatusOK)
		})
	})
}
