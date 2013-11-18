package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/martini"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func main() {
	fmt.Println("DataCon Server")
	m := martini.Classic()
	// m.Use("/", martini.Static("public/index.html"))
	m.Get("/", func(res http.ResponseWriter, req *http.Request) { // res and req are injected by Martini
		what := setupDatabase()
		what.Ping()
		http.ServeFile(res, req, "public/index.html")
		// res.WriteHeader(200) // HTTP 200
	})
	m.Run()
}

func setupDatabase() *sql.DB {
	fmt.Println("HI I AM A DATABASE!")
	con, err := sql.Open("mysql", "root:@10.0.0.2/DataCon")
	defer con.Close()
	con.Ping()
	con.Exec("SHOW TABLES", err)

	if err != nil {
		return con
	} else {
		panic(err)
		fmt.Println("OH DEAR!")
		return con
	}
}

func checkAuth(res http.ResponseWriter, req *http.Request) string {
	cookie, _ := req.Cookie("session")
	return cookie.String()
}
