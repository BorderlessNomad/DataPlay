package main

import (
	"github.com/jinzhu/gorm"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// generate some random tables and columns, amount dependent on correlation type
func GetRandomNameMap(m map[string]string, c cmeth) bool {

	m["guid1"] = NameToGuid(m["table1"])
	m["table2"] = RandomTableName()
	m["guid2"] = NameToGuid(m["table2"])
	cols1 := FetchTableCols(m["guid1"])
	cols2 := FetchTableCols(m["guid2"])
	m["dateCol1"] = RandomDateColumn(cols1)
	m["valCol1"] = RandomValueColumn(cols1)
	m["dateCol2"] = RandomDateColumn(cols2)
	m["valCol2"] = RandomValueColumn(cols2)
	allNames := true

	if m["table1"] == m["table2"] || m["table1"] == "" || m["table2"] == "" || m["valCol1"] == "" || m["valCol2"] == "" || m["dateCol1"] == "" || m["dateCol2"] == "" {
		allNames = false
	}
	if c == P {
		m["method"] = "Pearson"
	}
	if c == V {
		m["method"] = "Visual"
	}
	if c == S {
		m["table3"] = RandomTableName()
		m["guid3"] = NameToGuid(m["table3"])
		cols3 := FetchTableCols(m["guid3"])
		m["dateCol3"] = RandomDateColumn(cols3)
		m["valCol3"] = RandomValueColumn(cols3)
		m["method"] = "Spurious"
		if m["table1"] == m["table3"] || m["table2"] == m["table3"] || m["table3"] == "" || m["valCol3"] == "" || m["dateCol3"] == "" {
			allNames = false
		}
	}
	return allNames
}

// convert table name to guid
func NameToGuid(tablename string) string {
	var guid []string
	err := DB.Table("priv_onlinedata").Where("tablename = ?", tablename).Pluck("guid", &guid).Error
	if err != nil && err != gorm.RecordNotFound {
		return ""
	} else if err == gorm.RecordNotFound || len(guid) < 1 {
		return "No Record Found!"
	}

	return guid[0]
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

/// return date and value columns combined within struct from table
func ExtractDateVal(tablename string, dateCol string, valCol string) []DateVal {
	if tablename == "" || dateCol == "" || valCol == "" {
		return nil
	}

	var dates []time.Time
	var amounts []float64

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
		dVal = DayNum(v.Date)
		if dVal > high {
			high = dVal
			toDate = v.Date
		}
		if dVal < low {
			low = dVal
			fromDate = v.Date
		}
	}

	rng := DayNum(toDate) - DayNum(fromDate)
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
	if d.Date.Equal(from) || d.Date.After(from) && d.Date.Before(to) {
		return true
	}
	return false
}

// calculate intervals for date ranges
func Steps(a int, b int) (int, int) {
	stepNum := math.Ceil(float64(a) / float64(b))
	bucketNum := a / int(stepNum)
	return int(stepNum), bucketNum
}

// convert data values of type float to strings
func FloatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', -1, 64)
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

// transform date into day number (since 1900)
func DayNum(d time.Time) int {
	var date time.Time
	var days int

	for i := 1900; i < d.Year(); i++ {
		date = time.Date(i, 12, 31, 0, 0, 0, 0, time.UTC)
		days += date.YearDay()
	}

	days += d.YearDay()
	return days
}

func GetValues(vals []DateVal, from time.Time, to time.Time) ([]XYVal, bool) {
	values := make([]XYVal, 0)
	var tmpXY XYVal
	hasVals := false

	for _, v := range vals {
		if v.Between(from, to) {
			hasVals = true // at least 1 value within range exists
			tmpXY.X = (v.Date.String()[0:10])
			tmpXY.Y = FloatToString(v.Value)
			values = append(values, tmpXY)
		}
	}

	return values, hasVals
}

/// return true if there are too many zero values
func MostlyEmpty(slice []float64) bool {
	ct := 0

	for _, v := range slice {
		if v == 0.0 {
			ct++
		}
	}

	if ct > 5 {
		return true
	} else {
		return false
	}
}
