package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddCommentHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{
		"uid":     "1",
		"comment": "TESTING TESTING 1 2 3",
	}

	Convey("Should add comment", t, func() {
		result := AddCommentHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})

}
