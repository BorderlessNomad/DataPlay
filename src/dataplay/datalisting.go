package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/jinzhu/gorm"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Authresponse struct {
	Username string
	UserID   int
}

//This function is used to gather what is the username is
// This used to be used on the front page but now it is mainly used as a "noop" call to check if the user is logged in or not.
func CheckAuth(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	user := User{}
	err := DB.Where("uid = ?", GetUserID(res, req)).Find(&user).Error
	check(err)

	result := Authresponse{
		Username: user.Email,
		UserID:   user.Uid,
	}

	b, _ := json.Marshal(result)

	return string(b)
}

type SearchResult struct {
	Title        string
	GUID         string
	LocationData string
}

// This is the search function that is called though the API
func SearchForData(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	if prams["s"] == "" {
		http.Error(res, "There was no search request", http.StatusBadRequest)
		return ""
	}

	uid := GetUserID(res, req)
	Results := make([]SearchResult, 0)

	indices := []Index{}

	term := prams["s"] + "%" // e.g. "nhs" => "nhs%" (What about "%nhs"?)

	Logger.Println("Searching with Backward Wildcard", term)
	err := DB.Where("LOWER(title) LIKE LOWER(?)", term).Where("(owner = 0 OR owner = ?)", uid).Limit(10).Find(&indices).Error

	Results = ProcessSearchResults(indices, err)
	if len(Results) == 0 {
		term := "%" + prams["s"] + "%" // e.g. "nhs" => "%nhs%"

		Logger.Println("Searching with Forward + Backward Wildcard", term)
		err := DB.Where("LOWER(title) LIKE LOWER(?)", term).Where("(owner = 0 OR owner = ?)", uid).Limit(10).Find(&indices).Error
		Results = ProcessSearchResults(indices, err)
		if len(Results) == 0 {
			term := "%" + strings.Replace(prams["s"], " ", "%", -1) + "%" // e.g. "nh s" => "%nh%s%"

			Logger.Println("Searching with Forward + Backward + Trim Wildcard", term)

			err := DB.Where("LOWER(title) LIKE LOWER(?)", term).Where("(owner = 0 OR owner = ?)", uid).Limit(10).Find(&indices).Error
			Results = ProcessSearchResults(indices, err)

			if len(Results) == 0 && (len(prams["s"]) >= 3 && len(prams["s"]) < 20) {
				term := "%" + prams["s"] + "%" // e.g. "nhs" => "%nhs%"

				Logger.Println("Searching with Forward + Backward Wildcard in String Table", term)

				query := DB.Table("priv_stringsearch, priv_onlinedata, index")
				query = query.Select("DISTINCT ON (priv_onlinedata.guid) priv_onlinedata.guid, index.title")
				query = query.Where("(LOWER(value) LIKE LOWER(?) OR LOWER(x) LIKE LOWER(?))", term, term)
				query = query.Where("priv_stringsearch.tablename = priv_onlinedata.tablename")
				query = query.Where("priv_onlinedata.guid = index.guid")
				query = query.Where("(owner = ? OR owner = ?)", 0, uid)
				query = query.Order("priv_onlinedata.guid")
				query = query.Order("priv_stringsearch.count DESC")
				query = query.Limit(10)
				err := query.Find(&indices).Error

				Results = ProcessSearchResults(indices, err)
			}
		}
	}

	b, _ := json.Marshal(Results)

	return string(b)
}

