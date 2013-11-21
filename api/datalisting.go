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
	//This function is used to gather what is the
	session := manager.GetSession(res, req)
	database := msql.GetDB()
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
	res.Header().Set("Content-Type", "application/json")
	return string(b[:])
	// return session.Cookie()
}
