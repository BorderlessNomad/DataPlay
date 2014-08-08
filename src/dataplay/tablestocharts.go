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

// generate all the potentially valid charts that relate to a single tablename
func GetRelatedCharts(tableName string, offset int, count int) (RelatedCharts, *appError) {
	columns := FetchTableCols(tableName) //array column names
	guid, _ := GetRealTableName(tableName)
	charts := make([]TableData, 0) ///empty slice for adding all possible charts
	index := Index{}
	xyNames := XYPermutations(columns) // get all possible valid permuations of columns as X & Y

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
		pieSQL := fmt.Sprintf("SELECT %s AS x, COUNT(%s) AS y FROM %s GROUP BY %s", v.Name, v.Name, guid, v.Name)
		GetChartData("pie", pieSQL, xyPie, &charts, index)
	}

	for _, v := range xyNames { /// create column and line charts, stacked column if plotting varchar against varchar
		sql := fmt.Sprintf("SELECT %s AS x, %s AS y FROM  %s", v.X, v.Y, guid)
		if v.Xtype == "varchar" && v.Ytype == "varchar" {
			GetChartData("stacked column", sql, v, &charts, index)
		}
		GetChartData("column", sql, v, &charts, index)
		GetChartData("line", sql, v, &charts, index)
		GetChartData("row", sql, v, &charts, index)
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
	} else if offset == 0 && last > totalCharts {
		last = totalCharts
	}

	charts = charts[offset:last] // return marshalled slice
	return RelatedCharts{charts, totalCharts}, nil
}

func GetNewCorrelatedCharts(tableName string, searchDepth int, offset int, count int) (RelatedCharts, *appError) {
	charts := make([]TableData, 0) ///empty slice for adding all possible charts

	/// inject chart type for each, intelligent but random - if 3 then use bubble

	totalCharts := len(charts)
	if offset > totalCharts {
		return RelatedCharts{nil, 0}, &appError{nil, fmt.Sprintf("Offset value out of bounds (Max: %d)", totalCharts), http.StatusBadRequest}
	}

	last := offset + count
	if offset != 0 && last > totalCharts {
		return RelatedCharts{nil, 0}, &appError{nil, fmt.Sprintf("Count value out of bounds (Max: %d)", totalCharts-offset), http.StatusBadRequest}
	} else if offset == 0 && last > totalCharts {
		last = totalCharts
	}

	charts = charts[offset:last] // return marshalled slice
	return RelatedCharts{charts, totalCharts}, nil
}

func GetValidatedCorrelatedCharts(tableName string, offset int, count int) (RelatedCharts, *appError) {
	charts := make([]TableData, 0) ///empty slice for adding all possible charts

	//where table 1 is the same, highest user rating
	/// inject chart type for each, intelligent but random - if 3 then use bubble

	totalCharts := len(charts)
	if offset > totalCharts {
		return RelatedCharts{nil, 0}, &appError{nil, fmt.Sprintf("Offset value out of bounds (Max: %d)", totalCharts), http.StatusBadRequest}
	}

	last := offset + count
	if offset != 0 && last > totalCharts {
		return RelatedCharts{nil, 0}, &appError{nil, fmt.Sprintf("Count value out of bounds (Max: %d)", totalCharts-offset), http.StatusBadRequest}
	} else if offset == 0 && last > totalCharts {
		last = totalCharts
	}

	charts = charts[offset:last] // return marshalled slice
	return RelatedCharts{charts, totalCharts}, nil
}

// Generate all possible permutations of xy columns
func XYPermutations(columns []ColType) []XYVal {
	length := len(columns)
	var xyNames []XYVal
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
	return xyNames
}

// Get arrays of data for the types of charts requested
func GetChartData(chartType string, sql string, names XYVal, charts *[]TableData, ind Index) {
	var tmpTD TableData
	var tmpXY XYVal
	tmpTD.ChartType = chartType
	tmpTD.Title = SanitizeString(ind.Title)
	tmpTD.Desc = SanitizeString(ind.Notes)
	tmpTD.LabelX = names.X
	var dx, dy time.Time
	var fx, fy float64
	var vx, vy string
	pieSlices := 0

	rows, _ := DB.Raw(sql).Rows()

	defer rows.Close()

	if chartType == "pie" { // single column pie chart x = type, y = count
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
			} else if names.Xtype == "float" {
				rows.Scan(&fx, &fy)
				tmpXY.X = FloatToString(fx)
				tmpXY.Y = FloatToString(fy)
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else {
				tmpXY.X = ""
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			}
		}

		if pieSlices < 20 && pieSlices > 1 { // reject pies with too many slices or not enough
			*charts = append(*charts, tmpTD)
		}

	} else {
		tmpTD.LabelY = names.Y
		for rows.Next() {
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

		*charts = append(*charts, tmpTD)
	}
}

//////////////////////////////////////////////////////////////////////////
//////////// HTTP AND QUEUE FUNCTIONS TO CALL ABOVE METHODS///////////////
//////////////////////////////////////////////////////////////////////////

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

	if params["searchDepth"] == "" {
		searchDepth = 1000
	} else {
		searchDepth, err = strconv.Atoi(params["searchDepth"])
		if err != nil {
			http.Error(res, "Invalid searchDepth parameter.", http.StatusBadRequest)
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

	searchDepth, e := strconv.Atoi(params["searchDepth"])
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
