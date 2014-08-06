package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/jinzhu/gorm"
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type cmeth int

const ( //go version of enum
	P cmeth = iota
	S
	V
)

type RelatedCharts struct {
	Charts []TableData
	Count  int
}

type CorrelationData struct {
	Method    string    `json:"method"`
	ChartType []string  `json:"validCharts"`
	Table1    TableData `json:"table1"`
	Table2    TableData `json:"table2"`
	Table3    TableData `json:"table3"`
}

type TableData struct {
	ChartType string  `json:"chart,omitempty"`
	Title     string  `json:"title,omitempty"`
	Desc      string  `json:"desc,omitempty"`
	LabelX    string  `json:"xLabel,omitempty"`
	LabelY    string  `json:"yLabel,omitempty"`
	Values    []XYVal `json:"values,omitempty"`
}

type XYVal struct {
	X     string `json:"x,omitempty"`
	Y     string `json:"y,omitempty"`
	Xtype string `json:"-"`
	Ytype string `json:"-"`
}

type DateVal struct {
	Date  time.Time
	Value float64
}

type FromTo struct {
	From time.Time
	To   time.Time
}

// Take in table and column names and a threshold for the looping and get correlated tables
// Use 3:1:1 for Spurious to Pearson and Visual as Spurious less likely to find correlations
func GenerateCorrelations(table1 string, valCol1 string, dateCol1 string, thresh int) string {
	if table1 == "" || valCol1 == "" || dateCol1 == "" {
		return ""
	}
	m := make(map[string]string)
	m["table1"], m["dateCol1"], m["valCol1"] = table1, dateCol1, valCol1
	c := P

	for i := 0; i < thresh; i++ {
		r := i % 5

		if r == 0 {
			c = P
		} else if r == 1 {
			c = V
		} else {
			c = S
		}

		GenerateCorrelation(m, c)
	}

	return "win"
}

// Take in table and column values and a correlation type, generate some more random tables and check for a pre-existing correlation.
// If a correlation for the generated tables combination doesn't exist, attempt to generate a new correlation and return that
func GenerateCorrelation(m map[string]string, c cmeth) string {

	cor := Correlation{}
	var jsonData []string
	method := ""
	cd := new(CorrelationData)
	nameChk := GetRandomNames(m, c)

	if nameChk {

		if c == P {
			err := DB.Model(&cor).Where("tbl1 = ?", m["table1"]).Where("col1 = ?", m["valCol1"]).Where("tbl2 = ?", m["table2"]).Where("col2 = ?", m["valCol2"]).Where("method = ?", "Pearson").Pluck("json", &jsonData).Error
			check(err)
			method = "Pearson"
		} else if c == S {
			err := DB.Model(&cor).Where("tbl1 = ?", m["table1"]).Where("col1 = ?", m["valCol1"]).Where("tbl2 = ?", m["table2"]).Where("col2 = ?", m["valCol2"]).Where("tbl3 = ?", m["table3"]).Where("col3 = ?", m["valCol3"]).Where("method = ?", "Spurious").Pluck("json", &jsonData).Error
			check(err)
			method = "Spurious"
		} else if c == V {
			err := DB.Model(&cor).Where("tbl1 = ?", m["table1"]).Where("col1 = ?", m["valCol1"]).Where("tbl2 = ?", m["table2"]).Where("col2 = ?", m["valCol2"]).Where("method = ?", "Visual").Pluck("json", &jsonData).Error
			check(err)
			method = "Visual"
		}

		if jsonData == nil {
			cf := GetCoef(m, c, cd)

			if cf != 0 { //Save the correlation if one is generated
				m["method"] = method
				SaveCorrelation(m, c, cf, cd)
			}

			if c == P {
				err := DB.Model(&cor).Where("tbl1 = ?", m["table1"]).Where("col1 = ?", m["valCol1"]).Where("tbl2 = ?", m["table2"]).Where("col2 = ?", m["valCol2"]).Where("method = ?", "Pearson").Pluck("json", &jsonData).Error
				check(err)
			} else if c == S {
				err := DB.Model(&cor).Where("tbl1 = ?", m["table1"]).Where("col1 = ?", m["valCol1"]).Where("tbl2 = ?", m["table2"]).Where("col2 = ?", m["valCol2"]).Where("tbl3 = ?", m["table3"]).Where("col3 = ?", m["valCol3"]).Where("method = ?", "Spurious").Pluck("json", &jsonData).Error
				check(err)
			} else if c == V {
				err := DB.Model(&cor).Where("tbl1 = ?", m["table1"]).Where("col1 = ?", m["valCol1"]).Where("tbl2 = ?", m["table2"]).Where("col2 = ?", m["valCol2"]).Where("method = ?", "Visual").Pluck("json", &jsonData).Error
				check(err)
			}
		}
	}

	if len(jsonData) > 0 {
		return jsonData[0]
	} else {
		return ""
	}

}

