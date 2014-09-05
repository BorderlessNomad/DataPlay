package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/jinzhu/gorm"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Authresponse struct {
	Username string
	UserID   int
}

type mainDateVal struct {
	DateString string
	Count      int
}

//This function is used to gather what is the username is
// This used to be used on the front page but now it is mainly used as a "noop" call to check if the user is logged in or not.
func CheckAuth(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := params["session"]
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	user := User{}
	err1 := DB.Where("uid = ?", uid).Find(&user).Error
	check(err1)

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
	LocationData bool
}

func SearchForDataHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	result, error := SearchForData(params["s"], uid)
	if error != nil {
		http.Error(res, error.Message, error.Code)
		return ""
	}

	r, err1 := json.Marshal(result)
	if err1 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func SearchForDataQ(params map[string]string) string {
	if params["user"] == "" {
		return ""
	}

	uid, e := strconv.Atoi(params["user"])
	if e != nil {
		return ""
	}

	result, err := SearchForData(params["s"], uid)
	if err != nil {
		return ""
	}

	r, e := json.Marshal(result)
	if e != nil {
		return ""
	}

	return string(r)
}

/**
 * @brief Search a given term in database
 * @details This method searches for a matching title with following conditions,
 * 		Suffix wildcard
 * 		Prefix & suffix wildcard
 * 		Prefix, suffix & trimmed spaces with wildcard
 * 		Prefix & suffix on previously searched terms
 */
