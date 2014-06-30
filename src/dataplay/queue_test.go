package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
)

func TestQueueDecode(t *testing.T) {
	str := QueueEncode("TestFunction", map[string]string{"ValueX": "2", "ValueY": "3", "ValueZ": "4"})

	Convey("Should decode json string back into Message object and run method by name with passed params", t, func() {
		msg := QueueDecode(str)
		So(msg, ShouldEqual, "24")
	})

}

func (q *Queue) TestFunction(x string, y string, z string) string {
	ix, _ := strconv.Atoi(x)
	iy, _ := strconv.Atoi(y)
	iz, _ := strconv.Atoi(z)

	a := ix * iy * iz
	return strconv.Itoa(a)
}