// Generate a coefficient (if data allows), based on requested correlation type
func GetCoef(m map[string]string, c cmeth, cd *CorrelationData) float64 {
	if len(m) == 0 {
		return 0.0
	}

	var bucketRange []FromTo
	var xBuckets []float64
	var yBuckets []float64
	var zBuckets []float64
	var cf float64

	if c == P {
		charts := []string{"bar", "column", "line", "scatter"}
		(*cd).ChartType = charts
		x := ExtractDateVal(m["table1"], m["dateCol1"], m["valCol1"])
		y := ExtractDateVal(m["table2"], m["dateCol2"], m["valCol2"])
		fromX, toX, rngX := DetermineRange(x)
		fromY, toY, rngY := DetermineRange(y)

		if rngX == 0 || rngY == 0 {
			return 0
		}

		bucketRange = GetIntersect(&fromX, &toX, &rngX, fromY, toY, rngY)
		xBuckets = FillBuckets(x, bucketRange)
		yBuckets = FillBuckets(y, bucketRange)
		cf = Pearson(xBuckets, yBuckets)

	} else if c == S {
		charts := []string{"bubble", "line", "scatter", "stacked"}
		(*cd).ChartType = charts
		x := ExtractDateVal(m["table2"], m["dateCol2"], m["valCol2"])
		y := ExtractDateVal(m["table3"], m["dateCol3"], m["valCol3"])
		z := ExtractDateVal(m["table1"], m["dateCol1"], m["valCol1"])
		fromX, toX, rngX := DetermineRange(x)
		fromY, toY, rngY := DetermineRange(y)
		fromZ, toZ, rngZ := DetermineRange(z)

		if rngX == 0 || rngY == 0 || rngZ == 0 {
			return 0
		}

		_ = GetIntersect(&fromX, &toX, &rngX, fromY, toY, rngY)
		bucketRange = GetIntersect(&fromX, &toX, &rngX, fromZ, toZ, rngZ)
		xBuckets = FillBuckets(x, bucketRange)
		yBuckets = FillBuckets(y, bucketRange)
		zBuckets = FillBuckets(z, bucketRange)
		cf = Spurious(xBuckets, yBuckets, zBuckets)

	} else if c == V {
		charts := []string{"bar", "column", "line", "scatter"}
		(*cd).ChartType = charts
		x := ExtractDateVal(m["table1"], m["dateCol1"], m["valCol1"])
		y := ExtractDateVal(m["table2"], m["dateCol2"], m["valCol2"])
		fromX, toX, rngX := DetermineRange(x)
		fromY, toY, rngY := DetermineRange(y)

		if rngX == 0 || rngY == 0 {
			return 0
		}

		bucketRange = GetIntersect(&fromX, &toX, &rngX, fromY, toY, rngY)
		xBuckets = FillBuckets(x, bucketRange)
		yBuckets = FillBuckets(y, bucketRange)
		n := len(xBuckets)
		var high []float64
		var low []float64

		for i, v := range xBuckets {
			if v > yBuckets[i] {
				high = append(high, v)
				low = append(low, yBuckets[i])
			} else {
				high = append(high, yBuckets[i])
				low = append(low, v)
			}
		}

		highTotal := 0
		lowTotal := 0

		for i := 0; i < n; i++ {
			days := dayNum(bucketRange[i].To) - dayNum(bucketRange[i].From)
			highTotal = highTotal + days*int(high[i])
			lowTotal = lowTotal + days*int(low[i])
		}

		if highTotal == 0 || lowTotal == 0 {
			return 0
		}

		cf = 1 / (float64(highTotal) / float64(lowTotal))

	} else {
		return 0
	}

	labels := LabelGen(bucketRange)
	n := len(bucketRange)
	values1 := make([]XYVal, n)
	values2 := make([]XYVal, n)
	values3 := make([]XYVal, n)

	for i, v := range labels {
		values1[i].X = v
		values1[i].Y = strconv.FormatFloat(xBuckets[i], 'f', -1, 64)
		values2[i].X = v
		values2[i].Y = strconv.FormatFloat(yBuckets[i], 'f', -1, 64)

		if c == S {
			values3[i].X = v
			values3[i].Y = strconv.FormatFloat(zBuckets[i], 'f', -1, 64)
		}
	}

	(*cd).Table1.Values = values1
	(*cd).Table2.Values = values2
	if c == S {
		(*cd).Table3.Values = values3
	}
	return cf
}