func SearchForData(str string, uid int) ([]SearchResult, *appError) {
	if str == "" {
		return nil, &appError{nil, "There was no search request", http.StatusBadRequest}
	}

	AddSearchTerm(str) // add to search term count

	Indices := []Index{}

	term := str + "%" // e.g. "nhs" => "nhs%" (What about "%nhs"?)

	Logger.Println("Searching with Suffix Wildcard", term)

	err := DB.Where("LOWER(title) LIKE LOWER(?)", term).Where("(owner = 0 OR owner = ?)", uid).Limit(10).Find(&Indices).Error
	if err != nil && err != gorm.RecordNotFound {
		return nil, &appError{err, "Database query failed (SUFFIX)", http.StatusServiceUnavailable}
	}

	Results := ProcessSearchResults(Indices, err)
	if len(Results) == 0 {
		term := "%" + str + "%" // e.g. "nhs" => "%nhs%"

		Logger.Println("Searching with Prefix + Suffix Wildcard", term)

		err := DB.Where("LOWER(title) LIKE LOWER(?)", term).Where("(owner = 0 OR owner = ?)", uid).Limit(10).Find(&Indices).Error
		if err != nil && err != gorm.RecordNotFound {
			return nil, &appError{err, "Database query failed (PREFIX + SUFFIX)", http.StatusServiceUnavailable}
		}

		Results = ProcessSearchResults(Indices, err)
		if len(Results) == 0 {
			term := "%" + strings.Replace(str, " ", "%", -1) + "%" // e.g. "nh s" => "%nh%s%"

			Logger.Println("Searching with Prefix + Suffix + Trim Wildcard", term)

			err := DB.Where("LOWER(title) LIKE LOWER(?)", term).Where("(owner = 0 OR owner = ?)", uid).Limit(10).Find(&Indices).Error
			if err != nil && err != gorm.RecordNotFound {
				return nil, &appError{err, "Database query failed (PREFIX + SUFFIX + TRIM)", http.StatusServiceUnavailable}
			}

			Results = ProcessSearchResults(Indices, err)
			if len(Results) == 0 && (len(str) >= 3 && len(str) < 20) {
				term := "%" + str + "%" // e.g. "nhs" => "%nhs%"

				Logger.Println("Searching with Prefix + Suffix Wildcard in String Table", term)

				query := DB.Table("priv_stringsearch, priv_onlinedata, index")
				query = query.Select("DISTINCT ON (priv_onlinedata.guid) priv_onlinedata.guid, index.title")
				query = query.Where("(LOWER(value) LIKE LOWER(?) OR LOWER(x) LIKE LOWER(?))", term, term)
				query = query.Where("priv_stringsearch.tablename = priv_onlinedata.tablename")
				query = query.Where("priv_onlinedata.guid = index.guid")
				query = query.Where("(owner = ? OR owner = ?)", 0, uid)
				query = query.Order("priv_onlinedata.guid")
				query = query.Order("priv_stringsearch.count DESC")
				query = query.Limit(10)
				err := query.Find(&Indices).Error

				if err != nil && err != gorm.RecordNotFound {
					return nil, &appError{err, "Database query failed (PREFIX + SUFFIX + STRING)", http.StatusInternalServerError}
				}

				Results = ProcessSearchResults(Indices, err)
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

func DumpTableHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	if params["id"] == "" {
		http.Error(res, "Sorry! Could not complete this request (Hint, You didnt ask for a table to be dumped)", http.StatusBadRequest)
		return ""
	}
	result, error := DumpTable(params)
	if error != nil {
		http.Error(res, "No data", http.StatusBadRequest)
		return ""
	}

	r, err := json.Marshal(result)
	if err != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func DumpTableQ(params map[string]string) string {
	if params["id"] == "" {
		return ""
	}

	result, error := DumpTable(params)
	if error != nil {
		return ""
	}

	r, err := json.Marshal(result)
	if err != nil {
		return ""
	}

	return string(r)
}

// This function will empty a whole table out into JSON
// Due to what seems to be a golang bug, everything is outputted as a string.
func DumpTable(params map[string]string) ([]map[string]interface{}, *appError) {
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
			return nil, &appError{cE, "Please give valid numbers for offset and count", http.StatusBadRequest}
		}
	}

	var tablename string
	var err error

	tablename, err = GetRealTableName(params["id"])
	if err != nil || len(tablename) == 0 {
		return nil, &appError{err, "Unable to find that table", http.StatusBadRequest}
	}

	var rows *sql.Rows

	if UsingRanges {
		rows, err = DB.Raw(fmt.Sprintf("SELECT * FROM %s OFFSET %d LIMIT %d", tablename, offset, count)).Rows()
	} else {
		rows, err = DB.Raw(fmt.Sprintf("SELECT * FROM %s", tablename)).Rows()
	}

	if err != nil {
		return nil, &appError{err, "Database query failed (SELECT)", http.StatusInternalServerError}
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, &appError{err, "Database query failed (COLUMNS)", http.StatusInternalServerError}
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
			return nil, &appError{err, "Database query failed (ROWS)", http.StatusInternalServerError}
		}

		record := ScanRow(values, columns)
		array = append(array, record)
	}

	return array, nil
}

func DumpTableRangeHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	if params["id"] == "" {
		http.Error(res, "Sorry! Could not complete this request (Hint, You didnt ask for a table to be dumped)", http.StatusBadRequest)
		return ""
	}

	if params["x"] == "" || params["startx"] == "" || params["endx"] == "" {
		http.Error(res, "You did not provide enough infomation to make this kind of request :id/:x/:startx/:endx", http.StatusBadRequest)
		return ""
	}

	result, error := DumpTableRange(params)
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

func DumpTableRangeQ(params map[string]string) string {
	if params["id"] == "" {
		return ""
	}

	if params["x"] == "" || params["startx"] == "" || params["endx"] == "" {
		return ""
	}

	result, error := DumpTableRange(params)
	if error != nil {
		return ""
	}

	r, err := json.Marshal(result)
	if err != nil {
		return ""
	}

	return string(r)
}

// This function will empty a whole table out into JSON
// Due to what seems to be a golang bug, everything is outputted as a string.
func DumpTableRange(params map[string]string) ([]map[string]interface{}, *appError) {
	tablename, e := GetRealTableName(params["id"])
	if e != nil {
		return nil, &appError{e, "Unable to find that table", http.StatusBadRequest}
	}

	rows, err := DB.Raw("SELECT * FROM " + tablename).Rows()
	if err != nil {
		return nil, &appError{err, "Database query failed (SELECT)", http.StatusInternalServerError}
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, &appError{err, "Database query failed (COLUMNS)", http.StatusInternalServerError}
	}

	var xcol int
	xcol = 999
	startx, starte := strconv.ParseInt(params["startx"], 10, 64)
	endx, ende := strconv.ParseInt(params["endx"], 10, 64)
	if starte != nil || ende != nil {
		return nil, &appError{nil, "Please give valid numbers for start and end", http.StatusBadRequest}
	}

	for number, colname := range columns {
		if colname == params["x"] {
			xcol = number
			break
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
			return nil, &appError{err, "Database query failed (ROWS)", http.StatusInternalServerError}
		}

		xvalue := values[xcol].(int64)

		if xvalue >= startx && xvalue <= endx {
			record := ScanRow(values, columns)
			array = append(array, record)
		}
	}

	return array, nil
}

func DumpTableGroupedHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	if params["id"] == "" || params["x"] == "" || params["y"] == "" {
		http.Error(res, "You did not provide enough infomation to make this kind of request :id/:x/:y", http.StatusBadRequest)
		return ""
	}

	result, error := DumpTableGrouped(params)
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

func DumpTableGroupedQ(params map[string]string) string {
	if params["id"] == "" || params["x"] == "" || params["y"] == "" {
		return ""
	}

	result, error := DumpTableGrouped(params)
	if error != nil {
		return ""
	}

	r, err := json.Marshal(result)
	if err != nil {
		return ""
	}

	return string(r)
}

// This call with use the GROUP BY function in mysql to query and get the sum of things
// This is very useful for things like picharts
// /api/getdatagrouped/:id/:x/:y
func DumpTableGrouped(params map[string]string) ([]map[string]interface{}, *appError) {
	tablename, e := GetRealTableName(params["id"])
	if e != nil {
		return nil, &appError{e, "Unable to find that table", http.StatusBadRequest}
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
		return nil, &appError{nil, "Col X is invalid.", http.StatusBadRequest}
	}

	if !ValidY {
		return nil, &appError{nil, "Col Y is invalid.", http.StatusBadRequest}
	}

	q := ""
	if sumY {
		q = fmt.Sprintf("SELECT %[1]s, SUM(%[2]s) AS %[2]s FROM %[3]s GROUP BY %[1]s", params["x"], params["y"], tablename)
	} else {
		q = fmt.Sprintf("SELECT DISTINCT %[1]s, COUNT(%[1]s) AS %[2]s FROM %[3]s GROUP BY %[1]s ORDER BY %[1]s", params["x"], params["y"], tablename)
	}

	rows, e1 := DB.Raw(q).Rows()
	if e1 != nil {
		return nil, &appError{e1, "Database query failed (SELECT)", http.StatusInternalServerError}
	}

	columns, e2 := rows.Columns()
	if e1 != nil || e2 != nil {
		return nil, &appError{e2, "Database query failed (COLUMNS)", http.StatusInternalServerError}
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
			return nil, &appError{err, "Database query failed (ROWS)", http.StatusInternalServerError}
		}

		record := ScanRow(values, columns)
		array = append(array, record)
	}

	return array, nil
}

func DumpTablePredictionHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	if params["id"] == "" || params["x"] == "" || params["y"] == "" {
		http.Error(res, "You did not provide enough infomation to make this kind of request :id/:x/:y", http.StatusBadRequest)
		return ""
	}

	result, error := DumpTablePrediction(params)
	if result == nil {
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

func DumpTablePredictionQ(params map[string]string) string {
	if params["id"] == "" || params["x"] == "" || params["y"] == "" {
		return ""
	}

	result, error := DumpTablePrediction(params)
	if error != nil {
		return ""
	}

	r, err := json.Marshal(result)
	if err != nil {
		return ""
	}

	return string(r)
}

// This call will get a X,Y and a prediction of a value. that is asked for
func DumpTablePrediction(params map[string]string) ([]float64, *appError) {
	tablename, e := GetRealTableName(params["id"])
	if e != nil {
		return nil, &appError{e, "Unable to find that table", http.StatusBadRequest}
	}

	cls := FetchTableCols(params["id"])
	// Now we need to check that the rows that the client is asking for, are in the table.
	ValidX := false
	for _, clm := range cls {
		if clm.Name == params["x"] {
			ValidX = true
		}
	}
	if !ValidX {
		return nil, &appError{nil, "Col X is invalid.", http.StatusBadRequest}
	}

	ValidY := false
	for _, clm := range cls {
		if clm.Name == params["y"] {
			ValidY = true
		}
	}

	if !ValidY {
		return nil, &appError{nil, "Col Y is invalid.", http.StatusBadRequest}
	}

	rows, e1 := DB.Raw(fmt.Sprintf("SELECT %s, %s FROM %s", params["x"], params["y"], tablename)).Rows()
	if e1 != nil {
		return nil, &appError{e1, "Database query failed (SELECT)", http.StatusInternalServerError}

	}

	columns, e2 := rows.Columns()
	if e2 != nil {
		return nil, &appError{e2, "Database query failed (COLUMNS)", http.StatusInternalServerError}
	}

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	xarray := make([]float64, 0)
	yarray := make([]float64, 0)
	array := make([]map[string]interface{}, 0)

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, &appError{err, "Database query failed (ROWS)", http.StatusInternalServerError}
		}

		record := ScanRow(values, columns)
		f1, e := ConvertToFloat(record[columns[0]])
		if e != nil {
			return nil, &appError{e, "Could not parse the value into a float (" + columns[0] + ")", http.StatusBadRequest}
		}

		f2, e := ConvertToFloat(record[columns[1]])
		if e != nil {
			return nil, &appError{e, "Could not parse the value into a float (" + columns[1] + ")", http.StatusBadRequest}
		}

		xarray = append(xarray, f1)
		yarray = append(yarray, f2)
		array = append(array, record)
	}

	results := GetPolyResults(xarray, yarray)

	return results, nil
}

func DumpReducedTableHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	result, error := DumpReducedTable(params)
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

func DumpReducedTableQ(params map[string]string) string {
	result, error := DumpReducedTable(params)
	if error != nil {
		return ""
	}

	r, e := json.Marshal(result)
	if e != nil {
		return ""
	}

	return string(r)
}

/**
 * @brief This function will take a share of a table and return it as JSON
 */
func DumpReducedTable(params map[string]string) ([]map[string]interface{}, *appError) {
	if params["id"] == "" {
		return nil, &appError{nil, "Sorry! Could not compleate this request (Hint, You didnt ask for a table to be dumped)", http.StatusBadRequest}
	}

	tablename, e := GetRealTableName(params["id"])
	if e != nil {
		return nil, &appError{e, "Unable to find that table", http.StatusBadRequest}
	}

	var rows *sql.Rows
	var e1 error
	if params["x"] == "" || params["y"] == "" {
		rows, e1 = DB.Table(tablename).Rows()
	} else {
		x := params["x"]
		y := params["y"]

		cls := FetchTableCols(params["id"])

		ValidX := false
		ValidY := false
		sumY := false
		for _, clm := range cls {
			/* Check for existence of X & Y in Table */
			if !ValidX && clm.Name == x {
				ValidX = true
			} else if !ValidY && clm.Name == y {
				ValidY = true
				if !sumY && clm.Sqltype != "varchar" && clm.Sqltype != "date" {
					sumY = true
				}
			}

			if ValidX && ValidY {
				break
			}
		}

		if !ValidX {
			return nil, &appError{nil, "Col X is invalid.", http.StatusBadRequest}
		}
		if !ValidY {
			return nil, &appError{nil, "Col Y is invalid.", http.StatusBadRequest}
		}

		if sumY {
			// If Y is Int/Float we can SUM
			rows, e1 = DB.Table(tablename).Select("DISTINCT \"" + x + "\", SUM(\"" + y + "\") AS \"" + y + "\"").Group(x).Order(x).Rows()
		} else {
			// Just count X aginst Y
			rows, e1 = DB.Table(tablename).Select("DISTINCT \"" + x + "\", COUNT(\"" + x + "\") AS \"" + y + "\"").Group(x).Order(x).Rows()
		}
	}

	if e1 != nil {
		return nil, &appError{e1, "Database query failed (SELECT)", http.StatusInternalServerError}
	}

	columns, e2 := rows.Columns()
	if e2 != nil {
		return nil, &appError{e2, "Database query failed (COLUMNS)", http.StatusInternalServerError}
	}

	colLength := len(columns)
	values := make([]interface{}, colLength)
	scanArgs := make([]interface{}, colLength)
	for i := 0; i < colLength; i++ {
		scanArgs[i] = &values[i]
	}

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		err := rows.Scan(scanArgs...)

		if err != nil {
			return nil, &appError{err, "Database query failed (ROWS)", http.StatusInternalServerError}
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
			return nil, &appError{e, "Invalid percentage", http.StatusBadRequest}
		}

		Temp := (float64(Divider) / 100) * float64(DataLength)

		if Temp < 1 {
			Temp = 1
		}

		DataLength = DataLength / int(Temp)

		if params["min"] != "" {
			MinSpend, e := strconv.ParseInt(params["min"], 10, 64)
			if e != nil {
				return nil, &appError{e, "Invalid min", http.StatusBadRequest}
			}

			if DataLength > 0 && int(RealDL/DataLength) < int(MinSpend) {
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

	results_calc := make([]map[string]interface{}, 0)
	for i, result := range results {
		if i%DataLength == 0 {
			results_calc = append(results_calc, result)
		}
	}

	return results_calc, nil
}

/**
 * @brief Converts GUID ('friendly' name) into actual table inside database
 *
 * @param string GUID
 * @param http http.ResponseWriter
 *
 * @return string output, error
 */
func GetRealTableName(guid string) (out string, e error) {
	data := OnlineData{}
	err := DB.Select("tablename").Where("guid = ?", guid).Find(&data).Error

	if err == gorm.RecordNotFound {
		return "", fmt.Errorf("Could not find table")
	}

	return data.Tablename, err
}

/**
 * @brief Convert given interface's value to Float
 * @details Try and convert a given value to Float. Value can be in the form of Float,
 * Int, Un-signed Int, String etc.
 *
 * @param  val interface{}
 * @return float64, error
 */
func ConvertToFloat(val interface{}) (float64, error) {
	switch i := val.(type) {
	case float64:
		return float64(i), nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int16:
		return float64(i), nil
	case int8:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint16:
		return float64(i), nil
	case uint8:
		return float64(i), nil
	case int:
		return float64(i), nil
	case uint:
		return float64(i), nil
	case string:
		f, err := strconv.ParseFloat(i, 64)
		if err != nil {
			return math.NaN(), err
		}

		return f, err
	default:
		return math.NaN(), errors.New("ConvertToFloat: Unknown value is of incompatible type")
	}
}

func AddSearchTerm(str string) {
	searchterm := SearchTerm{}

	err := DB.Where("term = ?", str).Find(&searchterm).Error
	if err != nil && err != gorm.RecordNotFound {
		panic(err)
	} else if err == gorm.RecordNotFound {
		searchterm.Count = 0
		searchterm.Term = str
	}

	searchterm.Count++
	err = DB.Save(&searchterm).Error
}

// RUN ONCE AND POPULATE PRIMARY DATE FIELD IN PRIV_ONLINEDATA WITH MAIN TABLE DATE FOR USE IN SEARCH
func PrimaryDate() {
	var names []string

	DB.Table("priv_onlinedata").Pluck("guid", &names)

	for _, v := range names {
		cols := FetchTableCols(v)
		dateCol := RandomDateColumn(cols)
		table, _ := GetRealTableName(v)
		d := "DELETE FROM " + table + " WHERE " + dateCol + " = '0001-01-01 BC'" ////////TEMP FIX TO GET RID OF INVALID VALUES IN GOV DATA
		DB.Exec(d)
		if dateCol != "" {
			var dates []time.Time
			err := DB.Table(table).Pluck(dateCol, &dates).Error
			if err == nil {
				dv := make([]DateVal, 0)
				var d DateVal
				for _, v := range dates {

					d.Date = v
					dv = append(dv, d)
				}
				primaryDate := MainDate(dv)
				err := DB.Table("priv_onlinedata").Where("tablename = ?", table).Update("primarydate", primaryDate).Error
				check(err)
			}
		}
	}
}

func MainDate(d []DateVal) string {
	from, to, rng := DetermineRange(d)
	start, end, n := 0, 0, 0

	if rng > 366 { // get most popular year
		start = from.Year()
		end = to.Year()
		n = end - start
	} else { // get most popular month
		start = DayNum(from)
		end = DayNum(to)
		n = ((end - start) / 31) + 1
	}

	dv := make([]mainDateVal, n) // use date value for date and count

	if n > 0 {
		for _, v := range d {
			if rng > 366 {
				isit, i := stringInSlice(strconv.Itoa(v.Date.Year()), dv)
				if isit {
					dv[i].Count++
				} else {
					tmpdv := mainDateVal{DateString: strconv.Itoa(v.Date.Year()), Count: 1}
					dv = append(dv, tmpdv)
				}
			} else {
				isit, i := stringInSlice(v.Date.Month().String()+" "+strconv.Itoa(v.Date.Year()), dv)
				if isit {
					dv[i].Count++
				} else {
					str := v.Date.Month().String() + " " + strconv.Itoa(v.Date.Year())
					tmpdv := mainDateVal{DateString: str, Count: 1}
					dv = append(dv, tmpdv)
				}
			}
		}
	} else {
		return from.Month().String() + " " + strconv.Itoa(from.Year())
	}

	highest := 0
	maindate := ""

	for _, v := range dv {
		if v.Count > highest {
			highest = v.Count
			maindate = v.DateString
		}
	}

	return maindate
}

func stringInSlice(dateString string, list []mainDateVal) (bool, int) {
	for i, j := range list {
		if j.DateString == dateString {
			return true, i
		}
	}
	return false, 0
}
