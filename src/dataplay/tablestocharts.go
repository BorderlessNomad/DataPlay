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

type DataEntry struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

type PatternInfo struct {
	PatternID       string      `json:"patternid"`
	DiscoveredID    int         `json:"discoveredid"`
	Discoverer      string      `json:"discoveredby"`
	DiscoveryDate   time.Time   `json:"discoverydate"`
	Creditors       []string    `json:"creditedby, omitempty"`
	Discreditors    []string    `json:"discreditedby, omitempty"`
	PrimarySource   string      `json:"source1"`
	SecondarySource string      `json:"source2, omitempty"`
	Strength        string      `json:"statstrength, omitempty"`
	Observations    int         `json:"numobs"`
	UserCredited    bool        `json:"userhascredited"`
	UserDiscredited bool        `json:"userhasdiscredited"`
	ChartData       interface{} `json:"chartdata"`
}

type CorrelatedCharts struct {
	Charts []CorrelationData `json:"charts"`
	Count  int               `json:"count"`
}

type DiscoveredCharts struct {
	Charts []interface{} `json:"charts"`
	Count  int           `json:"count"`
}

// Get all data for single selected chart
func GetChart(tablename string, tablenum int, chartType string, uid int, coords ...string) (PatternInfo, *appError) {
	pattern := PatternInfo{}
	id := ""
	x, y, z := coords[0], coords[1], ""
	if len(coords) > 2 {
		z = coords[2]
		id = tablename + "/" + strconv.Itoa(tablenum) + "/" + chartType + "/" + x + "/" + y + "/" + z //unique id
	} else {
		id = tablename + "/" + strconv.Itoa(tablenum) + "/" + chartType + "/" + x + "/" + y //unique id
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
	if len(chart) == 0 {
		return pattern, &appError{err, "Not possible to plot this chart", http.StatusInternalServerError}
	}

	jByte, err1 := json.Marshal(chart[0])
	if err1 != nil {
		return pattern, &appError{err1, "Unable to parse JSON", http.StatusInternalServerError}
	}

	//if the table is as yet undiscovered then add to the discovered table as an initial discovery
	discovered := Discovered{}
	var err3 *appError
	err2 := DB.Where("relation_id = ?", id).Find(&discovered).Error
	if err2 == gorm.RecordNotFound {
		discovered, err3 = Discover(id, uid, jByte, false)
		if err3 != nil {
			return pattern, err3
		}
	}

	user := User{}
	err4 := DB.Where("uid = ?", discovered.Uid).Find(&user).Error
	if err4 != nil && err4 != gorm.RecordNotFound {
		return pattern, &appError{err4, "unable to retrieve user for related chart", http.StatusInternalServerError}
	}

	creditors := make([]string, 0)
	discreditors := make([]string, 0)
	creditingUsers := []struct {
		Credit
		Username string
	}{}

	query := DB.Select("DISTINCT uid, credflag, (SELECT priv_users.username FROM priv_users WHERE priv_users.uid = priv_credits.uid) as username")
	query = query.Where("discovered_id = ?", discovered.DiscoveredId)

	err5 := query.Find(&creditingUsers).Error

	if err5 != nil && err5 != gorm.RecordNotFound {
		return pattern, &appError{err5, "find creditors failed", http.StatusInternalServerError}
	} else {
		for _, vu := range creditingUsers {
			if vu.Credflag == true {
				creditors = append(creditors, vu.Username)
			} else if vu.Credflag == false {
				discreditors = append(discreditors, vu.Username)
			}
		}
	}

	var observation []Observation
	count := 0
	err6 := DB.Model(&observation).Where("discovered_id = ?", discovered.DiscoveredId).Count(&count).Error
	if err6 != nil {
		return pattern, &appError{err6, "observation count failed", http.StatusInternalServerError}
	}

	var td TableData
	err7 := json.Unmarshal(discovered.Json, &td)
	if err7 != nil {
		return pattern, &appError{err7, "json failed", http.StatusInternalServerError}
	}

	pattern.ChartData = td
	pattern.PatternID = discovered.RelationId
	pattern.DiscoveredID = discovered.DiscoveredId
	pattern.Discoverer = user.Username
	pattern.DiscoveryDate = discovered.Created
	pattern.Creditors = creditors
	pattern.Discreditors = discreditors
	pattern.PrimarySource = SanitizeString(index.Title)
	pattern.Observations = count

	// See if user has credited or discredited this chart
	cred := Credit{}
	err8 := DB.Where("discovered_id = ?", discovered.DiscoveredId).Where("uid = ?", uid).Find(&cred).Error
	if err8 == nil {
		if cred.Credflag == true {
			pattern.UserCredited = true
			pattern.UserDiscredited = false
		} else {
			pattern.UserDiscredited = true
			pattern.UserCredited = false
		}
	} else if err8 == gorm.RecordNotFound {
		pattern.UserCredited = false
		pattern.UserDiscredited = false
	} else if err8 != nil {
		return pattern, &appError{err8, "user credited failed", http.StatusInternalServerError}
	}

	return pattern, nil
}

// use the id relating to the record stored in the generated correlations table to return the json with the specific chart info
func GetChartCorrelated(cid int, uid int) (PatternInfo, *appError) {
	pattern := PatternInfo{}
	var chart []string
	var cd CorrelationData

	err := DB.Model(Correlation{}).Where("correlation_id = ?", cid).Pluck("json", &chart).Error
	if err != nil && err != gorm.RecordNotFound {
		return pattern, &appError{err, "Database query failed (ID)", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return pattern, &appError{err, "No related chart found", http.StatusNotFound}
	}

	//if undiscovered add to the discovered table as an initial discovery
	discovered := Discovered{}
	err1 := DB.Where("correlation_id = ?", cid).Find(&discovered).Error
	if err1 == gorm.RecordNotFound {
		Discover(strconv.Itoa(cid), uid, []byte(chart[0]), true)
	}

	user := User{}
	err2 := DB.Where("uid = ?", discovered.Uid).Find(&user).Error
	if err2 != nil && err2 != gorm.RecordNotFound {
		return pattern, &appError{err2, "unable to retrieve user for correlated chart", http.StatusInternalServerError}
	}

	creditors := make([]string, 0)
	discreditors := make([]string, 0)
	creditingUsers := []struct {
		Credit
		Username string
	}{}

	query := DB.Select("DISTINCT uid, credflag, (SELECT priv_users.username FROM priv_users WHERE priv_users.uid = priv_credits.uid) as username")
	query = query.Where("discovered_id = ?", discovered.DiscoveredId)
	err3 := query.Find(&creditingUsers).Error

	if err3 != nil && err3 != gorm.RecordNotFound {
		return pattern, &appError{err3, "find creditors failed", http.StatusInternalServerError}
	} else {
		for _, vu := range creditingUsers {
			if vu.Credflag == true {
				creditors = append(creditors, vu.Username)
			} else if vu.Credflag == false {
				discreditors = append(discreditors, vu.Username)
			}
		}
	}

	err4 := DB.Where("correlation_id = ?", cid).Find(&discovered).Error
	if err4 != nil {
		return pattern, &appError{err4, "Correlation failed", http.StatusInternalServerError}
	}

	var observation []Observation
	count := 0
	err5 := DB.Model(&observation).Where("discovered_id = ?", discovered.DiscoveredId).Count(&count).Error
	if err5 != nil {
		return pattern, &appError{err5, "observation count failed", http.StatusInternalServerError}
	}

	err6 := json.Unmarshal(discovered.Json, &cd)
	if err6 != nil {
		return pattern, &appError{err6, "Json failed", http.StatusInternalServerError}
	}

	correlation := Correlation{}
	err7 := DB.Where("correlation_id = ?", cid).Find(&correlation).Error
	if err7 != nil {
		return pattern, &appError{err7, "Correlation failed", http.StatusInternalServerError}
	}

	pattern.ChartData = cd
	pattern.PatternID = strconv.Itoa(discovered.CorrelationId)
	pattern.DiscoveredID = discovered.DiscoveredId
	pattern.Discoverer = user.Username
	pattern.DiscoveryDate = discovered.Created
	pattern.Creditors = creditors
	pattern.Discreditors = discreditors
	pattern.PrimarySource = cd.Table1.Title
	pattern.SecondarySource = cd.Table2.Title
	pattern.Strength = CalcStrength(correlation.Abscoef)
	pattern.Observations = count

	// See if user has credited or discredited this chart
	cred := Credit{}
	err8 := DB.Where("discovered_id = ?", discovered.DiscoveredId).Where("uid = ?", uid).Find(&cred).Error
	if err8 == nil {
		if cred.Credflag == true {
			pattern.UserCredited = true
			pattern.UserDiscredited = false
		} else {
			pattern.UserDiscredited = true
			pattern.UserCredited = false
		}
	} else if err8 == gorm.RecordNotFound {
		pattern.UserCredited = false
		pattern.UserDiscredited = false
	} else if err8 != nil {
		return pattern, &appError{err8, "user credited failed", http.StatusInternalServerError}
	}

	return pattern, nil
}

// save chart to valdiated table
func Discover(id string, uid int, json []byte, correlated bool) (Discovered, *appError) {
	discovered := Discovered{}

	if correlated {
		discovered.CorrelationId, _ = strconv.Atoi(id)
	} else {
		discovered.RelationId = id
	}

	discovered.DiscoveredId = 0
	discovered.Uid = uid
	discovered.Json = json
	discovered.Created = time.Now()
	discovered.Rating = 0
	discovered.Discredited = 0
	discovered.Credited = 0

	err = DB.Save(&discovered).Error
	if err != nil {
		return discovered, &appError{err, "could not save discovery", http.StatusInternalServerError}
	}
	return discovered, nil
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

	// @TODO: For future can see if these charts have been discovered yet or not
	// for i, v := range charts {
	// 	originid := tablename + "/" + strconv.Itoa(i) + "/" + v.ChartType + "/" + v.LabelX + "/" + v.LabelY
	// 	discovered := Discovered{}
	// 	err := DB.Where("relation_id = ?", originid).Find(&discovered).Error
	// 	if err == gorm.RecordNotFound {
	// 		v.Discovered = false
	// 	} else {
	// 		v.Discovered = true
	// 	}
	// }

	charts = charts[offset:last] // return marshalled slice
	return RelatedCharts{charts, totalCharts}, nil
}

// Look for new correlated charts, take the correlations and break them down into charting types, and return them along with their total count
// To return only existing charts use searchdepth = 0
func GetCorrelatedCharts(guid string, searchDepth int, offset int, count int) (CorrelatedCharts, *appError) {
	correlation := make([]Correlation, 0)
	charts := make([]CorrelationData, 0) ///empty slice for adding all possible charts
	od := OnlineData{}
	e := DB.Where("guid = ?", guid).Find(&od).Error
	if e != nil {
		return CorrelatedCharts{nil, 0}, &appError{nil, "Bad guid", http.StatusInternalServerError}
	}
	tableName := od.Tablename

	for i := 0; i < 10; i++ {
		go GenerateCorrelations(tableName, searchDepth)
	}

	err := DB.Where("tbl1 = ?", tableName).Order("abscoef DESC").Find(&correlation).Error
	if err != nil && err != gorm.RecordNotFound {
		return CorrelatedCharts{nil, 0}, &appError{nil, "Database query failed (TBL1)", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return CorrelatedCharts{nil, 0}, &appError{nil, "No correlated chart found", http.StatusNotFound}
	}

	for _, c := range correlation {
		var cd CorrelationData
		json.Unmarshal(c.Json, &cd)
		cd.CorrelationId = c.CorrelationId
		charts = append(charts, cd)
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

	charts = charts[offset:last] // return marshalled slice
	return CorrelatedCharts{charts, totalCharts}, nil
}

// As GetNew but get charts users have already voted on and return in an order based upon their absoulte ranking value
func GetDiscoveredCharts(tableName string, correlated bool, offset int, count int) (DiscoveredCharts, *appError) {
	discovered := make([]Discovered, 0)
	charts := make([]interface{}, 0)

	if correlated {
		c := Correlation{}
		d := Discovered{}

		joinStr := "LEFT JOIN " + c.TableName() + " ON " + d.TableName() + ".correlation_id = " + c.TableName() + ".correlation_id"
		whereStr := c.TableName() + ".tbl1 = ?"

		err := DB.Joins(joinStr).Where(whereStr, tableName).Order("rating DESC").Find(&discovered).Error
		if err != nil && err != gorm.RecordNotFound {
			return DiscoveredCharts{nil, 0}, &appError{nil, "Database query failed (JOIN)", http.StatusInternalServerError}
		} else if err == gorm.RecordNotFound {
			return DiscoveredCharts{nil, 0}, &appError{nil, "No valid chart found", http.StatusNotFound}
		}
	} else {
		tableName = tableName + "_%"
		err := DB.Model(Discovered{}).Where("relation_id LIKE ?", tableName).Order("rating DESC").Find(&discovered).Error
		if err != nil && err != gorm.RecordNotFound {
			return DiscoveredCharts{nil, 0}, &appError{nil, "Database query failed", http.StatusInternalServerError}
		} else if err == gorm.RecordNotFound {
			return DiscoveredCharts{nil, 0}, &appError{nil, "No valid chart found", http.StatusNotFound}
		}
	}

	for i, _ := range discovered {
		if correlated {
			var correlationData CorrelationData
			err1 := json.Unmarshal(discovered[i].Json, &correlationData)
			if err1 != nil {
				return DiscoveredCharts{nil, 0}, &appError{err1, "Json failed", http.StatusInternalServerError}
			}

			correlationData.CorrelationId = discovered[i].CorrelationId
			charts = append(charts, correlationData)
		} else {
			var tableData TableData
			err1 := json.Unmarshal(discovered[i].Json, &tableData)
			if err1 != nil {
				return DiscoveredCharts{nil, 0}, &appError{err1, "Json failed", http.StatusInternalServerError}
			}

			tableData.RelationId = discovered[i].RelationId
			charts = append(charts, tableData)
		}
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

	charts = charts[offset:last] // return slice
	return DiscoveredCharts{charts, totalCharts}, nil
}

// Get arrays of data for the types of charts requested (titles, descriptions, all the xy values etc)
// Determines what types of data are valid for any particular type of chart
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
			sql = fmt.Sprintf("SELECT %q AS x, SUM(%q) AS y FROM %q GROUP BY %q", names.X, names.X, guid, names.X)
			tmpTD.LabelY = "sum"
		} else {
			sql = fmt.Sprintf("SELECT %q AS x, COUNT(%q) AS y FROM %q GROUP BY %q", names.X, names.X, guid, names.X)
			tmpTD.LabelY = "count"
		}
	} else if chartType == "bubble" {
		sql = fmt.Sprintf("SELECT %q AS x, %q AS y, %q AS z FROM %q", names.X, names.Y, names.Z, guid)
	} else {
		sql = fmt.Sprintf("SELECT %q AS x, %q AS y FROM %q", names.X, names.Y, guid)
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

// Gnerate all possible permutations of xy columns
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

//determines strength of correlation
func CalcStrength(x float64) string {
	if x <= 0.1 {
		return "very low"
	} else if x <= 0.2 {
		return "low"
	} else if x <= 0.3 {
		return "quite low"
	} else if x <= 0.4 {
		return "medium"
	} else if x <= 0.6 {
		return "quite high"
	} else if x <= 0.7 {
		return "high"
	} else {
		return "very high"
	}
}

//////////////////////////////////////////////////////////////////////////
//////////// HTTP AND QUEUE FUNCTIONS TO CALL ABOVE METHODS///////////////
//////////////////////////////////////////////////////////////////////////

func GetChartHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	uid, err1 := GetUserID(session)
	if err1 != nil {
		http.Error(res, err1.Message, err1.Code)
		return ""
	}

	if params["tablename"] == "" {
		http.Error(res, "Invalid tablename", http.StatusBadRequest)
		return ""
	}

	tablenum, err := strconv.Atoi(params["tablenum"])
	if err != nil {
		http.Error(res, "Invalid tablenum parameter", http.StatusBadRequest)
		return ""
	}

	if params["type"] == "" {
		http.Error(res, "Invalid chart type", http.StatusBadRequest)
		return ""
	}

	if params["x"] == "" {
		http.Error(res, "Invalid x label", http.StatusBadRequest)
		return "I"
	}

	if params["y"] == "" {
		http.Error(res, "Invalid y label", http.StatusBadRequest)
		return ""
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
		return ""
	}

	return string(r)
}

func GetChartCorrelatedHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	cid, err := strconv.Atoi(params["cid"])
	if err != nil {
		http.Error(res, "Invalid id parameter", http.StatusBadRequest)
		return ""
	}

	uid, err1 := GetUserID(session)
	if err1 != nil {
		http.Error(res, err1.Message, err1.Code)
		return ""
	}

	result, error := GetChartCorrelated(cid, uid)
	if error != nil {
		http.Error(res, error.Message, error.Code)
		return ""
	}

	r, err2 := json.Marshal(result)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func GetRelatedChartsHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	var offset, count int
	var err error

	if params["offset"] == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(params["offset"])
		if err != nil {
			http.Error(res, "Invalid offset parameter", http.StatusBadRequest)
			return ""
		}
	}

	if params["count"] == "" {
		count = 1
	} else {
		count, err = strconv.Atoi(params["count"])
		if err != nil {
			http.Error(res, "Invalid count parameter", http.StatusBadRequest)
			return ""
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
		return ""
	}

	return string(r)
}

func GetCorrelatedChartsHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	var search, offset, count int
	var err error

	if params["offset"] == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(params["offset"])
		if err != nil {
			http.Error(res, "Invalid offset parameter", http.StatusBadRequest)
			return ""
		}
	}

	if params["count"] == "" {
		count = 3
	} else {
		count, err = strconv.Atoi(params["count"])
		if err != nil {
			http.Error(res, "Invalid count parameter", http.StatusBadRequest)
			return ""
		}
	}

	if params["search"] == "true" { ///default searchdepth when blank
		search = sd
	} else { // do not search when false so can return just what exist in table
		search = 0
	}

	result, error := GetCorrelatedCharts(params["tablename"], search, offset, count)
	if error != nil {
		http.Error(res, error.Message, error.Code)
		return error.Message
	}

	r, err1 := json.Marshal(result)
	if err1 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

func GetDiscoveredChartsHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
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
			return ""
		}
	}

	if params["count"] == "" {
		count = 3
	} else {
		count, err = strconv.Atoi(params["count"])
		if err != nil {
			http.Error(res, "Invalid count parameter", http.StatusBadRequest)
			return ""
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
		return ""
	}

	return string(r)
}

