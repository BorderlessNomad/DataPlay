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
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/martini-contrib/cors"
	"log"
	"math/rand"
	"net/http"
	"os"
	"playgen/database"
	"strconv"
	"strings"
	"time"
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
	rand.Seed(time.Now().Unix())
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

	m.Use(cors.Allow(&cors.Options{
		AllowAllOrigins: true,
		// AllowOrigins:     []string{"http://localhost:9000"},
		// AllowMethods: []string{"PUT", "PATCH"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowCredentials: true,
		AllowHeaders: []string{
			"Origin",
			"Accept",
			"Content-Type",
			"Authorization",
			"Accept-Encoding",
			"Content-Length",
			"Host",
			"Referer",
			"User-Agent",
			"X-CSRF-Token",
			"X-API-SESSION",
		},
	}))

	// m.Get("/", Authorisation)
	/* @todo convert to APIs */
	m.Get("/charts/:id", Charts)
	m.Get("/search/overlay", SearchOverlay)
	m.Get("/overlay/:id", Overlay)
	m.Get("/overview/:id", Overview)
	m.Get("/search", Search)
	m.Get("/maptest/:id", MapTest)

	/* APIs */
	m.Get("/api/ping", func(res http.ResponseWriter, req *http.Request) string {
		return "pong"
	})
	m.Post("/api/login", binding.Bind(UserForm{}), func(res http.ResponseWriter, req *http.Request, login UserForm) string {
		return HandleLogin(res, req, login)
	})
	m.Delete("/api/logout", HandleLogout)
	m.Post("/api/register", binding.Bind(UserForm{}), func(res http.ResponseWriter, req *http.Request, login UserForm) string {
		return HandleRegister(res, req, login)
	})
	m.Post("/api/user/check", binding.Bind(UserNameForm{}), func(res http.ResponseWriter, req *http.Request, username UserNameForm) string {
		return HandleCheckUsername(res, req, username)
	})
	m.Post("/api/user/forgot", binding.Bind(UserNameForm{}), func(res http.ResponseWriter, req *http.Request, username UserNameForm) string {
		return HandleForgotPassword(res, req, username)
	})
	m.Get("/api/user/reset/:token/:username", HandleResetPasswordCheck)
	m.Put("/api/user/reset/:token", binding.Bind(UserForm{}), func(res http.ResponseWriter, req *http.Request, params martini.Params, user UserForm) string {
		return HandleResetPassword(res, req, params, user)
	})

	m.Get("/api/user", GetUserDetails)
	m.Put("/api/user", binding.Bind(UserDetailsForm{}), func(res http.ResponseWriter, req *http.Request, user UserDetailsForm) string {
		return UpdateUserDetails(res, req, user)
	})

	m.Get("/api/visited", GetLastVisitedHttp)
	m.Post("/api/visited", binding.Bind(VisitedForm{}), func(res http.ResponseWriter, req *http.Request, visited VisitedForm) string {
		return TrackVisitedHttp(res, req, visited)
	})

	m.Get("/api/search/:s", SearchForDataHttp)
	m.Get("/api/getinfo/:id", GetEntry)
	m.Get("/api/getimportstatus/:id", CheckImportStatus)
	m.Get("/api/getdata/:id", DumpTableHttp)
	m.Get("/api/getdata/:id/:offset/:count", DumpTableHttp)
	m.Get("/api/getdata/:id/:x/:startx/:endx", DumpTableRangeHttp)
	m.Get("/api/getdatagrouped/:id/:x/:y", DumpTableGroupedHttp)
	m.Get("/api/getdatapred/:id/:x/:y", DumpTablePredictionHttp)
	m.Get("/api/getreduceddata/:id", DumpReducedTableHttp)
	m.Get("/api/getreduceddata/:id/:percent", DumpReducedTableHttp)
	m.Get("/api/getreduceddata/:id/:percent/:min", DumpReducedTableHttp)
	m.Get("/api/getreduceddata/:id/:x/:y/:percent/:min", DumpReducedTableHttp)
	m.Post("/api/setdefaults/:id", SetDefaults)
	m.Get("/api/getdefaults/:id", GetDefaults)
	m.Get("/api/identifydata/:id", IdentifyTable)
	m.Get("/api/findmatches/:id/:x/:y", AttemptToFindMatches)
	m.Get("/api/classifydata/:table/:col", SuggestColType)
	m.Get("/api/stringmatch/:word", FindStringMatches)
	m.Get("/api/stringmatch/:word/:x", FindStringMatches)
	m.Get("/api/relatedstrings/:guid", GetRelatedDatasetByStrings)

	// API v1.1
	m.Get("/api/chart/:tablename/:tablenum/:type/:x/:y", GetChartHttp)
	m.Get("/api/chart/:tablename/:tablenum/:type/:x/:y/:z", GetChartHttp)
	m.Get("/api/chartcorrelated/:cid", GetChartCorrelatedHttp)
	m.Get("/api/related/:tablename", GetRelatedChartsHttp)
	m.Get("/api/related/:tablename/:offset/:count", GetRelatedChartsHttp)
	m.Get("/api/correlated/:tablename", GetCorrelatedChartsHttp)
	m.Get("/api/correlated/:tablename/:searchdepth", GetCorrelatedChartsHttp)
	m.Get("/api/correlated/:tablename/:offset/:count/:searchdepth", GetCorrelatedChartsHttp)
	m.Get("/api/discovered/:tablename/:correlated", GetDiscoveredChartsHttp)
	m.Get("/api/discovered/:tablename/:correlated/:offset/:count", GetDiscoveredChartsHttp)
	m.Put("/api/chart/:rcid", ValidateChartHttp)
	m.Put("/api/chart/:rcid/:valflag", ValidateChartHttp)
	m.Put("/api/observations/:did/:x/:y/:comment", AddObservationHttp)
	m.Put("/api/observations/:oid", ValidateObservationHttp)
	m.Put("/api/observations/:oid/:valflag", ValidateObservationHttp)
	m.Get("/api/observations/:did", GetObservationsHttp)
	m.Get("/api/political/:type", GetPoliticalActivityHttp)
	m.Get("/api/profile/observations", GetProfileObservationsHttp)
	m.Get("/api/profile/discoveries", GetDiscoveriesHttp)
	m.Get("/api/profile/validated", GetValidatedDiscoveriesHttp)
	m.Get("/api/home/data", GetHomePageDataHttp)
	m.Get("/api/user/reputation", GetReputationHttp)
	m.Get("/api/user/discoveries", GetAmountDiscoveriesHttp)

	m.Use(JsonApiHandler)

	m.Use(SessionApiHandler)

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

	// m.Get("/", Authorisation)
	m.Post("/api/login", binding.Bind(UserForm{}), func(res http.ResponseWriter, req *http.Request, login UserForm) string {
		return HandleLogin(res, req, login)
	})
	m.Delete("/api/logout/:session", HandleLogout)
	m.Post("/api/register", binding.Bind(UserForm{}), func(res http.ResponseWriter, req *http.Request, login UserForm) string {
		return HandleRegister(res, req, login)
	})
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
	m.Get("/api/visited", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/visited", "GetLastVisitedQ")
	})
	m.Get("/api/getinfo/:id", GetEntry)
	m.Get("/api/getimportstatus/:id", CheckImportStatus)
	m.Post("/api/setdefaults/:id", SetDefaults)
	m.Get("/api/identifydata/:id", IdentifyTable)

	m.Get("/api/search/:s", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/search/:s", "SearchForDataQ")
	})
	m.Get("/api/getdata/:id", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/getdata/:id", "DumpTableQ")
	})
	m.Get("/api/getdata/:id/:offset/:count", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/getdata/:id/:offset/:count", "DumpTableQ")
	})
	m.Get("/api/getdata/:id/:x/:startx/:endx", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/getdata/:id/:x/:startx/:endx", "DumpTableRangeQ")
	})
	m.Get("/api/getdatagrouped/:id/:x/:y", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/getdatagrouped/:id/:x/:y", "DumpTableGroupedQ")
	})
	m.Get("/api/getdatapred/:id/:x/:y", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/getdatapred/:id/:x/:y", "DumpTablePredictionQ")
	})
	m.Get("/api/getreduceddata/:id", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/getreduceddata/:id", "DumpReducedTableQ")
	})
	m.Get("/api/getreduceddata/:id/:percent", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/getreduceddata/:id/:percent", "DumpReducedTableQ")
	})
	m.Get("/api/getreduceddata/:id/:percent/:min", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/getreduceddata/:id/:percent/:min", "DumpReducedTableQ")
	})
	m.Get("/api/getreduceddata/:id/:x/:y/:percent/:min", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/getreduceddata/:id/:x/:y/:percent/:min", "DumpReducedTableQ")
	})

	m.Get("/api/getdefaults/:id", GetDefaults)                     // Q
	m.Get("/api/findmatches/:id/:x/:y", AttemptToFindMatches)      // Q
	m.Get("/api/classifydata/:table/:col", SuggestColType)         // Q
	m.Get("/api/stringmatch/:word", FindStringMatches)             // Q
	m.Get("/api/stringmatch/:word/:x", FindStringMatches)          // Q
	m.Get("/api/relatedstrings/:guid", GetRelatedDatasetByStrings) // Q

	m.Use(JsonApiHandler)

	m.Use(LogRequest)

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

	defer DB.Close()

	myfuncs = make(funcs)
	myfuncs.registerCallback("GetLastVisitedQ", GetLastVisitedQ)
	myfuncs.registerCallback("SearchForDataQ", SearchForDataQ)
	myfuncs.registerCallback("DumpTableQ", DumpTableQ)
	myfuncs.registerCallback("DumpTableRangeQ", DumpTableRangeQ)
	myfuncs.registerCallback("DumpTableGroupedQ", DumpTableGroupedQ)
	myfuncs.registerCallback("DumpTablePredictionQ", DumpTablePredictionQ)
	myfuncs.registerCallback("DumpReducedTableQ", DumpReducedTableQ)

	/* Start request consumer in Listen mode */
	consumer := QueueConsumer{}
	consumer.Consume()
}

