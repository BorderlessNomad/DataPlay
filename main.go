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
	bcrypt "code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"github.com/codegangsta/martini"              // Worked at 890a2a52d2e59b007758538f9b845fa0ed7daccb
	"github.com/dre1080/martini-contrib/recovery" // Worked at efb5afbb743444c561125d607ba887554a0b9ee2
	"github.com/mattn/go-session-manager"         // Worked at 02b4822c40b5b3996ebbd8bd747d20587635c41b
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
	m.Get("/viewbookmark/:id", func(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) {
		checkAuth(res, req, monager)
		renderTemplate("public/bookmarked.html", nil, res)
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
	m.Get("/api/getdata/:id/:x/:startx/:endx", api.DumpTableRange)
	m.Get("/api/getdatagrouped/:id/:x/:y", api.DumpTableGrouped)
	m.Get("/api/getdatapred/:id/:x/:y", api.DumpTablePrediction)
	m.Get("/api/getcsvdata/:id/:x/:y", api.GetCSV)
	m.Get("/api/getreduceddata/:id", api.DumpReducedTable)
	m.Get("/api/getreduceddata/:id/:persent", api.DumpReducedTable)
	m.Get("/api/getreduceddata/:id/:persent/:min", api.DumpReducedTable)
	m.Post("/api/setbookmark/", api.SetBookmark)
	m.Get("/api/getbookmark/:id", api.GetBookmark)
	m.Post("/api/setdefaults/:id", api.SetDefaults)
	m.Get("/api/getdefaults/:id", api.GetDefaults)
	m.Get("/api/identifydata/:id", api.IdentifyTable)
	m.Get("/api/findmatches/:id/:x/:y", api.AttemptToFindMatches)
	m.Get("/api/classifydata/:table/:col", api.SuggestColType)
	m.Get("/api/stringmatch/:word", api.FindStringMatches)
	m.Get("/api/stringmatch/:word/:x", api.FindStringMatches)
	m.Get("/api/relatedstrings/:guid", api.GetRelatedDatasetByStrings)
	m.Use(ProabblyAPI)
	m.Use(recovery.Recovery())
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

func HandleLogin(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) {
	database := msql.GetDB()
	session := monager.GetSession(res, req)
	username := req.FormValue("username")
	password := req.FormValue("password")

	rows, e := database.Query("SELECT `password` FROM priv_users where email = ? LIMIT 1", username)
	check(e)
	rows.Next()
	var usrpassword string
	e = rows.Scan(&usrpassword)

	if usrpassword != "" && bcrypt.CompareHashAndPassword([]byte(usrpassword), []byte(password)) == nil {
		var uid int
		e := database.QueryRow("SELECT uid FROM priv_users where email = ? LIMIT 1", username).Scan(&uid)
		check(e)
		session.Value = fmt.Sprintf("%d", uid)
		http.Redirect(res, req, "/", http.StatusFound)
	} else {
		var md5test int
		e := database.QueryRow("SELECT count(*) FROM priv_users where email = ? AND password = MD5( ? ) LIMIT 1", username, password).Scan(&md5test)

		if e == nil {
			if md5test != 0 {
				// Ooooh, We need to upgrade this password!
				pwd, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				if e == nil {
					database.Exec("UPDATE `DataCon`.`priv_users` SET `password`= ? WHERE `email`=?", pwd, username)

					var uid int
					e := database.QueryRow("SELECT uid FROM priv_users where email = ? LIMIT 1", username).Scan(&uid)
					check(e)
					session.Value = fmt.Sprintf("%d", uid)

					http.Redirect(res, req, "/", http.StatusFound)
				}
				http.Redirect(res, req, fmt.Sprintf("/login?failed=3&r=%s", e), http.StatusFound)
			} else {
				http.Redirect(res, req, "/login?failed=1", http.StatusFound) // The user has failed this test as well :sad tuba:
			}
		} else {
			http.Redirect(res, req, "/login?failed=1", http.StatusFound) // Ditto to the above
		}
	}
}

func HandleRegister(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) string {
	database := msql.GetDB()
	session := monager.GetSession(res, req)
	username := req.FormValue("username")
	password := req.FormValue("password")

	rows, e := database.Query("SELECT COUNT(*) FROM priv_users where email = ? LIMIT 1", username)
	check(e)
	rows.Next()
	var doesusrexist int
	e = rows.Scan(&doesusrexist)

	if doesusrexist == 0 {
		pwd, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if e != nil {
			return "The password you entered is invalid."
		}
		r, e := database.Exec("INSERT INTO `DataCon`.`priv_users` (`email`, `password`) VALUES (?, ?);", username, pwd)
		if e != nil {
			return "Could not make the user you requested."
		}
		newid, _ := r.LastInsertId()
		session.Value = fmt.Sprintf("%d", newid)
		http.Redirect(res, req, "/", http.StatusFound)
		return ""
	} else {
		return "That username is already registered."
	}
	return ""
}

func checkAuth(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) {
	session := monager.GetSession(res, req)
	if !(session.Value != nil) {
		http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
		return
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
