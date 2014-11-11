package main

import (
	"github.com/ahirmayur/gorm"
	"github.com/codegangsta/martini"
	"net/http"
)

// The set Defaults function is there to save small bits of data that the client might set.
// things like "the key 'date' is a int" really does need to be stored. Thus these pair of calls
// allow the browser to put data next to the row and get it back with ease.
func SetDefaults(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	if params["id"] == "" {
		http.Error(res, "You didnt give me a id to store for", http.StatusBadRequest)
		return ""
	}

	jsondata := req.FormValue("data")

	onlinedata := OnlineData{}
	err := DB.Model(&onlinedata).Where("guid = ?", params["id"]).UpdateColumn("defaults", jsondata).Error
	check(err)

	return "{\"result\":\"OK\"}"
}

// The GetDefaults function is the retrival function for SetDefaults
func GetDefaults(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	if params["id"] == "" {
		http.Error(res, "You didnt give me a id to lookup", http.StatusBadRequest)
		return ""
	}

	var d string

	onlinedata := OnlineData{}
	err := DB.Select("defaults").Where("guid = ?", params["id"]).Find(&onlinedata).Error
	if err == gorm.RecordNotFound {
		d = "{}"
	} else if err == nil {
		d = onlinedata.Defaults
	} else {
		return ""
	}

	return d
}
