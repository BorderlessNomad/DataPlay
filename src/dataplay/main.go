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
	"os"
	"playgen/database"
	"strings"
)

var Logger *log.Logger = log.New(os.Stdout, "[API] ", log.Lshortfile)

type AuthHandler struct {
	http.Handler
	Users map[string]string
}

var DB database.Database
var isDBSetup bool

func DBSetup() error {
	if isDBSetup {
		return nil
	}
	isDBSetup = true
	DB.Setup()
	DB.ParseEnvironment()
	return DB.Connect()
}

/**
 * @details Application bootstrap
 *
 *   Checks database connection,
 *   Init templates,
 *   Init Martini API
 */
func main() {
	e := DBSetup()
	if e == nil {
		/* Database connection will be closed only when Server closes */
		defer DB.Close()
		fmt.Println("[Init] ---[ Welcome to DataCon Server ]---")
	} else {
		panic(fmt.Sprintf("[database] Unable to connect to the Database: %s\n", e))
		return
	}

	initTemplates() // Load all templates from the fs ready to serve to clients.

	m := martini.Classic()

	m.Get("/", Authorisation)
	m.Get("/login", Login)
	m.Get("/logout", Logout)
	m.Get("/register", Register)
	m.Get("/charts/:id", Charts)
	m.Get("/search/overlay", SearchOverlay)
	m.Get("/overlay/:id", Overlay)
	m.Get("/overview/:id", Overview)
	m.Get("/search", Search)
	m.Get("/maptest/:id", MapTest)
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
	if strings.HasPrefix(req.URL.Path, "/api") {
		CheckAuthRedirect(res, req) // Make everything in the API auth'd
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
