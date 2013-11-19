package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/martini"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strings"
)

func main() {
	what := setupDatabase()
	what.Ping()
	_, e := what.Exec("SHOW TABLES")
	check(e)
	fmt.Println("DataCon Server")
	m := martini.Classic()
	// m.Use("/", martini.Static("public/index.html"))
	m.Get("/", func(res http.ResponseWriter, req *http.Request) { // res and req are injected by Martini
		http.ServeFile(res, req, "public/index.html")
	})
	m.Get("/login", func(res http.ResponseWriter, req *http.Request) { // res and req are injected by Martini
		http.ServeFile(res, req, "public/signin.html")
	})
	m.Post("/noauth/login.json", HandleLogin)
	m.Use(checkAuth)
	m.Run()
}

func HandleLogin(res http.ResponseWriter, req *http.Request) {
	db := setupDatabase()
	res.Write()
}

func setupDatabase() *sql.DB {
	fmt.Println("[database] Asked for MySQL connection")
	con, err := sql.Open("mysql", "root:@tcp(10.0.0.2:3306)/DataCon")
	check(err)
	con.Ping()
	return con
}

func checkAuth(res http.ResponseWriter, req *http.Request) {
	cookie, _ := req.Cookie("session")
	database := setupDatabase()
	// fmt.Println(cookie.String())
	if req.RequestURI != "/login" && !strings.HasPrefix(req.RequestURI, "/assets") && !strings.HasPrefix(req.RequestURI, "/noauth") {
		if cookie.String() != "" {
			rows, e := database.Query("SELECT userid FROM priv_sessions where token = ? LIMIT 1", cookie.String())
			check(e)
			var userid int
			rows.Next()
			e = rows.Scan(&userid)
			if e != nil {
				http.Redirect(res, req, "/login", http.StatusMovedPermanently)
			}
		} else {
			http.Redirect(res, req, "/login", http.StatusMovedPermanently)
		}
	}
	// return cookie.String()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
