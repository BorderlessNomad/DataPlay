package main

import (
	"fmt"
	"github.com/codegangsta/martini"
)

func main() {
	fmt.Println("DataCon Server")
	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Use(martini.Static("../static/"))
	m.Run()
}
