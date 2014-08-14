package main

// import (
// 	. "github.com/smartystreets/goconvey/convey"
// 	"testing"
// )

// func TestRanking(t *testing.T) {
// 	n := 1 // id

// 	cor := Correlation{} // reset all
// 	err := DB.First(&cor, 1).Update("valid", 0).Error
// 	check(err)
// 	err = DB.First(&cor, 1).Update("invalid", 0).Error
// 	check(err)
// 	err = DB.First(&cor, 1).Update("rating", 0.0).Error
// 	check(err)

// 	for i := 0; i < 23; i++ {
// 		Valid(n)
// 	}

// 	for i := 0; i < 15; i++ {
// 		Invalid(n)
// 	}

// 	Convey("Should return ranking", t, func() {
// 		result := Ranking(n)
// 		So(result, ShouldEqual, 0.44717586998695963)
// 	})

// }
