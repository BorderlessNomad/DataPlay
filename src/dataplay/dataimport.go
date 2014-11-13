package main

// import (
// 	"encoding/json"
// 	"github.com/codegangsta/martini"
// 	"net/http"
// )

// type ImportResponse struct {
// 	State   string
// 	Request string
// }

// // This function checks to see if the data has been imported yet or still is in need of importing
// func CheckImportStatus(res http.ResponseWriter, req *http.Request, params martini.Params) string {
// 	if params["id"] == "" {
// 		http.Error(res, "There was no ID request", http.StatusBadRequest)
// 		return ""
// 	}

// 	onlinedata := OnlineData{}
// 	count := 0
// 	err := DB.Model(&onlinedata).Where("guid = ?", params["id"]).Count(&count).Error
// 	check(err)

// 	state := "offline"
// 	if count != 0 {
// 		state = "online"
// 	}

// 	result := ImportResponse{
// 		State:   state,
// 		Request: params["id"],
// 	}

// 	b, _ := json.Marshal(result)

// 	return string(b)
// }
