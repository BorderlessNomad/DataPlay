package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"net/http"
)

type CheckDataQualityResponce struct {
	Amount  int
	Request string
}

// This is a old legacy function that should not be used anymore, can be removed at any time
// anyone wants to take responsibility to fix the broken ness ( I don't expect there to be lots of it )
func CheckDataQuality(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	// 3 is the amount that is considered by the client as "high quality"
	returnobj := CheckDataQualityResponce{
		Amount:  3,
		Request: prams["id"],
	}
	b, _ := json.Marshal(returnobj)
	return string(b)
}
