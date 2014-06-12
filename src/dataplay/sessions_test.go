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

func TestClearSession(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	e_return := ClearSession(response, request)

	Convey("When there is no cookie to be found", t, func() {
		So(e_return, ShouldEqual, "No cookie found")
	})

}

func TestRandString(t *testing.T) {
	result := randString(5)

	Convey("When Random String length is 5", t, func() {
		So(len(result), ShouldEqual, 5)
		So(len(result), ShouldNotEqual, 6)
		So(result, ShouldNotContainSubstring, "!")

	})
}
