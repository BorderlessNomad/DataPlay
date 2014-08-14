package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRankValidations(t *testing.T) {
	Convey("Should return ranking", t, func() {
		result := RankValidations(23, 15)
		So(result, ShouldEqual, 0.44717586998695963)
	})
}
