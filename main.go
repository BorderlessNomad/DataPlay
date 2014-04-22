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
	"net/http"
	"net/url"
	"strings"
)

type AuthHandler struct {
	http.Handler
	Users map[string]string
}

func main() {
	what := GetDB()
	what.Ping() // Check that the database is actually there and isnt ~spooking~ around

	_, e := what.Exec("SHOW TABLES") // A null query to test functionaility of the SQL server
	check(e)
	what.Close() // Close down the SQL connection since it does nothing after this.
	fmt.Println("DataCon Server")
	initTemplates() // Load all templates from the fs ready to serve to clients.
	m := martini.Classic()
	m.Get("/", func(res http.ResponseWriter, req *http.Request) { // res and req are injected by Martini
		checkAuth(res, req)
		database := GetDB()
		defer database.Close()
		var uid string
		uid = fmt.Sprint(GetUserID(res, req))
		var username string
		database.QueryRow("select email from priv_users where uid = ?", uid).Scan(&username) // get the user's email so I can bake it into the page I am about to send
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
			TrackVisited(prams["id"], string(GetUserID(res, req))) // Make sure the tracking module knows about their visit.
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
			TrackVisited(prams["id"], string(GetUserID(res, req)))
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
	m.Use(ProabblyAPI)
	m.Use(martini.Static("node_modules"))
	m.Run()
}

// A HTTP middleware that Forces anything with /api
// to have a json doctype. Since it makes sence to
func ProabblyAPI(res http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.RequestURI, "/api") {
		checkAuth(res, req) // Make everything in the API auth'd
		res.Header().Set("Content-Type", "application/json")
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
