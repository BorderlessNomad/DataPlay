package api

import (
	// cache "../cache"
	// msql "../databasefuncs"
	"encoding/json"
	// "fmt"
	// goq "github.com/PuerkitoBio/goquery"
	"github.com/codegangsta/martini"
	// "github.com/mattn/go-session-manager"
	"net/http"
	// "strings"
)

type CheckDataQualityResponce struct {
	Amount  int
	Request string
}

func CheckDataQuality(res http.ResponseWriter, req *http.Request, prams martini.Params) string {

	returnobj := CheckDataQualityResponce{
		Amount:  3,
		Request: prams["id"],
	}
	b, _ := json.Marshal(returnobj)
	return string(b[:])
}
