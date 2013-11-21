package main

import (
	api "./api"
	msql "./databasefuncs"
	// "database/sql"
	"fmt"
	"github.com/codegangsta/martini"
	// _ "github.com/go-sql-driver/mysql"
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
	manager = session.NewSessionManager(logger)
	manager.SetPath("/") // Thanks for telling me this is needed?

	manager.OnStart(func(session *session.Session) {
		println("started new session", session.Id, session.Value)
	})
	manager.SetTimeout(120)

	what := msql.GetDB()
	what.Ping()

	_, e := what.Exec("SHOW TABLES")
	check(e)
	fmt.Println("DataCon Server")
	m := martini.Classic()
	m.Map(manager)
	m.Get("/", func(res http.ResponseWriter, req *http.Request) { // res and req are injected by Martini
		http.ServeFile(res, req, "public/index.html")
	})
	m.Get("/login", func(res http.ResponseWriter, req *http.Request) { // res and req are injected by Martini
		http.ServeFile(res, req, "public/signin.html")
	})
	m.Post("/noauth/login.json", HandleLogin)
	m.Get("/testtest", api.CheckAuth)
	m.Use(checkAuth)
	m.Run()
}

func ProabblyAPI(res http.ResponseWriter, req *http.Request) { // Forces anything with /api to have a json doctype. Since it makes sence to
	if strings.HasPrefix(req.RequestURI, "/api") {
		res.Header().Set("Content-Type", "application/json")
	}
}

func HandleLogin(res http.ResponseWriter, req *http.Request) {
	database := msql.GetDB()
	session := manager.GetSession(res, req)
	username := req.FormValue("username")
	password := req.FormValue("password")
	fmt.Println()
	rows, e := database.Query("SELECT COUNT(*) as count FROM priv_users where email = ? and password = MD5( ? ) LIMIT 1", username, password)
	check(e)
	rows.Next()
	var count int
	e = rows.Scan(&count)
	fmt.Println("SQL user", count)
	if count != 0 {
		var uid int
		e := database.QueryRow("SELECT uid FROM priv_users where email = ? and password = MD5( ? ) LIMIT 1", username, password).Scan(&uid)
		check(e)
		session.Value = fmt.Sprintf("%d", uid)
		// fmt.Println("hi", session.Id, session.Value, session)
		http.Redirect(res, req, "/?1=1", http.StatusFound)
	}
}

func checkAuth(res http.ResponseWriter, req *http.Request) {
	if req.RequestURI != "/login" && !strings.HasPrefix(req.RequestURI, "/assets") && !strings.HasPrefix(req.RequestURI, "/noauth") {
		session := manager.GetSession(res, req)

		if session.Value != nil {
			fmt.Println("The users sesion val is", session.Value, "so I am going let them pass")
		} else {
			fmt.Println("The users sesion val is", session.Value, "so I am going to send them to the login page")
			http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
