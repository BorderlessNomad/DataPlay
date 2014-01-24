package main

import (
	Rserve "../src"
	"fmt"
)

func main() {
	fmt.Println("Test")
	a := Rserve.New()
	a.Connect()
}
