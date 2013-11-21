package api

import (
	// msql "../databasefuncs"
	// "database/sql"
	// "fmt"
	"github.com/codegangsta/martini"
	"github.com/mattn/go-session-manager"
	"net/http"
)

type AuthResponce struct {
	Username string
	UserID   string
}

func CheckAuth(res http.ResponseWriter, req *http.Request, prams martini.Params, manager *session.SessionManager) string {
	//This function is used to gather what is the
	session := manager.GetSession(res, req)
	return session.Cookie()
}
