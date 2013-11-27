package api

import (
	msql "../databasefuncs"
	"github.com/codegangsta/martini"
	"net/http"
)

func SetDefaults(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	// Okay this is suppose to be a POST request that will spit back out a UUID

	if prams["id"] == "" {
		http.Error(res, "You didnt give me a id to store for", http.StatusBadRequest)
		return ""
	}

	database := msql.GetDB()
	defer database.Close()
	jsondata := req.FormValue("data")

	database.Exec("UPDATE `priv_onlinedata` SET `Defaults`=? WHERE  `GUID`=?", jsondata, prams["id"])
	return "{\"result\":\"OK\"}"
}

func GetDefaults(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	database := msql.GetDB()
	defer database.Close()

	if prams["id"] == "" {
		http.Error(res, "You didnt give me a id to lookup", http.StatusBadRequest)
		return ""
	}

	rows := database.QueryRow("SELECT Defaults FROM `priv_onlinedata` WHERE GUID = ?", prams["id"])
	var output string
	rows.Scan(&output)
	return output
}
