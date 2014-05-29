package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCheckAuth(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	Convey("On HTTP Request", t, func() {
		checkAuth(response, request)

		Convey("When authentication is successful", func() {
			// So(response.Code, ShouldBeIn, []int{200, 201, 301, 302, 303, 307})
			So(response.Code, ShouldEqual, http.StatusTemporaryRedirect)
		})

		Convey("When authentication is not successful", func() {
			// So(response.Code, ShouldNotBeIn, []int{200, 201, 301, 302, 303, 307})
			// So(response.Code, ShouldNotEqual, http.StatusTemporaryRedirect)
		})
	})
}

func TestHandleLogin(t *testing.T) {
	Convey("On HTTP Request", t, func() {
		handleLoginNoData(t)
		handleLoginInvalidData(t)
		handleLoginValidData(t)
	})
}

func handleLoginNoData(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", nil)
	response := httptest.NewRecorder()

	HandleLogin(response, request)

	Convey("When No data is provided", func() {
		So(response.Code, ShouldEqual, http.StatusNotFound)
	})
}

func handleLoginInvalidData(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=test&password=123456"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()

	HandleLogin(response, request)

	Convey("When Invalid data is provided", func() {
		So(response.Code, ShouldEqual, http.StatusNotFound)
	})
}

func handleLoginValidData(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=mayur@dataplay.com&password=whoru007"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()

	HandleLogin(response, request)

	Convey("When Correct data is provided", func() {
		So(response.Code, ShouldEqual, http.StatusFound)
	})
}

func TestHandleRegister(t *testing.T) {

}
