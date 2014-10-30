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
	Convey("Should search financial", t, func() {
		result, _ := SearchForData(1, "financial", params)
		So(result, ShouldNotBeNil)
	})
	Convey("Should search trust", t, func() {
		result, _ := SearchForData(1, "trust", params)
		So(result, ShouldNotBeNil)
	})
	Convey("Should search ambulance", t, func() {
		result, _ := SearchForData(1, "ambulance", params)
		So(result, ShouldNotBeNil)
	})
}
