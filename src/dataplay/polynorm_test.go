package main

// import (
// 	. "github.com/smartystreets/goconvey/convey"
// 	"testing"
// )

// func TestGetPolyResults(t *testing.T) {
// 	var xShort = []float64{1, 2, 3, 4}
// 	var yShort = []float64{1, 2, 3, 4}

// 	var xDiffLen = []float64{1, 2, 3, 4, 5, 6}
// 	var yDiffLen = []float64{1, 2, 3, 4, 5}

// 	var xCorrect = []float64{1.1, 3.7, 4.3, 2.6, 12.8}
// 	var yCorrect = []float64{4.6, 2.5, 3.9, 1.4, 8.5}

// 	result := []float64{1, 1, 1}
// 	fail := []float64{0, 0, 0}
// 	pass := []float64{4.568273715976804, -0.8052095636667054, 0.08719185121795317}

// 	Convey("When X & Y arrays size is below 5", t, func() {
// 		result = GetPolyResults(xShort, yShort)
// 		So(result, ShouldResemble, fail)
// 	})

// 	Convey("When X & Y arrays are different lengths", t, func() {
// 		result = GetPolyResults(xDiffLen, yDiffLen)
// 		So(result, ShouldResemble, fail)
// 	})

// 	Convey("When X & Y arrays are valid", t, func() {
// 		result = GetPolyResults(xCorrect, yCorrect)
// 		So(result, ShouldResemble, pass)
// 	})
// }
