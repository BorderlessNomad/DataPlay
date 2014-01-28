package main

import (
	Rserve "../src"
	"fmt"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stderr, "[test] ", log.Ldate)
	fmt.Println("Test")
	a := Rserve.New()
	a.AllowUnknownVersions = true
	e := a.Connect("10.0.0.2", 6311)
	if e != nil {
		logger.Fatal("Cannot connect to the RServer")
	}
	logger.Println("Server vesion is", a.ServerBanner)
	e = a.Eval("1+1")
	logger.Println(e)
}