//Create a json string containing all the data needed for generating a graph and then insert this and all the other correlation info into the correlations table
func SaveCorrelation(m map[string]string, c cmeth, cf float64, cd *CorrelationData) string {
	ind1 := Index{}
	ind2 := Index{}
	ind3 := Index{}

	guid1 := NameToGuid(m["table1"])
	guid2 := NameToGuid(m["table2"])
	guid3 := NameToGuid(m["table3"])

	err1 := DB.Model(&ind1).Where("guid= ?", guid1).Find(&ind1).Error
	check(err1)
	err2 := DB.Model(&ind2).Where("guid= ?", guid2).Find(&ind2).Error
	check(err2)

	if c == S {
		err3 := DB.Model(&ind3).Where("guid= ?", guid3).Find(&ind3).Error
		check(err3)
	}

	(*cd).Method = m["method"]
	(*cd).Table1.Title = ind1.Title
	(*cd).Table2.Title = ind2.Title
	(*cd).Table1.Desc = ind1.Notes
	(*cd).Table2.Desc = ind2.Notes
	(*cd).Table1.LabelX = m["dateCol1"]
	(*cd).Table2.LabelX = m["dateCol2"]
	(*cd).Table1.LabelY = m["valCol1"]
	(*cd).Table2.LabelY = m["valCol2"]

	if c == S {
		(*cd).Table3.Title = ind3.Title
		(*cd).Table3.Desc = ind3.Notes
		(*cd).Table3.LabelX = m["dateCol3"]
		(*cd).Table3.LabelY = m["valCol3"]
	}

	jv, _ := json.Marshal(*cd)

	correlation := Correlation{
		Tbl1:      m["table1"],
		Col1:      m["valCol1"],
		Tbl2:      m["table2"],
		Col2:      m["valCol2"],
		Tbl3:      m["table3"],
		Col3:      m["valCol3"],
		Method:    m["method"],
		Coef:      cf,
		Json:      string(jv),
		Rating:    0,
		Credit:    0,
		Discredit: 0,
	}

	err := DB.Save(&correlation).Error
	check(err)

	return string(jv)
}

// Determine if two sets of dates overlap - X values are referenced so they can be altered and passed back again for Spurious correlation
func GetIntersect(pFromX *time.Time, pToX *time.Time, pRngX *int, fromY time.Time, toY time.Time, rngY int) []FromTo {
	var bucketRange []FromTo
	fromX, toX, rngX := *pFromX, *pToX, *pRngX

	if rngX <= rngY && (fromX == fromY && toX == toY || fromX.After(fromY) && toX.Before(toY)) {
		bucketRange = CreateBuckets(fromX, toX, rngX)
	} else if rngY < rngX && fromY.After(fromX) && toY.Before(toX) {
		bucketRange = CreateBuckets(fromY, toY, rngY)
		*pFromX = fromY
		*pToX = toY
		*pRngX = rngY
	} else if fromX.Before(fromY) && toX.Before(fromY) || fromX.After(toY) && toX.After(toY) {
		return nil
	} else if fromX.Before(fromY) {
		rngYX := dayNum(toX) - dayNum(fromY)
		bucketRange = CreateBuckets(fromY, toX, rngYX)
		*pFromX = fromY
		*pRngX = rngYX
	} else {
		rngXY := dayNum(toY) - dayNum(fromX)
		bucketRange = CreateBuckets(fromX, toY, rngXY)
		*pToX = toY
		*pRngX = rngXY
	}

	return bucketRange
}

