package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pmylund/sortutil"
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// generate a random value column if one exists
func RandomValueColumn(cols []ColType) string {
	if cols == nil {
		return ""
	}

	rand.Seed(time.Now().UTC().UnixNano())
	columns := make([]string, 0)

	for i, _ := range cols {
		if (cols[i].Sqltype == "numeric" || cols[i].Sqltype == "float" || cols[i].Sqltype == "integer" || cols[i].Sqltype == "real") && cols[i].Name != "transaction_number" {
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
		isDateYear, _ := regexp.MatchString("date|year", strings.ToLower(v.Name)) //find a column of date type

		if isDateYear {
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

/// return date and value columns combined within struct from table
func ExtractDateVal(tablename string, dateCol string, valCol string) ([]DateVal, *appError) {
	if tablename == "" || dateCol == "" || valCol == "" {
		return nil, &appError{nil, "Invalid or empty data.", http.StatusBadRequest}
	}

	var dates []time.Time
	var amounts []float64

	err := DB.Table(fmt.Sprintf("%q", tablename)).Pluck(fmt.Sprintf("%q", dateCol), &dates).Error
	if err != nil && err != gorm.RecordNotFound {
		return nil, &appError{err, "Database query failed (DateCol).", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return nil, &appError{err, "No table for for DateCol.", http.StatusNotFound}
	}

	err = DB.Table(fmt.Sprintf("%q", tablename)).Pluck(fmt.Sprintf("%q", valCol), &amounts).Error
	if err != nil && err != gorm.RecordNotFound {
		return nil, &appError{err, "Database query failed (ValCol).", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return nil, &appError{err, "No table for for ValCol.", http.StatusNotFound}
	}

	result := make([]DateVal, len(dates))

	for i, v := range dates {
		result[i].Date = v
	}

	for i, v := range amounts {
		result[i].Value = v
	}

	return result, nil
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

// transforms date into day number (since 1900)
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

// returns a set of XY values that fall between 2 dates, adding dummy values as bookends if necessary
func GetValues(vals []DateVal, from time.Time, to time.Time) ([]XYVal, bool) {
	values := make([]XYVal, 0)
	var tmpXY XYVal
	hasVals := false
	fromchk, tochk := false, false // check whether to and from dates are included, if not add dummies at the end

	for _, v := range vals {
		if v.Between(from, to) {
			hasVals = true // at least 1 value within range exists
			if v.Date.Equal(from) {
				fromchk = true
			}

			if v.Date.Equal(to.AddDate(0, 0, -1)) {
				tochk = true
			}

			tmpXY.X = (v.Date.String()[0:10])
			tmpXY.Y = FloatToString(v.Value)
			values = append(values, tmpXY)
		}
	}

	if !fromchk {
		tmpXY.X = (from.String()[0:10])
		tmpXY.Y = "0" // @NOTE - Can also be neighbour value = values[0].Y
		values = append(values, tmpXY)
	}

	if !tochk {
		tmpXY.X = (to.AddDate(0, 0, -1).String()[0:10])
		tmpXY.Y = "0" // @NOTE - Can also be neighbour value = values[len(values)-1].Y
		values = append(values, tmpXY)
	}

	sortutil.AscByField(values, "X")

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

// Get guid from tablename
func GetGuid(tablename string) (string, error) {
	if tablename == "" {
		return "", fmt.Errorf("Invalid tablename")
	}

	data, err := GetOnlineDataByTablename(tablename)
	if err != nil && err != gorm.RecordNotFound {
		return "", fmt.Errorf("Internal Server Error.")
	} else if err == gorm.RecordNotFound {
		return "", fmt.Errorf("Could not find Table.")
	}

	return data.Guid, err
}

// Used to determine what type of correlation is performed (increasing number gives stronger weighting to Spurious as it's less likely to find results)
func RandomMethod() string {
	rand.Seed(time.Now().UTC().UnixNano())
	x := rand.Intn(7)
	if x == 0 {
		return "Pearson"
	} else if x == 1 {
		return "Visual"
	} else {
		return "Spurious"
	}
}
