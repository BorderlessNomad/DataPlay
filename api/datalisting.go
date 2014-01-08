package api

import (
	msql "../databasefuncs"
	"database/sql"
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
	Title        string
	GUID         string
	LocationData string
}

func SearchForData(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	// The three holy commands to setup HTTP handlers
	// session := manager.GetSession(res, req)
	database := msql.GetDB()
	defer database.Close()
	// End
	if prams["s"] == "" {
		http.Error(res, "There was no search request", http.StatusBadRequest)
		return ""
	}
	rows, e := database.Query("SELECT GUID,Title FROM `index` WHERE Title LIKE ? LIMIT 10", prams["s"]+"%")

	Results := make([]SearchResult, 1)
	Results = ProcessSearchResults(rows, e, database)
	if len(Results) == 1 {
		fmt.Println("falling back to overkill search")
		rows, e := database.Query("SELECT GUID,Title FROM `index` WHERE Title LIKE ? LIMIT 10", "%"+prams["s"]+"%")
		Results = ProcessSearchResults(rows, e, database)
		if len(Results) == 1 {
			fmt.Println("Going 100 persent mad search")
			query := strings.Replace(prams["s"], " ", "%", -1)
			rows, e := database.Query("SELECT GUID,Title FROM `index` WHERE Title LIKE ? LIMIT 10", "%"+query+"%")
			Results = ProcessSearchResults(rows, e, database)
		}
	}
	defer rows.Close()
	b, _ := json.Marshal(Results)
	return string(b[:])
}

func ProcessSearchResults(rows *sql.Rows, e error, database *sql.DB) []SearchResult {
	Results := make([]SearchResult, 1)
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
		Location := HasTableGotLocationData(id, database)
		SR := SearchResult{
			Title:        name,
			GUID:         id,
			LocationData: Location,
		}
		Results = append(Results, SR)
	}
	return Results
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
		return ""
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
		http.Error(res, "Could not find that data.", http.StatusNotFound)
		return ""
	}

	b, _ := json.Marshal(returner)
	return string(b[:])
}

func scanrow(values []interface{}, columns []string) map[string]interface{} {
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
	return record
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
		return
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
		record := scanrow(values, columns)
		array = append(array, record)
	}
	s, _ := json.Marshal(array)
	res.Write(s)
	io.WriteString(res, "\n")
}

