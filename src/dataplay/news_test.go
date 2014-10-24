package main

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSearchForNewsHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{
		"terms": "gold",
	}

	Convey("Should return news", t, func() {
		a := time.Now()
		result := SearchForNewsHttp(res, req, params)
		b := time.Now()
		fmt.Println("ROBOCOP MAIN", b.Sub(a).Seconds())
		So(result, ShouldNotBeNil)
	})
}
