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

const sd = 100

type RelatedCharts struct {
	Charts []TableData `json:"charts"`
	Count  int         `json:"count"`
}

type PatternInfo struct {
	PatternID       int         `json:"patternid"`
	Discoverer      string      `json:"discoveredby"`
	DiscoveryDate   time.Time   `json:"discoverydate"`
	Validators      []string    `json:"validatedby, omitempty"`
	PrimarySource   string      `json:"source1"`
	SecondarySource string      `json:"source2, omitempty"`
	Strength        float64     `json:"statstrength, omitempty"`
	Observations    int         `json:"numobs"`
	ChartData       interface{} `json:"chartdata"`
}

type CorrelatedCharts struct {
	Charts []CorrelationData `json:"charts"`
	Count  int               `json:"count"`
}

type DiscoveredCharts struct {
	Charts []string `json:"charts"`
	Count  int      `json:"count"`
}

// Get all data for single selected chart
func GetChart(tablename string, tablenum int, chartType string, uid int, coords ...string) (PatternInfo, *appError) {
	pattern := PatternInfo{}
	id := ""
	x, y, z := coords[0], coords[1], ""
	if len(coords) > 2 {
		z = coords[2]
		id = tablename + "_" + strconv.Itoa(tablenum) + "_" + chartType + "_" + x + "_" + y + "_" + z //unique id
	} else {
		id = tablename + "_" + strconv.Itoa(tablenum) + "_" + chartType + "_" + x + "_" + y //unique id
	}

	xyz := XYVal{X: x, Y: y, Z: z}

	guid, _ := GetRealTableName(tablename)
	index := Index{}
	err := DB.Where("guid = ?", tablename).Find(&index).Error
	if err != nil && err != gorm.RecordNotFound {
		return pattern, &appError{err, "Database query failed (GUID)", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return pattern, &appError{err, "No related chart found", http.StatusNotFound}
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

	chart := make([]TableData, 0)
	GenerateChartData(chartType, guid, xyz, &chart, index)

	jByte, err := json.Marshal(chart[0])
	if err != nil {
		return pattern, &appError{err, "Unable to parse JSON", http.StatusInternalServerError}
	}

	//if the table is as yet undiscovered then add to the discovered table as an initial discovery
	discovered := Discovered{}
	var errd *appError
	err1 := DB.Where("relation_id = ?", id).First(&discovered).Error
	if err1 == gorm.RecordNotFound {
		discovered, errd = Discover(id, uid, jByte, false)
		if errd != nil {
			return pattern, errd
		}
	}

	user := User{}
	err4 := DB.Where("uid = ?", discovered.Uid).First(&user).Error
	if err4 != nil && err4 != gorm.RecordNotFound {
		return pattern, &appError{err4, "unable to retrieve user", http.StatusInternalServerError}
	}

	validators := make([]string, 0)
	validatorsUsers := []struct {
		Validation
		Username string
	}{}
	err5 := DB.Select("DISTINCT uid, (SELECT priv_users.username FROM priv_users WHERE priv_users.uid = priv_validations.uid) as username").Where("discovered_id = ?", discovered.DiscoveredId).Where("valflag = ?", true).Find(&validatorsUsers).Error
	if err5 != nil && err5 != gorm.RecordNotFound {
		return pattern, &appError{err5, "find validators failed", http.StatusInternalServerError}
	} else {
		for _, vu := range validatorsUsers {
			validators = append(validators, vu.Username)
		}
	}

	var observation []Observation
	count := 0
	err7 := DB.Model(&observation).Where("discovered_id = ?", discovered.DiscoveredId).Count(&count).Error
	if err7 != nil {
		return pattern, &appError{err7, "observation count failed", http.StatusInternalServerError}
	}

	var td TableData
	err3 := json.Unmarshal(discovered.Json, &td)
	if err3 != nil {
		return pattern, &appError{err3, "json failed", http.StatusInternalServerError}
	}

	pattern.ChartData = td
	pattern.PatternID = discovered.DiscoveredId
	pattern.Discoverer = user.Username
	pattern.DiscoveryDate = discovered.Created
	pattern.Validators = validators
	pattern.PrimarySource = SanitizeString(index.Title)
	// pattern.SecondarySource = ""
	// pattern.Strength = 0
	pattern.Observations = count

	return pattern, nil
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

	//if undiscovered add to the discovered table as an initial discovery
	var discovered []Discovered
	err1 := DB.Where("correlation_id = ?", cid).Find(&discovered).Error
	if err1 == gorm.RecordNotFound {
		Discover(strconv.Itoa(cid), uid, []byte(chart[0]), true)
	}

	err2 := DB.Where("correlation_id = ?", cid).Find(&discovered).Error
	if err2 != nil {
		return result, &appError{err, "Validation failed", http.StatusInternalServerError}
	}

	err3 := json.Unmarshal(discovered[0].Json, &result)
	check(err3)
	result.CorrelationId = discovered[0].CorrelationId
	return result, nil
}

// save chart to valdiated table
func Discover(id string, uid int, json []byte, correlated bool) (Discovered, *appError) {
	val := Discovered{}

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
	if err != nil {
		return val, &appError{err, "unable to create discovery", http.StatusInternalServerError}
	}

	return val, nil
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
		originid := tablename + "_" + strconv.Itoa(i) + "_" + v.ChartType + "_" + v.LabelX + "_" + v.LabelY
		discovered := Discovered{}
		err := DB.Where("relation_id = ?", originid).Find(&discovered).Error
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
func GetCorrelatedCharts(tableName string, offset int, count int, searchDepth int) (CorrelatedCharts, *appError) {
	correlation := make([]Correlation, 0)
	charts := make([]CorrelationData, 0) ///empty slice for adding all possible charts
	var cd CorrelationData

	GenerateCorrelations(tableName, searchDepth)
	fmt.Println("ROBOCOP1", time.Now())
	err := DB.Where("tbl1 = ?", tableName).Order("abscoef DESC").Find(&correlation).Error
	if err != nil && err != gorm.RecordNotFound {
		return CorrelatedCharts{nil, 0}, &appError{nil, "Database query failed (TBL1)", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return CorrelatedCharts{nil, 0}, &appError{nil, "No correlated chart found", http.StatusNotFound}
	}

	fmt.Println("ROBOCOP2", time.Now())
	for _, c := range correlation {
		json.Unmarshal(c.Json, &cd)
		cd.CorrelationId = c.CorrelationId

		if cd.Method == "Pearson" {
			cd.ChartType = "bar"
			charts = append(charts, cd)
			cd.ChartType = "column"
			charts = append(charts, cd)
			cd.ChartType = "line"
			charts = append(charts, cd)
			cd.ChartType = "scatter"
			charts = append(charts, cd)

		} else if cd.Method == "Spurious" {
			cd.ChartType = "line"
			charts = append(charts, cd)
			cd.ChartType = "scatter"
			charts = append(charts, cd)
			cd.ChartType = "stacked"
			charts = append(charts, cd)

		} else if cd.Method == "Visual" {
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
		return CorrelatedCharts{nil, 0}, &appError{nil, fmt.Sprintf("Offset value out of bounds (Max: %d)", totalCharts), http.StatusBadRequest}
	}

	last := offset + count
	if offset != 0 && last > totalCharts {
		return CorrelatedCharts{nil, 0}, &appError{nil, fmt.Sprintf("Count value out of bounds (Max: %d)", totalCharts-offset), http.StatusBadRequest}
	} else if offset == 0 && (last > totalCharts || count == 0) {
		last = totalCharts
	}

	fmt.Println("ROBOCOP3", time.Now())
	for _, c := range charts {
		originid := strconv.Itoa(c.CorrelationId)
		discovered := Discovered{}
		err := DB.Where("correlation_id = ?", originid).Find(&discovered).Error
		if err == gorm.RecordNotFound {
			c.Discovered = false
		} else {
			c.Discovered = true
		}
	}

	fmt.Println("ROBOCOP4", time.Now())
	charts = charts[offset:last] // return marshalled slice
	return CorrelatedCharts{charts, totalCharts}, nil
}

// As GetNew but get charts users have already voted on and return in an order based upon their absoulte ranking value
func GetDiscoveredCharts(tableName string, correlated bool, offset int, count int) (DiscoveredCharts, *appError) {
	discovered := make([]Discovered, 0)
	charts := make([]string, 0)
	var vd []byte

	if correlated {
		c := Correlation{}
		v := Discovered{}
		selectStr := v.TableName() + ".json"
		joinStr := "LEFT JOIN " + c.TableName() + " ON " + v.TableName() + ".correlation_id = " + c.TableName() + ".correlation_id"
		whereStr := c.TableName() + ".tbl1 = ?"
		err := DB.Select(selectStr).Joins(joinStr).Where(whereStr, tableName).Order("rating DESC").Find(&discovered).Error
		if err != nil && err != gorm.RecordNotFound {
			return DiscoveredCharts{nil, 0}, &appError{nil, "Database query failed (JOIN)", http.StatusInternalServerError}
		} else if err == gorm.RecordNotFound {
			return DiscoveredCharts{nil, 0}, &appError{nil, "No valid chart found", http.StatusNotFound}
		}
	} else {
		tableName = tableName + "_%"
		err := DB.Model(Discovered{}).Select("json").Where("relation_id LIKE ?", tableName).Order("rating DESC").Find(&discovered).Error
		if err != nil && err != gorm.RecordNotFound {
			return DiscoveredCharts{nil, 0}, &appError{nil, "Database query failed", http.StatusInternalServerError}
		} else if err == gorm.RecordNotFound {
			return DiscoveredCharts{nil, 0}, &appError{nil, "No valid chart found", http.StatusNotFound}
		}
	}

	for _, v := range discovered {
		vd = v.Json
		charts = append(charts, string(vd))
	}

	totalCharts := len(charts)
	if offset > totalCharts {
		return DiscoveredCharts{nil, 0}, &appError{nil, fmt.Sprintf("Offset value out of bounds (Max: %d)", totalCharts), http.StatusBadRequest}
	}

	last := offset + count
	if offset != 0 && last > totalCharts {
		return DiscoveredCharts{nil, 0}, &appError{nil, fmt.Sprintf("Count value out of bounds (Max: %d)", totalCharts-offset), http.StatusBadRequest}
	} else if offset == 0 && last > totalCharts {
		last = totalCharts
	}

	charts = charts[offset:last] // return marshalled slice
	return DiscoveredCharts{charts, totalCharts}, nil
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

	rows, err := DB.Raw(sql).Rows()
	if err == nil {
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
	} else {
		*charts = append(*charts, tmpTD)
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

	var result PatternInfo
	var err2 *appError

	if params["z"] == "" {
		result, err2 = GetChart(params["tablename"], tablenum, params["type"], uid, params["x"], params["y"])
		if err2 != nil {
			http.Error(res, err2.Message, err2.Code)
			return err2.Message
		}
	} else {
		result, err2 = GetChart(params["tablename"], tablenum, params["type"], uid, params["x"], params["y"], params["z"])
		if err2 != nil {
			http.Error(res, err2.Message, err2.Code)
			return err2.Message
		}
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

	cid, err := strconv.Atoi(params["cid"])
	if err != nil {
		http.Error(res, "Invalid id parameter", http.StatusBadRequest)
		return "Invalid id parameter"
	}

	uid, err1 := GetUserID(session)
	if err1 != nil {
		http.Error(res, err1.Message, err1.Code)
		return "Could not validate user"
	}

	result, error := GetChartCorrelated(cid, uid)
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

	var search, offset, count int
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

	if params["search"] == "true" { ///default searchdepth when blank
		search = sd
	} else if params["search"] == "false" { // do not search when 0 so can return just what exist in table
		search = 0
	}

	result, error := GetCorrelatedCharts(params["tablename"], offset, count, search)
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

func GetDiscoveredChartsHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
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

	result, error := GetDiscoveredCharts(params["tablename"], correlated, offset, count)
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

	var result PatternInfo
	var err2 *appError

	if params["z"] == "" {
		result, err2 = GetChart(params["tablename"], tablenum, params["type"], uid, params["x"], params["y"])
		if err2 != nil {
			return err2.Message
		}
	} else {
		result, err2 = GetChart(params["tablename"], tablenum, params["type"], uid, params["x"], params["y"], params["z"])
		if err2 != nil {
			return err2.Message
		}
	}

	r, e := json.Marshal(result)
	if e != nil {
		return e.Error()
	}

	return string(r)
}

func GetChartCorrelatedQ(params map[string]string) string {
	cid, err := strconv.Atoi(params["cid"])
	if err != nil {
		return err.Error()
	}

	uid, err1 := strconv.Atoi(params["uid"])
	if err1 != nil {
		return err1.Error()
	}

	result, err2 := GetChartCorrelated(cid, uid)
	if err2 != nil {
		return err2.Message
	}

	r, err3 := json.Marshal(result)
	if err3 != nil {
		return err3.Error()
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

	search := 0

	if params["search"] == "true" { ///default searchdepth when true
		search = sd
	} else if params["search"] == "false" { // do not search when 0 so can return just what exist in table
		search = 0
	}

	result, err := GetCorrelatedCharts(params["tablename"], offset, count, search)
	if err != nil {
		return err.Message
	}

	r, e := json.Marshal(result)
	if e != nil {
		return e.Error()
	}

	return string(r)
}

func GetDiscoveredChartsQ(params map[string]string) string {
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

	result, err := GetDiscoveredCharts(params["tablename"], correlated, offset, count)
	if err != nil {
		return err.Message
	}

	r, e := json.Marshal(result)
	if e != nil {
		return e.Error()
	}

	return string(r)
}
