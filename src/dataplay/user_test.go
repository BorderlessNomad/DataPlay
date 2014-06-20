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

func TestCheckAuthRedirect(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	Convey("On HTTP Request", t, func() {
		CheckAuthRedirect(response, request)

		Convey("When authentication is successful", func() {
			So(response.Code, ShouldEqual, http.StatusTemporaryRedirect)
		})

		Convey("When authentication is unsuccessful", func() {
			So(response.Code, ShouldNotBeIn, []int{200, 201})
		})
	})
}

func TestHandleLogin(t *testing.T) {
	Convey("On HTTP Request 1", t, func() {
		handleLoginWithoutData(t)
	})
	Convey("On HTTP Request 2", t, func() {
		handleLoginInvalidData(t)
	})
	Convey("On HTTP Request 3", t, func() {
		handleLoginValidData(t)
	})
	Convey("On HTTP Request 4", t, func() {
		handleLoginValidDataMD5(t)
	})
}

func handleLoginInvalidData(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=random&password=123456"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()

	Convey("When Invalid data is provided", func() {
		HandleLogin(response, request)
		So(response.Code, ShouldEqual, http.StatusNotFound)
	})
}

func handleLoginValidData(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=mayur@dataplay.com&password=whoru007"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()

	Convey("When Correct data is provided", func() {
		HandleLogin(response, request)
		So(response.Code, ShouldEqual, http.StatusFound)
	})
}

func handleLoginValidDataMD5(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=glyn@dataplay.com&password=123456"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()

	//switch stored password to MD5 version before test
	DB.Model(&User{}).Where("email = ?", "glyn@dataplay.com").UpdateColumn("password", "e10adc3949ba59abbe56e057f20f883e")

	Convey("When user has old MD5 password", func() {
		HandleLogin(response, request)
		So(response.Code, ShouldEqual, http.StatusFound)
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

	Convey("When User does not exist", func() {
		HandleRegister(response, request)
		So(response.Code, ShouldEqual, http.StatusFound)
	})
}

func handleRegisterExisitingData(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=mayur@dataplay.com&password=whoru007"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()

	Convey("When User already exists", func() {
		HandleRegister(response, request)
		So(response.Code, ShouldEqual, http.StatusConflict)
	})
}

func handleLoginWithoutData(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=&password="))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()

	HandleLogin(response, request)

	Convey("When No data is provided", func() {
		So(response.Code, ShouldEqual, http.StatusNotFound)
	})
}