func ProcessSearchResults(rows []Index, e error) []SearchResult {
	if e != nil && e != gorm.RecordNotFound {
		check(e)
	}

	Results := make([]SearchResult, 0)

	for _, row := range rows {
		Location := HasTableGotLocationData(row.Guid)

		result := SearchResult{
			Title:        row.Title,
			GUID:         row.Guid,
			LocationData: Location,
		}

		Results = append(Results, result)
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

// This function gets the extended infomation FROM the index, things like the notes are used
// in the "wiki" section of the page.
func GetEntry(res http.ResponseWriter, req *http.Request, prams martini.Params) string {
	if prams["id"] == "" {
		http.Error(res, "There was no ID request", http.StatusBadRequest)
		return ""
	}

	index := Index{}
	err := DB.Where("LOWER(guid) LIKE LOWER(?)", prams["id"]+"%").Find(&index).Error
	if err == gorm.RecordNotFound {
		return "[]"
	} else if err != nil {
		panic(err)
		http.Error(res, "Could not find that data.", http.StatusNotFound)
		return ""
	}

	result := DataEntry{
		GUID:     index.Guid,
		Name:     index.Name,
		Title:    index.Title,
		Notes:    index.Notes,
		Ckan_url: strings.Replace(index.CkanUrl, "//", "/", -1),
	}

	b, _ := json.Marshal(result)

	return string(b)
}

// This function casts everything into what it /Should/ Be
// But due to a obscureity in mysql / go / DB.SQL\sql
// everything wants to be a []byte. So I just cast them to that
// then make them strings.
func ScanRow(values []interface{}, columns []string) map[string]interface{} {
	record := make(map[string]interface{})

	for i, col := range values {
		if col != nil {
			switch t := col.(type) {
			default:
				Logger.Printf("Unexpected type %T\n", t)
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

type Dataresponse struct {
	Results []interface{}
	Name    string
}

// This function will empty a whole table out into JSON
// Due to what seems to be a golang bug, everything is outputted as a string.
func DumpTable(res http.ResponseWriter, req *http.Request, prams martini.Params) {
	if prams["id"] == "" {
		http.Error(res, "Sorry! Could not complete this request (Hint, You didnt ask for a table to be dumped)", http.StatusBadRequest)
		return
	}

	var offset int64 = 0
	var count int64 = 0

	UsingRanges := true
	if prams["offset"] == "" || prams["count"] == "" {
		UsingRanges = false
	} else {
		var oE, cE error
		offset, oE = strconv.ParseInt(prams["offset"], 10, 64)
		count, cE = strconv.ParseInt(prams["count"], 10, 64)

		if oE != nil || cE != nil {
			http.Error(res, "Please give valid numbers for offset and count", http.StatusBadRequest)
			return
		}
	}

	tablename, e := getRealTableName(prams["id"], res)
	if e != nil {
		return
	}

	var rows *sql.Rows
	var err error

	if UsingRanges {
		rows, err = DB.Raw(fmt.Sprintf("SELECT * FROM %s OFFSET %d LIMIT %d", tablename, offset, count)).Rows()
	} else {
		rows, err = DB.Raw(fmt.Sprintf("SELECT * FROM %s", tablename)).Rows()
	}

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
		err = rows.Scan(scanArgs...) // This may look like a typo, But it is infact not. This is what you use for interfaces.
		if err != nil {
			panic(err)
		}
		record := ScanRow(values, columns)
		array = append(array, record)
	}
	s, _ := json.Marshal(array)
	res.Write(s)
	io.WriteString(res, "\n")
}

// This function will empty a whole table out into JSON
// Due to what seems to be a golang bug, everything is outputted as a string.
func DumpTableRange(res http.ResponseWriter, req *http.Request, prams martini.Params) {

	// :id/:x/:startx/:endx

	if prams["id"] == "" {
		http.Error(res, "Sorry! Could not compleate this request (Hint, You didnt ask for a table to be dumped)", http.StatusBadRequest)
		return
	}

	if prams["x"] == "" || prams["startx"] == "" || prams["endx"] == "" {
		http.Error(res, "You did not provide enough infomation to make this kind of request :id/:x/:startx/:endx", http.StatusBadRequest)
		return
	}

	tablename, e := getRealTableName(prams["id"], res)
	if e != nil {
		return
	}

	rows, err := DB.Raw("SELECT * FROM " + tablename).Rows()
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

		xvalue := values[xcol].(int64)

		if xvalue >= startx && xvalue <= endx {
			record := ScanRow(values, columns)
			array = append(array, record)
		}
	}
	s, _ := json.Marshal(array)
	res.Write(s)
	io.WriteString(res, "\n")
}

// This call with use the GROUP BY function in mysql to query and get the sum of things
// This is very useful for things like picharts
// /api/getdatagrouped/:id/:x/:y
func DumpTableGrouped(res http.ResponseWriter, req *http.Request, prams martini.Params) {
	if prams["id"] == "" || prams["x"] == "" || prams["y"] == "" {
		http.Error(res, "You did not provide enough infomation to make this kind of request :id/:x/:y", http.StatusBadRequest)
		return
	}

	tablename, e := getRealTableName(prams["id"], res)
	if e != nil {
		return
	}

	cls := FetchTableCols(prams["id"])
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

	rows, e1 := DB.Raw(fmt.Sprintf("SELECT %[1]s, SUM(%[2]s) AS %[2]s FROM %[3]s GROUP BY %[1]s", prams["x"], prams["y"], tablename)).Rows()
	// You may think the above might have some security downsides, It could but what you
	// are proabs thinking is not true, if a user wants to SQL inject as any of the %s's
	// then the table col name will also have to be the SQLi, and frankly, if a user
	// does that then I have no idea what that user should expect, apart FROM broken queries
	// =
	// This could also be filtered at the import level as a form as "moron detection"
	if e1 != nil {
		panic(e1)
		http.Error(res, "Could not query the data FROM the datastore E1", http.StatusInternalServerError)
		return
	}

	columns, e2 := rows.Columns()
	if e1 != nil || e2 != nil {
		panic(e2)
		http.Error(res, "Could not query the data FROM the datastore E2", http.StatusInternalServerError)
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

		record := ScanRow(values, columns)
		array = append(array, record)
	}

	s, _ := json.Marshal(array)
	res.Write(s)
	io.WriteString(res, "\n")
}

// This call will get a X,Y and a prediction of a value. that is asked for
func DumpTablePrediction(res http.ResponseWriter, req *http.Request, prams martini.Params) {
	// /api/getdatapred/:id/:x/:y

	if prams["id"] == "" || prams["x"] == "" || prams["y"] == "" {
		http.Error(res, "You did not provide enough infomation to make this kind of request :id/:x/:y", http.StatusBadRequest)
		return
	}

	tablename, e := getRealTableName(prams["id"], res)
	if e != nil {
		return
	}

	cls := FetchTableCols(prams["id"])
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

	rows, e1 := DB.Raw(fmt.Sprintf("SELECT %s, %s FROM %s", prams["x"], prams["y"], tablename)).Rows()

	if e1 != nil {
		http.Error(res, "Could not query the data FROM the datastore", http.StatusInternalServerError)
		return
	}

	columns, e2 := rows.Columns()
	if e2 != nil {
		http.Error(res, "Could not query the data FROM the datastore", http.StatusInternalServerError)
		return
	}
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	array := make([]map[string]interface{}, 0)
	xarray := make([]float64, 0)
	yarray := make([]float64, 0)

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			panic(err)
		}

		record := ScanRow(values, columns)
		/*Going to if both things are float's else I can't predict them*/
		f1, e := strconv.ParseFloat(record[columns[0]].(string), 64)
		if e != nil {
			http.Error(res, "Could not parse one of the values into a float, therefore cannot run Poly Prediction over it", http.StatusBadRequest)
			return
		}
		f2, e := strconv.ParseFloat(record[columns[1]].(string), 64)
		if e != nil {
			http.Error(res, "Could not parse one of the values into a float, therefore cannot run Poly Prediction over it", http.StatusBadRequest)
			return
		}
		xarray = append(xarray, f1)
		yarray = append(yarray, f2)
		array = append(array, record)
	}
	wat := GetPolyResults(xarray, yarray)
	s, _ := json.Marshal(wat)
	res.Write(s)
	io.WriteString(res, "\n")
}

// This function will take a share of a table and return it as JSON
// Due to what seems to be a golang bug, everything is outputted as a string.
func DumpReducedTable(res http.ResponseWriter, req *http.Request, prams martini.Params) {
	if prams["id"] == "" {
		http.Error(res, "Sorry! Could not compleate this request (Hint, You didnt ask for a table to be dumped)", http.StatusBadRequest)
		return
	}

	tablename, e := getRealTableName(prams["id"], res)
	if e != nil {
		return
	}

	rows, e1 := DB.Raw("SELECT * FROM " + tablename).Rows()

	if e1 != nil {
		http.Error(res, "Could not read that table", http.StatusInternalServerError)
		return
	}

	columns, e2 := rows.Columns()
	if e2 != nil {
		http.Error(res, "Could not read that table", http.StatusInternalServerError)
		return
	}

	DataLength := 0
	err := DB.Table(tablename).Count(&DataLength).Error

	if err != nil && err != gorm.RecordNotFound {
		check(err)
	}

	RealDL := DataLength
	if prams["percent"] == "" {
		DataLength = DataLength / 25
	} else {
		percent := prams["percent"]
		Divider, e := strconv.ParseInt(percent, 10, 64)
		if e != nil {
			http.Error(res, "Invalid percentage", http.StatusBadRequest)
			return // Halt!
		}

		Temp := (float64(Divider) / 100) * float64(DataLength)

		if Temp < 1 {
			Temp = 1
		}

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
		DataLength = 1 // In the case that the percentage returnes a super small amount, then
		// force it to be 1, and return it all
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
		err := rows.Scan(scanArgs...)

		if err != nil {
			panic(err)
		}

		if RowsScanned%DataLength == 0 {
			record := ScanRow(values, columns)
			array = append(array, record)
		}

		RowsScanned++
	}

	s, _ := json.Marshal(array)
	res.Write(s)
	io.WriteString(res, "\n")
}

/**
 * @brief Converts GUID ('friendly' name) into actual table inside database
 *
 * @param string GUID
 * @param http http.ResponseWriter
 *
 * @return string output, error
 */
func getRealTableName(guid string, res ...http.ResponseWriter) (out string, e error) {
	Data := OnlineData{}
	err := DB.Select("tablename").Where("guid = ?", guid).Find(&Data).Error
	if err == gorm.RecordNotFound {
		if res != nil {
			http.Error(res[0], "Could not find that table", http.StatusNotFound)
		}

		return "", fmt.Errorf("Could not find table")
	}

	return Data.Tablename, err
}
