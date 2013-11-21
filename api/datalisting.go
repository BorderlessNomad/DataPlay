package api

import (
	msql "../databasefuncs"
	// "database/sql"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/mattn/go-session-manager"
	"net/http"
	"strconv"
)

type AuthResponce struct {
	Username string
	UserID   int64
}

func CheckAuth(res http.ResponseWriter, req *http.Request, prams martini.Params, manager *session.SessionManager) string {
	//This function is used to gather what is the username is

	// The three holy commands to setup HTTP handlers
	session := manager.GetSession(res, req)
	database := msql.GetDB()
	defer database.Close()
	// End

	var uid string
	uid = fmt.Sprint(session.Value)
	intuid, _ := strconv.ParseInt(uid, 10, 16)
	var username string
	database.QueryRow("select email from priv_users where uid = ?", uid).Scan(&username)

	returnobj := AuthResponce{
		Username: username,
		UserID:   intuid,
	}
	b, _ := json.Marshal(returnobj)
	return string(b[:])
}

type SearchResult struct {
	Title string
	GUID  string
}

func SearchForData(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	// The three holy commands to setup HTTP handlers
	// session := manager.GetSession(res, req)
	database := msql.GetDB()
	defer database.Close()
	// End
	if prams["s"] == "" {
		http.Error(res, "There was no search request", http.StatusBadRequest)
	}
	rows, e := database.Query("SELECT GUID,Title FROM `index` WHERE Title LIKE ? LIMIT 10", prams["s"]+"%")

	if e != nil {
		panic(e)
	}
	Results := make([]SearchResult, 1)
	for rows.Next() {
		var id string
		var name string

		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}

		SR := SearchResult{
			Title: name,
			GUID:  id,
		}
		Results = append(Results, SR)
	}

	defer rows.Close()
	b, _ := json.Marshal(Results)
	return string(b[:])
}

type DataEntry struct {
	GUID     string
	Name     string
	Title    string
	Notes    string
	Ckan_url string
}

func GetEntry(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	database := msql.GetDB()
	defer database.Close()
	if prams["id"] == "" {
		http.Error(res, "There was no ID request", http.StatusBadRequest)
	}
	var GUID string
	var Name string
	var Title string
	var Notes string
	var ckan_url string
	e := database.QueryRow("SELECT * FROM `index` WHERE GUID LIKE ? LIMIT 10", prams["id"]+"%").Scan(&GUID, &Name, &Title, &Notes, &ckan_url)

	returner := DataEntry{
		GUID:     GUID,
		Name:     Name,
		Title:    Title,
		Notes:    Notes,
		Ckan_url: ckan_url,
	}
	if e != nil {
		panic(e)
	}

	b, _ := json.Marshal(returner)
	return string(b[:])
}
