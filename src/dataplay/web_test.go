package main

import (
	 //. "github.com/smartystreets/goconvey/convey"
	 "net/http"
	 "net/http/httptest"
	 "testing"
)

	// func TestAuthorisation (t *testing.T) {
	// 	request, _ := http.NewRequest("GET", "/", nil)
	// 	response := httptest.NewRecorder()

	// 	Authorisation(response, request)
	// }

	func TestLogin (t *testing.T) {
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()

		Login(response, request)
	}

	func TestLogout (t *testing.T) {
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()

		Logout(response, request)
	}

	func TestRegister (t *testing.T) {
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()

		Register(response, request)
	}

	func TestCharts (t *testing.T) {
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()
		prams := map[string]string{
			"id": "",
		}

		Charts(response, request, prams)
	}

	func TestSearchOverlay (t *testing.T) {
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()

		SearchOverlay(response, request)
	}

	func TestOverlay (t *testing.T) {
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()

		Overlay(response, request)
	}

	func TestOverview (t *testing.T) {
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()
		prams := map[string]string{
			"id": "",
		}

		Overview(response, request, prams)
	}

	func TestSearch (t *testing.T) {
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()

		Search(response, request)
	}

	func TestMapTest (t *testing.T) {
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()

		MapTest(response, request)
	}


