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

func TestExtractColumnWithExpression(t *testing.T) {
	name := "houseprices"
	cols := FetchTableCols(name)
	result := ExtractColumnWithExpression("date", name, cols)
	Convey("Should return date column", t, func() {
		So(result, ShouldNotBeNil)
	})
}
