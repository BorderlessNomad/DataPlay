package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"net/http"
)

func main() {
	fmt.Println("DataCon Serverr")
	m := martini.Classic()
	// m.Use("/", martini.Static("public/index.html"))
	m.Get("/", func(res http.ResponseWriter, req *http.Request) { // res and req are injected by Martini
		http.ServeFile(res, req, "public/index.html")
		// res.WriteHeader(200) // HTTP 200
	})
	m.Run()
}
