package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mattn/go-session-manager"
	"log"
	"net/http"
	"os"
	"strings"
)

type AuthHandler struct {
	http.Handler
	Users map[string]string
}

var manager *session.SessionManager

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	manager := session.NewSessionManager(logger)

	manager.OnStart(func(session *session.Session) {
		println("started new session")
	})

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
	// handler := seshcookie.NewSessionHandler(
	// 	&AuthHandler{http.FileServer(contentDir), userDb},
	// 	"session key, preferably a sequence of data from /dev/urandom",
	// 	nil)
	// m.Use(handler)
	m.Use(checkAuth)
	m.Run()
}

func HandleLogin(res http.ResponseWriter, req *http.Request) {
	database := setupDatabase()
	session := manager.GetSession(res, req)

	// res.Write()
	username := req.FormValue("username")
	password := req.FormValue("password")
	rows, e := database.Query("SELECT COUNT(*) as count FROM priv_users where email = ? and password = MD5(?) LIMIT 1", username, password)
	check(e)
	rows.Next()
	var count int
	e = rows.Scan(&count)
	if count != 0 {
		// session := seshcookie.Session.Get(req)
		// loggedin, _ := session["loggedin"].(int)
		// loggedin := 1
		// session["loggedin"] = loggedin
		session.Value = "a"
	}
}

func setupDatabase() *sql.DB {
	fmt.Println("[database] Asked for MySQL connection")
	con, err := sql.Open("mysql", "root:@tcp(10.0.0.2:3306)/DataCon")
	check(err)
	con.Ping()
	return con
}

func checkAuth(res http.ResponseWriter, req *http.Request) {
	// cookie, _ := req.Cookie("session")
	// database := setupDatabase()
	// fmt.Println(cookie.String())
	if req.RequestURI != "/login" && !strings.HasPrefix(req.RequestURI, "/assets") && !strings.HasPrefix(req.RequestURI, "/noauth") {
		session := manager.GetSession(res, req)
		if session.Value != nil {
			http.Redirect(res, req, "/", http.StatusMovedPermanently)
		} else {
			http.Redirect(res, req, "/login", http.StatusMovedPermanently)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
