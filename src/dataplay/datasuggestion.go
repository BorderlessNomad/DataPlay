package main

import (
	"encoding/json"
	// "fmt"
	"github.com/jinzhu/gorm"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

type DateAmt struct {
	Date   time.Time
	Amount float64
}

type FromTo struct {
	From time.Time
	To   time.Time
}

/**
 * @brief Gets (or generates if one does not exist) a JSON string containing the details of the correlation between a random numeric column of the
 * passed table and a random numeric column of another randomly selected table from the database
 */
func GetCorrelation(table1 string) string {
	if table1 == "" {
		return ""
	}

	c := Correlation{}
	m := make(map[string]string)
	m["table1"] = table1
	GetNames(m)
	var coef []float64 // check if correlation already exists for this pairing first @TODO: Add more cols, date etc???
	err := DB.Model(&c).Where("tbl1 = ?", m["table1"]).Where("col1 = ?", m["amtCol1"]).Where("tbl2 = ?", m["table2"]).Where("col2 = ?", m["amtCol2"]).Where("method = ?", "Pearson").Pluck("coef", &coef).Error
	check(err)

	if coef == nil {
		cf := GetCoef(m)
		correlation := Correlation{
			Tbl1:   m["table1"],
			Col1:   m["amtCol1"],
			Tbl2:   m["table2"],
			Col2:   m["amtCol2"],
			Tbl3:   "",
			Col3:   "",
			Method: "Pearson",
			Coef:   cf,
		}

		jv, _ := json.Marshal(correlation)
		correlation.Json = string(jv)
		err := DB.Save(&correlation).Error // save newly generated row in correlations table
		check(err)
	}

	var result []string //query again and result now exists!
	err = DB.Model(&c).Where("tbl1 = ?", m["table1"]).Where("col1 = ?", m["amtCol1"]).Where("tbl2 = ?", m["table2"]).Where("col2 = ?", m["amtCol2"]).Where("method = ?", "Pearson").Pluck("json", &result).Error
	check(err)
	return result[0]
}

/**
 * @brief Get Random appropriate table and column names
 */
func GetNames(m map[string]string) {

	m["table2"] = RandomTableName() // get random 2nd table name
	guid1 := NameToGuid(m["table1"])
	guid2 := NameToGuid(m["table2"])
	columnNames1 := FetchTableCols(guid1)           // get all column names in table 1
	columnNames2 := FetchTableCols(guid2)           // get all columns names in table 2
	m["amtCol1"] = RandomAmountColumn(columnNames1) // get name of random numeric column from table 1
	m["amtCol2"] = RandomAmountColumn(columnNames2) // get name of random numeric column from table 2
	m["dateCol1"] = RandomDateColumn(columnNames1)  // get name of random date column from table 1
	m["dateCol2"] = RandomDateColumn(columnNames2)  // get name of random date column from table 2
}

/**
 * @brief Bulk of the algorithm, take in map of column and table names and spit out correlation coefficient based on them
 */
func GetCoef(m map[string]string) float64 {
	if len(m) == 0 {
		return 0.0
	}
	var cf float64

	x := ExtractDateAmt(m["table1"], m["dateCol1"], m["amtCol1"]) // get the chosen random dates and amounts from table 1
	y := ExtractDateAmt(m["table2"], m["dateCol2"], m["amtCol2"]) // get the chosen random dates and amounts from table 2
	fromDateX, toDateX, rngX := DetermineRange(x)                 // get the date range for table 1
	fromDateY, toDateY, rngY := DetermineRange(y)                 // get the date range for table 2

	// choose whichever range is smaller to be the template range as long as it overlaps the other range (if there's no overlap return 0)
	// then create standard range of discrete date periods (buckets) and place dated data from x and y into these buckets

	if rngX == rngY || (rngX < rngY && (fromDateX.After(fromDateY) && fromDateX.Before(fromDateY))) || (fromDateX.After(fromDateY) || fromDateX.Before(fromDateY)) {
		bucketRange := CreateBuckets(x, fromDateX, toDateX, rngX) // create appropriate bucket range
		xBuckets := FillBuckets(x, bucketRange)                   // put table 1 values into buckets
		yBuckets := FillBuckets(y, bucketRange)                   // put table 2 values into buckets
		cf = Pearson(xBuckets, yBuckets)                          // calculate coefficient of table 1 and table 2 values
		return cf
	} else if rngY < rngX && (fromDateX.After(fromDateY) && fromDateX.Before(fromDateY)) || (fromDateX.After(fromDateY) || fromDateX.Before(fromDateY)) {
		bucketRange := CreateBuckets(y, fromDateY, toDateY, rngY)
		xBuckets := FillBuckets(x, bucketRange)
		yBuckets := FillBuckets(y, bucketRange)
		cf = Pearson(xBuckets, yBuckets)
		return cf
	} else {
		return 0
	}
}

/**
 * @brief Takes a bunch of column names and types and returns a random amount column of a numeric type
 */
func RandomAmountColumn(cols []ColType) string {
	if cols == nil {
		return ""
	}

	rand.Seed(time.Now().UTC().UnixNano())
	columns := make([]string, 0)

	for i, _ := range cols {
		if cols[i].Sqltype == "integer" || cols[i].Sqltype == "float" {
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

/**
 * @brief Takes a bunch of column names and types and returns a random amount date column
 * @TODO: Add date type check
 */
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

/**
 * @brief Returns a random table name from the database schema
 */
func RandomTableName() string {
	var name []string
	err := DB.Table("priv_onlinedata").Order("random()").Limit(1).Pluck("tablename", &name).Error
	if err != nil && err != gorm.RecordNotFound {
		return ""
	}
	return name[0]
}

/**
 * @brief Converts table name to GUID
 */
func NameToGuid(tablename string) string {
	var guid []string
	err := DB.Table("priv_onlinedata").Where("tablename = ?", tablename).Pluck("guid", &guid).Error
	if err != nil && err != gorm.RecordNotFound {
		return ""
	}
	return guid[0]
}

/**
 * @brief Extracts date column and amount column from specified table and returns slice of DateAmt structs
 */
func ExtractDateAmt(tablename string, dateCol string, amtCol string) []DateAmt {
	if tablename == "" || dateCol == "" || amtCol == "" {
		return nil
	}

	var dates []time.Time
	var amounts []float64

	DB.Table(tablename).Pluck(dateCol, &dates)
	//DB.Raw("SELECT " + dateCol + " FROM " + tablename).Scan(&dates)

	err = DB.Table(tablename).Pluck(amtCol, &amounts).Error
	if err != nil && err != gorm.RecordNotFound {
		check(err)
	}

	result := make([]DateAmt, len(amounts))

	for i, v := range amounts {
		result[i].Amount = v
	}

	for i, v := range dates {
		result[i].Date = v
	}

	return result
}

/**
 * @brief Returns the date range (from date, to date and the intervening difference between those dates in days) of an array of dates
 */
func DetermineRange(Dates []DateAmt) (time.Time, time.Time, int) {
	lim := 5 // less dates than this gives nothing worth plotting
	var fromDate time.Time
	var toDate time.Time

	if Dates == nil || len(Dates) < lim {
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

/**
 * @brief Creates a series of dated buckets (each bucket represents an individual date or a range of dates)
 */
func CreateBuckets(Dates []DateAmt, fromDate time.Time, toDate time.Time, rng int) []FromTo {
	if Dates == nil || rng == 0 {
		return nil
	}

	lim := 10
	bucketAmt := 0

	if rng >= lim { /// no more than 10 buckets
		bucketAmt = lim
	} else {
		bucketAmt = rng
	}

	result := make([]FromTo, bucketAmt)
	step := rng / bucketAmt // get steps between dates - rounds down so will never go over date range in loop
	date := fromDate        // set starting date
	i := 0

	for ; i < bucketAmt; i++ {
		result[i].From = date
		result[i].To = date.AddDate(0, 0, step)
		date = result[i].To
	}

	result[i-1].To = toDate /// catch any dates that were rounded off
	return result
}

/**
 * @brief Takes array of dates and amount values and drops and sums them into a discrete range of values which is returned
 */
func FillBuckets(dateAmt []DateAmt, bucketRange []FromTo) []float64 {
	if dateAmt == nil || bucketRange == nil {
		return nil
	}

	bucket := make([]float64, len(bucketRange))

	for i, _ := range dateAmt {
		for j, _ := range bucketRange {
			if dateAmt[i].Between(bucketRange[j].From, bucketRange[j].To) {
				bucket[j] = bucket[j] + dateAmt[i].Amount
			}
		}
	}
	return bucket
}

/**
 * @brief Determine if date lies between 2 dates
 */
func (d DateAmt) Between(from time.Time, to time.Time) bool {
	if d.Date == from || d.Date == to || (d.Date.After(from) && d.Date.Before(to)) {
		return true
	}
	return false
}

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

/**
 * @brief Return number of days in month
 */
func daysInMonth(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

/**
 * @brief Return number of days in year
 */
func daysInYear(y int) int {
	d1 := time.Date(y, 1, 1, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(y+1, 1, 1, 0, 0, 0, 0, time.UTC)
	return int(d2.Sub(d1) / (24 * time.Hour))
}
