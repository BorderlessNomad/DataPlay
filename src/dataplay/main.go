package main

/**
 * Hi there. You can find the SQL layout in "layout.sql"
 * You can get the data from http://data.gov.uk/data/dumps
 *
 * For those who want to build this in the future (Hi future!)
 * This was written in Go 1.2.2/1.3
 */

import (
	"flag"
	"fmt"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
	"os"
	"playgen/database"
	"strings"
)

var (
	mode = flag.Int("mode", 1, "1=Node (default), 2=Master, 3=Standalone")
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

func init() {
	flag.Parse()
}

/**
 * @details Application bootstrap
 *
 *   Checks database connection,
 *   init templates,
 *   init Martini API
 */
func main() {
	/**
	 * This application run in 3 types of mode
	 *
	 * 1 = Node (default): Acts as a simple compute instance i.e. no API,
	 * template handling also not exposed to Public. It continuously listens for
	 * incoming requests by means of QueueConsumer. Multiple instances of this mode
	 * can be spawned and killed depending on overall Queue lenght and load, latency etc on system
	 * as whole.
	 *
	 * 2 = Master: APIs are exposed to public and so does everything else.
	 * However only minor calculations are performed by the machine itself (federated system)
	 * major computations are passed to Queue Manager (QueueProducer). Ideally only single instance
	 * of Master should be running unless we configure application to handle load balancing and distributed
	 * Queue (Channels).
	 *
	 * 3 = Single: This is for development purpose where we don't utilize Queue management and instead
	 * evething runs in a single Box (VM).
	 */

	if *mode == 3 {
		initClassicMode()
	} else if *mode == 2 {
		initMasterMode()
	} else {
		initNodeMode()
	}
}

func initClassicMode() {
	fmt.Println("[init] starting in Classic mode")

	e := DBSetup()
	if e != nil {
		panic(fmt.Sprintf("[database] Unable to connect to the Database: %s\n", e))
		return
	}

	/* Database connection will be closed only when Server closes */
	defer DB.Close()

	// // MigrateColumns() // @FUTURE DO NOT UNCOMMENT THIS LINE

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
	m.Get("/api/getdata/:id/:offset/:count", DumpTable)
	m.Get("/api/getdata/:id/:x/:startx/:endx", DumpTableRange)
	m.Get("/api/getdatagrouped/:id/:x/:y", DumpTableGrouped)
	m.Get("/api/getdatapred/:id/:x/:y", DumpTablePrediction)
	m.Get("/api/getreduceddata/:id", DumpReducedTable)
	m.Get("/api/getreduceddata/:id/:percent", DumpReducedTable)
	m.Get("/api/getreduceddata/:id/:percent/:min", DumpReducedTable)
	m.Get("/api/getreduceddata/:id/:x/:y/:percent/:min", DumpReducedTable)
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

func initMasterMode() {
	fmt.Println("[init] starting in Master mode")
	// Logic for Master (QueueProducer)
}

func initNodeMode() {
	fmt.Println("[init] starting in Node mode")
	// Logic for Node (QueueConsumer)
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
