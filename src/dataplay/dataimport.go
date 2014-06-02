package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"net/http"
)

type CheckImportResponce struct {
	State   string
	Request string
}

// This function checks to see if the data has been imported yet or still is in need of importing
func CheckImportStatus(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	if prams["id"] == "" {
		http.Error(res, "There was no ID request", http.StatusBadRequest)
		return ""
	}

	var count int
	Database.DB.QueryRow("SELECT COUNT(*) FROM `priv_onlinedata` WHERE GUID = ?", prams["id"]).Scan(&count)
	var state string

	if count != 0 { // If we have any hits from that query then we have that dataset in the system
		state = "online"
	} else {
		state = "offline"
	}

	returnobj := CheckImportResponce{
		State:   state,
		Request: prams["id"],
	}
	b, _ := json.Marshal(returnobj)

	return string(b)
}
