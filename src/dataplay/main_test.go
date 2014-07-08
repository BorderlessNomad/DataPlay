package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(t *testing.T) {
	*mode = 1
	go main()
}

func TestMain2(t *testing.T) {
	*mode = 2
	go func() {
		main()
	}()
}

func TestJsonApiHandler(t *testing.T) {
	request, _ := http.NewRequest("GET", "http://dataplay.com/api/user", nil)
	response := httptest.NewRecorder()
	JsonApiHandler(response, request)
}

func TestSendToQueue(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	reqstr, methstr := "", ""
	params := map[string]string{"": ""}
	result := sendToQueue(response, request, params, reqstr, methstr)
	Convey("Should return response channel", t, func() {
		So(result, ShouldNotBeNil)
	})
}
