package main

import (
	"github.com/codegangsta/martini"
	"net/http"
)

func Charts(res http.ResponseWriter, req *http.Request, params martini.Params) {
	CheckAuthRedirect(res, req)

	if IsUserLoggedIn(res, req) {
		session := params["session"]
		if len(session) <= 0 {
			http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		}

		_, err := GetUserID(session)
		if err != nil {
			http.Error(res, err.Message, err.Code)
		}

		// TrackVisited(params["id"], uid) // Make sure the tracking module knows about their visit.
	}

	RenderTemplate("public/charts.html", nil, res)
	return
}

func SearchOverlay(res http.ResponseWriter, req *http.Request) {
	CheckAuthRedirect(res, req)

	RenderTemplate("public/search.html", nil, res)
	return
}

func Overlay(res http.ResponseWriter, req *http.Request) {
	CheckAuthRedirect(res, req)

	RenderTemplate("public/overlay.html", nil, res)
	return
}

func Overview(res http.ResponseWriter, req *http.Request, params martini.Params) {
	CheckAuthRedirect(res, req)
	if IsUserLoggedIn(res, req) {
		session := params["session"]
		if len(session) <= 0 {
			http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		}

		_, err := GetUserID(session)
		if err != nil {
			http.Error(res, err.Message, err.Code)
		}

		// TrackVisited(params["id"], uid)
	}

	RenderTemplate("public/overview.html", nil, res)
	return
}

func Search(res http.ResponseWriter, req *http.Request) {
	CheckAuthRedirect(res, req)

	RenderTemplate("public/search.html", nil, res)
	return
}

func MapTest(res http.ResponseWriter, req *http.Request) {
	CheckAuthRedirect(res, req)

	RenderTemplate("public/maptest.html", nil, res)
	return
}
