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
	Charts []TableData
	Count  int
}

type RelatedCorrelatedCharts struct {
	Charts []CorrelationData
	Count  int
}

func GetChart(tableName string, chartType string, coords ...string) (TableData, *appError) {
	guid, _ := GetRealTableName(tableName)
	index := Index{}
	chart := make([]TableData, 0)
	x, y, z := coords[0], coords[1], ""
	if len(coords) > 2 {
		z = coords[2]
	}
	xyz := XYVal{X: x, Y: y, Z: z}

	err := DB.Where("guid = ?", tableName).Find(&index).Error
	if err != nil && err != gorm.RecordNotFound {
		return chart[0], &appError{err, "Database query failed", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return chart[0], &appError{err, "No related chart found", http.StatusNotFound}
	}

	columns := FetchTableCols(tableName)
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

	GetChartData(chartType, guid, xyz, &chart, index)
	result := chart[0]
	return result, nil
}

// generate all the potentially valid charts that relate to a single tablename, add apt charting types, and return them along with their total count
func GetRelatedCharts(tableName string, offset int, count int) (RelatedCharts, *appError) {
	columns := FetchTableCols(tableName) //array column names
	guid, _ := GetRealTableName(tableName)
	charts := make([]TableData, 0) ///empty slice for adding all possible charts
	index := Index{}
	xyNames := XYPermutations(columns, false) // get all possible valid permuations of columns as X & Y

	err := DB.Where("guid = ?", tableName).Find(&index).Error
	if err != nil && err != gorm.RecordNotFound {
		return RelatedCharts{nil, 0}, &appError{err, "Database query failed", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return RelatedCharts{nil, 0}, &appError{err, "No related chart found", http.StatusNotFound}
	}

	var xyPie XYVal

	for _, v := range columns { // create single column pie charts
		xyPie.X = v.Name
		xyPie.Xtype = v.Sqltype
		GetChartData("pie", guid, xyPie, &charts, index)
	}

	for _, v := range xyNames { /// create all other types of chart

		if v.Xtype == "varchar" && v.Ytype == "varchar" { // stacked or scatter charts if string v string values
			GetChartData("stacked column", guid, v, &charts, index)
			GetChartData("scatter", guid, v, &charts, index)
			// column and row charts for all that are not string v string values and are not date v string or string v date values
		} else if !(v.Xtype == "varchar" && v.Ytype == "date") || !(v.Xtype == "date" && v.Ytype == "varchar") {
			GetChartData("column", guid, v, &charts, index)
			GetChartData("row", guid, v, &charts, index)
		}

		if v.Xtype == "date" || (v.Xtype == "float" || v.Xtype == "integer") { // line chart only where x is a date or number value
			GetChartData("line", guid, v, &charts, index)
		}
	}

	if len(columns) > 2 { // if there's more than 2 columns grab a 3rd variable for bubble charts
		xyNames = XYPermutations(columns, true) // set z flag to true to get all possible valid permuations of columns as X, Y & Z
		for _, v := range xyNames {
			GetChartData("bubble", guid, v, &charts, index)
		}
	}

	// for i := range charts { // shuffle charts into random order
	// 	j := rand.Intn(i + 1)
	// 	charts[i], charts[j] = charts[j], charts[i]
	// }

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

	charts = charts[offset:last] // return marshalled slice
	return RelatedCharts{charts, totalCharts}, nil
}

// Look for new correlated charts, take the correlations and break them down into charting types, and return them along with their total count
func GetNewCorrelatedCharts(tableName string, searchDepth int, offset int, count int) (RelatedCorrelatedCharts, *appError) {
	charts := make([]CorrelationData, 0) ///empty slice for adding all possible charts
	GenerateCorrelations(tableName, searchDepth)
	sql := fmt.Sprintf("SELECT method, json, abscoef FROM priv_correlation WHERE tbl1 = %s AND rating = 0 ORDER BY abscoef DESC", tableName)
	var meth string
	var js []byte
	var acoef float64
	var cd CorrelationData

	rows, _ := DB.Raw(sql).Rows()
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&meth, &js, &acoef)
		err := json.Unmarshal(js, &cd)
		check(err)

		if meth == "Pearson" {
			cd.Chart = "bar"
			charts = append(charts, cd)
			cd.Chart = "column"
			charts = append(charts, cd)
			cd.Chart = "line"
			charts = append(charts, cd)
			cd.Chart = "scatter"
			charts = append(charts, cd)

		} else if meth == "Spurious" {
			cd.Chart = "line"
			charts = append(charts, cd)
			cd.Chart = "scatter"
			charts = append(charts, cd)
			cd.Chart = "stacked"
			charts = append(charts, cd)

		} else if meth == "Visual" {
			cd.Chart = "bar"
			charts = append(charts, cd)
			cd.Chart = "column"
			charts = append(charts, cd)
			cd.Chart = "line"
			charts = append(charts, cd)
			cd.Chart = "scatter"
			charts = append(charts, cd)
		} else {
			cd.Chart = "unknown"
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
	} else if offset == 0 && last > totalCharts {
		last = totalCharts
	}

	charts = charts[offset:last] // return marshalled slice
	return RelatedCorrelatedCharts{charts, totalCharts}, nil
}

// As GetNew but get charts users have already voted on and return in an order based upon their absoulte ranking value
func GetValidatedCorrelatedCharts(tableName string, offset int, count int) (RelatedCorrelatedCharts, *appError) {
	charts := make([]CorrelationData, 0) ///empty slice for adding all possible charts
	sql := fmt.Sprintf("SELECT method, json, abscoef FROM priv_correlation WHERE tbl1 = %s AND rating != 0 ORDER BY abscoef DESC", tableName)
	var meth string
	var js []byte
	var acoef float64
	var cd CorrelationData

	rows, _ := DB.Raw(sql).Rows()
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&meth, &js, &acoef)
		err := json.Unmarshal(js, &cd)
		check(err)

		if meth == "Pearson" {
			cd.Chart = "bar"
			charts = append(charts, cd)
			cd.Chart = "column"
			charts = append(charts, cd)
			cd.Chart = "line"
			charts = append(charts, cd)
			cd.Chart = "scatter"
			charts = append(charts, cd)

		} else if meth == "Spurious" {
			cd.Chart = "line"
			charts = append(charts, cd)
			cd.Chart = "scatter"
			charts = append(charts, cd)
			cd.Chart = "stacked"
			charts = append(charts, cd)

		} else if meth == "Visual" {
			cd.Chart = "bar"
			charts = append(charts, cd)
			cd.Chart = "column"
			charts = append(charts, cd)
			cd.Chart = "line"
			charts = append(charts, cd)
			cd.Chart = "scatter"
			charts = append(charts, cd)
		} else {
			cd.Chart = "unknown"
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
	} else if offset == 0 && last > totalCharts {
		last = totalCharts
	}

	charts = charts[offset:last] // return marshalled slice
	return RelatedCorrelatedCharts{charts, totalCharts}, nil
}

