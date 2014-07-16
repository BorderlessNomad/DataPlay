package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestGetCorrelation(t *testing.T) {
	for i := 0; i < 100; i++ {
		Convey("Should return JSON string with correlation", t, func() {
			table := RandomTableName()
			result := GetCorrelation(table)
			So(result, ShouldNotBeNil)
		})
	}
}

func TestGetCoef(t *testing.T) {
	tst := make(map[string]string)

	Convey("Should return nothing when passed empty map", t, func() {
		result := GetCoef(tst)
		So(result, ShouldEqual, 0)
	})

	tst["table1"] = "gold"
	tst["table2"] = "gold"
	tst["amtCol1"] = "price"
	tst["amtCol2"] = "price"
	tst["dateCol1"] = "date"
	tst["dateCol2"] = "date"

	Convey("Should return coefficient value when passed map", t, func() {
		result := GetCoef(tst)
		So(result, ShouldEqual, 0.9999762331129333)
	})

}

func TestRandomAmountColumn(t *testing.T) {
	test := make([]ColType, 4)
	test[0].Name = "num_a"
	test[1].Name = "num_b"
	test[2].Name = "c"
	test[3].Name = "d"
	test[0].Sqltype = "integer"
	test[1].Sqltype = "float"
	test[2].Sqltype = "varchar"
	test[3].Sqltype = "int"

	Convey("Should return random date column", t, func() {
		result := RandomAmountColumn(test)
		So(result, ShouldStartWith, "num")
	})
}

func TestRandomDateColumn(t *testing.T) {
	test := make([]ColType, 4)
	test[0].Name = "num"
	test[1].Name = "ddatea"
	test[2].Name = "dATE"
	test[3].Name = "date"

	Convey("Should return random date column", t, func() {
		result := RandomDateColumn(test)
		So(result, ShouldStartWith, "d")
	})
}

func TestRandomTableName(t *testing.T) {
	Convey("Should return random db table name", t, func() {
		result := RandomTableName()
		So(result, ShouldNotBeNil)
	})
}

func TestExtractDateAmt(t *testing.T) {
	Convey("Should return extracted date and amoutn cols", t, func() {
		result := ExtractDateAmt("gold", "date", "price")
		So(result, ShouldNotBeNil)
	})
}

func TestDetermineRange(t *testing.T) {
	testX := make([]DateAmt, 7)
	testY := make([]DateAmt, 8)
	testZ := make([]DateAmt, 1)
	testX[0].Date = time.Date(2013, 10, 1, 0, 0, 0, 0, time.UTC)
	testX[1].Date = time.Date(2013, 10, 7, 0, 0, 0, 0, time.UTC)
	testX[2].Date = time.Date(2013, 10, 2, 0, 0, 0, 0, time.UTC)
	testX[3].Date = time.Date(2013, 10, 5, 0, 0, 0, 0, time.UTC)
	testX[4].Date = time.Date(2013, 10, 6, 0, 0, 0, 0, time.UTC)
	testX[5].Date = time.Date(2013, 10, 3, 0, 0, 0, 0, time.UTC)
	testX[6].Date = time.Date(2013, 10, 4, 0, 0, 0, 0, time.UTC)

	testY[0].Date = time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC)
	testY[1].Date = time.Date(2013, 10, 2, 0, 0, 0, 0, time.UTC)
	testY[2].Date = time.Date(2013, 10, 3, 0, 0, 0, 0, time.UTC)
	testY[3].Date = time.Date(2013, 10, 4, 0, 0, 0, 0, time.UTC)
	testY[4].Date = time.Date(2013, 10, 8, 0, 0, 0, 0, time.UTC)
	testY[5].Date = time.Date(2013, 12, 31, 0, 0, 0, 0, time.UTC)
	testY[6].Date = time.Date(2013, 10, 6, 0, 0, 0, 0, time.UTC)
	testY[7].Date = time.Date(2013, 10, 7, 0, 0, 0, 0, time.UTC)

	testZ[0].Date = time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC)

	_, _, result := DetermineRange(testX)
	Convey("Should return correct range", t, func() {
		So(result, ShouldEqual, 6)
	})

	Convey("Should return correct range", t, func() {
		_, _, result = DetermineRange(testY)
		So(result, ShouldEqual, 364)
	})

	Convey("Should return correct from date", t, func() {
		from, _, _ := DetermineRange(testY)
		So(from, ShouldResemble, testY[0].Date)
	})

	Convey("Should return correct to date", t, func() {
		_, to, _ := DetermineRange(testY)
		So(to, ShouldResemble, testY[5].Date)
	})

	Convey("Should return nothing if range is not large enough", t, func() {
		_, _, result = DetermineRange(testZ)
		So(result, ShouldEqual, 0)
	})
}

