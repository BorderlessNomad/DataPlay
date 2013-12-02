package api

import (
	msql "../databasefuncs"
	"fmt"
	"github.com/codegangsta/martini"
	"net/http"
	"strconv"
)

func SetBookmark(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	// Okay this is suppose to be a POST request that will spit back out a UUID
	database := msql.GetDB()
	defer database.Close()
	jsondata := req.FormValue("data")
	_, e := database.Query("INSERT INTO `priv_shares` (`jsoninfo`) VALUES (?);", jsondata)
	if e != nil {
		panic(e)
	}
	r := database.QueryRow("SELECT `shareid` FROM priv_shares ORDER BY `shareid` DESC LIMIT 1")

	var id int
	r.Scan(&id)
	return fmt.Sprintf("%d", id)
}

func GetBookmark(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	database := msql.GetDB()
	defer database.Close()

	if prams["id"] == "" {
		http.Error(res, "You didnt give me a id to lookup", 404)
	}
	var input int64
	var e error
	input, e = strconv.ParseInt(prams["id"], 10, 32)
	if e != nil {
		http.Error(res, "Thats not a number.", 500)
	}
	rows := database.QueryRow("SELECT jsoninfo FROM `priv_shares` WHERE shareid = ?", input)
	var output string
	rows.Scan(&output)
	return output
}
