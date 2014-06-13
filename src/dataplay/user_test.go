package main

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

		Convey("When authentication is unsuccessful", func() {
			So(response.Code, ShouldNotBeIn, []int{200, 201})
			// So(response.Code, ShouldNotEqual, http.StatusTemporaryRedirect)
		})
	})
}

func TestHandleLogin(t *testing.T) {
	Convey("On HTTP Request", t, func() {
		handleLoginNoData(t)
	})
	Convey("On HTTP Request", t, func() {
		handleLoginInvalidData(t)
	})
	Convey("On HTTP Request", t, func() {
		handleLoginValidData(t)
	})
	Convey("On HTTP Request", t, func() {
		handleLoginValidDataMD5(t)
	})
}

func handleLoginNoData(t *testing.T) {
	// request, _ := http.NewRequest("POST", "/", nil)
	// response := httptest.NewRecorder()

	// HandleLogin(response, request)

	Convey("When No data is provided", func() {
		// So(response.Code, ShouldEqual, http.StatusNotFound)
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

func handleLoginValidDataMD5(t *testing.T) {
	// request, _ := http.NewRequest("POST", "/", strings.NewReader("username=glyn@dataplay.com&password=123456"))
	// request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	// response := httptest.NewRecorder()

	// //switch password to MD5
	// DB.SQL.Exec("UPDATE `DataCon`.`priv_users` SET `password`= ? WHERE `email`=?", "e10adc3949ba59abbe56e057f20f883e", "glyn@dataplay.com")

	// HandleLogin(response, request)

	Convey("When user has old MD5 password", func() {
		// So(response.Code, ShouldEqual, http.StatusFound)
	})
}

func TestHandleLogout(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	Convey("On HTTP Request", t, func() {
		HandleLogout(response, request)

		Convey("Should Logout", func() {
			So(response.Code, ShouldEqual, http.StatusTemporaryRedirect)
		})
	})
}

func TestHandleRegister(t *testing.T) {
	Convey("On HTTP Request 1", t, func() {
		handleRegisterValidData(t)
	})
	Convey("On HTTP Request 2", t, func() {
		handleRegisterExisitingData(t)
	})
}

func handleRegisterValidData(t *testing.T) {
	time := time.Now()
	testuser := fmt.Sprintf("testuser_%d", time.Unix())
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username="+testuser+"@dataplay.com&password=123456"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()

	HandleRegister(response, request)

	Convey("When User does not exist", func() {
		So(response.Code, ShouldEqual, http.StatusFound)
	})
}

func handleRegisterExisitingData(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=mayur@dataplay.com&password=whoru007"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()

	HandleRegister(response, request)

	Convey("When User already exists", func() {
		So(response.Code, ShouldEqual, http.StatusConflict)
	})
}
