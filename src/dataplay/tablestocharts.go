package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

type RelatedCharts struct {
	Charts []TableData `json:"charts"`
	Count  int         `json:"count"`
}

type RelatedCorrelatedCharts struct {
	Charts []CorrelationData `json:"charts"`
	Count  int               `json:"count"`
}

type ValidatedCharts struct {
	Charts []string `json:"charts"`
	Count  int      `json:"count"`
}

// Get all data for single selected chart
func GetChart(tablename string, tablenum int, chartType string, uid int, coords ...string) (TableData, *appError) {
	guid, _ := GetRealTableName(tablename)
	index := Index{}
	chart := make([]TableData, 0)
	x, y, z := coords[0], coords[1], ""
	if len(coords) > 2 {
		z = coords[2]
	}
	xyz := XYVal{X: x, Y: y, Z: z}

	err := DB.Where("guid = ?", tablename).Find(&index).Error
	if err != nil && err != gorm.RecordNotFound {
		return chart[0], &appError{err, "Database query failed (GUID)", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return chart[0], &appError{err, "No related chart found", http.StatusNotFound}
	}

	columns := FetchTableCols(tablename)
	for _, v := range columns {
		if v.Name == x {
			xyz.Xtype = v.Sqltype
		}
		if v.Name == y {
			xyz.Ytype = v.Sqltype
		}
		if v.Name == z {
			xyz.Ztype = v.Sqltype
		}
	}

	GenerateChartData(chartType, guid, xyz, &chart, index)
	id := tablename + "_" + strconv.Itoa(tablenum) //unique id

	jByte, err := json.Marshal(chart[0])
	if err != nil {
		return chart[0], &appError{err, "Unable to parse JSON", http.StatusInternalServerError}
	}

	//if the table is as yet undiscovered then add to the validated table as an initial discovery
	var validated []Validated
	err1 := DB.Where("relation_id = ?", id).Find(&validated).Error
	if err1 == gorm.RecordNotFound {
		Discover(id, uid, jByte, false)
	}

	err2 := DB.Where("relation_id = ?", id).Find(&validated).Error
	if err2 != nil {
		return chart[0], &appError{err, "Validation failed", http.StatusInternalServerError}
	}

	var result TableData
	err3 := json.Unmarshal(validated[0].Json, &result)
	check(err3)
	result.RelationId = validated[0].RelationId
	return result, nil
}

// use the id relating to the record stored in the generated correlations table to return the json with the specific chart info
func GetChartCorrelated(cid int, uid int) (CorrelationData, *appError) {
	var chart []string
	var result CorrelationData
	err := DB.Model(Correlation{}).Where("correlation_id = ?", cid).Pluck("json", &chart).Error

	if err != nil && err != gorm.RecordNotFound {
		return result, &appError{err, "Database query failed (ID)", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return result, &appError{err, "No related chart found", http.StatusNotFound}
	}

	//if undiscovered add to the validated table as an initial discovery
	var validated []Validated
	err1 := DB.Where("correlation_id = ?", cid).Find(&validated).Error
	if err1 == gorm.RecordNotFound {
		Discover(strconv.Itoa(cid), uid, []byte(chart[0]), true)
	}

	err2 := DB.Where("correlation_id = ?", cid).Find(&validated).Error
	if err2 != nil {
		return result, &appError{err, "Validation failed", http.StatusInternalServerError}
	}

	err3 := json.Unmarshal(validated[0].Json, &result)
	check(err3)
	result.CorrelationId = validated[0].CorrelationId
	return result, nil
}

// save chart to valdiated table
func Discover(id string, uid int, json []byte, correlated bool) {
	val := Validated{}

	if correlated {
		val.CorrelationId, _ = strconv.Atoi(id)
	} else {
		val.RelationId = id
	}

	val.Uid = uid
	val.Json = json
	val.Created = time.Now()
	val.Rating = 0
	val.Invalid = 0
	val.Valid = 0
	err := DB.Save(&val).Error
	check(err)
}

// generate all the potentially valid charts that relate to a single tablename, add apt charting types,
// and return them along with their total count and whether they've been discovered
func GetRelatedCharts(tablename string, offset int, count int) (RelatedCharts, *appError) {
	columns := FetchTableCols(tablename) //array column names
	guid, _ := GetRealTableName(tablename)
	charts := make([]TableData, 0) ///empty slice for adding all possible charts
	index := Index{}
	xyNames := XYPermutations(columns, false) // get all possible valid permuations of columns as X & Y

	err := DB.Where("guid = ?", tablename).Find(&index).Error
	if err != nil && err != gorm.RecordNotFound {
		return RelatedCharts{nil, 0}, &appError{err, "Database query failed (GUID)", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return RelatedCharts{nil, 0}, &appError{err, "No related chart found", http.StatusNotFound}
	}

	var xyPie XYVal

	for _, v := range columns { // create single column pie charts
		xyPie.X = v.Name
		xyPie.Xtype = v.Sqltype
		GenerateChartData("pie", guid, xyPie, &charts, index)
	}

	for _, v := range xyNames { /// create all other types of chart

		if v.Xtype == "varchar" && v.Ytype == "varchar" { // stacked or scatter charts if string v string values
			GenerateChartData("stacked column", guid, v, &charts, index)
			GenerateChartData("scatter", guid, v, &charts, index)
			// column and row charts for all that are not string v string values and are not date v string or string v date values
		} else if !(v.Xtype == "varchar" && v.Ytype == "date") || !(v.Xtype == "date" && v.Ytype == "varchar") {
			GenerateChartData("row", guid, v, &charts, index)
			// no string values for y axis on column charts
			if v.Ytype != "varchar" && v.Ytype != "date" {
				GenerateChartData("column", guid, v, &charts, index)
			}
		}

		if v.Xtype != "varchar" && (v.Ytype != "date" || v.Ytype != "varchar") { // line chart cannot be based on strings or have date on the Y axis
			GenerateChartData("line", guid, v, &charts, index)
		}
	}

	if len(columns) > 2 { // if there's more than 2 columns grab a 3rd variable for bubble charts
		xyNames = XYPermutations(columns, true) // set z flag to true to get all possible valid permuations of columns as X, Y & Z
		for _, v := range xyNames {
			GenerateChartData("bubble", guid, v, &charts, index)
		}
	}

	totalCharts := len(charts)
	if offset > totalCharts {
		return RelatedCharts{nil, 0}, &appError{nil, fmt.Sprintf("Offset value out of bounds (Max: %d)", totalCharts), http.StatusBadRequest}
	}

	last := offset + count
	if offset != 0 && last > totalCharts {
		return RelatedCharts{nil, 0}, &appError{nil, fmt.Sprintf("Count value out of bounds (Max: %d)", totalCharts-offset), http.StatusBadRequest}
	} else if offset == 0 && (last > totalCharts || count == 0) {
		last = totalCharts
	}

	for i, v := range charts {
		originid := tablename + "_" + strconv.Itoa(i)
		validated := Validated{}
		err := DB.Where("relation_id = ?", originid).Find(&validated).Error
		if err == gorm.RecordNotFound {
			v.Discovered = false
		} else {
			v.Discovered = true
		}
	}

	charts = charts[offset:last] // return marshalled slice
	return RelatedCharts{charts, totalCharts}, nil
}

// Look for new correlated charts, take the correlations and break them down into charting types, and return them along with their total count
// To return only existing charts use searchdepth = 0
func GetCorrelatedCharts(tableName string, offset int, count int, searchDepth int) (RelatedCorrelatedCharts, *appError) {
	correlation := make([]Correlation, 0)
	charts := make([]CorrelationData, 0) ///empty slice for adding all possible charts
	var cd CorrelationData

	GenerateCorrelations(tableName, searchDepth)
	err := DB.Where("tbl1 = ?", tableName).Order("abscoef DESC").Find(&correlation).Error
	if err != nil && err != gorm.RecordNotFound {
		return RelatedCorrelatedCharts{nil, 0}, &appError{nil, "Database query failed (TBL1)", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return RelatedCorrelatedCharts{nil, 0}, &appError{nil, "No correlated chart found", http.StatusNotFound}
	}

	for _, c := range correlation {
		err := json.Unmarshal(c.Json, &cd)
		check(err)

		if c.Method == "Pearson" {
			cd.ChartType = "bar"
			charts = append(charts, cd)
			cd.ChartType = "column"
			charts = append(charts, cd)
			cd.ChartType = "line"
			charts = append(charts, cd)
			cd.ChartType = "scatter"
			charts = append(charts, cd)

		} else if c.Method == "Spurious" {
			cd.ChartType = "line"
			charts = append(charts, cd)
			cd.ChartType = "scatter"
			charts = append(charts, cd)
			cd.ChartType = "stacked"
			charts = append(charts, cd)

		} else if c.Method == "Visual" {
			cd.ChartType = "bar"
			charts = append(charts, cd)
			cd.ChartType = "column"
			charts = append(charts, cd)
			cd.ChartType = "line"
			charts = append(charts, cd)
			cd.ChartType = "scatter"
			charts = append(charts, cd)
		} else {
			cd.ChartType = "unknown"
			charts = append(charts, cd)
		}
	}

	totalCharts := len(charts)
	if offset > totalCharts {
		return RelatedCorrelatedCharts{nil, 0}, &appError{nil, fmt.Sprintf("Offset value out of bounds (Max: %d)", totalCharts), http.StatusBadRequest}
	}

	last := offset + count
	if offset != 0 && last > totalCharts {
		return RelatedCorrelatedCharts{nil, 0}, &appError{nil, fmt.Sprintf("Count value out of bounds (Max: %d)", totalCharts-offset), http.StatusBadRequest}
	} else if offset == 0 && (last > totalCharts || count == 0) {
		last = totalCharts
	}

	for _, c := range charts {
		originid := strconv.Itoa(c.CorrelationId)
		validated := Validated{}
		err := DB.Where("correlation_id = ?", originid).Find(&validated).Error
		if err == gorm.RecordNotFound {
			c.Discovered = false
		} else {
			c.Discovered = true
		}
	}

	charts = charts[offset:last] // return marshalled slice
	return RelatedCorrelatedCharts{charts, totalCharts}, nil
}

// As GetNew but get charts users have already voted on and return in an order based upon their absoulte ranking value
func GetValidatedCharts(tableName string, correlated bool, offset int, count int) (ValidatedCharts, *appError) {
	validated := make([]Validated, 0)
	charts := make([]string, 0)
	var vd []byte

	if correlated {
		err := DB.Select("priv_validatedtables.json").Joins("LEFT JOIN priv_correlation ON priv_validatedtables.correlation_id = priv_correlation.correlation_id").Where("priv_correlation.tbl1 = ?", tableName).Order("priv_validatedtables.rating DESC").Find(&validated).Error
		if err != nil && err != gorm.RecordNotFound {
			return ValidatedCharts{nil, 0}, &appError{nil, "Database query failed (JOIN)", http.StatusInternalServerError}
		} else if err == gorm.RecordNotFound {
			return ValidatedCharts{nil, 0}, &appError{nil, "No valid chart found", http.StatusNotFound}
		}
	} else {
		tableName = tableName + "_%"
		err := DB.Select("priv_validatedtables.json").Where("priv_validatedtables.relation_id LIKE ?", tableName).Order("priv_validatedtables.rating DESC").Find(&validated).Error
		if err != nil && err != gorm.RecordNotFound {
			return ValidatedCharts{nil, 0}, &appError{nil, "Database query failed", http.StatusInternalServerError}
		} else if err == gorm.RecordNotFound {
			return ValidatedCharts{nil, 0}, &appError{nil, "No valid chart found", http.StatusNotFound}
		}
	}

	for _, v := range validated {
		vd = v.Json
		charts = append(charts, string(vd))
	}

	totalCharts := len(charts)
	if offset > totalCharts {
		return ValidatedCharts{nil, 0}, &appError{nil, fmt.Sprintf("Offset value out of bounds (Max: %d)", totalCharts), http.StatusBadRequest}
	}

	last := offset + count
	if offset != 0 && last > totalCharts {
		return ValidatedCharts{nil, 0}, &appError{nil, fmt.Sprintf("Count value out of bounds (Max: %d)", totalCharts-offset), http.StatusBadRequest}
	} else if offset == 0 && last > totalCharts {
		last = totalCharts
	}

	charts = charts[offset:last] // return marshalled slice
	return ValidatedCharts{charts, totalCharts}, nil
}

// Get arrays of data for the types of charts requested (titles, descriptions, all the xy values etc)
func GenerateChartData(chartType string, guid string, names XYVal, charts *[]TableData, ind Index) {
	var tmpTD TableData
	var tmpXY XYVal
	tmpTD.ChartType = chartType
	tmpTD.Title = SanitizeString(ind.Title)
	tmpTD.Desc = SanitizeString(ind.Notes)
	tmpTD.LabelX = names.X
	var dx, dy time.Time
	var fx, fy, fz float64
	var vx, vy string
	pieSlices, rowAmt := 0, 0
	sql := ""

	if chartType == "pie" {
		if names.Xtype == "float" {
			sql = fmt.Sprintf("SELECT %s AS x, SUM(%s) AS y FROM %s GROUP BY %s", names.X, names.X, guid, names.X)
			tmpTD.LabelY = "sum"
		} else {
			sql = fmt.Sprintf("SELECT %s AS x, COUNT(%s) AS y FROM %s GROUP BY %s", names.X, names.X, guid, names.X)
			tmpTD.LabelY = "count"
		}
	} else if chartType == "bubble" {
		sql = fmt.Sprintf("SELECT %s AS x, %s AS y, %s AS z FROM  %s", names.X, names.Y, names.Z, guid)
	} else {
		sql = fmt.Sprintf("SELECT %s AS x, %s AS y FROM  %s", names.X, names.Y, guid)
	}

	rows, _ := DB.Raw(sql).Rows()
	defer rows.Close()

	if chartType == "bubble" {
		tmpTD.LabelY = names.Y
		tmpTD.LabelZ = names.Z
		for rows.Next() {
			if (names.Xtype == "float" || names.Xtype == "integer") && (names.Ytype == "float" || names.Ytype == "integer") && (names.Ztype == "float" || names.Ztype == "integer") {
				rows.Scan(&fx, &fy, &fz)
				tmpXY.X = FloatToString(fx)
				tmpXY.Y = FloatToString(fy)
				tmpXY.Z = FloatToString(fz)
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if (names.Xtype == "float" || names.Xtype == "integer") && names.Ytype == "varchar" && (names.Ztype == "float" || names.Ztype == "integer") {
				rows.Scan(&fx, &vy, &fz)
				tmpXY.X = FloatToString(fx)
				tmpXY.Y = vy
				tmpXY.Z = FloatToString(fz)
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if names.Xtype == "varchar" && (names.Ytype == "float" || names.Ytype == "integer") && (names.Ztype == "float" || names.Ztype == "integer") {
				rows.Scan(&vx, &fy, &fz)
				tmpXY.X = vx
				tmpXY.Y = FloatToString(fy)
				tmpXY.Z = FloatToString(fz)
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else {
				tmpXY.X = ""
				tmpXY.Y = ""
				tmpXY.Z = ""
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			}
		}
		if ValueCheck(tmpTD) && NegCheck(tmpTD) {
			*charts = append(*charts, tmpTD)
		}

	} else if chartType == "pie" { // single column pie chart x = type, y = count
		for rows.Next() {
			pieSlices++
			if names.Xtype == "varchar" {
				rows.Scan(&vx, &fy)
				tmpXY.X = vx
				tmpXY.Y = FloatToString(fy)
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if names.Xtype == "date" {
				rows.Scan(&dx, &fy)
				tmpXY.X = (dx.String()[0:10])
				tmpXY.Y = FloatToString(fy)
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if names.Xtype == "float" || names.Xtype == "integer" {
				rows.Scan(&fx, &fy)
				tmpXY.X = FloatToString(fx)
				tmpXY.Y = FloatToString(fy)
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else {
				tmpXY.X = ""
				tmpXY.Y = ""
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			}
		}

		if pieSlices <= 20 && pieSlices > 1 { // reject pies with too many slices or not enough
			if ValueCheck(tmpTD) {
				*charts = append(*charts, tmpTD)
			}
		}

	} else { // for all other types of chart
		tmpTD.LabelY = names.Y
		for rows.Next() {
			if chartType == "row" {
				rowAmt++
			}
			if names.Xtype == "date" && names.Ytype == "date" {
				rows.Scan(&dx, &dy)
				tmpXY.X = (dx.String()[0:10])
				tmpXY.Y = (dy.String()[0:10])
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if names.Xtype == "date" && (names.Ytype == "float" || names.Ytype == "integer") {
				rows.Scan(&dx, &fy)
				tmpXY.X = (dx.String()[0:10])
				tmpXY.Y = FloatToString(fy)
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if names.Xtype == "date" && names.Ytype == "varchar" {
				rows.Scan(&dx, &vy)
				tmpXY.X = (dx.String()[0:10])
				tmpXY.Y = vy
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if (names.Xtype == "float" || names.Xtype == "integer") && names.Ytype == "date" {
				rows.Scan(&fx, &dy)
				tmpXY.X = FloatToString(fx)
				tmpXY.Y = (dy.String()[0:10])
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if (names.Xtype == "float" || names.Xtype == "integer") && (names.Ytype == "float" || names.Ytype == "integer") {
				rows.Scan(&fx, &fy)
				tmpXY.X = FloatToString(fx)
				tmpXY.Y = FloatToString(fy)
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if (names.Xtype == "float" || names.Xtype == "integer") && names.Ytype == "varchar" {
				rows.Scan(&fx, &vy)
				tmpXY.X = FloatToString(fx)
				tmpXY.Y = vy
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if names.Xtype == "varchar" && names.Ytype == "date" {
				rows.Scan(&vx, &dy)
				tmpXY.X = vx
				tmpXY.Y = (dy.String()[0:10])
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if names.Xtype == "varchar" && (names.Ytype == "float" || names.Ytype == "integer") {
				rows.Scan(&vx, &fy)
				tmpXY.X = vx
				tmpXY.Y = FloatToString(fy)
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if names.Xtype == "varchar" && names.Ytype == "varchar" {
				rows.Scan(&vx, &vy)
				tmpXY.X = vx
				tmpXY.Y = vy
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else {
				tmpXY.X = ""
				tmpXY.Y = ""
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			}
		}

		if chartType == "row" {
			if rowAmt <= 20 && rowAmt > 1 {
				if ValueCheck(tmpTD) {
					*charts = append(*charts, tmpTD)
				}
			}
		} else {
			if ValueCheck(tmpTD) {
				*charts = append(*charts, tmpTD)
			}
		}
	}
}

// Generate all possible permutations of xy columns
func XYPermutations(columns []ColType, bubble bool) []XYVal {
	length := len(columns)
	var xyNames []XYVal
	var xyzNames []XYVal
	var tmpXY XYVal

	for i := 0; i < length; i++ {
		for j := 0; j < length; j++ {
			if columns[i] != columns[j] {
				tmpXY.X = columns[i].Name
				tmpXY.Y = columns[j].Name
				tmpXY.Xtype = columns[i].Sqltype
				tmpXY.Ytype = columns[j].Sqltype
				xyNames = append(xyNames, tmpXY)
			}
		}
	}
	// if bubble chart add xyz permutations
	if bubble {
		for _, v := range xyNames {
			for k := 0; k < length; k++ {
				if columns[k].Name != v.X && columns[k].Name != v.Y {
					tmpXY.X = v.X
					tmpXY.Y = v.Y
					tmpXY.Z = columns[k].Name
					tmpXY.Xtype = v.Xtype
					tmpXY.Ytype = v.Ytype
					tmpXY.Ztype = columns[k].Sqltype
					xyzNames = append(xyzNames, tmpXY)
				}
			}
		}
		return xyzNames
	}

	return xyNames
}

// checks whether data values for x or y are singular and if so returns false so that the corresponding chart is not added
func ValueCheck(t TableData) bool {
	lastXval, lastYval := t.Values[0].X, t.Values[0].Y
	xChk, yChk := false, false

	for _, v := range t.Values {
		if lastXval != v.X {
			xChk = true
		}
		if lastYval != v.Y {
			yChk = true
		}
	}

	if xChk && yChk {
		return true
	} else {
		return false
	}
}

// checks whether any X axis values are negative as bubble won't plot if they are
func NegCheck(t TableData) bool {
	for _, v := range t.Values {
		x, _ := strconv.Atoi(v.X)
		if x < 0 {
			return false
		}
	}
	return true
}

//////////////////////////////////////////////////////////////////////////
//////////// HTTP AND QUEUE FUNCTIONS TO CALL ABOVE METHODS///////////////
//////////////////////////////////////////////////////////////////////////

func GetChartHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	uid, err1 := GetUserID(session)
	if err1 != nil {
		http.Error(res, err1.Message, err1.Code)
		return "Could not validate user"
	}

	if params["tablename"] == "" {
		http.Error(res, "Invalid tablename", http.StatusBadRequest)
		return "Invalid tablename"
	}

	tablenum, err := strconv.Atoi(params["tablenum"])
	if err != nil {
		http.Error(res, "Invalid tablenum parameter", http.StatusBadRequest)
		return "Invalid tablenum parameter"
	}

	if params["type"] == "" {
		http.Error(res, "Invalid chart type", http.StatusBadRequest)
		return "Invalid chart type"
	}

	if params["x"] == "" {
		http.Error(res, "Invalid x label", http.StatusBadRequest)
		return "Invalid x label"
	}

	if params["y"] == "" {
		http.Error(res, "Invalid y label", http.StatusBadRequest)
		return "Invalid y label"
	}

	result, err2 := GetChart(params["tablename"], tablenum, params["type"], uid, params["x"], params["y"], params["z"])
	if err2 != nil {
		http.Error(res, err2.Message, err2.Code)
		return err2.Message
	}

	r, err3 := json.Marshal(result)
	if err3 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return "Unable to parse JSON"
	}

	return string(r)
}

func GetChartCorrelatedHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(res, "Invalid id parameter", http.StatusBadRequest)
		return "Invalid id parameter"
	}

	uid, err1 := GetUserID(session)
	if err1 != nil {
		http.Error(res, err1.Message, err1.Code)
		return "Could not validate user"
	}

	result, error := GetChartCorrelated(id, uid)
	if error != nil {
		http.Error(res, error.Message, error.Code)
		return ""
	}

	r, err2 := json.Marshal(result)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return "Unable to parse JSON"
	}

	return string(r)
}

func GetRelatedChartsHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	var offset, count int
	var err error

	if params["offset"] == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(params["offset"])
		if err != nil {
			http.Error(res, "Invalid offset parameter", http.StatusBadRequest)
			return "Invalid offset parameter"
		}
	}

	if params["count"] == "" {
		count = 3
	} else {
		count, err = strconv.Atoi(params["count"])
		if err != nil {
			http.Error(res, "Invalid count parameter", http.StatusBadRequest)
			return "Invalid count parameter"
		}
	}

	result, error := GetRelatedCharts(params["tablename"], offset, count)
	if error != nil {
		http.Error(res, error.Message, error.Code)
		return error.Message
	}

	r, err1 := json.Marshal(result)
	if err1 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return "Unable to parse JSON"
	}

	return string(r)
}

func GetCorrelatedChartsHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	var searchDepth, offset, count int
	var err error

	if params["offset"] == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(params["offset"])
		if err != nil {
			http.Error(res, "Invalid offset parameter", http.StatusBadRequest)
			return "Invalid offset parameter"
		}
	}

	if params["count"] == "" {
		count = 3
	} else {
		count, err = strconv.Atoi(params["count"])
		if err != nil {
			http.Error(res, "Invalid count parameter", http.StatusBadRequest)
			return "Invalid count parameter"
		}
	}

	if params["searchdepth"] == "" { ///default searchdepth when blank
		searchDepth = 100
	} else if params["searchdepth"] == "0" { // do not search when 0 so can return just what exist in table
		searchDepth = 0
	} else {
		searchDepth, err = strconv.Atoi(params["searchdepth"])
		if err != nil {
			http.Error(res, "Invalid searchdepth parameter", http.StatusBadRequest)
			return "Invalid searchdepth parameter"
		}
	}

	result, error := GetCorrelatedCharts(params["tablename"], offset, count, searchDepth)
	if error != nil {
		http.Error(res, error.Message, error.Code)
		return error.Message
	}

	r, err1 := json.Marshal(result)
	if err1 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return "Unable to parse JSON"
	}

	return string(r)
}

func GetValidatedChartsHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return "Missing session parameter"
	}

	var offset, count int
	var err error

	correlated, e := strconv.ParseBool(params["correlated"])
	if e != nil {
		return e.Error()
	}

	if params["offset"] == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(params["offset"])
		if err != nil {
			http.Error(res, "Invalid offset parameter", http.StatusBadRequest)
			return "Invalid offset parameter"
		}
	}

	if params["count"] == "" {
		count = 3
	} else {
		count, err = strconv.Atoi(params["count"])
		if err != nil {
			http.Error(res, "Invalid count parameter", http.StatusBadRequest)
			return "Invalid count parameter"
		}
	}

	result, error := GetValidatedCharts(params["tablename"], correlated, offset, count)
	if error != nil {
		http.Error(res, error.Message, error.Code)
		return error.Message
	}

	r, err1 := json.Marshal(result)
	if err1 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return "Unable to parse JSON"
	}

	return string(r)
}

func GetChartQ(params map[string]string) string {
	if params["tablename"] == "" {
		return "no tablename"
	}

	tablenum, err := strconv.Atoi(params["tablenum"])
	if err != nil {
		return "no tablenum"
	}

	if params["type"] == "" {
		return "no type"
	}

	uid, err1 := strconv.Atoi(params["uid"])
	if err1 != nil {
		return "invalid uid"
	}

	if params["x"] == "" {
		return "no x coordinate"
	}

	if params["y"] == "" {
		return "no y coordinate"
	}

	result, err2 := GetChart(params["tablename"], tablenum, params["type"], uid, params["x"], params["y"], params["z"])
	if err2 != nil {
		return err2.Message
	}

	r, e := json.Marshal(result)
	if e != nil {
		return e.Error()
	}

	return string(r)
}

func GetChartCorrelatedQ(params map[string]string) string {
	id, e := strconv.Atoi(params["id"])
	if e != nil {
		return e.Error()
	}

	uid, e := strconv.Atoi(params["uid"])
	if err != nil {
		return e.Error()
	}

	result, err1 := GetChartCorrelated(id, uid)
	if err1 != nil {
		return err1.Message
	}

	r, e := json.Marshal(result)
	if e != nil {
		return e.Error()
	}

	return string(r)
}

func GetRelatedChartsQ(params map[string]string) string {
	if params["user"] == "" {
		return "no user"
	}

	offset, e := strconv.Atoi(params["offset"])
	if e != nil {
		return e.Error()
	}

	count, e := strconv.Atoi(params["count"])
	if e != nil {
		return e.Error()
	}

	result, err := GetRelatedCharts(params["tablename"], offset, count)
	if err != nil {
		return err.Message
	}

	r, e := json.Marshal(result)
	if e != nil {
		return e.Error()
	}

	return string(r)
}

func GetCorrelatedChartsQ(params map[string]string) string {
	offset, e := strconv.Atoi(params["offset"])
	if e != nil {
		return e.Error()
	}

	count, e := strconv.Atoi(params["count"])
	if e != nil {
		return e.Error()
	}

	searchDepth, e := strconv.Atoi(params["searchdepth"])
	if e != nil {
		return e.Error()
	}

	result, err := GetCorrelatedCharts(params["tablename"], offset, count, searchDepth)
	if err != nil {
		return err.Message
	}

	r, e := json.Marshal(result)
	if e != nil {
		return e.Error()
	}

	return string(r)
}

func GetValidatedChartsQ(params map[string]string) string {
	correlated, e := strconv.ParseBool(params["correlated"])
	if e != nil {
		return e.Error()
	}

	offset, e := strconv.Atoi(params["offset"])
	if e != nil {
		return e.Error()
	}

	count, e := strconv.Atoi(params["count"])
	if e != nil {
		return e.Error()
	}

	result, err := GetValidatedCharts(params["tablename"], correlated, offset, count)
	if err != nil {
		return err.Message
	}

	r, e := json.Marshal(result)
	if e != nil {
		return e.Error()
	}

	return string(r)
}
