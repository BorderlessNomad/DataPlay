package main

import (
	"fmt"
	"github.com/codegangsta/martini"
	"net/http"
	"net/url"
)

func Authorisation(res http.ResponseWriter, req *http.Request) { // res and req are injected by Martini
	CheckAuthRedirect(res, req)

	user := User{}
	err := DB.Where("uid = ?", GetUserID(res, req)).Find(&user).Error
	if err != nil {
		panic(err)
	}

	custom := map[string]string{
		"username": user.Email,
	}

	RenderTemplate("public/home.html", custom, res)
	return
}

func Login(res http.ResponseWriter, req *http.Request) {
	failedstr := ""
	queryprams, _ := url.ParseQuery(req.URL.String())

	if queryprams.Get("/login?failed") != "" {
		failedstr = "Incorrect User Name or Password" // They are wrong
		if queryprams.Get("/login?failed") == "2" {
			failedstr = "Your password has been upgraded, please login again." // This should not show anymore, we auto redirect
		} else if queryprams.Get("/login?failed") == "3" {
			failedstr = "Failed to login you in, Sorry!" // somehting went wrong in password upgrade.
		}
	}

	custom := map[string]string{
		"fail": failedstr,
	}

	RenderTemplate("public/signin.html", custom, res)
	return
}

func Logout(res http.ResponseWriter, req *http.Request) {
	HandleLogout(res, req)

	failedstr := ""
	custom := map[string]string{
		"fail": failedstr,
	}

	RenderTemplate("public/signin.html", custom, res)
	return
}

func Register(res http.ResponseWriter, req *http.Request) {
	failedstr := ""
	custom := map[string]string{
		"fail": failedstr,
	}

	RenderTemplate("public/register.html", custom, res)
	return
}

func Charts(res http.ResponseWriter, req *http.Request, prams martini.Params) {
	CheckAuthRedirect(res, req)

	if IsUserLoggedIn(res, req) {
		TrackVisited(prams["id"], fmt.Sprint(GetUserID(res, req))) // Make sure the tracking module knows about their visit.
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

func Overview(res http.ResponseWriter, req *http.Request, prams martini.Params) {
	CheckAuthRedirect(res, req)
	if IsUserLoggedIn(res, req) {
		TrackVisited(prams["id"], fmt.Sprint(GetUserID(res, req)))
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
