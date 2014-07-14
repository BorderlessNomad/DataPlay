package main

import (
	"github.com/codegangsta/martini"
	"github.com/jinzhu/gorm"
	"net/http"
)

// Deprecated
func Authorisation(res http.ResponseWriter, req *http.Request) { // res and req are injected by Martini
	CheckAuthRedirect(res, req)

	user := User{}
	err := DB.Where("uid = ?", GetUserID(res, req)).Find(&user).Error
	if err != nil && err != gorm.RecordNotFound {
		panic(err)
	}

	custom := map[string]string{
		"username": user.Email,
	}

	RenderTemplate("public/home.html", custom, res)
	return
}

func Charts(res http.ResponseWriter, req *http.Request, params martini.Params) {
	CheckAuthRedirect(res, req)

	if IsUserLoggedIn(res, req) {
		TrackVisited(params["id"], GetUserID(res, req)) // Make sure the tracking module knows about their visit.
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
		TrackVisited(params["id"], GetUserID(res, req))
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
