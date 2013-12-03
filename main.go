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
	m := martini.Classic()
	m.Map(manager)
	m.Get("/", func(res http.ResponseWriter, req *http.Request) { // res and req are injected by Martini
		http.ServeFile(res, req, "public/index.html")
	})
	m.Get("/login", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, "public/signin.html")
	})
	m.Get("/view/:id", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, "public/displaydataset.html")
	})
	m.Get("/viewbookmark/:id", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, "public/bookmarked.html")
	})
	m.Get("/search/overlay", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, "public/search.html")
	})
	m.Get("/overlay/:id", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, "public/overlay.html")
	})
	m.Get("/import/:id", func(res http.ResponseWriter, req *http.Request, prams martini.Params) {
		http.ServeFile(res, req, "public/processing.html")

		if prams["id"] == "" {
			http.Error(res, "There was no ID request", http.StatusBadRequest)
		}
		database := msql.GetDB()
		defer database.Close()
		var ckan_url string
		database.QueryRow("SELECT ckan_url FROM `index` WHERE GUID = ? LIMIT 1", prams["id"]).Scan(&ckan_url)
		url := strings.Replace(strings.Replace(ckan_url, "//", "/", -1), "http:/", "http://", 1)
		// Because the data.gov.uk dataset is braindead it adds too many "/"'s on the end of the host,
		// because the system can't cope with that I have to filter them all out.
		api.ImportAllDatasets(url, prams["id"])
	})
	m.Post("/noauth/login.json", HandleLogin)
	m.Get("/api/user", api.CheckAuth)
	m.Get("/api/search/:s", api.SearchForData)
	m.Get("/api/getinfo/:id", api.GetEntry)
	m.Get("/api/getquality/:id", api.CheckDataQuality)
	m.Get("/api/getimportstatus/:id", api.CheckImportStatus)
	m.Get("/api/getdata/:id", api.DumpTable)
	m.Post("/api/setbookmark/", api.SetBookmark)
	m.Get("/api/getbookmark/:id", api.GetBookmark)
	m.Post("/api/setdefaults/:id", api.SetDefaults)
	m.Get("/api/getdefaults/:id", api.GetDefaults)

	m.Use(checkAuth)
	m.Use(ProabblyAPI)
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
	}
}

func checkAuth(res http.ResponseWriter, req *http.Request, monager *session.SessionManager) {
	if req.RequestURI != "/login" && !strings.HasPrefix(req.RequestURI, "/assets") && !strings.HasPrefix(req.RequestURI, "/noauth") {
		session := monager.GetSession(res, req)

		if session.Value != nil {
		} else {
			http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
