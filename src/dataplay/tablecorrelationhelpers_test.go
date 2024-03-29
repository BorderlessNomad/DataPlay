package main

// import (
// 	. "github.com/smartystreets/goconvey/convey"
// 	"testing"
// 	"time"
// )

// func TestRandomValueColumn(t *testing.T) {
// 	test := make([]ColType, 4)
// 	test[0].Name = "num_a"
// 	test[1].Name = "num_b"
// 	test[2].Name = "c"
// 	test[3].Name = "d"
// 	test[0].Sqltype = "integer"
// 	test[1].Sqltype = "float"
// 	test[2].Sqltype = "varchar"
// 	test[3].Sqltype = "int"

// 	Convey("Should return random date column", t, func() {
// 		result := RandomValueColumn(test)
// 		So(result, ShouldStartWith, "num")
// 	})
// }

// func TestRandomDateColumn(t *testing.T) {
// 	test := make([]ColType, 4)
// 	test[0].Name = "num"
// 	test[1].Name = "ddatea"
// 	test[2].Name = "dATE"
// 	test[3].Name = "date"

// 	Convey("Should return random date column", t, func() {
// 		result := RandomDateColumn(test)
// 		So(result, ShouldStartWith, "d")
// 	})
// }

// func TestRandomTableName(t *testing.T) {
// 	Convey("Should return random db table name", t, func() {
// 		result := RandomTableName()
// 		So(result, ShouldNotBeNil)
// 	})
// }

// func TestExtractDateVal(t *testing.T) {
// 	Convey("Should return extracted date and amoutn cols", t, func() {
// 		result, _ := ExtractDateVal("gold", "date", "price")
// 		So(result, ShouldNotBeNil)
// 	})
// }

// func TestDetermineRange(t *testing.T) {
// 	testX := make([]DateVal, 7)
// 	testY := make([]DateVal, 8)
// 	testZ := make([]DateVal, 1)

// 	testX[0].Date = time.Date(2013, 10, 1, 0, 0, 0, 0, time.UTC)
// 	testX[1].Date = time.Date(2013, 10, 7, 0, 0, 0, 0, time.UTC)
// 	testX[2].Date = time.Date(2013, 10, 2, 0, 0, 0, 0, time.UTC)
// 	testX[3].Date = time.Date(2013, 10, 5, 0, 0, 0, 0, time.UTC)
// 	testX[4].Date = time.Date(2013, 10, 6, 0, 0, 0, 0, time.UTC)
// 	testX[5].Date = time.Date(2013, 10, 3, 0, 0, 0, 0, time.UTC)
// 	testX[6].Date = time.Date(2013, 10, 4, 0, 0, 0, 0, time.UTC)

// 	testY[0].Date = time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC)
// 	testY[1].Date = time.Date(2013, 10, 2, 0, 0, 0, 0, time.UTC)
// 	testY[2].Date = time.Date(2013, 10, 3, 0, 0, 0, 0, time.UTC)
// 	testY[3].Date = time.Date(2013, 10, 4, 0, 0, 0, 0, time.UTC)
// 	testY[4].Date = time.Date(2013, 10, 8, 0, 0, 0, 0, time.UTC)
// 	testY[5].Date = time.Date(2013, 12, 31, 0, 0, 0, 0, time.UTC)
// 	testY[6].Date = time.Date(2013, 10, 6, 0, 0, 0, 0, time.UTC)
// 	testY[7].Date = time.Date(2013, 10, 7, 0, 0, 0, 0, time.UTC)

// 	testZ[0].Date = time.Date(2013, 1, 1, 0, 0, 0, 0, time.UTC)

// 	_, _, result := DetermineRange(testX)
// 	Convey("Should return correct range", t, func() {
// 		So(result, ShouldEqual, 6)
// 	})

// 	Convey("Should return correct range", t, func() {
// 		_, _, result = DetermineRange(testY)
// 		So(result, ShouldEqual, 364)
// 	})

// 	Convey("Should return correct from date", t, func() {
// 		from, _, _ := DetermineRange(testY)
// 		So(from, ShouldResemble, testY[0].Date)
// 	})

// 	Convey("Should return correct to date", t, func() {
// 		_, to, _ := DetermineRange(testY)
// 		So(to, ShouldResemble, testY[5].Date)
// 	})

// 	Convey("Should return nothing if range is not large enough", t, func() {
// 		_, _, result = DetermineRange(testZ)
// 		So(result, ShouldEqual, 0)
// 	})
// }

