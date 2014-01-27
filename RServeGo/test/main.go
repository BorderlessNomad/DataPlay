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
	e := a.Connect("192.168.1.122", 6311)
	if e != nil {
		logger.Fatal("Cannot connect to the RServer")
	}
	logger.Println("Server vesion is", a.ServerBanner)
}
