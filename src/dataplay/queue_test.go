package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRunMethod(t *testing.T) {

	Convey("When calling method by string name", t, func() {
		r := RunMethod("RunMethodTestFunction", 2, 3, 4)
		result := r[0].Int()
		So(result, ShouldEqual, 24)
	})
}

func (q *Queue) RunMethodTestFunction(x int, y int, z int) int {
	return x * y * z
}
