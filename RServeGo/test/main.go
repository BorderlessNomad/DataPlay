package main

import (
	Rserve "../src"
	"fmt"
)

func main() {
	fmt.Println("Test")
	a := Rserve.New()
	e := a.Connect("localhost", 6311)
	fmt.Println(e)
}