// generate some random tables and columns, amount dependent on correlation type
func GetRandomNames(m map[string]string, c cmeth) bool {
	allNames := true

	m["table2"] = RandomTableName()
	guid2 := NameToGuid(m["table2"])
	columnNames2 := FetchTableCols(guid2)
	m["valCol2"] = RandomValueColumn(columnNames2)
	m["dateCol2"] = RandomDateColumn(columnNames2)

	if m["table1"] == "" || m["table2"] == "" || m["valCol1"] == "" || m["valCol2"] == "" || m["dateCol1"] == "" || m["dateCol2"] == "" {
		allNames = false
	}

	if c == S {
		m["table3"] = RandomTableName()
		guid3 := NameToGuid(m["table3"])
		columnNames3 := FetchTableCols(guid3)
		m["valCol3"] = RandomValueColumn(columnNames3)
		m["dateCol3"] = RandomDateColumn(columnNames3)
		if m["table3"] == "" || m["valCol3"] == "" || m["dateCol3"] == "" {
			allNames = false
		}
	}

	return allNames
}

// generate a random value column if one exists
func RandomValueColumn(cols []ColType) string {
	if cols == nil {
		return ""
	}

	rand.Seed(time.Now().UTC().UnixNano())
	columns := make([]string, 0)

	for i, _ := range cols {
		if (cols[i].Sqltype == "numeric" || cols[i].Sqltype == "float" || cols[i].Sqltype == "integer") && cols[i].Name != "transaction_number" {
			columns = append(columns, cols[i].Name)
		}
	}

	n := len(columns)

	if n > 0 {
		x := rand.Intn(n)
		return columns[x]
	} else {
		return ""
	}
}

// generate a random date column if one exists
func RandomDateColumn(cols []ColType) string {
	if cols == nil {
		return ""
	}

	rand.Seed(time.Now().UTC().UnixNano())
	columns := make([]string, 0)

	for _, v := range cols {
		isDate, _ := regexp.MatchString("date", strings.ToLower(v.Name)) //find a column of date type
		if isDate {
			columns = append(columns, v.Name)
		}
	}

	n := len(columns)

	if n > 0 {
		x := rand.Intn(n)
		return columns[x]
	} else {
		return ""
	}
}

// generate a random table name
func RandomTableName() string {
	var name []string
	err := DB.Table("priv_onlinedata").Order("random()").Limit(1).Pluck("tablename", &name).Error
	if err != nil && err != gorm.RecordNotFound {
		return ""
	}
	return name[0]
}

// convert table name to guid
func NameToGuid(tablename string) string {
	var guid []string
	err := DB.Table("priv_onlinedata").Where("tablename = ?", tablename).Pluck("guid", &guid).Error
	if err != nil && err != gorm.RecordNotFound {
		return ""
	}
	return guid[0]
}

/// return date and value columns combined within struct from table
func ExtractDateVal(tablename string, dateCol string, valCol string) []DateVal {
	if tablename == "" || dateCol == "" || valCol == "" {
		return nil
	}

	var dates []time.Time
	var amounts []float64

	d := "DELETE FROM " + tablename + " WHERE " + dateCol + " = '0001-01-01 BC'" ////////TEMP FIX TO GET RID OF INVALID VALUES IN GOV DATA
	DB.Exec(d)

	err = DB.Table(tablename).Pluck(dateCol, &dates).Error
	if err != nil && err != gorm.RecordNotFound {
		check(err)
	}

	err = DB.Table(tablename).Pluck(valCol, &amounts).Error
	if err != nil && err != gorm.RecordNotFound {
		check(err)
	}

	result := make([]DateVal, len(dates))

	for i, v := range dates {
		result[i].Date = v
	}

	for i, v := range amounts {
		result[i].Value = v
	}

	return result
}

