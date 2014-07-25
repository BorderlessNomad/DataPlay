package main

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type cmeth int

const (
	P cmeth = iota
	S
	V
)

type CorrelationData struct {
	Title1  string
	Title2  string
	Title3  string
	Desc1   string
	Desc2   string
	Desc3   string
	LabelX1 string
	LabelX2 string
	LabelX3 string
	LabelY1 string
	LabelY2 string
	LabelY3 string
	Vals1   []XYVal
	Vals2   []XYVal
	Vals3   []XYVal
}

type XYVal struct {
	XVal string
	YVal string
}

type DateVal struct {
	Date  time.Time
	Value float64
}

type FromTo struct {
	From time.Time
	To   time.Time
}

func GetCorrelation(table1 string, valCol1 string, dateCol1 string, c cmeth) string {
	if table1 == "" || valCol1 == "" || dateCol1 == "" {
		return ""
	}

	m := make(map[string]string)
	m["table1"], m["dateCol1"], m["valCol1"] = table1, dateCol1, valCol1
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

			if cf != 0 {
				m["method"] = method
				_ = Validate(m, cf, cd)
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
		values1[i].XVal = v
		values1[i].YVal = strconv.FormatFloat(xBuckets[i], 'f', -1, 64)
		values2[i].XVal = v
		values2[i].YVal = strconv.FormatFloat(yBuckets[i], 'f', -1, 64)

		if c == S {
			values3[i].XVal = v
			values3[i].YVal = strconv.FormatFloat(zBuckets[i], 'f', -1, 64)
		}
	}

	(*cd).Vals1 = values1
	(*cd).Vals2 = values2
	(*cd).Vals3 = values3
	return cf
}

func Validate(m map[string]string, cf float64, cd *CorrelationData) string {
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
	err3 := DB.Model(&ind3).Where("guid= ?", guid3).Find(&ind3).Error
	check(err3)

	(*cd).Title1 = ind1.Title
	(*cd).Title2 = ind2.Title
	(*cd).Title3 = ind3.Title
	(*cd).Desc1 = ind1.Notes
	(*cd).Desc2 = ind2.Notes
	(*cd).Desc3 = ind3.Notes
	(*cd).LabelX1 = m["dateCol1"]
	(*cd).LabelX2 = m["dateCol2"]
	(*cd).LabelX3 = m["dateCol3"]
	(*cd).LabelY1 = m["valCol1"]
	(*cd).LabelY2 = m["valCol2"]
	(*cd).LabelY3 = m["valCol3"]

	jv, _ := json.Marshal(*cd)

	correlation := Correlation{
		Tbl1:   m["table1"],
		Col1:   m["valCol1"],
		Tbl2:   m["table2"],
		Col2:   m["valCol2"],
		Tbl3:   m["table3"],
		Col3:   m["valCol3"],
		Method: m["method"],
		Coef:   cf,
		Json:   string(jv),
	}

	err := DB.Save(&correlation).Error
	check(err)

	return string(jv)
}

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

/// return date column from table
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

// sum dated values that fall on or between relevant dates and return array of new values
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

// return true if date is between 2 other dates
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
