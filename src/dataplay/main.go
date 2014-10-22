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
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/cors"
	"log"
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
	fmt.Println("[init] starting in Classic mode")

	e := DBSetup()
	if e != nil {
		fmt.Sprintf("[database] Unable to connect to the Database: %s\n", e)
	}

	/* Database connection will be closed only when Server closes */
	defer DB.Close()

	m := initAPI()

	m.Get("/api/chart/:tablename/:tablenum/:type/:x/:y", GetChartHttp)
	m.Get("/api/chart/:tablename/:tablenum/:type/:x/:y/:z", GetChartHttp)
	m.Get("/api/chartcorrelated/:cid", GetChartCorrelatedHttp)
	m.Get("/api/correlated/:tablename/:search", GetCorrelatedChartsHttp)
	m.Get("/api/correlated/:tablename/:search/:offset/:count", GetCorrelatedChartsHttp)
	m.Get("/api/discovered/:tablename/:correlated", GetDiscoveredChartsHttp)
	m.Get("/api/discovered/:tablename/:correlated/:offset/:count", GetDiscoveredChartsHttp)
	m.Get("/api/getdata/:id", DumpTableHttp)
	m.Get("/api/getdata/:id/:offset/:count", DumpTableHttp)
	m.Get("/api/getdata/:id/:x/:startx/:endx", DumpTableRangeHttp)
	m.Get("/api/getdatagrouped/:id/:x/:y", DumpTableGroupedHttp)
	m.Get("/api/getdatapred/:id/:x/:y", DumpTablePredictionHttp)
	m.Get("/api/getreduceddata/:id", DumpReducedTableHttp)
	m.Get("/api/getreduceddata/:id/:percent", DumpReducedTableHttp)
	m.Get("/api/getreduceddata/:id/:percent/:min", DumpReducedTableHttp)
	m.Get("/api/getreduceddata/:id/:percent/:min/:x/:y", DumpReducedTableHttp)
	m.Get("/api/news/search/:terms", SearchForNewsHttp)
	m.Get("/api/observations/:did", GetObservationsHttp)
	m.Get("/api/political/:type", GetPoliticalActivityHttp)
	m.Get("/api/related/:tablename", GetRelatedChartsHttp)
	m.Get("/api/related/:tablename/:offset/:count", GetRelatedChartsHttp)
	m.Get("/api/relatedstrings/:guid", GetRelatedDatasetByStrings)
	m.Get("/api/search/:keyword", SearchForDataHttp)
	m.Get("/api/search/:keyword/:offset", SearchForDataHttp)
	m.Get("/api/search/:keyword/:offset/:count", SearchForDataHttp)
	m.Get("/api/user/activitystream", GetActivityStreamHttp)
	m.Get("/api/visited", GetLastVisitedHttp)

	m.Use(JsonApiHandler)

	m.Use(ApiSessionHandler)

	m.Run()
}

func initMasterMode() {
	fmt.Println("[init] starting in Master mode")

	e := DBSetup()
	if e != nil {
		fmt.Sprintf("[database] Unable to connect to the Database: %s\n", e)
	}

	/* Database connection will be closed only when Server closes */
	defer DB.Close()

	m := initAPI()

	m.Get("/api/chart/:tablename/:tablenum/:type/:x/:y", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/chart/:tablename/:tablenum/:type/:x/:y", "GetChartQ")
	})
	m.Get("/api/chart/:tablename/:tablenum/:type/:x/:y/:z", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/chart/:tablename/:tablenum/:type/:x/:y/:z", "GetChartQ")
	})
	m.Get("/api/chartcorrelated/:cid", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/chartcorrelated/:cid", "GetChartCorrelatedQ")
	})
	m.Get("/api/correlated/:tablename/:search", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/correlated/:tablename/:search", "GetCorrelatedChartsQ")
	})
	m.Get("/api/correlated/:tablename/:search/:offset/:count", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/correlated/:tablename/:search/:offset/:count", "GetCorrelatedChartsQ")
	})
	m.Get("/api/discovered/:tablename/:correlated", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/discovered/:tablename/:correlated", "GetDiscoveredChartsQ")
	})
	m.Get("/api/discovered/:tablename/:correlated/:offset/:count", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/discovered/:tablename/:correlated/:offset/:count", "GetDiscoveredChartsQ")
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
	m.Get("/api/getreduceddata/:id/:percent/:min/:x/:y", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/getreduceddata/:id/:percent/:min/:x/:y", "DumpReducedTableQ")
	})
	m.Get("/api/news/search/:terms", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/news/search/:terms", "SearchForNewsQ")
	})
	m.Get("/api/observations/:did", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/observations/:did", "GetObservationsQ")
	})
	m.Get("/api/political/:type", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/political/:type", "GetPoliticalActivityQ")
	})
	m.Get("/api/related/:tablename", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/related/:tablename", "GetRelatedChartsQ")
	})
	m.Get("/api/related/:tablename/:offset/:count", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/related/:tablename/:offset/:count", "GetRelatedChartsQ")
	})
	m.Get("/api/search/:keyword", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/search/:keyword", "SearchForDataQ")
	})
	m.Get("/api/search/:keyword/:offset", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/search/:keyword/:offset", "SearchForDataQ")
	})
	m.Get("/api/search/:keyword/:offset/:count", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/search/:keyword/:offset/:count", "SearchForDataQ")
	})
	m.Get("/api/user/activitystream", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/user/activitystream", "GetActivityStreamQ")
	})
	m.Get("/api/visited", func(res http.ResponseWriter, req *http.Request, params martini.Params) string {
		return sendToQueue(res, req, params, "/api/visited", "GetLastVisitedQ")
	})

	m.Use(JsonApiHandler)

	m.Use(LogRequest)

	m.Use(ApiSessionHandler)

	m.Run()
}