// return the starting point and end point of a date range and the length of days in between
func DetermineRange(Dates []DateVal) (time.Time, time.Time, int) {
	var fromDate time.Time
	var toDate time.Time

	if Dates == nil {
		return toDate, fromDate, 0
	}

	dVal, high, low := 0, 0, 100000000

	for _, v := range Dates {
		dVal = dayNum(v.Date)
		if dVal > high {
			high = dVal
			toDate = v.Date
		}
		if dVal < low {
			low = dVal
			fromDate = v.Date
		}
	}

	rng := dayNum(toDate) - dayNum(fromDate)
	return fromDate, toDate, rng
}

// create an array of date ranges
func CreateBuckets(fromDate time.Time, toDate time.Time, rng int) []FromTo {
	if rng == 0 {
		return nil
	}

	lim, max := 10, 0

	if rng >= lim { /// no more than 10 buckets
		max = lim
	} else {
		max = rng
	}

	days, bNum := Steps(rng, max) // get days between dates and number of buckets
	extra := 0
	if days > 1 {
		extra = 1
	}
	result := make([]FromTo, bNum+extra)
	date := fromDate // set starting date

	for i := 0; i < bNum; i++ {
		result[i].From = date                   // current date becomes from date
		result[i].To = date.AddDate(0, 0, days) // step amount to to date
		date = result[i].To
	}
	if extra == 1 {
		result[bNum].From = date
		result[bNum].To = toDate.AddDate(0, 0, 1)
	}
	return result
}

// sum dated values that fall on or between relevant dates and return array of these summed values
func FillBuckets(dateVal []DateVal, bucketRange []FromTo) []float64 {
	if dateVal == nil || bucketRange == nil {
		return nil
	}

	buckets := make([]float64, len(bucketRange))

	for _, v := range dateVal {
		for j, w := range bucketRange {
			if v.Between(w.From, w.To) {
				buckets[j] += float64(v.Value)
				break
			}
		}
	}
	return buckets
}

// return true if date is between 2 other dates (from inclusive, up to exclusive)
func (d DateVal) Between(from time.Time, to time.Time) bool {
	if d.Date == from || (d.Date.After(from) && d.Date.Before(to)) {
		return true
	}
	return false
}

// transform date into day number (since 1900)
func dayNum(d time.Time) int {
	var date time.Time
	var days int

	for i := 1900; i < d.Year(); i++ {
		date = time.Date(i, 12, 31, 0, 0, 0, 0, time.UTC)
		days += date.YearDay()
	}

	days += d.YearDay()
	return days
}

// calculate intervals for date ranges
func Steps(a int, b int) (int, int) {
	stepNum := math.Ceil(float64(a) / float64(b))
	bucketNum := a / int(stepNum)
	return int(stepNum), bucketNum
}

// generates array of labels based on from and to dates
func LabelGen(dv []FromTo) []string {
	result := make([]string, len(dv))

	for i, v := range dv {
		if v.From != v.To {
			result[i] = strconv.Itoa(v.From.Day()) + " " + strings.ToUpper(v.From.Month().String()[0:3]) + " " + strconv.Itoa(v.From.Year()) + " to " + strconv.Itoa(v.To.Day()) + " " + strings.ToUpper(v.To.Month().String()[0:3]) + " " + strconv.Itoa(v.To.Year())
		} else {
			result[i] = strconv.Itoa(v.From.Day()) + " " + strings.ToUpper(v.From.Month().String()[0:3]) + " " + strconv.Itoa(v.From.Year())
		}
	}

	return result
}