func GetChartQ(params map[string]string) string {
	uid, err1 := GetUserID(params["session"])
	if err1 != nil {
		return "invalid uid"
	}

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
	} else { // do not search when 0 so can return just what exist in table
		search = 0
	}

	result, err := GetCorrelatedCharts(params["tablename"], search, offset, count)
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

// Get charts and info awaiting user credit
func GetAwaitingCreditHttp(res http.ResponseWriter, req *http.Request) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	discovered := []Discovered{}
	charts := make([]interface{}, 0)

	query := DB.Select("json, priv_discovered.correlation_id, priv_discovered.relation_id, priv_discovered.discovered_id")
	query = query.Joins("LEFT JOIN priv_credits ON priv_discovered.discovered_id = priv_credits.discovered_id")
	query = query.Where("priv_credits.uid != ?", uid) //@todo check this in practice
	query = query.Order("random()")
	err1 := query.Find(&discovered).Error
	if err1 != nil && err1 != gorm.RecordNotFound {
		http.Error(res, "Failed to find discovered charts", http.StatusBadRequest)
		return ""
	}

	for _, d := range discovered {
		if d.RelationId == "" {
			var correlationData CorrelationData
			err2 := json.Unmarshal(d.Json, &correlationData)
			if err2 != nil {
				http.Error(res, "Failed to unmarshal json for discovered charts", http.StatusBadRequest)
				return ""
			}

			correlationData.CorrelationId = d.CorrelationId
			correlationData.Discovered = true

			type correlationExtender struct {
				CorrelationData
				Source  string `json:"source_title"`
				SourceX string `json:"source_X"`
				SourceY string `json:"source_Y"`
				SourceZ string `json:"source_Z, omitempty"`
			}

			var correlation Correlation
			err3 := DB.Where("correlation_id = ?", d.CorrelationId).Find(&correlation).Error
			if err3 != nil {
				http.Error(res, "Failed to find correlation data", http.StatusBadRequest)
				return ""
			}

			var c correlationExtender
			c.CorrelationData = correlationData
			c.Source = correlation.Tbl1
			c.SourceX = correlation.Col1
			c.SourceY = correlation.Col2
			c.SourceZ = correlation.Col3

			charts = append(charts, c)

		} else {
			var tableData TableData
			err2 := json.Unmarshal(d.Json, &tableData)
			if err2 != nil {
				http.Error(res, "Failed to unmarshal json for discovered charts", http.StatusBadRequest)
				return ""
			}

			tableData.RelationId = d.RelationId
			charts = append(charts, tableData)
		}
	}

	resp := map[string]interface{}{
		"charts": charts,
		"total":  len(charts),
	}
	response, _ := json.Marshal(resp)

	return string(response)
}

