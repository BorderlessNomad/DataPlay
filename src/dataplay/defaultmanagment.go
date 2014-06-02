package main

import (
	"github.com/codegangsta/martini"
	"net/http"
)

// The set Defaults function is there to save small bits of data that the client might set.
// things like "the key 'date' is a int" really does need to be stored. Thus these pair of calls
// allow the browser to put data next to the row and get it back with ease.
func SetDefaults(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	if prams["id"] == "" {
		http.Error(res, "You didnt give me a id to store for", http.StatusBadRequest)
		return ""
	}

	jsondata := req.FormValue("data")

	Database.DB.Exec("UPDATE `priv_onlinedata` SET `Defaults`=? WHERE  `GUID`=? LIMIT 1", jsondata, prams["id"])

	return "{\"result\":\"OK\"}"
}

// The GetDefaults function is the retrival function for SetDefaults
func GetDefaults(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	if prams["id"] == "" {
		http.Error(res, "You didnt give me a id to lookup", http.StatusBadRequest)
		return ""
	}

	rows := Database.DB.QueryRow("SELECT Defaults FROM `priv_onlinedata` WHERE GUID = ? LIMIT 1", prams["id"])
	var output string
	rows.Scan(&output)

	return output
}
