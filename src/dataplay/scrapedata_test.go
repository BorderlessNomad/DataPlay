package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDateScrape(t *testing.T) {
	Convey("DateScrape", t, func() {
		res := DateScrapeA()
		DateScrapeB(res)
	})
}
