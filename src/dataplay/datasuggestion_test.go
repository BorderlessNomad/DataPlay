package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRandomTableGUID(t *testing.T) {
	Convey("Should return random db table name", t, func() {
		result := RandomTableGUID()
		So(result, ShouldNotBeBlank)
	})
}

func TestGetCorrelation(t *testing.T) {
	Convey("Should return JSON string with correlation", t, func() {
		result := GetCorrelation("gdp")
		So(result, ShouldNotBeNil)
	})
}

func TestSelectRandomValidColumn(t *testing.T) {
	cols := FetchTableCols("popu")
	column := SelectRandomValidColumn(cols)
	Convey("Should return JSON string with correlation", t, func() {
		So(column, ShouldNotBeBlank)
	})
}

// func TestAFunkySituation(t *testing.T) {
// 	cols := FetchTableCols("761e568dd9534b4671eed5dcbd94a6da64e65083cea29e40b9b6a051bde")
// 	fmt.Println("222222", cols)
// 	column := SelectRandomValidColumn(cols)
// 	fmt.Println("222222", column)
// 	Convey("WHOWIDDYNAT", t, func() {
// 		So(column, ShouldNotBeBlank)
// 	})
// }