var myfuncs funcs

func initNodeMode() {
	fmt.Println("[init] starting in Node mode")

	e := DBSetup()
	if e != nil {
		fmt.Sprintf("[database] Unable to connect to the Database: %s\n", e)
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
	myfuncs.registerCallback("GetChartQ", GetChartQ)
	myfuncs.registerCallback("GetChartCorrelatedQ", GetChartCorrelatedQ)
	myfuncs.registerCallback("GetRelatedChartsQ", GetRelatedChartsQ)
	myfuncs.registerCallback("GetCorrelatedChartsQ", GetCorrelatedChartsQ)
	myfuncs.registerCallback("GetDiscoveredChartsQ", GetDiscoveredChartsQ)
	myfuncs.registerCallback("GetObservationsQ", GetObservationsQ)
	myfuncs.registerCallback("GetPoliticalActivityQ", GetPoliticalActivityQ)
	myfuncs.registerCallback("GetActivityStreamQ", GetActivityStreamQ)
	myfuncs.registerCallback("SearchForNewsQ", SearchForNewsQ)

	/* Start request consumer in Listen mode */
	consumer := QueueConsumer{}
	consumer.Consume()
}

func initAPI() *martini.ClassicMartini { // initialise martini and add in common methods
	m := martini.Classic()

	m.Use(cors.Allow(&cors.Options{
		AllowAllOrigins: true,
		// AllowOrigins:     []string{"http://localhost:9000"},
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

	m.Get("/api/ping", func(res http.ResponseWriter, req *http.Request) string { return "pong" })

	m.Delete("/api/admin/observations/:id", DeleteObservationHttp)
	m.Delete("/api/logout", HandleLogout)
	m.Delete("/api/logout/:session", HandleLogout)

	m.Get("/api/admin/observations/get/:order/:offset/:count/:flagged", GetObservationsTableHttp)
	m.Get("/api/admin/user/get/:order/:offset/:count", GetUserTableHttp)
	m.Get("/api/chart/awaitingcredit", GetAwaitingCreditHttp)
	m.Get("/api/chart/toprated", GetTopRatedChartsHttp)
	m.Get("/api/chartinfo/:tablename", GetChartInfoHttp)
	m.Get("/api/classifydata/:table/:col", SuggestColType)
	m.Get("/api/correlated/:tablename", GetCorrelatedChartsHttp)
	m.Get("/api/findmatches/:id/:x/:y", AttemptToFindMatches)
	m.Get("/api/getdefaults/:id", GetDefaults)
	m.Get("/api/getimportstatus/:id", CheckImportStatus)
	m.Get("/api/home/data", GetHomePageDataHttp)
	m.Get("/api/identifydata/:id", IdentifyTable)
	m.Get("/api/profile/credited", GetCreditedDiscoveriesHttp)
	m.Get("/api/profile/discoveries", GetDiscoveriesHttp)
	m.Get("/api/profile/observations", GetProfileObservationsHttp)
	m.Get("/api/recentobservations", GetRecentObservationsHttp)
	m.Get("/api/stringmatch/:word", FindStringMatches)
	m.Get("/api/stringmatch/:word/:x", FindStringMatches)
	m.Get("/api/tweets/:searchterms", GetTweetsHttp)
	m.Get("/api/user", GetUserDetails)
	m.Get("/api/user/discoveries", GetAmountDiscoveriesHttp)
	m.Get("/api/user/experts", GetDataExpertsHttp)
	m.Get("/api/user/reputation", GetReputationHttp)
	m.Get("/api/user/reset/:token/:username", HandleResetPasswordCheck)

	m.Post("/api/login", binding.Bind(UserForm{}), func(res http.ResponseWriter, req *http.Request, login UserForm) string {
		return HandleLogin(res, req, login)
	})
	m.Post("/api/register", binding.Bind(UserForm{}), func(res http.ResponseWriter, req *http.Request, login UserForm) string {
		return HandleRegister(res, req, login)
	})
	m.Post("/api/setdefaults/:id", SetDefaults)
	m.Post("/api/login/social", binding.Bind(UserSocialForm{}), func(res http.ResponseWriter, req *http.Request, login UserSocialForm) string {
		return HandleSocialLogin(res, req, login)
	})
	m.Post("/api/observations/flag/:id", FlagObservationHttp)
	m.Post("/api/user/check", binding.Bind(UserNameForm{}), func(res http.ResponseWriter, req *http.Request, username UserNameForm) string {
		return HandleCheckUsername(res, req, username)
	})
	m.Post("/api/user/forgot", binding.Bind(UserNameForm{}), func(res http.ResponseWriter, req *http.Request, username UserNameForm) string {
		return HandleForgotPassword(res, req, username)
	})
	m.Post("/api/visited", binding.Bind(VisitedForm{}), func(res http.ResponseWriter, req *http.Request, visited VisitedForm) string {
		return TrackVisitedHttp(res, req, visited)
	})

	m.Put("/api/admin/user/edit", binding.Bind(UserEdit{}), func(res http.ResponseWriter, req *http.Request, userEdit UserEdit) string {
		return EditUserHttp(res, req, userEdit)
	})
	m.Put("/api/chart", binding.Bind(CreditRequest{}), func(res http.ResponseWriter, req *http.Request, params martini.Params, credit CreditRequest) string {
		return CreditChartHttp(res, req, params, credit)
	})
	m.Put("/api/chart/:credflag", binding.Bind(CreditRequest{}), func(res http.ResponseWriter, req *http.Request, params martini.Params, credit CreditRequest) string {
		return CreditChartHttp(res, req, params, credit)
	})
	m.Put("/api/observations", binding.Bind(ObservationComment{}), func(res http.ResponseWriter, req *http.Request, observation ObservationComment) string {
		return AddObservationHttp(res, req, observation)
	})
	m.Put("/api/observations/:oid", CreditObservationHttp)
	m.Put("/api/observations/:oid/:credflag", CreditObservationHttp)
	m.Put("/api/user", binding.Bind(UserDetailsForm{}), func(res http.ResponseWriter, req *http.Request, user UserDetailsForm) string {
		return UpdateUserDetails(res, req, user)
	})
	m.Put("/api/user/reset/:token", binding.Bind(UserForm{}), func(res http.ResponseWriter, req *http.Request, params martini.Params, user UserForm) string {
		return HandleResetPassword(res, req, params, user)
	})

	return m
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

	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	params["user"] = strconv.Itoa(uid)
	params["session"] = session
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

func ApiSessionHandler(res http.ResponseWriter, req *http.Request) {
	if strings.HasPrefix(req.URL.Path, "/api/admin") {
		AdminApiSessionHandler(res, req)
	} else {
		UserApiSessionHandler(res, req)
	}
}

func AdminApiSessionHandler(res http.ResponseWriter, req *http.Request) {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
	}

	user := User{}
	err1 := DB.Where("uid = ?", uid).First(&user).Error
	if err1 != nil && err1 != gorm.RecordNotFound {
		http.Error(res, "Database query failed (User).", http.StatusInternalServerError)
	} else if err1 == gorm.RecordNotFound {
		http.Error(res, "No such user found!", http.StatusNotFound)
	}

	if user.Usertype != UserTypeAdmin {
		http.Error(res, "User is not authorised to perform this action.", http.StatusUnauthorized)
	}
}

func UserApiSessionHandler(res http.ResponseWriter, req *http.Request) {
	noAuthPaths := map[string]bool{
		"/api/login":          true,
		"/api/register":       true,
		"/api/user/check":     true,
		"/api/user/forgot":    true,
		"/api/user/reset":     true,
		"/api/home/data":      true,
		"/api/chart/toprated": true,
	}

	pathTrimmed := strings.TrimLeft(req.URL.Path, "/")
	path := strings.Split(pathTrimmed, "/")
	pathA := "/"
	if len(path) > 0 {
		pathA = pathA + path[0]

		if len(path) > 1 {
			pathA = pathA + "/" + path[1]
		}
	}
	pathB := pathA
	if len(path) > 2 {
		pathB = pathA + "/" + path[2]
	}

	if pathA == "/" || pathA == "/favicon.ico" {
		res.WriteHeader(http.StatusOK)
	} else if (!noAuthPaths[pathA] && !noAuthPaths[pathB]) && (len(req.Header.Get("X-API-SESSION")) < 64 || req.Header.Get("X-API-SESSION") == "false") {
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

	c.Next() // Execute requested method

	endTime := time.Since(startTime)
	executionTime := endTime.Nanoseconds() / 1000 // nanoseconds/1000 = microsecond (us)

	urlData := strings.Split(req.URL.Path, "/")

	// Send data for storage
	go StoreMonitoringData(urlData[1], urlData[2], req.URL.Path, req.Method, rw.Status(), executionTime)
}

/**
 * @details Error Handler
 *
 * @param error
 * @print error
 */
func check(e error) {
	if e != nil {
		fmt.Println(e.Error())
	}
}
