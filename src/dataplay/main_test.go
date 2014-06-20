package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(t *testing.T) {
	go func() {
		main()
	}()
}

func TestJsonApiHandler(t *testing.T) {
	request, _ := http.NewRequest("GET", "http://dataplay.com/api/user", nil)
	response := httptest.NewRecorder()
	JsonApiHandler(response, request)
}
