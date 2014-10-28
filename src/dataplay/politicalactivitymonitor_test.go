package main

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetPoliticalActivityHttp(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-SESSION", "00TK6wuwwj1DmVDtn8mmveDMVYKxAJKLVdghTynDXBd62wDqGUGlAmEykcnaaO66")
	res := httptest.NewRecorder()
	params := map[string]string{
		"type": "d",
	}
	a := time.Now()

	Convey("Should return departments PoliticalActivity", t, func() {
		result := GetPoliticalActivityHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})

	params["type"] = "e"

	b := time.Now()
	fmt.Println("POLITICAL DEPT TIME", b.Sub(a).Seconds())

	Convey("Should return events PoliticalActivity", t, func() {
		result := GetPoliticalActivityHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})

	params["type"] = "r"

	c := time.Now()
	fmt.Println("POLITICAL EVENT TIME", c.Sub(b).Seconds())

	Convey("Should return regions PoliticalActivity", t, func() {
		result := GetPoliticalActivityHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})

	params["type"] = "p"

	d := time.Now()
	fmt.Println("POLITICAL REGION TIME", d.Sub(c).Seconds())

	Convey("Should return popular PoliticalActivity", t, func() {
		result := GetPoliticalActivityHttp(res, req, params)
		So(result, ShouldNotBeNil)
	})

	z := time.Now()
	fmt.Println("POLITICAL POPULAR TIME", z.Sub(d).Seconds())
}
