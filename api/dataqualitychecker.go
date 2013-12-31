package api

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"net/http"
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