var responseChannel chan string

/**
 * @details Send requet to Queue for remote execution in parallel mode.
 * Request & responses are async however output will be sent to ResponseWriter
 * as soon as it is received via singleton channel.
 */
func sendToQueue(res http.ResponseWriter, req *http.Request, params martini.Params, request string, method string) string {
	responseChannel = make(chan string, 1)

	q := Queue{}
	go q.Response()

	session := params["session"]
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	params["user"] = strconv.Itoa(uid)
	message := q.Encode(method, params)

	fmt.Println("Sending request to Queue", request, params, message)

	go q.send(message)

	return <-responseChannel
}

/**
 * @details A HTTP middleware that Forces anything with /api to have a json doctype.
 *
 * @param http.ResponseWriter
 * @param *http.Request
 */
func JsonApiHandler(res http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.URL.Path, "/api") {
		// CheckAuthRedirect(res, req) // Make everything in the API auth'd
		res.Header().Set("Content-Type", "application/json")
	}
}

func SessionApiHandler(res http.ResponseWriter, req *http.Request) {
	noAuthPaths := map[string]bool{
		"/api/login":       true,
		"/api/register":    true,
		"/api/user/check":  true,
		"/api/user/forgot": true,
		"/api/user/reset":  true,
	}

	pathTrimmed := strings.TrimLeft(req.URL.Path, "/")
	path := strings.Split(pathTrimmed, "/")
	pathA := "/" + path[0] + "/" + path[1]
	pathB := pathA
	if len(path) > 2 {
		pathB = pathA + "/" + path[2]
	}

	if (!noAuthPaths[pathA] && !noAuthPaths[pathB]) && (len(req.Header.Get("X-API-SESSION")) <= 0 || req.Header.Get("X-API-SESSION") == "false") {
		res.WriteHeader(http.StatusUnauthorized)
	}
}

/**
 * @brief Log incoming requests
 * @details Log all requests ending on '/api' for performance monitoring,
 * this method will start a timer before procesing data and will report at
 * the end of execution process. Data is stored in Redis via StoreMonitoringData
 * in sync manner.
 *
 * @param http [description]
 * @param http [description]
 * @param martini [description]
 * @return [description]
 */

func LogRequest(res http.ResponseWriter, req *http.Request, c martini.Context) {
	// Do not proceed if request is not for "/api"
	if !strings.HasPrefix(req.URL.Path, "/api") {
		return
	}

	startTime := time.Now()

	rw := res.(martini.ResponseWriter)

	// tm := TimeMachine(100, 500)
	// time.Sleep(tm * time.Millisecond)

	c.Next() // Execute requested method

	endTime := time.Since(startTime)
	executionTime := endTime.Nanoseconds() / 1000 // nanoseconds/1000 = microsecond (us)

	urlData := strings.Split(req.URL.Path, "/")

	// Send data for storage
	go StoreMonitoringData(urlData[1], urlData[2], req.URL.Path, req.Method, rw.Status(), executionTime)
}

func TimeMachine(min, max int) time.Duration {
	return time.Duration(rand.Intn(max-min) + min)
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
