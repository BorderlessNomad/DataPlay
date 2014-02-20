package api

import (
	msql "../databasefuncs"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/mattn/go-session-manager"
	"net/http"
	"strconv"
)

type BookmarkInternalMetadata struct {
	owner  int
	parent string
}

func SetBookmark(res http.ResponseWriter, req *http.Request, prams martini.Params, session *session.Session) string {
	// Okay this is suppose to be a POST request that will spit back out a UUID.
	database := msql.GetDB()
	defer database.Close()
	jsondata := req.FormValue("data")
	var ownerid int64
	ownerid, _ = strconv.ParseInt(string(session.Value.(string)), 10, 0)
	BMData := BookmarkInternalMetadata{
		owner:  int(ownerid),
		parent: "test",
	}
	privatejson, _ := json.Marshal(BMData)
	_, e := database.Query("INSERT INTO `priv_shares` (`jsoninfo`,`privateinfo`) VALUES (?);", jsondata, string(privatejson))
	if e != nil {
		http.Error(res, "Could not save bookmark", http.StatusInternalServerError)
		return ""
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
		return ""
	}
	var input int64
	var e error
	input, e = strconv.ParseInt(prams["id"], 10, 32)
	if e != nil {
		http.Error(res, "Thats not a number.", 500)
		return ""
	}
	rows := database.QueryRow("SELECT jsoninfo FROM `priv_shares` WHERE shareid = ?", input)
	var output string
	rows.Scan(&output)
	return output
}
