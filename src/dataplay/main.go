package main

// Hi there. You can find the SQL layout in "layout.sql"
// You can get the data from http://data.gov.uk/data/dumps
// Not actually sure they want you to do that. But it works for now.
// For those who want to build this in the future (Hi future!)
// This was written in Go 1.1.1 or 1.2
// You will also need to run "go get" and hope to god the packages
// still exist.
import (
	"fmt"
	"github.com/codegangsta/martini" // Worked at 890a2a52d2e59b007758538f9b845fa0ed7daccb
	"log"
	"net/http"
	"net/url"
	"os"
	"playgen/database"
	"strings"
)

var Logger *log.Logger = log.New(os.Stdout, "[API] ", log.Lshortfile)

type AuthHandler struct {
	http.Handler
	Users map[string]string
}

var Database struct {
	database.Database
	enabled bool
}

/**
 * @details Application bootstrap
 *
 *   Checks database connection,
 *   Init templates,
 *   Init Martini API
 */
func main() {
	Database.Setup("playgen", "aDam3ntiUm", "10.0.0.2", 5432, "dataplay")
	Database.ParseEnvironment()
	e := Database.Connect()
	if e == nil {
		/* Database connection will be closed only when Server closes */
		defer Database.DB.Close()
		fmt.Println("[Init] ---[ Welcome to DataCon Server ]---")
	} else {
		panic(fmt.Sprintf("[database] Unable to connect to the Database: %s\n", e))
		return
	}

	initTemplates() // Load all templates from the fs ready to serve to clients.

	m := martini.Classic()

	m.Get("/", func(res http.ResponseWriter, req *http.Request) { // res and req are injected by Martini
		checkAuth(res, req)

		var uid, username string
		uid = fmt.Sprint(GetUserID(res, req))

		Database.DB.QueryRow("SELECT email FROM priv_users WHERE uid = $1", uid).Scan(&username) // get the user's email so I can bake it into the page I am about to send
		custom := map[string]string{
			"username": username,
		}

		renderTemplate("public/home.html", custom, res)
	})

	m.Get("/login", func(res http.ResponseWriter, req *http.Request) {
		failedstr := ""
		queryprams, _ := url.ParseQuery(req.URL.String())
		if queryprams.Get("/login?failed") != "" {
			failedstr = "Incorrect User Name or Password" // They are wrong
			if queryprams.Get("/login?failed") == "2" {
				failedstr = "You're password has been upgraded, please login again." // This should not show anymore, we auto redirect
			} else if queryprams.Get("/login?failed") == "3" {
				failedstr = "Failed to login you in, Sorry!" // somehting went wrong in password upgrade.
			}
		}

		custom := map[string]string{
			"fail": failedstr,
		}

		renderTemplate("public/signin.html", custom, res)
	})

	m.Get("/logout", func(res http.ResponseWriter, req *http.Request) {
		HandleLogout(res, req)

		failedstr := ""
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

	m.Get("/charts/:id", func(res http.ResponseWriter, req *http.Request, prams martini.Params) {
		checkAuth(res, req)

		if IsUserLoggedIn(res, req) {
			TrackVisited(prams["id"], fmt.Sprint(GetUserID(res, req))) // Make sure the tracking module knows about their visit.
		}

		renderTemplate("public/charts.html", nil, res)
	})

	m.Get("/search/overlay", func(res http.ResponseWriter, req *http.Request) {
		checkAuth(res, req)

		renderTemplate("public/search.html", nil, res)
	})

	m.Get("/overlay/:id", func(res http.ResponseWriter, req *http.Request) {
		checkAuth(res, req)

		renderTemplate("public/overlay.html", nil, res)
	})

	m.Get("/overview/:id", func(res http.ResponseWriter, req *http.Request, prams martini.Params) {
		checkAuth(res, req)
		if IsUserLoggedIn(res, req) {
			TrackVisited(prams["id"], fmt.Sprint(GetUserID(res, req)))
		}

		renderTemplate("public/overview.html", nil, res)
	})

	m.Get("/search", func(res http.ResponseWriter, req *http.Request) {
		checkAuth(res, req)

		renderTemplate("public/search.html", nil, res)
	})

	m.Get("/maptest/:id", func(res http.ResponseWriter, req *http.Request) {
		checkAuth(res, req)

		renderTemplate("public/maptest.html", nil, res)
	})

	m.Post("/noauth/login.json", HandleLogin)
	m.Post("/noauth/logout.json", HandleLogout)
	m.Post("/noauth/register.json", HandleRegister)
	m.Get("/api/user", CheckAuth)
	m.Get("/api/visited", GetLastVisited)
	m.Get("/api/search/:s", SearchForData)
	m.Get("/api/getinfo/:id", GetEntry)
	m.Get("/api/getimportstatus/:id", CheckImportStatus)
	m.Get("/api/getdata/:id", DumpTable)
	m.Get("/api/getdata/:id/:top/:bot", DumpTable)
	m.Get("/api/getdata/:id/:x/:startx/:endx", DumpTableRange)
	m.Get("/api/getdatagrouped/:id/:x/:y", DumpTableGrouped)
	m.Get("/api/getdatapred/:id/:x/:y", DumpTablePrediction)
	m.Get("/api/getcsvdata/:id/:x/:y", GetCSV)
	m.Get("/api/getreduceddata/:id", DumpReducedTable)
	m.Get("/api/getreduceddata/:id/:persent", DumpReducedTable)
	m.Get("/api/getreduceddata/:id/:persent/:min", DumpReducedTable)
	m.Post("/api/setdefaults/:id", SetDefaults)
	m.Get("/api/getdefaults/:id", GetDefaults)
	m.Get("/api/identifydata/:id", IdentifyTable)
	m.Get("/api/findmatches/:id/:x/:y", AttemptToFindMatches)
	m.Get("/api/classifydata/:table/:col", SuggestColType)
	m.Get("/api/stringmatch/:word", FindStringMatches)
	m.Get("/api/stringmatch/:word/:x", FindStringMatches)
	m.Get("/api/relatedstrings/:guid", GetRelatedDatasetByStrings)

	m.Use(JsonApiHandler)

	m.Use(martini.Static("../node_modules")) //Why?

	m.Run()
}

/**
 * @details A HTTP middleware that Forces anything with /api to have a json doctype. Since it makes sence to
 *
 * @param http.ResponseWriter
 * @param *http.Request
 */
func JsonApiHandler(res http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.RequestURI, "/api") {
		checkAuth(res, req) // Make everything in the API auth'd
		res.Header().Set("Content-Type", "application/json")
	}
}

/**
 * @details Error Handler
 *
 * @param error
 * @return panic
 */
func check(e error) {
	if e != nil {
		panic(e)
	}
}