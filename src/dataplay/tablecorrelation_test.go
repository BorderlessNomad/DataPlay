package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCalculateCoefficient(t *testing.T) {
	tstM := make(map[string]string)
	tstCd := new(CorrelationData)

	Convey("Should return nothing when passed empty map", t, func() {
		result := CalculateCoefficient(tstM, P, tstCd)
		So(result, ShouldEqual, 0)
	})

	tstM["table1"] = "gold"
	tstM["table2"] = "gold"
	tstM["valCol1"] = "price"
	tstM["valCol2"] = "price"
	tstM["dateCol1"] = "date"
	tstM["dateCol2"] = "date"

	Convey("Should return coefficient value of approx 1 when passed same map for table 1 and 2", t, func() {
		result := CalculateCoefficient(tstM, P, tstCd)
		So(result, ShouldEqual, 1)
	})

	tstM["table3"] = "gold"
	tstM["valCol3"] = "price"
	tstM["dateCol3"] = "date"
	Convey("Should return coefficient value of approx 1 when passed same map for table 1, 2 and 3", t, func() {
		result := CalculateCoefficient(tstM, S, tstCd)
		So(result, ShouldEqual, 0.29639139506866363)
	})
}

func TestVisualCorrelation(t *testing.T) {
	for i := 0; i < 10; i++ {
		Convey("Should return true if Visual Correlation found", t, func() {
			table1 := RandomTableName()
			guid, _ := GetRealTableName(table1)
			colNames := FetchTableCols(guid)
			tstM := map[string]string{
				"table1":   table1,
				"valCol1":  RandomValueColumn(colNames),
				"dateCol1": RandomDateColumn(colNames),
			}
			AttemptCorrelation(tstM, V)
		})
	}
}

func TestPearsonCorrelation(t *testing.T) {
	for i := 0; i < 10; i++ {
		Convey("Should return true if Pearson Correlation found", t, func() {
			table1 := RandomTableName()
			guid, _ := GetRealTableName(table1)
			colNames := FetchTableCols(guid)
			tstM := map[string]string{
				"table1":   table1,
				"valCol1":  RandomValueColumn(colNames),
				"dateCol1": RandomDateColumn(colNames),
			}
			AttemptCorrelation(tstM, P)
		})
	}
}

func TestSpuriousCorrelation(t *testing.T) {
	for i := 0; i < 10; i++ {
		Convey("Should return true if Spurious Correlation found", t, func() {
			table1 := RandomTableName()
			guid, _ := GetRealTableName(table1)
			colNames := FetchTableCols(guid)
			tstM := map[string]string{
				"table1":   table1,
				"valCol1":  RandomValueColumn(colNames),
				"dateCol1": RandomDateColumn(colNames),
			}
			AttemptCorrelation(tstM, S)
		})
	}
}