func TestCreateBuckets(t *testing.T) {
	test := make([]DateAmt, 20)
	test[0].Date = time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC)
	test[1].Date = time.Date(2014, 12, 31, 0, 0, 0, 0, time.UTC)
	test[2].Date = time.Date(2013, 8, 2, 0, 0, 0, 0, time.UTC)
	test[3].Date = time.Date(2013, 9, 5, 0, 0, 0, 0, time.UTC)
	test[4].Date = time.Date(2014, 3, 6, 0, 0, 0, 0, time.UTC)
	test[5].Date = time.Date(2013, 10, 3, 0, 0, 0, 0, time.UTC)
	test[6].Date = time.Date(2013, 2, 4, 0, 0, 0, 0, time.UTC)
	test[7].Date = time.Date(2013, 10, 20, 0, 0, 0, 0, time.UTC)
	test[8].Date = time.Date(2013, 8, 23, 0, 0, 0, 0, time.UTC)
	test[9].Date = time.Date(2013, 10, 5, 0, 0, 0, 0, time.UTC)
	test[10].Date = time.Date(2013, 8, 6, 0, 0, 0, 0, time.UTC)
	test[11].Date = time.Date(2013, 6, 7, 0, 0, 0, 0, time.UTC)
	test[12].Date = time.Date(2014, 10, 2, 0, 0, 0, 0, time.UTC)
	test[13].Date = time.Date(2013, 11, 15, 0, 0, 0, 0, time.UTC)
	test[14].Date = time.Date(2014, 3, 6, 0, 0, 0, 0, time.UTC)
	test[15].Date = time.Date(2013, 10, 3, 0, 0, 0, 0, time.UTC)
	test[16].Date = time.Date(2013, 7, 14, 0, 0, 0, 0, time.UTC)
	test[17].Date = time.Date(2014, 1, 7, 0, 0, 0, 0, time.UTC)
	test[18].Date = time.Date(2013, 4, 12, 0, 0, 0, 0, time.UTC)
	test[19].Date = time.Date(2013, 9, 15, 0, 0, 0, 0, time.UTC)
	from := time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2013, 12, 31, 0, 0, 0, 0, time.UTC)

	Convey("Should return range of dated FromTo buckets", t, func() {
		result := CreateBuckets(test, from, to, 730)
		So(result, ShouldNotBeNil)
	})
}

func TestFillBuckets(t *testing.T) {

	testDA := make([]DateAmt, 5)
	testDA[0].Date = time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC)
	testDA[1].Date = time.Date(2014, 1, 2, 0, 0, 0, 0, time.UTC)
	testDA[2].Date = time.Date(2014, 1, 3, 0, 0, 0, 0, time.UTC)
	testDA[3].Date = time.Date(2014, 2, 28, 0, 0, 0, 0, time.UTC)
	testDA[4].Date = time.Date(2014, 12, 31, 0, 0, 0, 0, time.UTC)
	testDA[0].Amount = 1.2
	testDA[1].Amount = 0.8
	testDA[2].Amount = 3.7
	testDA[3].Amount = 6.3
	testDA[4].Amount = 5.0

	testBkt := make([]FromTo, 3)
	testBkt[0].From = time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC)
	testBkt[0].To = time.Date(2014, 1, 2, 0, 0, 0, 0, time.UTC)
	testBkt[1].From = time.Date(2014, 1, 3, 0, 0, 0, 0, time.UTC)
	testBkt[1].To = time.Date(2014, 3, 1, 0, 0, 0, 0, time.UTC)
	testBkt[2].From = time.Date(2014, 3, 1, 0, 0, 0, 0, time.UTC)
	testBkt[2].To = time.Date(2014, 12, 31, 0, 0, 0, 0, time.UTC)

	Convey("Should return bucket", t, func() {
		result := FillBuckets(testDA, testBkt)
		chk := []float64{2.0, 10.0, 5.0}
		So(result, ShouldResemble, chk)
	})

}

func TestDaysInMonth(t *testing.T) {
	date := time.Date(2016, 2, 11, 0, 0, 0, 0, time.UTC)

	Convey("Should return days in February 2016", t, func() {
		result := daysInMonth(date.Month(), date.Year())
		So(result, ShouldEqual, 29)
	})
}

func TestDaysInYear(t *testing.T) {
	date := time.Date(2016, 8, 4, 0, 0, 0, 0, time.UTC)

	Convey("Should return days in 2016", t, func() {
		result := daysInYear(date.Year())
		So(result, ShouldEqual, 366)
	})
}

func TestBetween(t *testing.T) {
	from := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2010, 12, 31, 0, 0, 0, 0, time.UTC)
	date := DateAmt{Date: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC), Amount: 0}

	Convey("Date should be on start date", t, func() {
		result := date.Between(from, to)
		So(result, ShouldEqual, true)
	})

	Convey("Date should be between start and end dates", t, func() {
		date = DateAmt{Date: time.Date(2010, 6, 1, 0, 0, 0, 0, time.UTC), Amount: 0}
		result := date.Between(from, to)
		So(result, ShouldEqual, true)
	})

	Convey("Date should be outside start and end dates", t, func() {
		date = DateAmt{Date: time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC), Amount: 0}
		result := date.Between(from, to)
		So(result, ShouldEqual, false)
	})
}

func TestDayNum(t *testing.T) {
	Convey("Should return date as day number since 01/01/1900", t, func() {
		date := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
		result := dayNum(date)
		So(result, ShouldEqual, 42004)
	})
}
