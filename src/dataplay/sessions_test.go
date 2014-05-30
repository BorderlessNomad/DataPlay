package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIsUserLoggedIn(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=mayur@dataplay.com&password=whoru007"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()

	Convey("On HTTP Request", t, func() {
		HandleLogin(response, request)
		handleLoggedIn(t)

		HandleLogout(response, request)
		handleLoggedOut(t)
	})
}

func handleLoggedIn(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	status := IsUserLoggedIn(response, request)

	Convey("When User is Logged In", func() {
		So(status, ShouldEqual, false) //Since we don't have cookies in Simulation
		So(response.Code, ShouldEqual, http.StatusOK)
	})
}

func handleLoggedOut(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	status := IsUserLoggedIn(response, request)

	Convey("When User is Logged Out", func() {
		So(status, ShouldEqual, false) //Since we don't have cookies in Simulation
		So(response.Code, ShouldEqual, http.StatusOK)
	})
}
