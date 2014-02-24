package main

// Hi there. You can find the SQL layout in "layout.sql"
// You can get the data from http://data.gov.uk/data/dumps
// Not actually sure they want you to do that. But it works for now.
// For those who want to build this in the future (Hi future!)
// This was written in Go 1.1.1 or 1.2
// You will also need to run "go get" and hope to god the packages
// still exist.
import (
	api "./api"
	msql "./databasefuncs"
	"fmt"
	"github.com/codegangsta/martini"      // Worked at 890a2a52d2e59b007758538f9b845fa0ed7daccb
	"github.com/mattn/go-session-manager" // Worked at 02b4822c40b5b3996ebbd8bd747d20587635c41b
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
	logger := log.New(os.Stdout, "[Sessions] ", log.Ldate|log.Ltime)
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
	what.Close() // Close down the SQL connection since it does nothing after this.
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
	})
	m.Get("/login", func(res http.ResponseWriter, req *http.Request) {
		failedstr := ""
		queryprams, _ := url.ParseQuery(req.URL.String())
		if queryprams.Get("/login?failed") != "" {
			failedstr = "Incorrect User Name or Password"
			if queryprams.Get("/login?failed") == "2" {
				failedstr = "You're password has been upgraded, please login again."
			} else if queryprams.Get("/login?failed") == "3" {
				failedstr = "Failed to login you in, Sorry!"
			}
		}
		custom := map[string]string{
			"fail": failedstr,
		}
		renderTemplate("public/signin.html", custom, res)
	})
	m.Get("/register", func(res http.ResponseWriter, req *http.Request) {
		failedstr := ""
		custom := map[string]string{
			"fail": failedstr,
		}
		renderTemplate("public/register.html", custom, res)
	})
	m.Get("/charts/:id", func(res http.ResponseWriter, req *http.Request, prams martini.Params, monager *session.SessionManager) {
		checkAuth(res, req, monager)
		session := monager.GetSession(res, req)
		if session.Value != nil {
			api.TrackVisited(prams["id"], session.Value.(string))
		}
		renderTemplate("public/charts.html", nil, res)
	})
	m.Get("/search/overlay", func(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) {
		checkAuth(res, req, monager)
		renderTemplate("public/search.html", nil, res)
	})
	m.Get("/overlay/:id", func(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) {
		checkAuth(res, req, monager)
		renderTemplate("public/overlay.html", nil, res)
	})
	m.Get("/overview/:id", func(res http.ResponseWriter, req *http.Request, prams martini.Params, monager *session.SessionManager) {
		checkAuth(res, req, monager)
		session := monager.GetSession(res, req)
		if session.Value != nil {
			api.TrackVisited(prams["id"], session.Value.(string))
		}
		renderTemplate("public/overview.html", nil, res)
	})
	m.Get("/search", func(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) {
		checkAuth(res, req, monager)
		renderTemplate("public/search.html", nil, res)
	})
	m.Get("/maptest/:id", func(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) {
		checkAuth(res, req, monager)
		renderTemplate("public/maptest.html", nil, res)
	})

	m.Post("/noauth/login.json", HandleLogin)
	m.Post("/noauth/register.json", HandleRegister)
	m.Get("/api/user", api.CheckAuth)
	m.Get("/api/visited", api.GetLastVisited)
	m.Get("/api/search/:s", api.SearchForData)
	m.Get("/api/getinfo/:id", api.GetEntry)
	m.Get("/api/getimportstatus/:id", api.CheckImportStatus)
	m.Get("/api/getdata/:id", api.DumpTable)
	m.Get("/api/getdata/:id/:top/:bot", api.DumpTable)
	m.Get("/api/getdata/:id/:x/:startx/:endx", api.DumpTableRange)
	m.Get("/api/getdatagrouped/:id/:x/:y", api.DumpTableGrouped)
	m.Get("/api/getdatapred/:id/:x/:y", api.DumpTablePrediction)
	m.Get("/api/getcsvdata/:id/:x/:y", api.GetCSV)
	m.Get("/api/getreduceddata/:id", api.DumpReducedTable)
	m.Get("/api/getreduceddata/:id/:persent", api.DumpReducedTable)
	m.Get("/api/getreduceddata/:id/:persent/:min", api.DumpReducedTable)
	m.Post("/api/setdefaults/:id", api.SetDefaults)
	m.Get("/api/getdefaults/:id", api.GetDefaults)
	m.Get("/api/identifydata/:id", api.IdentifyTable)
	m.Get("/api/findmatches/:id/:x/:y", api.AttemptToFindMatches)
	m.Get("/api/classifydata/:table/:col", api.SuggestColType)
	m.Get("/api/stringmatch/:word", api.FindStringMatches)
	m.Get("/api/stringmatch/:word/:x", api.FindStringMatches)
	m.Get("/api/relatedstrings/:guid", api.GetRelatedDatasetByStrings)
	m.Use(ProabblyAPI)
	m.Use(martini.Static("node_modules"))
	m.Run()
}

func ProabblyAPI(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) {
	// Forces anything with /api to have a json doctype. Since it makes sence to
	if strings.HasPrefix(req.RequestURI, "/api") {
		checkAuth(res, req, monager) // Make everything in the API auth'd
		res.Header().Set("Content-Type", "application/json")
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