// func TestCreateBuckets(t *testing.T) {
// 	t1 := time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC)
// 	t2 := time.Date(2014, 1, 2, 0, 0, 0, 0, time.UTC)
// 	t3 := time.Date(2014, 1, 3, 0, 0, 0, 0, time.UTC)
// 	t4 := time.Date(2014, 1, 4, 0, 0, 0, 0, time.UTC)
// 	t5 := time.Date(2014, 1, 5, 0, 0, 0, 0, time.UTC)
// 	t6 := time.Date(2014, 1, 6, 0, 0, 0, 0, time.UTC)
// 	t7 := time.Date(2014, 1, 7, 0, 0, 0, 0, time.UTC)

// 	chk := make([]FromTo, 6)
// 	chk[0].From = t1
// 	chk[0].To = t2
// 	chk[1].From = t2
// 	chk[1].To = t3
// 	chk[2].From = t3
// 	chk[2].To = t4
// 	chk[3].From = t4
// 	chk[3].To = t5
// 	chk[4].From = t5
// 	chk[4].To = t6
// 	chk[5].From = t6
// 	chk[5].To = t7

// 	Convey("Should return range of dated FromTo buckets", t, func() {
// 		result := CreateBuckets(t1, t6, 6)
// 		So(result, ShouldResemble, chk)
// 	})
// }

// func TestFillBuckets(t *testing.T) {

// 	testDA := make([]DateVal, 5)
// 	testDA[0].Date = time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC)
// 	testDA[1].Date = time.Date(2014, 1, 2, 0, 0, 0, 0, time.UTC)
// 	testDA[2].Date = time.Date(2014, 1, 3, 0, 0, 0, 0, time.UTC)
// 	testDA[3].Date = time.Date(2014, 2, 28, 0, 0, 0, 0, time.UTC)
// 	testDA[4].Date = time.Date(2014, 12, 31, 0, 0, 0, 0, time.UTC)
// 	testDA[0].Value = 1.3
// 	testDA[1].Value = 1.7
// 	testDA[2].Value = 3.4
// 	testDA[3].Value = 6.6
// 	testDA[4].Value = 5.0

// 	testBkt := make([]FromTo, 3)
// 	testBkt[0].From = time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC)
// 	testBkt[0].To = time.Date(2014, 1, 3, 0, 0, 0, 0, time.UTC)
// 	testBkt[1].From = time.Date(2014, 1, 3, 0, 0, 0, 0, time.UTC)
// 	testBkt[1].To = time.Date(2014, 3, 1, 0, 0, 0, 0, time.UTC)
// 	testBkt[2].From = time.Date(2014, 3, 1, 0, 0, 0, 0, time.UTC)
// 	testBkt[2].To = time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)

// 	Convey("Should return bucket", t, func() {
// 		result := FillBuckets(testDA, testBkt)
// 		chk := []float64{3.0, 10.0, 5.0}
// 		So(result, ShouldResemble, chk)
// 	})

// }

// func TestDaysInMonth(t *testing.T) {
// 	date := time.Date(2016, 2, 11, 0, 0, 0, 0, time.UTC)

// 	Convey("Should return days in February 2016", t, func() {
// 		result := daysInMonth(date.Month(), date.Year())
// 		So(result, ShouldEqual, 29)
// 	})
// }

// func TestDaysInYear(t *testing.T) {
// 	date := time.Date(2016, 8, 4, 0, 0, 0, 0, time.UTC)

// 	Convey("Should return days in 2016", t, func() {
// 		result := daysInYear(date.Year())
// 		So(result, ShouldEqual, 366)
// 	})
// }

// func TestBetween(t *testing.T) {
// 	from := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
// 	to := time.Date(2010, 12, 31, 0, 0, 0, 0, time.UTC)
// 	date := DateVal{Date: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC), Value: 0}

// 	Convey("Date should be on start date", t, func() {
// 		result := date.Between(from, to)
// 		So(result, ShouldEqual, true)
// 	})

// 	Convey("Date should be between start and end dates", t, func() {
// 		date = DateVal{Date: time.Date(2010, 6, 1, 0, 0, 0, 0, time.UTC), Value: 0}
// 		result := date.Between(from, to)
// 		So(result, ShouldEqual, true)
// 	})

// 	Convey("Date should be outside start and end dates", t, func() {
// 		date = DateVal{Date: time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC), Value: 0}
// 		result := date.Between(from, to)
// 		So(result, ShouldEqual, false)
// 	})
// }

// func TestDayNum(t *testing.T) {
// 	Convey("Should return date as day number since 01/01/1900", t, func() {
// 		date := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
// 		result := DayNum(date)
// 		So(result, ShouldEqual, 42004)
// 	})
// }
