package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"net/http"
)

type ImportResponce struct {
	State   string
	Request string
}

// This function checks to see if the data has been imported yet or still is in need of importing
func CheckImportStatus(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	if prams["id"] == "" {
		http.Error(res, "There was no ID request", http.StatusBadRequest)
		return ""
	}

	onlinedata := OnlineData{}
	count := 0
	err := DB.Model(&onlinedata).Where("guid = ?", prams["id"]).Count(&count).Error
	check(err)

	state := "offline"
	if count != 0 {
		state = "online"
	}

	result := ImportResponce{
		State:   state,
		Request: prams["id"],
	}

	b, _ := json.Marshal(result)

	return string(b)
}
