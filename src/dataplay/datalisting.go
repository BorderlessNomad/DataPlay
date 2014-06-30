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
func CheckAuth(res http.ResponseWriter, req *http.Request, params martini.Params) string {
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

func SearchForDataHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	uid := GetUserID(res, req)

	result, error := SearchForData(params["s"], uid)
	if error != nil {
		http.Error(res, error.Message, error.Code)
		return ""
	}

	r, err := json.Marshal(result)
	if err != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

/**
 * @brief Search a given term in database
 * @details This method searches for a matching title with following conditions,
 * 		Postfix wildcard
 * 		Prefix & postfix wildcard
 * 		Prefix, postfix & trimmed spaces with wildcard
 * 		Prefix & postfix on previously searched terms
 */
func SearchForData(str string, uid int) ([]SearchResult, *appError) {
	Results := make([]SearchResult, 0)

	if str == "" {
		return Results, &appError{nil, "There was no search request", http.StatusBadRequest}
	}

	indices := []Index{}

	term := str + "%" // e.g. "nhs" => "nhs%" (What about "%nhs"?)

	Logger.Println("Searching with Postfix Wildcard", term)

	err := DB.Where("LOWER(title) LIKE LOWER(?)", term).Where("(owner = 0 OR owner = ?)", uid).Limit(10).Find(&indices).Error
	if err != nil && err != gorm.RecordNotFound {
		return Results, &appError{err, "Database query failed", http.StatusServiceUnavailable}
	}

	Results = ProcessSearchResults(indices, err)
	if len(Results) == 0 {
		term := "%" + str + "%" // e.g. "nhs" => "%nhs%"

		Logger.Println("Searching with Prefix + Postfix Wildcard", term)

		err := DB.Where("LOWER(title) LIKE LOWER(?)", term).Where("(owner = 0 OR owner = ?)", uid).Limit(10).Find(&indices).Error
		if err != nil && err != gorm.RecordNotFound {
			return Results, &appError{err, "Database query failed", http.StatusServiceUnavailable}
		}

		Results = ProcessSearchResults(indices, err)
		if len(Results) == 0 {
			term := "%" + strings.Replace(str, " ", "%", -1) + "%" // e.g. "nh s" => "%nh%s%"

			Logger.Println("Searching with Prefix + Postfix + Trim Wildcard", term)

			err := DB.Where("LOWER(title) LIKE LOWER(?)", term).Where("(owner = 0 OR owner = ?)", uid).Limit(10).Find(&indices).Error
			if err != nil && err != gorm.RecordNotFound {
				return Results, &appError{err, "Database query failed", http.StatusServiceUnavailable}
			}

			Results = ProcessSearchResults(indices, err)
			if len(Results) == 0 && (len(str) >= 3 && len(str) < 20) {
				term := "%" + str + "%" // e.g. "nhs" => "%nhs%"

				Logger.Println("Searching with Prefix + Postfix Wildcard in String Table", term)

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

				if err != nil && err != gorm.RecordNotFound {
					return Results, &appError{err, "Database query failed", http.StatusInternalServerError}
				}

				Results = ProcessSearchResults(indices, err)
			}
		}
	}

	return Results, nil
}

func ProcessSearchResults(rows []Index, e error) []SearchResult {
	if e != nil && e != gorm.RecordNotFound {
		check(e)
	}

	Results := make([]SearchResult, 0)

	for _, row := range rows {
		Location := HasTableGotLocationData(row.Guid)

		result := SearchResult{
			Title:        SanitizeString(row.Title),
			GUID:         SanitizeString(row.Guid),
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
func GetEntry(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	if params["id"] == "" {
		http.Error(res, "There was no ID request", http.StatusBadRequest)
		return ""
	}

	index := Index{}
	err := DB.Where("LOWER(guid) LIKE LOWER(?)", params["id"]+"%").Find(&index).Error
	if err == gorm.RecordNotFound {
		return "[]"
	} else if err != nil {
		panic(err)
		http.Error(res, "Could not find that data.", http.StatusNotFound)
		return ""
	}

	result := DataEntry{
		GUID:     index.Guid,
		Name:     SanitizeString(index.Name),
		Title:    SanitizeString(index.Title),
		Notes:    SanitizeString(index.Notes),
		Ckan_url: strings.Replace(index.CkanUrl, "//", "/", -1),
	}

	b, _ := json.Marshal(result)

	return string(b)
}

func SanitizeString(str string) string {
	return strings.Replace(str, "Ã‚Â£", "£", -1)
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
				record[columns[i]] = t.Format("2006-01-02")
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
func DumpTable(res http.ResponseWriter, req *http.Request, params martini.Params) {
	if params["id"] == "" {
		http.Error(res, "Sorry! Could not complete this request (Hint, You didnt ask for a table to be dumped)", http.StatusBadRequest)
		return
	}

	var offset int64 = 0
	var count int64 = 0

	UsingRanges := true
	if params["offset"] == "" || params["count"] == "" {
		UsingRanges = false
	} else {
		var oE, cE error
		offset, oE = strconv.ParseInt(params["offset"], 10, 64)
		count, cE = strconv.ParseInt(params["count"], 10, 64)

		if oE != nil || cE != nil {
			http.Error(res, "Please give valid numbers for offset and count", http.StatusBadRequest)
			return
		}
	}

	tablename, e := getRealTableName(params["id"], res)
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
func DumpTableRange(res http.ResponseWriter, req *http.Request, params martini.Params) {

	// :id/:x/:startx/:endx

	if params["id"] == "" {
		http.Error(res, "Sorry! Could not compleate this request (Hint, You didnt ask for a table to be dumped)", http.StatusBadRequest)
		return
	}

	if params["x"] == "" || params["startx"] == "" || params["endx"] == "" {
		http.Error(res, "You did not provide enough infomation to make this kind of request :id/:x/:startx/:endx", http.StatusBadRequest)
		return
	}

	tablename, e := getRealTableName(params["id"], res)
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
	startx, starte := strconv.ParseInt(params["startx"], 10, 64)
	endx, ende := strconv.ParseInt(params["endx"], 10, 64)
	if starte != nil || ende != nil {
		http.Error(res, "You didnt pass me proper numbers to start with.", http.StatusBadRequest)
		return
	}

	for number, colname := range columns {
		if colname == params["x"] {
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
func DumpTableGrouped(res http.ResponseWriter, req *http.Request, params martini.Params) {
	if params["id"] == "" || params["x"] == "" || params["y"] == "" {
		http.Error(res, "You did not provide enough infomation to make this kind of request :id/:x/:y", http.StatusBadRequest)
		return
	}

	tablename, e := getRealTableName(params["id"], res)
	if e != nil {
		return
	}

	cls := FetchTableCols(params["id"])
	ValidX := false
	ValidY := false
	sumY := false

	/* Check for existence of X & Y in Table */
	for _, clm := range cls {
		if clm.Name == params["x"] {
			ValidX = true
		} else if clm.Name == params["y"] {
			if clm.Sqltype != "varchar" && clm.Sqltype != "date" {
				sumY = true
			}
			ValidY = true
		}

		if ValidX && ValidY {
			break
		}
	}

	if !ValidX {
		http.Error(res, "Col X is invalid.", http.StatusBadRequest)
		return
	}

	if !ValidY {
		http.Error(res, "Col Y is invalid.", http.StatusBadRequest)
		return
	}

	q := ""
	if sumY {
		q = fmt.Sprintf("SELECT %[1]s, SUM(%[2]s) AS %[2]s FROM %[3]s GROUP BY %[1]s", params["x"], params["y"], tablename)
	} else {
		q = fmt.Sprintf("SELECT DISTINCT %[1]s, COUNT(%[1]s) AS %[2]s FROM %[3]s GROUP BY %[1]s ORDER BY %[1]s", params["x"], params["y"], tablename)
	}

	rows, e1 := DB.Raw(q).Rows()
	if e1 != nil {
		check(e1)
		http.Error(res, "Could not query the data FROM the datastore E1", http.StatusInternalServerError)
		return
	}

	columns, e2 := rows.Columns()
	if e1 != nil || e2 != nil {
		check(e2)
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
func DumpTablePrediction(res http.ResponseWriter, req *http.Request, params martini.Params) {
	// /api/getdatapred/:id/:x/:y

	if params["id"] == "" || params["x"] == "" || params["y"] == "" {
		http.Error(res, "You did not provide enough infomation to make this kind of request :id/:x/:y", http.StatusBadRequest)
		return
	}

	tablename, e := getRealTableName(params["id"], res)
	if e != nil {
		return
	}

	cls := FetchTableCols(params["id"])
	// Now we need to check that the rows that the client is asking for, are in the table.
	Valid := false
	for _, clm := range cls {
		if clm.Name == params["x"] {
			Valid = true
		}
	}
	if !Valid {
		http.Error(res, "Col X is invalid.", http.StatusBadRequest)
		return
	}
	Valid = false
	for _, clm := range cls {
		if clm.Name == params["y"] {
			Valid = true
		}
	}
	if !Valid {
		http.Error(res, "Col Y is invalid.", http.StatusBadRequest)
		return
	}

	rows, e1 := DB.Raw(fmt.Sprintf("SELECT %s, %s FROM %s", params["x"], params["y"], tablename)).Rows()

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

/**
 * @brief This function will take a share of a table and return it as JSON
 */
func DumpReducedTable(res http.ResponseWriter, req *http.Request, params martini.Params) {
	if params["id"] == "" {
		http.Error(res, "Sorry! Could not compleate this request (Hint, You didnt ask for a table to be dumped)", http.StatusBadRequest)
		return
	}

	tablename, e := getRealTableName(params["id"], res)
	if e != nil {
		return
	}

	var rows *sql.Rows
	var e1 error
	if params["x"] == "" || params["y"] == "" {
		rows, e1 = DB.Table(tablename).Rows()
	} else {
		x := params["x"]
		y := params["y"]

		cls := FetchTableCols(params["id"])
		sumY := false

		/* Check for existence of X & Y in Table */
		for _, clm := range cls {
			if clm.Name == y && clm.Sqltype != "varchar" && clm.Sqltype != "date" {
				sumY = true
			}
		}

		if sumY {
			rows, e1 = DB.Table(tablename).Select("DISTINCT \"" + x + "\", SUM(\"" + y + "\") AS \"" + y + "\"").Group(x).Order(x).Rows()
		} else {
			rows, e1 = DB.Table(tablename).Select("DISTINCT \"" + x + "\", COUNT(\"" + x + "\") AS \"" + y + "\"").Group(x).Order(x).Rows()
		}
	}

	if e1 != nil {
		panic(e1)
		http.Error(res, "Could not read that table", http.StatusInternalServerError)
		return
	}

	columns, e2 := rows.Columns()
	if e2 != nil {
		http.Error(res, "Could not read that table", http.StatusInternalServerError)
		return
	}

	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	results := make([]map[string]interface{}, 0)
	results_calc := make([]map[string]interface{}, 0)
	for rows.Next() {
		err := rows.Scan(scanArgs...)

		if err != nil {
			panic(err)
		}

		record := ScanRow(values, columns)
		results = append(results, record)
	}

	DataLength := len(results)
	RealDL := DataLength
	if params["percent"] == "" {
		DataLength = DataLength / 25
	} else {
		percent := params["percent"]
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

		if params["min"] != "" {
			MinSpend, e := strconv.ParseInt(params["min"], 10, 64)
			if e != nil {
				http.Error(res, "Invalid Min", http.StatusBadRequest)
				return // Halt!
			}

			if int(RealDL/DataLength) < int(MinSpend) {
				DataLength = RealDL / int(MinSpend)
			}
		}
	}

	/**
	 * In the case that the percentage returns a super small amount,
	 * then force it to be 1, and return it all
	 */
	if DataLength < 1 {
		DataLength = 1
	}

	for i, result := range results {
		if i%DataLength == 0 {
			results_calc = append(results_calc, result)
		}
	}

	s, _ := json.Marshal(results_calc)
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
	data := OnlineData{}
	err := DB.Select("tablename").Where("guid = ?", guid).Find(&data).Error
	if err == gorm.RecordNotFound {
		if res != nil {
			http.Error(res[0], "Could not find that table", http.StatusNotFound)
		}

		return "", fmt.Errorf("Could not find table")
	}

	return data.Tablename, err
}
