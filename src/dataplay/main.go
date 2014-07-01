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
	"strconv"
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
	m.Get("/api/search/:s", SearchForDataHttp)
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

	e := DBSetup()
	if e != nil {
		panic(fmt.Sprintf("[database] Unable to connect to the Database: %s\n", e))
		return
	}

	/* Database connection will be closed only when Server closes */
	defer DB.Close()

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
	m.Get("/api/getinfo/:id", GetEntry)
	m.Get("/api/getimportstatus/:id", CheckImportStatus)
	m.Post("/api/setdefaults/:id", SetDefaults)
	m.Get("/api/identifydata/:id", IdentifyTable)

	m.Get("/api/search/:s", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		sendToQueue("/api/search/:s", "SearchForDataQ", res, req, params)
	})
	// m.Get("/api/getdata/:id", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
	// 	sendToQueue("/api/getdata/:id", "DumpTable", res, req, params)
	// })
	// m.Get("/api/getdata/:id/:offset/:count", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
	// 	sendToQueue("/api/getdata/:id/:offset/:count", "DumpTable", res, req, params)
	// })
	// m.Get("/api/getdata/:id/:x/:startx/:endx", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
	// 	sendToQueue("/api/getdata/:id/:x/:startx/:endx", "DumpTableRange", res, req, params)
	// })
	// m.Get("/api/getdatagrouped/:id/:x/:y", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
	// 	sendToQueue("/api/getdatagrouped/:id/:x/:y", "DumpTableGrouped", res, req, params)
	// })
	// m.Get("/api/getdatapred/:id/:x/:y", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
	// 	sendToQueue("/api/getdatapred/:id/:x/:y", "DumpTablePrediction", res, req, params)
	// })

	m.Get("/api/getreduceddata/:id", DumpReducedTable)          // Q
	m.Get("/api/getreduceddata/:id/:percent", DumpReducedTable) // Q// Q
	m.Get("/api/getreduceddata/:id/:percent/:min", DumpReducedTable)
	m.Get("/api/getreduceddata/:id/:x/:y/:percent/:min", DumpReducedTable) // Q

	m.Get("/api/getdefaults/:id", GetDefaults)                     // Q
	m.Get("/api/findmatches/:id/:x/:y", AttemptToFindMatches)      // Q
	m.Get("/api/classifydata/:table/:col", SuggestColType)         // Q
	m.Get("/api/stringmatch/:word", FindStringMatches)             // Q
	m.Get("/api/stringmatch/:word/:x", FindStringMatches)          // Q
	m.Get("/api/relatedstrings/:guid", GetRelatedDatasetByStrings) // Q

	m.Use(JsonApiHandler)

	m.Use(martini.Static("../node_modules")) //Why?

	m.Run()
}

var myfuncs funcs

func initNodeMode() {
	fmt.Println("[init] starting in Node mode")

	e := DBSetup()
	if e != nil {
		panic(fmt.Sprintf("[database] Unable to connect to the Database: %s\n", e))
		return
	}

	/* Database connection will be closed only when Server closes */
	defer DB.Close()

	// Logic for Node (QueueConsumer)
	myfuncs = make(funcs)
	myfuncs.registerCallback("SearchForDataQ", SearchForDataQ)
	// myfuncs.registerCallback("DumpTable", DumpTable)
	// myfuncs.registerCallback("DumpTableRange", DumpTableRange)
	// myfuncs.registerCallback("DumpTableGrouped", DumpTableGrouped)
	// myfuncs.registerCallback("DumpTablePrediction", DumpTablePrediction)
	consumer := QueueConsumer{}
	consumer.Consume()
}

func sendToQueue(request string, method string, res http.ResponseWriter, req *http.Request, params martini.Params) {
	q := Queue{}
	params["user"] = strconv.Itoa(GetUserID(res, req))
	message := q.Encode(method, params)
	//fmt.Println("Sending request to Queue", request, params, message)
	q.send(message)
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
