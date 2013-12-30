package main

// Hi there. You can find the SQL layout in "layout.sql"
// You can get the data from http://data.gov.uk/data/dumps
// Not actually sure they want you to do that. But it works for now.
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
	"net/url"
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
	manager.SetTimeout(120 * 60)

	what := msql.GetDB()
	what.Ping()

	_, e := what.Exec("SHOW TABLES")
	check(e)
	fmt.Println("DataCon Server")
	initTemplates()
	m := martini.Classic()
	m.Map(manager)
	m.Get("/", func(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) { // res and req are injected by Martini
		checkAuth(res, req, monager)
		session := monager.GetSession(res, req)
		database := msql.GetDB()
		defer database.Close()
		var uid string
		uid = fmt.Sprint(session.Value)
		var username string
		database.QueryRow("select email from priv_users where uid = ?", uid).Scan(&username)
		custom := map[string]string{
			"username": username,
		}
		renderTemplate("public/home.html", custom, res)
		//ApplyTemplate("public/home.html", username, res)
	})
	m.Get("/login", func(res http.ResponseWriter, req *http.Request) {
		failedstr := ""
		queryprams, _ := url.ParseQuery(req.URL.String())
		if queryprams.Get("/login?failed") != "" {
			failedstr = "Incorrect User Name or Password"
		}
		custom := map[string]string{
			"fail": failedstr,
		}
		renderTemplate("public/signin.html", custom, res)
		//ApplyTemplate("public/signin.html", failedstr, res)
	})
	m.Get("/charts/:id", func(res http.ResponseWriter, req *http.Request, prams martini.Params, monager *session.SessionManager) {
		checkAuth(res, req, monager)
		session := monager.GetSession(res, req)
		api.TrackVisited(prams["id"], session.Value.(string))
		renderTemplate("public/charts.html", nil, res)
		//http.ServeFile(res, req, "public/charts.html")
	})
	m.Get("/viewbookmark/:id", func(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) {
		checkAuth(res, req, monager)
		renderTemplate("public/bookmarked.html", nil, res)
		//http.ServeFile(res, req, "public/bookmarked.html")
	})
	m.Get("/search/overlay", func(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) {
		checkAuth(res, req, monager)
		renderTemplate("public/search.html", nil, res)
		//http.ServeFile(res, req, "public/search.html")
	})
	m.Get("/overlay/:id", func(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) {
		checkAuth(res, req, monager)
		renderTemplate("public/overlay.html", nil, res)
		//http.ServeFile(res, req, "public/overlay.html")
	})
	m.Get("/grid/:id", func(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) {
		checkAuth(res, req, monager)
		renderTemplate("public/grid.html", nil, res)
		//http.ServeFile(res, req, "public/grid.html")
	})
	m.Get("/overview/:id", func(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) {
		checkAuth(res, req, monager)
		renderTemplate("public/overview.html", nil, res)
		//http.ServeFile(res, req, "public/overview.html")
	})
	m.Get("/search", func(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) {
		checkAuth(res, req, monager)
		renderTemplate("public/search.html", nil, res)
	})

	m.Post("/noauth/login.json", HandleLogin)
	m.Get("/api/user", api.CheckAuth)
	m.Get("/api/visited", api.GetLastVisited)
	m.Get("/api/search/:s", api.SearchForData)
	m.Get("/api/getinfo/:id", api.GetEntry)
	m.Get("/api/getimportstatus/:id", api.CheckImportStatus)
	m.Get("/api/getdata/:id", api.DumpTable)
	m.Get("/api/getdata/:id/:x/:startx/:endx", api.DumpTableRange)
	m.Get("/api/getcsvdata/:id/:x/:y", api.GetCSV)
	m.Get("/api/getreduceddata/:id", api.DumpReducedTable)
	m.Post("/api/setbookmark/", api.SetBookmark)
	m.Get("/api/getbookmark/:id", api.GetBookmark)
	m.Post("/api/setdefaults/:id", api.SetDefaults)
	m.Get("/api/getdefaults/:id", api.GetDefaults)
	m.Get("/api/identifydata/:id", api.IdentifyTable)
	m.Get("/api/classifydata/:table/:col", api.SuggestColType)
	//m.Use(checkAuth)
	m.Use(ProabblyAPI)
	m.Use(martini.Static("node_modules")) 
	m.Run()
}

func ProabblyAPI(res http.ResponseWriter, req *http.Request) { // Forces anything with /api to have a json doctype. Since it makes sence to
	if strings.HasPrefix(req.RequestURI, "/api") {
		res.Header().Set("Content-Type", "application/json")
	}
}

func HandleLogin(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) {
	database := msql.GetDB()
	session := monager.GetSession(res, req)
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
		http.Redirect(res, req, "/", http.StatusFound)
	} else {
		http.Redirect(res, req, "/login?failed=1", http.StatusFound)
	}
}

// func checkAuth(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) {
// 	if !strings.HasPrefix(req.RequestURI, "/login") && !strings.HasPrefix(req.RequestURI, "/assets") && !strings.HasPrefix(req.RequestURI, "/lib") && !strings.HasPrefix(req.RequestURI, "/noauth") {
// 		session := monager.GetSession(res, req)

// 		if session.Value != nil {
// 		} else {
// 			http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
// 		}
// 	} 
// }

func checkAuth(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) {
	session := monager.GetSession(res, req)
	if session.Value == nil {
		http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