// returns number of days in month
func daysInMonth(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// returns number of days in year
func daysInYear(y int) int {
	d1 := time.Date(y, 1, 1, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(y+1, 1, 1, 0, 0, 0, 0, time.UTC)
	return int(d2.Sub(d1) / (24 * time.Hour))
}

// increment user credited total for correlation
func Credit(id int) {
	cor := Correlation{}
	err := DB.Where("id= ?", id).Find(&cor).Error
	check(err)
	err = DB.Model(&cor).Update("credit", cor.Credit+1).Error
	check(err)
}

// increment user discredited total for correlation
func Discredit(id int) {
	cor := Correlation{}
	err := DB.Where("id= ?", id).Find(&cor).Error
	check(err)
	err = DB.Model(&cor).Update("discredit", cor.Discredit+1).Error
	check(err)
}

// given a small fraction of ratings there is a strong (95%) chance that the "real", final positive rating will be this value
// eg: gives expected (not necessarily current as there may have only been a few votes so far) value of positive ratings / total ratings
func Ranking(id int) float64 {
	cor := Correlation{}
	err := DB.Where("id= ?", id).Find(&cor).Error
	check(err)
	pos := float64(cor.Credit)
	tot := float64(cor.Credit + cor.Discredit)

	if tot == 0 {
		return 0
	}

	z := 1.96
	phat := pos / tot
	cor.Rating = (phat + z*z/(2*tot) - z*math.Sqrt((phat*(1-phat)+z*z/(4*tot))/tot)) / (1 + z*z/tot)
	err = DB.Save(&cor).Error
	check(err)
	return cor.Rating
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

func GetRelatedCharts(tableName string, offset int, count int) (RelatedCharts, *appError) {
	columns := FetchTableCols(tableName) //array column names
	guid := NameToGuid(tableName)
	charts := make([]TableData, 0)     ///empty slice for adding all possible charts
	xyNames := XYPermutations(columns) // get all possible valid permuations of columns as X & Y
	index := Index{}

	err := DB.Where("guid= ?", guid).Find(&index).Error
	if err != nil && err != gorm.RecordNotFound {
		return RelatedCharts{nil, 0}, &appError{err, "Database query failed", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return RelatedCharts{nil, 0}, &appError{err, "No Chart found", http.StatusNotFound}
	}

	var xyPie XYVal

	for _, v := range columns { // create single column pie charts
		xyPie.X = v.Name
		single := fmt.Sprintf("SELECT %s AS x, COUNT(%s) AS y FROM %s GROUP BY %s", v.Name, v.Name, guid, v.Name)

		GetChartData("pie", single, xyPie, &charts, index)
	}

	for _, v := range xyNames { /// creare column and line charts, stacked column if plotting varchar against varchar
		double := fmt.Sprintf("SELECT %s AS x, %s AS y FROM  %s", v.X, v.Y, guid)
		if v.Xtype == "varchar" && v.Ytype == "varchar" {
			GetChartData("stacked column", double, v, &charts, index)
		}

		GetChartData("column", double, v, &charts, index)
		GetChartData("line", double, v, &charts, index)
	}

	// for i := range charts { // shuffle charts into random order
	// 	j := rand.Intn(i + 1)
	// 	charts[i], charts[j] = charts[j], charts[i]
	// }

	chartLength := len(charts)
	charts = charts[offset : offset+count] // return marshalled slice

	return RelatedCharts{charts, chartLength}, nil
}

func GetChartData(chartType string, sql string, names XYVal, charts *[]TableData, ind Index) {
	var tmpTD TableData
	var tmpXY XYVal
	tmpTD.ChartType = chartType
	tmpTD.Title = ind.Title
	tmpTD.Desc = ind.Notes
	tmpTD.LabelX = names.X
	if chartType != "pie" {
		tmpTD.LabelY = names.Y
	}
	rows, _ := DB.Raw(sql).Rows()
	defer rows.Close()
	pieSlices := 0

	for rows.Next() {
		pieSlices++
		rows.Scan(&tmpXY.X, &tmpXY.Y)
		tmpTD.Values = append(tmpTD.Values, tmpXY)
	}

	if chartType != "pie" || (chartType == "pie" && pieSlices < 20) { // drop pie charts with too many slices
		*charts = append(*charts, tmpTD)
	}
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

func GetCreditedCorrelatedCharts() {
	//where table 1 is the same, highest user rating
	/// inject chart type for each, intelligent but random - if 3 then use bubble
	/// return JSON with title etc

}
func GetNewCorrelatedCharts() {
	// where table 1 is the same and no correlation value exists
	/// inject chart type for each, intelligent but random - if 3 then use bubble
	/// return JSON with title etc
}
