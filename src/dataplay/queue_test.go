package main

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
)

func TestQueueDecode(t *testing.T) {
	//register test function in function map first
	myfuncs = make(funcs)
	myfuncs.registerCallback("QueueTestFunction", QueueTestFunction)
	fmt.Println("myfuncs", myfuncs)

	q := Queue{}
	str := q.Encode("QueueTestFunction", map[string]string{"X": "2", "Y": "3", "Z": "4"})
	bstr := []byte(str)

	Convey("Should decode json string back into Message object and run method by name with passed params", t, func() {
		msg := q.Decode(bstr)
		So(msg, ShouldEqual, "24")
	})
}

func QueueTestFunction(params map[string]string) string {
	ix, _ := strconv.Atoi(params["X"])
	iy, _ := strconv.Atoi(params["Y"])
	iz, _ := strconv.Atoi(params["Z"])
	a := ix * iy * iz
	return strconv.Itoa(a)
}

func TestQueueTestFunction(t *testing.T) {
	params := make(map[string]string)
	result := RunMethodByName("nosuchfunction", params)
	Convey("Should fail to find function", t, func() {
		So(result, ShouldEqual, "[]")
	})
}
