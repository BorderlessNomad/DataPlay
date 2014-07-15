package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestRandomTableGUID(t *testing.T) {
	Convey("Should return random db table name", t, func() {
		result := RandomTableGUID()
		So(result, ShouldNotBeBlank)
	})
}

// func TestGetCorrelation(t *testing.T) {
// 	Convey("Should return JSON string with correlation", t, func() {
// 		result := GetCorrelation("gold")
// 		So(result, ShouldNotBeNil)
// 	})
// }

func TestExtractDateAmt(t *testing.T) {
	Convey("Should return extracted date and amoutn cols", t, func() {
		result := ExtractDateAmt("gold", "date", "price")
		So(result, ShouldNotBeNil)
	})
}

func TestDetermineRange(t *testing.T) {
	testX := make([]DateAmt, 7)
	testY := make([]DateAmt, 8)
	testX[0].Date = time.Date(2013, 10, 01, 0, 0, 0, 0, time.UTC)
	testX[1].Date = time.Date(2013, 10, 07, 0, 0, 0, 0, time.UTC)
	testX[2].Date = time.Date(2013, 10, 02, 0, 0, 0, 0, time.UTC)
	testX[3].Date = time.Date(2013, 10, 05, 0, 0, 0, 0, time.UTC)
	testX[4].Date = time.Date(2013, 10, 06, 0, 0, 0, 0, time.UTC)
	testX[5].Date = time.Date(2013, 10, 03, 0, 0, 0, 0, time.UTC)
	testX[6].Date = time.Date(2013, 10, 04, 0, 0, 0, 0, time.UTC)

	testY[0].Date = time.Date(2011, 10, 01, 0, 0, 0, 0, time.UTC)
	testY[1].Date = time.Date(2013, 10, 02, 0, 0, 0, 0, time.UTC)
	testY[2].Date = time.Date(2013, 10, 03, 0, 0, 0, 0, time.UTC)
	testY[3].Date = time.Date(2013, 10, 04, 0, 0, 0, 0, time.UTC)
	testY[4].Date = time.Date(2013, 10, 8, 0, 0, 0, 0, time.UTC)
	testY[5].Date = time.Date(2013, 10, 05, 0, 0, 0, 0, time.UTC)
	testY[6].Date = time.Date(2013, 10, 06, 0, 0, 0, 0, time.UTC)
	testY[7].Date = time.Date(2013, 10, 07, 0, 0, 0, 0, time.UTC)

	Convey("Should return correct range", t, func() {
		_, _, result1 := DetermineRange(testX)
		So(result1, ShouldEqual, 6)
	})
	Convey("Should return correct range", t, func() {
		_, _, result2 := DetermineRange(testY)
		So(result2, ShouldEqual, 737)
	})
}
