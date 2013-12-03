package api

import (
	msql "../databasefuncs"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/mattn/go-session-manager"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AuthResponce struct {
	Username string
	UserID   int64
}

func CheckAuth(res http.ResponseWriter, req *http.Request, prams martini.Params, manager *session.SessionManager) string {
	//This function is used to gather what is the username is

	// This is used on the home page where it says "Welcome ..." and then replaces the ... with the user name, its also a nice check to see
	// if you are still logged in or not

	session := manager.GetSession(res, req)
	database := msql.GetDB()
	defer database.Close()

	var uid string
	uid = fmt.Sprint(session.Value)
	intuid, _ := strconv.ParseInt(uid, 10, 16)
	var username string
	database.QueryRow("select email from priv_users where uid = ?", uid).Scan(&username)

	returnobj := AuthResponce{
		Username: username,
		UserID:   intuid,
	}
	b, _ := json.Marshal(returnobj)
	return string(b[:])
}

type SearchResult struct {
	Title string
	GUID  string
}

func SearchForData(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	// The three holy commands to setup HTTP handlers
	// session := manager.GetSession(res, req)
	database := msql.GetDB()
	defer database.Close()
	// End
	if prams["s"] == "" {
		http.Error(res, "There was no search request", http.StatusBadRequest)
	}
	rows, e := database.Query("SELECT GUID,Title FROM `index` WHERE Title LIKE ?", prams["s"]+"%")

	if e != nil {
		panic(e)
	}
	Results := make([]SearchResult, 1)
	for rows.Next() {
		var id string
		var name string

		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}

		SR := SearchResult{
			Title: name,
			GUID:  id,
		}
		Results = append(Results, SR)
	}
	if len(Results) == 1 {
		fmt.Println("falling back to overkill search")
		rows, e := database.Query("SELECT GUID,Title FROM `index` WHERE Title LIKE ? LIMIT 10", "%"+prams["s"]+"%")

		if e != nil {
			panic(e)
		}
		for rows.Next() {
			var id string
			var name string

			err := rows.Scan(&id, &name)
			if err != nil {
				panic(err)
			}

			SR := SearchResult{
				Title: name,
				GUID:  id,
			}
			Results = append(Results, SR)
		}
		if len(Results) == 1 {
			fmt.Println("Going 100 persent mad search")
			query := strings.Replace(prams["s"], " ", "%", -1)
			rows, e := database.Query("SELECT GUID,Title FROM `index` WHERE Title LIKE ? LIMIT 10", "%"+query+"%")

			if e != nil {
				panic(e)
			}
			for rows.Next() {
				var id string
				var name string

				err := rows.Scan(&id, &name)
				if err != nil {
					panic(err)
				}

				SR := SearchResult{
					Title: name,
					GUID:  id,
				}
				Results = append(Results, SR)
			}
		}
	}
	defer rows.Close()
	b, _ := json.Marshal(Results)
	return string(b[:])
}

type DataEntry struct {
	GUID     string
	Name     string
	Title    string
	Notes    string
	Ckan_url string
}

func GetEntry(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	// This function gets the extended infomation from the index, things like the notes are used
	// in the "wiki" section of the page.
	database := msql.GetDB()
	defer database.Close()
	if prams["id"] == "" {
		http.Error(res, "There was no ID request", http.StatusBadRequest)
	}
	var GUID string
	var Name string
	var Title string
	var Notes string
	var ckan_url string
	e := database.QueryRow("SELECT * FROM `index` WHERE GUID LIKE ? LIMIT 10", prams["id"]+"%").Scan(&GUID, &Name, &Title, &Notes, &ckan_url)
	strings.Replace(ckan_url, "//", "/", -1)

	returner := DataEntry{
		GUID:     GUID,
		Name:     Name,
		Title:    Title,
		Notes:    Notes,
		Ckan_url: ckan_url,
	}
	if e != nil {
		panic(e)
	}

	b, _ := json.Marshal(returner)
	return string(b[:])
}

type DataResponce struct {
	Results []interface{}
	Name    string
}

func DumpTable(res http.ResponseWriter, req *http.Request, prams martini.Params) {
	// This function will empty a whole table out into JSON
	// Due to what seems to be a golang bug, everything is outputted as a string.

	if prams["id"] == "" {
		http.Error(res, "u wot (Hint, You didnt ask for a table to be dumped)", http.StatusBadRequest)
	}
	database := msql.GetDB()
	defer database.Close()

	var tablename string
	database.QueryRow("SELECT TableName FROM `priv_onlinedata` WHERE GUID = ? LIMIT 1", prams["id"]).Scan(&tablename)
	if tablename == "" {
		http.Error(res, "Could not find that table", http.StatusNotFound)
		return
	}
	rows, err := database.Query("SELECT * FROM " + tablename)
	if err != nil {
		panic(err)
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	array := make([]map[string]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err)
		}

		record := make(map[string]interface{})

		for i, col := range values {
			if col != nil {

				switch t := col.(type) {
				default:
					fmt.Printf("Unexpected type %T\n", t)
				case bool:
					record[columns[i]] = col.(bool)
				case int:
					record[columns[i]] = col.(int)
				case int64:
					record[columns[i]] = col.(int64)
				case float64:
					record[columns[i]] = col.(float64)
				case string:
					record[columns[i]] = col.(string)
				case []byte: // -- all cases go HERE!
					record[columns[i]] = string(col.([]byte))
				case time.Time:
				}
			}
		}
		array = append(array, record)
	}
	s, _ := json.Marshal(array)
	res.Write(s)
	io.WriteString(res, "\n")
}

func DumpReducedTable(res http.ResponseWriter, req *http.Request, prams martini.Params) {
	// This function will empty a whole table out into JSON
	// Due to what seems to be a golang bug, everything is outputted as a string.

	if prams["id"] == "" {
		http.Error(res, "u wot (Hint, You didnt ask for a table to be dumped)", http.StatusBadRequest)
	}
	database := msql.GetDB()
	defer database.Close()

	var tablename string
	database.QueryRow("SELECT TableName FROM `priv_onlinedata` WHERE GUID = ? LIMIT 1", prams["id"]).Scan(&tablename)
	if tablename == "" {
		http.Error(res, "Could not find that table", http.StatusNotFound)
		return
	}
	rows, err := database.Query("SELECT * FROM " + tablename)
	if err != nil {
		panic(err)
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	var DataLength int
	database.QueryRow("SELECT COUNT(*) FROM " + tablename).Scan(&DataLength)
	DataLength = DataLength / 25
	if DataLength < 1 {
		DataLength = 1
	}
	var RowsScanned int
	RowsScanned = 0
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	array := make([]map[string]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err)
		}
		if RowsScanned%DataLength == 0 {
			record := make(map[string]interface{})

			for i, col := range values {
				if col != nil {

					switch t := col.(type) {
					default:
						fmt.Printf("Unexpected type %T\n", t)
					case bool:
						record[columns[i]] = col.(bool)
					case int:
						record[columns[i]] = col.(int)
					case int64:
						record[columns[i]] = col.(int64)
					case float64:
						record[columns[i]] = col.(float64)
					case string:
						record[columns[i]] = col.(string)
					case []byte: // -- all cases go HERE!
						record[columns[i]] = string(col.([]byte))
					case time.Time:
					}
				}
			}
			array = append(array, record)
		}
		RowsScanned++
	}
	s, _ := json.Marshal(array)
	res.Write(s)
	io.WriteString(res, "\n")
}
