package main

import (
	"github.com/codegangsta/martini"
	"github.com/jinzhu/gorm"
	"net/http"
	"net/url"
)

func Authorisation(res http.ResponseWriter, req *http.Request) { // res and req are injected by Martini
	CheckAuthRedirect(res, req)

	uid := GetUserID(res, req)

	user := User{}
	err := DB.Where("uid = ?", uid).Find(&user).Error
	if err != nil && err != gorm.RecordNotFound {
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
	queryparams, _ := url.ParseQuery(req.URL.String())

	if queryparams.Get("/login?failed") != "" {
		failedstr = "Incorrect User Name or Password" // They are wrong
		if queryparams.Get("/login?failed") == "2" {
			failedstr = "Your password has been upgraded, please login again." // This should not show anymore, we auto redirect
		} else if queryparams.Get("/login?failed") == "3" {
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