func DumpTableRange(res http.ResponseWriter, req *http.Request, prams martini.Params) {
	// This function will empty a whole table out into JSON
	// Due to what seems to be a golang bug, everything is outputted as a string.

	// :id/:x/:startx/:endx

	if prams["id"] == "" {
		http.Error(res, "u wot (Hint, You didnt ask for a table to be dumped)", http.StatusBadRequest)
		return
	}

	if prams["x"] == "" || prams["startx"] == "" || prams["endx"] == "" {
		http.Error(res, "You did not provide enough infomation to make this kind of request :id/:x/:startx/:endx", http.StatusBadRequest)
		return
	}

	database := msql.GetDB()
	defer database.Close()

	var tablename string
	database.QueryRow("SELECT TableName FROM `priv_onlinedata` WHERE GUID = ? LIMIT 1", prams["id"]).Scan(&tablename)
	if tablename == "" {
		http.Error(res, "Could not find that table", http.StatusNotFound)
		return
	}
	rows, err := database.Query("SELECT * FROM `" + tablename + "`")
	if err != nil {
		panic(err)
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	var xcol int
	xcol = 999
	startx, starte := strconv.ParseInt(prams["startx"], 10, 64)
	endx, ende := strconv.ParseInt(prams["endx"], 10, 64)
	if starte != nil || ende != nil {
		http.Error(res, "You didnt pass me proper numbers to start with.", http.StatusBadRequest)
		return
	}

	for number, colname := range columns {
		if colname == prams["x"] {
			xcol = number
		}
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

		xvalue, e := strconv.ParseInt(string(values[xcol].([]byte)), 10, 0)

		if e != nil {
			http.Error(res, "Read loop error D: Looks like int this is a imposter.", http.StatusInternalServerError)
			return
		}
		if xvalue >= startx && xvalue <= endx {
			record := scanrow(values, columns)
			array = append(array, record)
		}
	}
	s, _ := json.Marshal(array)
	res.Write(s)
	io.WriteString(res, "\n")
}

func DumpTableGrouped(res http.ResponseWriter, req *http.Request, prams martini.Params) {
	// This call with use the GROUP BY function in mysql to query and get the sum of things
	// This is very useful for things like picharts
	// /api/getdatagrouped/:id/:x/:y

	if prams["id"] == "" || prams["x"] == "" || prams["y"] == "" {
		http.Error(res, "You did not provide enough infomation to make this kind of request :id/:x/:y", http.StatusBadRequest)
		return
	}

	database := msql.GetDB()
	defer database.Close()

	var tablename string
	database.QueryRow("SELECT TableName FROM `priv_onlinedata` WHERE GUID = ? LIMIT 1", prams["id"]).Scan(&tablename)
	if tablename == "" {
		http.Error(res, "Could not find that table", http.StatusNotFound)
		return
	}
	cls := FetchTableCols(prams["id"], database)
	// Now we need to check that the rows that the client is asking for, are in the table.
	Valid := false
	for _, clm := range cls {
		if clm.Name == prams["x"] {
			Valid = true
		}
	}
	if !Valid {
		http.Error(res, "Col X is invalid.", http.StatusBadRequest)
		return
	}
	Valid = false
	for _, clm := range cls {
		if clm.Name == prams["y"] {
			Valid = true
		}
	}
	if !Valid {
		http.Error(res, "Col Y is invalid.", http.StatusBadRequest)
		return
	}
	rows, e1 := database.Query(fmt.Sprintf("SELECT `%s`,SUM(%s) AS %s FROM `%s` GROUP BY %s", prams["x"], prams["y"], prams["y"], tablename, prams["x"]))
	columns, e2 := rows.Columns()
	if e1 != nil || e2 != nil {
		http.Error(res, "Could not query the data from the datastore", http.StatusInternalServerError)
		return
	}

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	array := make([]map[string]interface{}, 0)
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			panic(err)
		}

		record := scanrow(values, columns)
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
		return
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
	RealDL := DataLength
	if prams["persent"] == "" {
		DataLength = DataLength / 25
	} else {
		Persent := prams["persent"]
		Divider, e := strconv.ParseInt(Persent, 10, 64)
		if e != nil {
			http.Error(res, "Invalid Persentage", http.StatusBadRequest)
			return // Halt!
		}
		Temp := (float64(Divider) / 100) * float64(DataLength)
		DataLength = DataLength / int(Temp)
		if prams["min"] != "" {
			MinSpend, e := strconv.ParseInt(prams["min"], 10, 64)
			if e != nil {
				http.Error(res, "Invalid Min", http.StatusBadRequest)
				return // Halt!
			}
			if int(RealDL/DataLength) < int(MinSpend) {
				DataLength = RealDL / int(MinSpend)
			}
		}
	}
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
			record := scanrow(values, columns)
			array = append(array, record)
		}
		RowsScanned++
	}
	s, _ := json.Marshal(array)
	res.Write(s)
	io.WriteString(res, "\n")
}

func GetCSV(res http.ResponseWriter, req *http.Request, prams martini.Params) {
	// This function will empty a whole table out into JSON
	// Due to what seems to be a golang bug, everything is outputted as a string.

	if prams["id"] == "" {
		http.Error(res, "u wot (Hint, You didnt ask for a table to be dumped)", http.StatusBadRequest)
		return
	}
	if prams["x"] == "" || prams["y"] == "" {
		http.Error(res, "I don't have a x and y to make the CSV for.", http.StatusBadRequest)
		return
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
	// We need to find the Columns to relay back.

	var xcol int
	var ycol int
	xcol = 999
	ycol = 999
	for number, colname := range columns {
		if colname == prams["x"] {
			xcol = number
		} else if colname == prams["y"] {
			ycol = number
		}
	}
	if xcol == 999 || ycol == 999 {
		http.Error(res, "Could not find some of the columns that you asked for.", http.StatusNotFound)
		return
	}

	var output string
	output = "\"name\",\"word\",\"count\"\n"
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

		output = output + fmt.Sprintf("\"%s\",\"%s\",%s\n", values[xcol], values[xcol], values[ycol])
		record := scanrow(values, columns)
		array = append(array, record)
	}
	// s, _ := json.Marshal(array)
	res.Write([]byte(output))
	// io.WriteString(output, "\n")
}
