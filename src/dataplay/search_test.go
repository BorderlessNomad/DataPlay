package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDataDictionary(t *testing.T) {
	Convey("Create DataDictionary", t, func() {
		DataDict()
	})
}

func TestSearchForData(t *testing.T) {

	params := map[string]string{
		"offset": "0",
		"count":  "5",
	}
	Convey("Should search", t, func() {
		result, _ := SearchForData(1, "financial", params)
		So(result, ShouldEqual, "")
	})
}
