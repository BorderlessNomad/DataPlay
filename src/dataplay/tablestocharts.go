package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/jinzhu/gorm"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type RelatedCharts struct {
	Charts []TableData `json:"charts"`
	Count  int         `json:"count"`
}

type DataEntry struct {
	Name       string `json:"name"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	SourceUrl  string `json:"sourceurl"`
	SourceName string `json:"sourcename"`
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
	Coefficient     float64     `json:"coefficient, omitempty"`
	Strength        string      `json:"statstrength, omitempty"`
	PrimOverview    string      `json:"overview1, omitempty"`
	SecoOverview    string      `json:"overview2, omitempty"`
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

var RelatedChartsCollection = make(map[string]map[int]TableData)

// given a small fraction of ratings there is a strong (95%) chance that the "real", final positive rating will be this value
// eg: gives expected (not necessarily current as there may have only been a few votes so far) value of positive ratings / total ratings
func RankCredits(credit int, discredit int) float64 {
	pos := float64(credit)
	tot := float64(credit + discredit)

	if tot == 0 {
		return 0
	}

	z := 1.96
	phat := pos / tot
	result := (phat + z*z/(2*tot) - z*math.Sqrt((phat*(1-phat)+z*z/(4*tot))/tot)) / (1 + z*z/tot)
	return result
}

// increment user discovered total for chart and rerank, return discovered id
func CreditChart(rcid string, uid int, credflag bool) (string, *appError) {
	t := time.Now()
	discovered := Discovered{}
	credit := Credit{}

	if strings.ContainsAny(rcid, "_") { // if a relation id
		rcid = strings.Replace(rcid, "_", "/", -1)
		err := DB.Where("relation_id = ?", rcid).Find(&discovered).Error
		if err != nil && err != gorm.RecordNotFound {
			return "", &appError{err, ", database query failed (relation_id)", http.StatusInternalServerError}
		}
	} else { // if a correlation id of type int
		cid, e := strconv.Atoi(rcid)
		if e != nil {
			return "", &appError{e, ", could not convert id to int", http.StatusInternalServerError}
		}

		err := DB.Where("correlation_id = ?", cid).Find(&discovered).Error
		if err != nil && err != gorm.RecordNotFound {
			return "", &appError{err, ", database query failed (correlation_id)", http.StatusInternalServerError}
		}
	}

	if credflag {
		discovered.Credited++
		Reputation(discovered.Uid, discCredit) // add points for discovery credit
		AddActivity(uid, "cc", t, discovered.DiscoveredId, 0)
	} else {
		discovered.Discredited++
		Reputation(discovered.Uid, discDiscredit) // remove points for discovery discredit
		AddActivity(uid, "dc", t, discovered.DiscoveredId, 0)
	}
	discovered.Rating = RankCredits(discovered.Credited, discovered.Discredited)
	err1 := DB.Save(&discovered).Error
	if err1 != nil {
		return "", &appError{err1, ", database query failed - credit chart (Save discovered)", http.StatusInternalServerError}
	}
	credit.DiscoveredId = discovered.DiscoveredId
	credit.Uid = uid
	credit.Created = t
	credit.ObservationId = 0 // not an observation
	credit.Credflag = credflag

	creditchk := Credit{}

	err2 := DB.Where("discovered_id = ?", credit.DiscoveredId).Where("uid = ?", credit.Uid).Where("observation_id = ?", credit.ObservationId).Find(&creditchk).Error
	if err2 == gorm.RecordNotFound {
		err3 := DB.Save(&credit).Error
		if err3 != nil {
			return "", &appError{err3, ", database query failed (Save credit)", http.StatusInternalServerError}
		}
	} else {
		credit.CreditId = creditchk.CreditId
		err4 := DB.Model(&creditchk).Update("credflag", credflag).Error
		if err4 != nil {
			return "", &appError{err4, ", database query failed (Update credit)", http.StatusInternalServerError}
		}
	}

	return strconv.Itoa(discovered.DiscoveredId), nil
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

	onlineData, err := GetOnlineDataByGuid(tablename)
	if err != nil {
		return pattern, &appError{err, "Database query failed (GUID)", http.StatusInternalServerError}
	}

	index, err := GetTableIndex(tablename)
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

		if chartType == "column" && x == "boroughs" && IsDateYear(v.Name) {
			xyz.Y = v.Name
			xyz.Ytype = v.Sqltype
			xyz.Z = y
		}
	}

	GenerateChartData(chartType, onlineData.Tablename, xyz, index, false)

	var chart TableData
	hashKey := onlineData.Tablename + "|" + chartType + "|" + xyz.X + "|" + xyz.Y + "|" + xyz.Z
	if _, ok := RelatedChartsCollection[hashKey]; ok {
		for _, ch := range RelatedChartsCollection[hashKey] {
			chart = ch
			break
		}
	}

	if len(chart.Values) < 1 {
		return pattern, &appError{err, "Not possible to plot this chart.", http.StatusBadRequest}
	}

	jByte, err1 := json.Marshal(chart)
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
		return pattern, &appError{err4, "Unable to retrieve user for related chart", http.StatusInternalServerError}
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
		return pattern, &appError{err5, "Unable to Find Crediting user", http.StatusInternalServerError}
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
	if err6 != nil && err6 != gorm.RecordNotFound {
		return pattern, &appError{err6, "Observation count failed.", http.StatusInternalServerError}
	}

	pattern.ChartData = chart
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
	if err8 != nil && err8 != gorm.RecordNotFound {
		return pattern, &appError{err8, "User credited failed.", http.StatusInternalServerError}
	} else if err8 == gorm.RecordNotFound {
		pattern.UserCredited = false
		pattern.UserDiscredited = false
	} else if cred.Credflag == true {
		pattern.UserCredited = true
		pattern.UserDiscredited = false
	} else {
		pattern.UserCredited = false
		pattern.UserDiscredited = true
	}

	return pattern, nil
}

// use the id relating to the record stored in the generated correlations table to return the json with the specific chart info
func GetChartCorrelated(cid int, uid int) (PatternInfo, *appError) {
	pattern := PatternInfo{}
	var cd CorrelationData

	correlation := Correlation{}
	err := DB.Where("correlation_id = ?", cid).Find(&correlation).Error
	if err != nil && err != gorm.RecordNotFound {
		return pattern, &appError{err, "Database query failed (ID)", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return pattern, &appError{err, "No correlated chart found", http.StatusNotFound}
	}

	//if undiscovered add to the discovered table as an initial discovery
	discovered := Discovered{}
	err1 := DB.Where("correlation_id = ?", cid).Find(&discovered).Error
	if err1 == gorm.RecordNotFound {
		Discover(strconv.Itoa(cid), uid, correlation.Json, true)
	}

	user := User{}
	err2 := DB.Where("uid = ?", discovered.Uid).Find(&user).Error
	if err2 != nil && err2 != gorm.RecordNotFound {
		return pattern, &appError{err2, "Unable to retrieve User for correlated chart.", http.StatusInternalServerError}
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

	pattern.ChartData = cd
	pattern.PatternID = strconv.Itoa(discovered.CorrelationId)
	pattern.DiscoveredID = discovered.DiscoveredId
	pattern.Discoverer = user.Username
	pattern.DiscoveryDate = discovered.Created
	pattern.Creditors = creditors
	pattern.Discreditors = discreditors
	pattern.PrimarySource = cd.Table1.Title
	pattern.SecondarySource = cd.Table2.Title
	pattern.Coefficient = correlation.Coef
	pattern.Strength = CalcStrength(correlation.Abscoef)
	pattern.Observations = count
	pattern.PrimOverview = correlation.Tbl1
	pattern.SecoOverview = correlation.Tbl2

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

	err := DB.Save(&discovered).Error
	if err != nil {
		return discovered, &appError{err, "could not save discovery", http.StatusInternalServerError}
	}
	return discovered, nil
}

// generate all the potentially valid charts that relate to a single tablename, add apt charting types,
// and return them along with their total count and whether they've been discovered
func GetRelatedCharts(tablename string, offset int, count int) (RelatedCharts, *appError) {
	columns := FetchTableCols(tablename)      //array column names
	xyNames := XYPermutations(columns, false) // get all possible valid permuations of columns as X & Y
	onlineData, err := GetOnlineDataByGuid(tablename)
	if err != nil {
		return RelatedCharts{nil, 0}, &appError{err, "Database query failed (GUID)", http.StatusInternalServerError}
	}

	index, err := GetTableIndex(tablename)
	if err != nil && err != gorm.RecordNotFound {
		return RelatedCharts{nil, 0}, &appError{err, "Database query failed (GUID)", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return RelatedCharts{nil, 0}, &appError{err, "No related chart found", http.StatusNotFound}
	}

	var xyPie XYVal

	for _, v := range columns { // create single column pie charts
		xyPie.X = v.Name
		xyPie.Xtype = v.Sqltype
		GenerateChartData("pie", onlineData.Tablename, xyPie, index, true)
	}

	for _, v := range xyNames { // create all other types of chart
		if v.Xtype != "varchar" && v.Ytype != "varchar" {
			GenerateChartData("column dated", onlineData.Tablename, v, index, true)
		}

		if v.Xtype == "varchar" && v.Ytype == "varchar" { // stacked or scatter charts if string v string values
			// GenerateChartData("stacked column", onlineData.Tablename, v, index, true)
			// GenerateChartData("scatter", onlineData.Tablename, v, index, true)
		}

		if !(v.Xtype == "varchar" && v.Ytype == "date") || !(v.Xtype == "date" && v.Ytype == "varchar") { // column and row charts for all that are not string v string values and are not date v string or string v date values
			// GenerateChartData("row", onlineData.Tablename, v, index, true)
			if v.Ytype != "varchar" && v.Ytype != "date" { // no string values for y axis on column charts
				GenerateChartData("column", onlineData.Tablename, v, index, true)
			}
		}

		if v.Ytype != "varchar" && (v.Xtype == "date" || (v.Xtype != "varchar" && IsDateYear(v.X))) { // line chart cannot be based on strings or have date on the Y axis
			GenerateChartData("line", onlineData.Tablename, v, index, true)
		}
	}

	if len(columns) > 2 { // if there's more than 2 columns grab a 3rd variable for bubble charts
		xyNames = XYPermutations(columns, true) // set z flag to true to get all possible valid permuations of columns as X, Y & Z
		for _, v := range xyNames {
			if v.Xtype == "date" && IsNumeric(v.Ytype) && IsNumeric(v.Ztype) {
				GenerateChartData("bubble", onlineData.Tablename, v, index, true)
			}
		}
	}

	uniqueCharts := make([]TableData, 0)
	check := make(map[string]bool)
	for key, charts := range RelatedChartsCollection {
		hashKey := strings.Split(key, "|")
		if hashKey[0] == onlineData.Tablename {
			for _, chart := range charts {
				if !check[hashKey[1]+chart.Title] {
					check[hashKey[1]+chart.Title] = true
					uniqueCharts = append(uniqueCharts, chart)
				}
			}
		}
	}

	totalCharts := len(uniqueCharts)
	if offset > totalCharts {
		return RelatedCharts{nil, 0}, &appError{nil, fmt.Sprintf("Offset value out of bounds (Max: %d)", totalCharts), http.StatusBadRequest}
	}

	last := offset + count
	if offset != 0 && last > totalCharts {
		return RelatedCharts{nil, 0}, &appError{nil, fmt.Sprintf("Count value out of bounds (Max: %d)", totalCharts-offset), http.StatusBadRequest}
	} else if offset == 0 && (last > totalCharts || count == 0) {
		last = totalCharts
	}

	// randomise order
	sort.Sort(MixRepeatably(uniqueCharts))

	uniqueCharts = uniqueCharts[offset:last] // return marshalled slice

	return RelatedCharts{uniqueCharts, totalCharts}, nil
}

// Look for new correlated charts, take the correlations and break them down into charting types, and return them along with their total count
// To return only existing charts use searchdepth = 0
func GetCorrelatedCharts(guid string, searchDepth int, offset int, count int, reduce bool) (CorrelatedCharts, *appError) {
	correlation := make([]Correlation, 0)
	charts := make([]CorrelationData, 0) ///empty slice for adding all possible charts

	onlineData, e := GetOnlineDataByGuid(guid)
	if e != nil {
		return CorrelatedCharts{nil, 0}, &appError{nil, "Invalid or Empty GUID", http.StatusBadRequest}
	}

	// for i := 0; i < 30; i++ {
	// 	go GenerateCorrelations(tableName, searchDepth)
	// }
	// // time.Sleep(5 * time.Second) // WHY??????

	// @todo Mayur Run once and generate all possible correlations
	go GenerateCorrelations(onlineData.Tablename, searchDepth)

	err := DB.Where("tbl1 = ?", onlineData.Tablename).Order("abscoef DESC").Find(&correlation).Error
	if err != nil && err != gorm.RecordNotFound {
		return CorrelatedCharts{nil, 0}, &appError{nil, "Database query failed (TBL1)", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		DB.Order("random()").Find(&correlation)
	}

	for _, c := range correlation {
		var cd CorrelationData
		json.Unmarshal(c.Json, &cd)
		cd.CorrelationId = c.CorrelationId
		cd.Coefficient = c.Coef
		cd.Strength = CalcStrength(c.Abscoef)

		if reduce {
			cd.Table1.Values = ReduceXYValues(cd.Table1.Values)
			cd.Table2.Values = ReduceXYValues(cd.Table2.Values)
		}

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

	origins := make([]int, 0)
	for _, c := range charts {
		c.Discovered = false
		origins = append(origins, c.CorrelationId)
	}

	if len(origins) == 0 {
		return CorrelatedCharts{nil, 0}, &appError{nil, "Unable to Generate Correlations (Origins).", http.StatusInternalServerError}
	}

	discoveries := []Discovered{}
	err1 := DB.Where("correlation_id IN (?)", origins).Find(&discoveries).Error
	if err1 != nil && err1 != gorm.RecordNotFound {
		return CorrelatedCharts{nil, 0}, &appError{err1, "Database query failed (Discoveries).", http.StatusInternalServerError}
	}

	for _, discovery := range discoveries {
		for _, c := range charts {
			if c.CorrelationId == discovery.CorrelationId {
				c.Discovered = true
				break
			}
		}
	}

	// sort by coefficient
	sort.Sort(SortByCoefficient(charts))

	if count > 30 {
		last = offset + 30
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
func GenerateChartData(chartType string, guid string, names XYVal, ind Index, reduce bool) {
	var tmpTD TableData
	var tmpXY XYVal
	tmpTD.ChartType = chartType
	tmpTD.Title = SanitizeString(ind.Title)
	tmpTD.Desc = SanitizeString(ind.Notes)
	tmpTD.LabelX = names.X
	tmpTD.LabelYLong = SanitizeString(ind.Notes)

	var dx, dy time.Time
	var fx, fy, fz float64
	var vx, vy string
	rowAmt := 0

	hashKey := guid + "|" + chartType + "|" + names.X + "|" + names.Y + "|" + names.Z
	if _, ok := RelatedChartsCollection[hashKey]; ok {
		return
	}

	if chartType == "bubble" {
		hashKeyBubble := guid + "|" + chartType + "|" + names.X + "|" + names.Z + "|" + names.Y
		if _, ok := RelatedChartsCollection[hashKeyBubble]; ok {
			return
		}
	}

	sql := ""
	validChartType := ""
	if chartType == "pie" {
		if IsNumeric(names.Xtype) {
			sql = fmt.Sprintf("SELECT * FROM %q ORDER BY %q", guid, names.X)
		} else if IsDateYear(names.X) && names.Ytype == "" {
			validChartType = "pie dated"
			y, _ := strconv.Atoi(names.Y)
			sql = fmt.Sprintf("SELECT * FROM %q WHERE EXTRACT(year FROM %q) = %d", guid, names.X, y)
		}
	} else if chartType == "bubble" && len(names.Z) > 0 {
		sql = fmt.Sprintf("SELECT %q AS x, %q AS y, %q AS z FROM %q", names.X, names.Y, names.Z, guid)
	} else if (chartType == "column" && names.X == "boroughs") || (chartType == "column dated" && IsDateYear(names.X) && IsNumeric(names.Ytype)) {
		validChartType = "column dated"
		if y, err := strconv.Atoi(names.Z); err == nil {
			sql = fmt.Sprintf("SELECT * FROM %q WHERE EXTRACT(year FROM %q) = %d", guid, names.Y, y)
		} else {
			sql = fmt.Sprintf("SELECT * FROM %q", guid)
		}

	} else if IsDateYear(names.X) {
		sql = fmt.Sprintf("SELECT %q AS x, %q AS y FROM %q", names.X, names.Y, guid)
	}

	rows, err := DB.Raw(sql).Rows()
	if err != nil {
		return
	}

	defer rows.Close()

	// columnNames, _ := rows.Columns()
	columnNames := FetchTableCols(guid)
	columns := make([]interface{}, len(columnNames))
	columnPointers := make([]interface{}, len(columnNames))
	for i := 0; i < len(columnNames); i++ {
		columnPointers[i] = &columns[i]
	}

	RelatedChartsCollection[hashKey] = make(map[int]TableData, 0)

	if chartType == "bubble" {
		tmpTD.LabelY = names.Y
		tmpTD.LabelZ = names.Z

		for rows.Next() {
			if IsNumeric(names.Xtype) && IsNumeric(names.Ytype) && IsNumeric(names.Ztype) {
				rows.Scan(&fx, &fy, &fz)

				tmpXY.X = FloatToString(fx)
				tmpXY.Y = FloatToString(fy)
				tmpXY.Z = FloatToString(fz)

				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if IsNumeric(names.Xtype) && names.Ytype == "varchar" && IsNumeric(names.Ztype) {
				rows.Scan(&fx, &vy, &fz)

				tmpXY.X = FloatToString(fx)
				tmpXY.Y = vy
				tmpXY.Z = FloatToString(fz)

				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if names.Xtype == "varchar" && IsNumeric(names.Ytype) && IsNumeric(names.Ztype) {
				rows.Scan(&vx, &fy, &fz)

				tmpXY.X = vx
				tmpXY.Y = FloatToString(fy)
				tmpXY.Z = FloatToString(fz)

				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if names.Xtype == "date" && IsNumeric(names.Ytype) && IsNumeric(names.Ztype) {
				rows.Scan(&dx, &fy, &fz)

				tmpXY.X = strconv.Itoa(dx.Year())
				tmpXY.Y = FloatToString(fy)
				tmpXY.Z = FloatToString(fz)

				tmpTD.Title = tmpTD.Desc + " in " + tmpTD.LabelY + " & " + tmpTD.LabelZ
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			}
		}

		if ValueCheck(tmpTD) && NegCheck(tmpTD) {
			RelatedChartsCollection[hashKey][0] = tmpTD
		}
	} else if chartType == "pie" { // single column pie chart x = type, y = count
		for rows.Next() {
			if IsNumeric(names.Xtype) || (names.Xtype == "date" && names.Ytype == "") {
				rows.Scan(columnPointers...)

				tmpTDPie := tmpTD
				var date time.Time

				for k, v := range columnNames {
					if !IsDateYear(v.Name) {
						length := len(tmpTDPie.Values)
						if length > 0 && tmpTDPie.Values[length-1].X == v.Name {
							return
						}

						tmpXY.X = v.Name
						switch value := columns[k].(type) {
						case int64:
							tmpXY.Y = strconv.FormatInt(value, 10)
						case float64:
							tmpXY.Y = FloatToString(value)
						case string:
							tmpXY.Y = value
						default:
							continue
						}

						tmpTDPie.Values = append(tmpTDPie.Values, tmpXY)
					} else if v.Sqltype == "date" {
						date = columns[k].(time.Time)

						tmpTDPie.LabelX = v.Name
						tmpTDPie.LabelY = strconv.Itoa(date.Year())
						tmpTDPie.Title = tmpTD.Desc + " in " + v.Name + " " + tmpTDPie.LabelY
					}
				}

				if ValueCheck(tmpTDPie) && NegCheck(tmpTDPie) {
					RelatedChartsCollection[hashKey][date.Year()] = tmpTDPie
				}
			}
		}
	} else if chartType == "column dated" || validChartType == "column dated" { // single column pie chart x = type, y = count
		for rows.Next() {
			if names.Ytype == "date" || (names.Xtype == "date" && IsNumeric(names.Ytype)) {
				rows.Scan(columnPointers...)

				var tmpTDColDated TableData
				var tmpXY XYVal
				var date time.Time
				var year int
				tmpTDColDated.ChartType = "column"

				for k, v := range columnNames {
					if !IsDateYear(v.Name) {
						length := len(tmpTDColDated.Values)
						if length > 0 && tmpTDColDated.Values[length-1].X == v.Name {
							return
						}

						tmpXY.X = v.Name
						switch value := columns[k].(type) {
						case int64:
							tmpXY.Y = strconv.FormatInt(value, 10)
						case float64:
							tmpXY.Y = FloatToString(value)
						case string:
							tmpXY.Y = value
						default:
							continue
						}

						tmpTDColDated.Values = append(tmpTDColDated.Values, tmpXY)
					} else if v.Sqltype == "date" {
						date = columns[k].(time.Time)
						year = date.Year()

						tmpTDColDated.LabelX = "boroughs" // @todo: mayur
						tmpTDColDated.LabelY = strconv.Itoa(year)
						tmpTDColDated.LabelYLong = tmpTD.LabelYLong
						tmpTDColDated.Title = tmpTD.Desc + " in " + v.Name + " " + strconv.Itoa(year)
					}
				}

				if ValueCheck(tmpTDColDated) && NegCheck(tmpTDColDated) {
					RelatedChartsCollection[hashKey][year] = tmpTDColDated
				}
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
			} else if names.Xtype == "date" && IsNumeric(names.Ytype) {
				rows.Scan(&dx, &fy)
				tmpXY.X = (dx.String()[0:10])
				tmpXY.Y = FloatToString(fy)
				tmpTD.Title = tmpTD.Desc + " in " + names.Y
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if names.Xtype == "date" && names.Ytype == "varchar" {
				rows.Scan(&dx, &vy)
				tmpXY.X = (dx.String()[0:10])
				tmpXY.Y = vy
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if IsNumeric(names.Xtype) && names.Ytype == "date" {
				rows.Scan(&fx, &dy)
				tmpXY.X = FloatToString(fx)
				tmpXY.Y = (dy.String()[0:10])
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if IsNumeric(names.Xtype) && IsNumeric(names.Ytype) && names.Xtype != names.Ytype {
				rows.Scan(&fx, &fy)
				tmpXY.X = FloatToString(fx)
				tmpXY.Y = FloatToString(fy)
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if IsNumeric(names.Xtype) && names.Ytype == "varchar" {
				rows.Scan(&fx, &vy)
				tmpXY.X = FloatToString(fx)
				tmpXY.Y = vy
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if names.Xtype == "varchar" && names.Ytype == "date" {
				rows.Scan(&vx, &dy)
				tmpXY.X = vx
				tmpXY.Y = (dy.String()[0:10])
				tmpTD.Values = append(tmpTD.Values, tmpXY)
			} else if names.Xtype == "varchar" && IsNumeric(names.Ytype) {
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
					RelatedChartsCollection[hashKey][0] = tmpTD
				}
			}
		} else if ValueCheck(tmpTD) {
			if chartType == "column" && len(tmpTD.Values) > 50 {
				return
			}

			if reduce {
				tmpTD.Values = ReduceXYValues(tmpTD.Values)
			}

			RelatedChartsCollection[hashKey][0] = tmpTD
		}
	}

	err = rows.Err() // get any error encountered during iteration
}

func IsDateYear(str string) bool {
	if strings.Contains(strings.ToLower(str), "date") || strings.Contains(strings.ToLower(str), "year") {
		return true
	}

	return false
}

// Gnerate all possible permutations of xy columns
func XYPermutations(columns []ColType, bubble bool) []XYVal {
	length := len(columns)
	maxIter := 100
	var xyNames []XYVal
	var xyzNames []XYVal
	var tmpXY XYVal

	randomColumns := make([]ColType, 0)
	checkedColumns := make([]bool, length)

	for k, v := range columns {
		if v.Sqltype == "date" {
			randomColumns = append(randomColumns, v)
			checkedColumns[k] = true
		}
	}

	if length > maxIter {
		for i := 0; i < maxIter; i++ {
			randomIndex := rand.Intn(length - 1)
			if checkedColumns[randomIndex] {
				i--
				continue
			}

			randomColumns = append(randomColumns, columns[randomIndex])
			checkedColumns[randomIndex] = true
		}

		columns = randomColumns
		length = len(columns)
	}

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
	if len(t.Values) <= 0 {
		return false
	}

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
	}

	return false
}

// checks whether any X axis values are negative as bubble won't plot if they are
func NegCheck(t TableData) bool {
	for _, v := range t.Values {
		x, _ := strconv.Atoi(v.X)
		if x < 0 {
			return false
		}

		y, _ := strconv.Atoi(v.Y)
		if y < 0 {
			return false
		}

		z, _ := strconv.Atoi(v.Z)
		if z < 0 {
			return false
		}
	}

	return true
}

/**
 * @brief Determines strength of correlation
 */
func CalcStrength(x float64) string {
	if x < 0.25 {
		return "very low"
	} else if x < 0.33 {
		return "low"
	} else if x < 0.5 {
		return "medium"
	} else if x < 0.75 {
		return "high"
	}

	return "very high"
}

//////////////////////////////////////////////////////////////////////////
//////////// HTTP FUNCTIONS TO CALL ABOVE METHODS///////////////
//////////////////////////////////////////////////////////////////////////
func CreditChartHttp(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	credflag := false
	rcid := ""

	if params["credflag"] == "" { // if no credflag then skip credit and just return discovered id
		http.Error(res, "Missing credflag", http.StatusBadRequest)
		return ""
	} else {
		credflag, _ = strconv.ParseBool(params["credflag"])
	}

	if params["rcid"] == "" {
		http.Error(res, "No Relation/Correlation ID provided.", http.StatusBadRequest)
		return ""
	} else {
		rcid = params["rcid"]
	}

	uid, err1 := GetUserID(session)
	if err1 != nil {
		http.Error(res, err1.Message, err1.Code)
		return ""
	}

	result, err2 := CreditChart(rcid, uid, credflag)
	if err2 != nil {
		msg := ""
		if credflag {
			msg = "Could not credit chart" + err2.Message
		} else {
			msg = "Could not discredit chart" + err2.Message
		}

		http.Error(res, err2.Message+msg, http.StatusBadRequest)
		return msg
	}

	return result
}

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
			return ""
		}
	} else {
		result, err2 = GetChart(params["tablename"], tablenum, params["type"], uid, params["x"], params["y"], params["z"])
		if err2 != nil {
			http.Error(res, err2.Message, err2.Code)
			return ""
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
		search = 50
	} else { // do not search when false so can return just what exist in table
		search = 0
	}

	result, error := GetCorrelatedCharts(params["tablename"], search, offset, count, true)
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
	query = query.Limit(5)
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

			c.CorrelationData.Table1.Values = ReduceXYValues(c.CorrelationData.Table1.Values)
			c.CorrelationData.Table2.Values = ReduceXYValues(c.CorrelationData.Table2.Values)

			charts = append(charts, c)
		} else {
			var tableData TableData
			err2 := json.Unmarshal(d.Json, &tableData)
			if err2 != nil {
				http.Error(res, "Failed to unmarshal json for discovered charts", http.StatusBadRequest)
				return ""
			}

			tableData.Values = ReduceXYValues(tableData.Values)

			tableData.RelationId = d.RelationId
			charts = append(charts, tableData)
		}
	}

	for i := range charts {
		j := rand.Intn(i + 1)
		charts[i], charts[j] = charts[j], charts[i]
	}

	resp := map[string]interface{}{
		"charts": charts,
		"total":  len(charts),
	}
	response, _ := json.Marshal(resp)

	return string(response)
}

func ReduceXYValues(originalValues []XYVal) []XYVal {
	if len(originalValues) > 100 {
		values := make([]XYVal, 0)
		rate := len(originalValues) / 100
		for k, v := range originalValues {
			if k == 0 || k == len(originalValues)-1 || k%rate == 0 {
				values = append(values, v)
			}
		}

		return values
	}

	return originalValues
}

func GetTopRatedChartsHttp(res http.ResponseWriter, req *http.Request) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter", http.StatusBadRequest)
		return ""
	}

	discovered := []Discovered{}
	charts := make([]interface{}, 0)
	err := DB.Order("rating DESC").Limit(6).Find(&discovered).Error
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

			correlationData.Table1.Values = ReduceXYValues(correlationData.Table1.Values)
			correlationData.Table2.Values = ReduceXYValues(correlationData.Table2.Values)

			charts = append(charts, correlationData)
		} else {
			var tableData TableData
			err1 := json.Unmarshal(d.Json, &tableData)
			if err1 != nil {
				http.Error(res, "Failed to unmarshal json for discovered charts", http.StatusBadRequest)
				return ""
			}

			tableData.RelationId = d.RelationId

			tableData.Values = ReduceXYValues(tableData.Values)

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

	table := strings.ToLower(strings.Trim(params["tablename"], " "))
	index, err := GetTableIndex(table)
	if err != nil && err != gorm.RecordNotFound {
		http.Error(res, fmt.Sprintf("Database Query Failed (%q)", table), http.StatusInternalServerError)
		return ""
	} else if err == gorm.RecordNotFound {
		table = "%" + table + "%"
		err = DB.Where("LOWER(guid) LIKE LOWER(?)", table).Find(&index).Error
		if err != nil && err != gorm.RecordNotFound {
			http.Error(res, fmt.Sprintf("Database Query Failed (%q)", table), http.StatusInternalServerError)
			return ""
		} else if err == gorm.RecordNotFound {
			http.Error(res, "No matching data found", http.StatusBadRequest)
			return ""
		}
	}

	result := DataEntry{
		Name:       SanitizeString(index.Name),
		Title:      SanitizeString(index.Title),
		Desc:       SanitizeString(index.Desc),
		SourceUrl:  SanitizeString(index.SourceUrl),
		SourceName: SanitizeString(index.SourceName),
	}

	r, _ := json.Marshal(result)

	return string(r)
}

func IsNumeric(str string) bool {
	if str == "float" || str == "integer" || str == "real" {
		return true
	}

	return false
}
