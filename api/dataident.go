package api

import (
	msql "../databasefuncs"
	// "fmt"
	"github.com/codegangsta/martini"
	"net/http"
	// "strconv"
)

// type IdentifyResponce struct {
// 	State   string
// 	Request string
// }

func IdentifyTable(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	// This function checks to see if the data has been imported yet or still is in need of importing
	database := msql.GetDB()
	defer database.Close()
	if prams["id"] == "" {
		http.Error(res, "There was no ID request", http.StatusBadRequest)
	}

	var tablename string
	database.QueryRow("SELECT TableName FROM `priv_onlinedata` WHERE GUID = ? LIMIT 1", prams["id"]).Scan(&tablename)
	if tablename == "" {
		http.Error(res, "Could not find that table", http.StatusNotFound)
		return ""
	}

	var createcode string
	database.QueryRow("SHOW CREATE TABLE "+tablename).Scan(&tablename, &createcode)
	return createcode

	// returnobj := CheckImportResponce{
	// 	State:   state,
	// 	Request: prams["id"],
	// }
	// b, _ := json.Marshal(returnobj)
	// return string(b[:])
}

// Okay so, The best case is that the system has already tagged the
// data with the correct SQL feilds, So at first we will try and ident
// that and see what happens, We can find this below:

// // CREATE TABLE `hips` (
// //   `Hospital` varchar(50) NOT NULL,
// //   `60t69` int(11) DEFAULT NULL,
// //   `70t79` int(11) DEFAULT NULL,
// //   `80t89` int(11) DEFAULT NULL,
// //   `90p` int(11) DEFAULT NULL,
// //   PRIMARY KEY (`Hospital`)
// // ) ENGINE=InnoDB DEFAULT CHARSET=latin1
