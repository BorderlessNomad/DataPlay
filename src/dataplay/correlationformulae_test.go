package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestPearson(t *testing.T) {
	var empty = []float64{}
	var data1 = []float64{29.8, 30.1, 30.5, 30.6, 31.3, 31.7, 32.6, 33.1, 32.7, 32.8}
	var data2 = []float64{327, 456, 509, 497, 596, 573, 661, 741, 809, 717}

	Convey("Should return 0 for empty data set", t, func() {
		result := Pearson(data1, empty)
		So(result, ShouldEqual, 0)
	})

	Convey("Should return valid correlation coefficient for sets of equal size", t, func() {
		result := Pearson(data1, data2)
		So(result, ShouldEqual, 0.9470910552716333)
	})

	Convey("Should return correlation coefficient of 0 for empty sets", t, func() {
		var empty = []float64{0.0}
		result := Pearson(empty, empty)
		So(result, ShouldEqual, 0.0)
	})

	Convey("Should return correlation coefficient of 0 when set 1 has 10 values and set 2 has more than 10 values", t, func() {
		data2 = append(data2, 300)
		result := Pearson(data1, data2)
		So(result, ShouldEqual, 0)
	})
}

func TestSpurious(t *testing.T) {
	var data1 = []float64{5249, 6402, 6854, 802, 2997, 9770, 8938, 1496, 8668, 5253, 1061, 3096, 5083, 4420, 9294, 250, 8648, 8602, 2440, 2267}
	var data2 = []float64{314, 196, 195, 244, 184, 249, 182, 232, 127, 141, 85, 189, 121, 338, 112, 225, 110, 168, 195, 262}
	var data3 = []float64{40036, 87900, 41390, 85953, 59604, 22848, 78542, 51792, 8811, 13540, 67289, 43760, 87331, 45984, 37737, 62219, 54737, 8169, 12550, 87735}

	Convey("Should return valid correlation coefficient for sets of equal size", t, func() {
		result := Spurious(data1, data2, data3)
		So(result, ShouldEqual, 0.7672911757618174)
	})

	Convey("Should return correlation coefficient of 0 for empty sets", t, func() {
		var empty = []float64{0.0}
		result := Spurious(empty, empty, empty)
		So(result, ShouldEqual, 0.0)
	})
}

func TestSpearman(t *testing.T) {
	var data1 = []float64{56, 75, 45, 71, 62, 64, 58, 80, 76, 61}
	var data2 = []float64{66, 70, 40, 60, 65, 56, 59, 77, 67, 63}

	Convey("Should return valid correlation coefficient for sets of equal size with unique data", t, func() {
		result := Spearman(data1, data2)
		So(result, ShouldEqual, 0.6727272727272727)
	})

	var tiedData1 = []float64{43, 75, 45, 71, 61, 64, 58, 80, 76, 61}
	var tiedData2 = []float64{66, 70, 40, 60, 12, 56, 59, 77, 67, 63}
	Convey("Should return valid correlation coefficient for sets of equal size with tied data", t, func() {
		result := Spearman(tiedData1, tiedData2)
		So(result, ShouldEqual, 0.5957474328064633)
	})
}

func TestVisual(t *testing.T) {
	t1 := time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2014, 1, 2, 0, 0, 0, 0, time.UTC)
	t3 := time.Date(2014, 1, 3, 0, 0, 0, 0, time.UTC)
	t4 := time.Date(2014, 1, 4, 0, 0, 0, 0, time.UTC)
	t5 := time.Date(2014, 1, 5, 0, 0, 0, 0, time.UTC)
	t6 := time.Date(2014, 1, 6, 0, 0, 0, 0, time.UTC)
	t7 := time.Date(2014, 1, 7, 0, 0, 0, 0, time.UTC)

	ft := make([]FromTo, 6)
	ft[0].From = t1
	ft[0].To = t2
	ft[1].From = t2
	ft[1].To = t3
	ft[2].From = t3
	ft[2].To = t4
	ft[3].From = t4
	ft[3].To = t5
	ft[4].From = t5
	ft[4].To = t6
	ft[5].From = t6
	ft[5].To = t7

	var data1 = []float64{65, 70, 40, 60, 65, 56}
	var data2 = []float64{66, 70, 40, 60, 65, 56}

	Convey("Should return valid correlation coefficient for sets of equal size with tied data", t, func() {
		result := Visual(data1, data2, ft)
		So(result, ShouldEqual, 0.9971988795518206)
	})
}

func TestMean(t *testing.T) {
	var data = []float64{29.8, 30.1, 30.5, 30.6, 31.3, 31.7, 32.6, 33.1, 32.7, 32.8}
	Convey("Should return mean", t, func() {
		result := Mean(data)
		So(result, ShouldEqual, 31.52)
	})
}

func TestVariation(t *testing.T) {
	var data = []float64{29.8, 30.1, 30.5, 30.6, 31.3, 31.7, 32.6, 33.1, 32.7, 32.8}
	Convey("Should return coefficient of variation", t, func() {
		result := Variation(data)
		So(result, ShouldEqual, 0.037047361870545685)
	})
}

func TestStandDev(t *testing.T) {
	var data = []float64{29.8, 30.1, 30.5, 30.6, 31.3, 31.7, 32.6, 33.1, 32.7, 32.8}
	Convey("Should return standard deviation", t, func() {
		result := StandDev(data)
		So(result, ShouldEqual, 1.1677328461596)
	})
}

func TestSgn(t *testing.T) {
	a := 0.0

	Convey("Should return standard sgn(a) = 0", t, func() {
		result := Sgn(a)
		So(result, ShouldEqual, 0)
	})

	Convey("Should return standard sgn(a) = 1", t, func() {
		a = 2
		result := Sgn(a)
		So(result, ShouldEqual, 1)
	})

	Convey("Should return standard sgn(a) = -1", t, func() {
		a = -2
		result := Sgn(a)
		So(result, ShouldEqual, -1)
	})
}
