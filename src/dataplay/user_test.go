package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCheckAuthRedirect(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	Convey("On HTTP Request 0", t, func() {
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
	// Convey("On HTTP Request 4", t, func() {
	// 	handleLoginValidDataMD5(t)
	// })
}

func handleLoginInvalidData(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=random&password=123456"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()
	u := UserForm{}
	u.Username = "nonsense"
	u.Password = "99999999"

	Convey("When Invalid data is provided", func() {
		HandleLogin(response, request, u)
		So(response.Code, ShouldEqual, http.StatusNotFound)
	})
}

func handleLoginValidData(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=mayur@dataplay.com&password=whoru007"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()
	u := UserForm{}
	u.Username = "glyn"
	u.Password = "123456"

	Convey("When Correct data is provided", func() {
		HandleLogin(response, request, u)
		So(response.Code, ShouldEqual, http.StatusOK)
	})
}

// func handleLoginValidDataMD5(t *testing.T) {
// 	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=glyn@dataplay.com&password=123456"))
// 	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
// 	response := httptest.NewRecorder()
// 	u := UserForm{}
// 	u.Username = "glyn@dataplay.com"
// 	u.Password = "123456"
// 	//switch stored password to MD5 version before test
// 	DB.Model(&User{}).Where("email = ?", "glyn@dataplay.com").UpdateColumn("password", "e10adc3949ba59abbe56e057f20f883e")

// 	Convey("When user has old MD5 password", func() {
// 		HandleLogin(response, request, u)
// 		So(response.Code, ShouldEqual, http.StatusFound)
// 	})
// }

func TestHandleRegister(t *testing.T) {
	// Convey("On HTTP Request 5", t, func() {
	// 	handleRegisterValidData(t)
	// })
	Convey("On HTTP Request 6", t, func() {
		handleRegisterExisitingData(t)
	})
}

// func handleRegisterValidData(t *testing.T) {
// 	time := time.Now()
// 	testuser := fmt.Sprintf("testuser_%d", time.Unix())
// 	request, _ := http.NewRequest("POST", "/", strings.NewReader("username="+testuser+"@dataplay.com&password=123456"))
// 	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
// 	response := httptest.NewRecorder()
// 	u := UserForm{}
// 	u.Username = testuser
// 	u.Password = "123456"

// 	Convey("When User does not already exist", func() {
// 		HandleRegister(response, request, u)
// 		So(response.Code, ShouldEqual, http.StatusOK)
// 	})
// }

func handleRegisterExisitingData(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=mayur@dataplay.com&password=whoru007"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()
	u := UserForm{}
	u.Username = "glyn"
	u.Password = "123456"
	u.Email = "glyn@dataplay.com"

	Convey("When User already exists", func() {
		HandleRegister(response, request, u)
		So(response.Code, ShouldEqual, http.StatusConflict)
	})
}

func handleLoginWithoutData(t *testing.T) {
	request, _ := http.NewRequest("POST", "/", strings.NewReader("username=&password="))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()
	u := UserForm{}
	u.Username = ""
	u.Password = ""

	HandleLogin(response, request, u)

	Convey("When No data is provided", func() {
		So(response.Code, ShouldEqual, http.StatusBadRequest)
	})
}

func TestHandleLogout(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	request.Header.Set("X-API-SESSION", "00pEM6oZTuo88GDrdtgpIYUzuw4LOVhkxdE9ywBTjpf44EcX03o037VYlvqacTbk")
	params := map[string]string{
		"s": "",
	}

	Convey("On HTTP Request 7", t, func() {
		HandleLogout(response, request, params)

		Convey("Should Logout", func() {
			So(response.Code, ShouldEqual, http.StatusOK)
		})
	})
}

func TestGetCreditedDiscoveriesHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()

	Convey("Should return credited discoveries", t, func() {
		result := GetCreditedDiscoveriesHttp(res, req)
		So(result, ShouldNotBeNil)
	})
}

func TestGetDiscoveriesHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()

	Convey("Should return discoveries", t, func() {
		result := GetDiscoveriesHttp(res, req)
		So(result, ShouldNotBeNil)
	})
}

func TestGetDataExpertsHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()

	Convey("Should return discoveries", t, func() {
		result := GetDataExpertsHttp(res, req)
		So(result, ShouldNotBeNil)
	})
}

func TestGetActivityStreamHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()

	Convey("Should return discoveries", t, func() {
		result := GetActivityStreamHttp(res, req)
		So(result, ShouldNotBeNil)
	})
}

func TestGetHomePageDataHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()

	Convey("Should return discoveries", t, func() {
		result := GetHomePageDataHttp(res, req)
		So(result, ShouldNotBeNil)
	})
}