// Get arrays of data for the types of charts requested (titles, descriptions, all the xy values etc)
func GetChartData(chartType string, guid string, names XYVal, charts *[]TableData, ind Index) {
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
			} else if (names.Xtype == "float" || names.Xtype == "integer") && names.Ytype == "date" && (names.Ztype == "float" || names.Ztype == "integer") {
				rows.Scan(&fx, &dy, &fz)
				tmpXY.X = FloatToString(fx)
				tmpXY.Y = (dy.String()[0:10])
				tmpXY.Z = FloatToString(fz)
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if (names.Xtype == "float" || names.Xtype == "integer") && names.Ytype == "varchar" && (names.Ztype == "float" || names.Ztype == "integer") {
				rows.Scan(&fx, &vy, &fz)
				tmpXY.X = FloatToString(fx)
				tmpXY.Y = vy
				tmpXY.Z = FloatToString(fz)
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if names.Xtype == "date" && (names.Ytype == "float" || names.Ytype == "integer") && (names.Ztype == "float" || names.Ztype == "integer") {
				rows.Scan(&dx, &fy, &fz)
				tmpXY.Y = (dx.String()[0:10])
				tmpXY.Y = FloatToString(fy)
				tmpXY.Z = FloatToString(fz)
			} else if names.Xtype == "varchar" && (names.Ytype == "float" || names.Ytype == "integer") && (names.Ztype == "float" || names.Ztype == "integer") {
				rows.Scan(&vx, &fy, &fz)
				tmpXY.Y = vx
				tmpXY.Y = FloatToString(fy)
				tmpXY.Z = FloatToString(fz)
			} else {
				tmpXY.X = ""
				tmpXY.Y = ""
				tmpXY.Z = ""
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			}
		}

		*charts = append(*charts, tmpTD)

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
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			}
		}

		if pieSlices <= 20 && pieSlices > 1 { // reject pies with too many slices or not enough
			*charts = append(*charts, tmpTD)
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
				*charts = append(*charts, tmpTD)
			}
		} else {
			*charts = append(*charts, tmpTD)
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

//////////////////////////////////////////////////////////////////////////
//////////// HTTP AND QUEUE FUNCTIONS TO CALL ABOVE METHODS///////////////
//////////////////////////////////////////////////////////////////////////

func GetChartHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	if params["tablename"] == "" {
		http.Error(res, "Invalid tablename.", http.StatusBadRequest)
		return ""
	}

	if params["type"] == "" {
		http.Error(res, "Invalid chart type.", http.StatusBadRequest)
		return ""
	}

	if params["x"] == "" {
		http.Error(res, "Invalid x label.", http.StatusBadRequest)
		return ""
	}

	if params["y"] == "" {
		http.Error(res, "Invalid y label.", http.StatusBadRequest)
		return ""
	}

	result, error := GetChart(params["tablename"], params["type"], params["x"], params["y"], params["z"])
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

func GetRelatedChartsHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	var offset, count int
	var err error

	if params["offset"] == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(params["offset"])
		if err != nil {
			http.Error(res, "Invalid offset parameter.", http.StatusBadRequest)
			return ""
		}
	}

	if params["count"] == "" {
		count = 3
	} else {
		count, err = strconv.Atoi(params["count"])
		if err != nil {
			http.Error(res, "Invalid count parameter.", http.StatusBadRequest)
			return ""
		}
	}

	result, error := GetRelatedCharts(params["tablename"], offset, count)
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

func GetNewCorrelatedChartsHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	var searchDepth, offset, count int
	var err error

	if params["offset"] == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(params["offset"])
		if err != nil {
			http.Error(res, "Invalid offset parameter.", http.StatusBadRequest)
			return ""
		}
	}

	if params["count"] == "" {
		count = 3
	} else {
		count, err = strconv.Atoi(params["count"])
		if err != nil {
			http.Error(res, "Invalid count parameter.", http.StatusBadRequest)
			return ""
		}
	}

	if params["searchdepth"] == "" {
		searchDepth = 1000
	} else {
		searchDepth, err = strconv.Atoi(params["searchdepth"])
		if err != nil {
			http.Error(res, "Invalid searchdepth parameter.", http.StatusBadRequest)
			return ""
		}
	}

	result, error := GetNewCorrelatedCharts(params["tablename"], searchDepth, offset, count)
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

func GetValidatedCorrelatedChartsHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	var offset, count int
	var err error

	if params["offset"] == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(params["offset"])
		if err != nil {
			http.Error(res, "Invalid offset parameter.", http.StatusBadRequest)
			return ""
		}
	}

	if params["count"] == "" {
		count = 3
	} else {
		count, err = strconv.Atoi(params["count"])
		if err != nil {
			http.Error(res, "Invalid count parameter.", http.StatusBadRequest)
			return ""
		}
	}

	result, error := GetValidatedCorrelatedCharts(params["tablename"], offset, count)
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

func GetChartQ(params map[string]string) string {
	if params["tablename"] == "" {
		return ""
	}

	if params["type"] == "" {
		return ""
	}

	if params["x"] == "" {
		return ""
	}

	if params["y"] == "" {
		return ""
	}

	result, err := GetChart(params["tablename"], params["type"], params["x"], params["y"], params["z"])
	if err != nil {
		return ""
	}

	r, e := json.Marshal(result)
	if e != nil {
		return ""
	}

	return string(r)
}

func GetRelatedChartsQ(params map[string]string) string {
	if params["user"] == "" {
		return ""
	}

	offset, e := strconv.Atoi(params["offset"])
	if e != nil {
		return ""
	}

	count, e := strconv.Atoi(params["count"])
	if e != nil {
		return ""
	}

	result, err := GetRelatedCharts(params["tablename"], offset, count)
	if err != nil {
		return ""
	}

	r, e := json.Marshal(result)
	if e != nil {
		return ""
	}

	return string(r)
}

func GetNewCorrelatedChartsQ(params map[string]string) string {
	if params["user"] == "" {
		return ""
	}

	offset, e := strconv.Atoi(params["offset"])
	if e != nil {
		return ""
	}

	count, e := strconv.Atoi(params["count"])
	if e != nil {
		return ""
	}

	searchDepth, e := strconv.Atoi(params["searchdepth"])
	if e != nil {
		return ""
	}

	result, err := GetNewCorrelatedCharts(params["tablename"], searchDepth, offset, count)
	if err != nil {
		return ""
	}

	r, e := json.Marshal(result)
	if e != nil {
		return ""
	}

	return string(r)
}

func GetValidatedCorrelatedChartsQ(params map[string]string) string {
	if params["user"] == "" {
		return ""
	}

	offset, e := strconv.Atoi(params["offset"])
	if e != nil {
		return ""
	}

	count, e := strconv.Atoi(params["count"])
	if e != nil {
		return ""
	}

	result, err := GetValidatedCorrelatedCharts(params["tablename"], offset, count)
	if err != nil {
		return ""
	}

	r, e := json.Marshal(result)
	if e != nil {
		return ""
	}

	return string(r)
}