func GetTopRatedChartsHttp(res http.ResponseWriter, req *http.Request) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	discovered := []Discovered{}
	charts := make([]interface{}, 0)
	err = DB.Order("rating DESC").Limit(6).Find(&discovered).Error
	if err != nil {
		http.Error(res, "Failed to find discovered charts", http.StatusBadRequest)
		return ""
	}

	for _, d := range discovered {
		if d.RelationId == "" {
			var correlationData CorrelationData
			err1 := json.Unmarshal(d.Json, &correlationData)
			if err1 != nil {
				http.Error(res, "Failed to unmarshal json for discovered charts", http.StatusBadRequest)
				return ""
			}

			correlationData.CorrelationId = d.CorrelationId
			charts = append(charts, correlationData)
		} else {
			var tableData TableData
			err1 := json.Unmarshal(d.Json, &tableData)
			if err1 != nil {
				http.Error(res, "Failed to unmarshal json for discovered charts", http.StatusBadRequest)
				return ""
			}

			tableData.RelationId = d.RelationId
			charts = append(charts, tableData)
		}
	}

	r, err2 := json.Marshal(charts)
	if err2 != nil {
		http.Error(res, "Unable to parse JSON", http.StatusInternalServerError)
		return ""
	}

	return string(r)
}

// This function gets the extended infomation FROM the index, things like the notes are used
// in the "wiki" section of the page.
func GetChartInfoHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	if params["tablename"] == "" {
		http.Error(res, "There was no ID request", http.StatusBadRequest)
		return ""
	}

	index := Index{}
	err := DB.Where("LOWER(guid) LIKE LOWER(?)", params["tablename"]+"%").Find(&index).Error
	if err == gorm.RecordNotFound {
		return "[]"
	} else if err != nil {
		http.Error(res, "Could not find that data.", http.StatusNotFound)
		return ""
	}

	result := DataEntry{
		Name:  SanitizeString(index.Name),
		Title: SanitizeString(index.Title),
	}

	r, _ := json.Marshal(result)

	return string(r)
}
