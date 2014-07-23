package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"math"
	"math/rand"
	"regexp"
	// "strconv"
	"strings"
	"time"
)

type DateVal struct {
	Date  time.Time
	Value float64
}

type FromTo struct {
	From time.Time
	To   time.Time
}

func GetVisualCorrelation(table1 string, valCol1 string, dateCol1 string) float64 {
	if table1 == "" || valCol1 == "" || dateCol1 == "" {
		return 0
	}

	result := 0.0
	m := make(map[string]string)
	m["table1"], m["dateCol1"], m["valCol1"] = table1, dateCol1, valCol1
	// c := Correlation{}

	nameChk := GetRandomNames(m, false)
	if nameChk {
		// var coef []float64
		// err := DB.Model(&c).Where("tbl1 = ?", m["table1"]).Where("col1 = ?", m["valCol1"]).Where("tbl2 = ?", m["table2"]).Where("col2 = ?", m["valCol2"]).Where("tbl3 = ?", m["table3"]).Where("col3 = ?", m["valCol3"]).Where("method = ?", "Spurious").Pluck("coef", &coef).Error
		// check(err)
		// if coef == nil {
		x := ExtractDateVal(m["table1"], m["dateCol1"], m["valCol1"])
		y := ExtractDateVal(m["table2"], m["dateCol2"], m["valCol2"])
		fromX, toX, rngX := DetermineRange(x)
		fromY, toY, rngY := DetermineRange(y)
		if rngX == 0 || rngY == 0 {
			return 0
		}

		var bucketRange []FromTo

		if rngX <= rngY && (fromX == fromY && toX == toY || fromX.After(fromY) && toX.Before(toY)) {
			bucketRange = CreateBuckets(fromX, toX, rngX)
			fmt.Println("\nxxxFROM", fromX)
			fmt.Println("\nxxxTO", toX)
			fmt.Println("\nrng", rngX)
		} else if rngY < rngX && fromY.After(fromX) && toY.Before(toX) {
			bucketRange = CreateBuckets(fromY, toY, rngY)
			fmt.Println("\nxxxFROM", fromY)
			fmt.Println("\nxxxTO", toY)
			fmt.Println("\nrng", rngY)
		} else if fromX.Before(fromY) && toX.Before(fromY) || fromX.After(toY) && toX.After(toY) {
			return 0
		} else if fromX.Before(fromY) {
			rngYX := dayNum(toX) - dayNum(fromY)
			bucketRange = CreateBuckets(fromY, toX, rngYX)
			fmt.Println("\nxxxFROM", fromY)
			fmt.Println("\nxxxTO", toX)
			fmt.Println("\nrng", rngYX)
		} else {
			rngXY := dayNum(toY) - dayNum(fromX)
			bucketRange = CreateBuckets(fromX, toY, rngXY)
			fmt.Println("\nxxxFROM", fromX)
			fmt.Println("\nxxxTO", toY)
			fmt.Println("\nrng", rngXY)
		}

		xBuckets := FillBuckets(x, bucketRange)
		yBuckets := FillBuckets(y, bucketRange)
		n := len(xBuckets)
		fmt.Println("\nxxx RANGE", bucketRange)
		fmt.Println("\nzzz X BUCKETS", xBuckets)
		fmt.Println("\nzzz Y BUCKETS", yBuckets)
		fmt.Println("\nzzz N", n)

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
		fmt.Println("xxxHIGH", high)
		fmt.Println("xxxLOW", low)
		highTotal := 0
		lowTotal := 0

		for i := 0; i < n; i++ {
			days := dayNum(bucketRange[i].To) - dayNum(bucketRange[i].From)
			highTotal = highTotal + days*int(high[i])
			lowTotal = lowTotal + days*int(low[i])
		}
		fmt.Println("aaahighafter", highTotal)
		fmt.Println("aaalowafter", lowTotal)
		result = 1 / (float64(highTotal) / float64(lowTotal))
		fmt.Println("xxxresult", result)
	}
	return result
}

/**
 * @brief Gets (or generates if one does not exist) a JSON string containing the details of the correlation between a random numeric column of the
 * passed table and a random numeric column of another randomly selected table from the database
 */
func GetCorrelation(table1 string, valCol1 string, dateCol1 string, spur bool) string {
	if table1 == "" || valCol1 == "" || dateCol1 == "" {
		return ""
	}

	m := make(map[string]string)
	m["table1"], m["dateCol1"], m["valCol1"] = table1, dateCol1, valCol1
	c := Correlation{}
	result := ""

	if spur {
		nameChk := GetRandomNames(m, spur)
		if nameChk {
			var coef []float64 // check if correlation already exists for this pairing first
			err := DB.Model(&c).Where("tbl1 = ?", m["table1"]).Where("col1 = ?", m["valCol1"]).Where("tbl2 = ?", m["table2"]).Where("col2 = ?", m["valCol2"]).Where("tbl3 = ?", m["table3"]).Where("col3 = ?", m["valCol3"]).Where("method = ?", "Spurious").Pluck("coef", &coef).Error
			check(err)
			if coef == nil {
				cf := GetCoefP(m)
				correlation := Correlation{
					Tbl1:   m["table1"],
					Col1:   m["valCol1"],
					Tbl2:   m["table2"],
					Col2:   m["valCol2"],
					Tbl3:   m["table3"],
					Col3:   m["valCol3"],
					Method: "Spurious",
					Coef:   cf,
				}
				jv, _ := json.Marshal(correlation)
				correlation.Json = string(jv)
				err := DB.Save(&correlation).Error // save newly generated row in correlations table
				check(err)
			}
			var r []string //query again and result now exists!
			err = DB.Model(&c).Where("tbl1 = ?", m["table1"]).Where("col1 = ?", m["valCol1"]).Where("tbl2 = ?", m["table2"]).Where("col2 = ?", m["valCol2"]).Where("tbl3 = ?", m["table3"]).Where("col3 = ?", m["valCol3"]).Where("method = ?", "Spurious").Pluck("json", &r).Error
			check(err)
			result = r[0]
		}
	} else {
		nameChk := GetRandomNames(m, spur)
		if nameChk {
			var coef []float64 // check if correlation already exists for this pairing first
			err := DB.Model(&c).Where("tbl1 = ?", m["table1"]).Where("col1 = ?", m["valCol1"]).Where("tbl2 = ?", m["table2"]).Where("col2 = ?", m["valCol2"]).Where("method = ?", "Pearson").Pluck("coef", &coef).Error
			check(err)
			if coef == nil {
				cf := GetCoefP(m)
				correlation := Correlation{
					Tbl1:   m["table1"],
					Col1:   m["valCol1"],
					Tbl2:   m["table2"],
					Col2:   m["valCol2"],
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
			var r []string //query again and result now exists!
			err = DB.Model(&c).Where("tbl1 = ?", m["table1"]).Where("col1 = ?", m["valCol1"]).Where("tbl2 = ?", m["table2"]).Where("col2 = ?", m["valCol2"]).Where("method = ?", "Pearson").Pluck("json", &r).Error
			check(err)
			result = r[0]
		}
	}
	return result
}

/**
 * @brief Get Random appropriate table and column names
 */
func GetRandomNames(m map[string]string, spur bool) bool {
	allNames := true

	m["table2"] = RandomTableName() // get random 2nd table name
	guid2 := NameToGuid(m["table2"])
	columnNames2 := FetchTableCols(guid2)           // get all columns names in table 2
	m["valCol2"] = RandomAmountColumn(columnNames2) // get name of random numeric column from table 2
	m["dateCol2"] = RandomDateColumn(columnNames2)  // get name of random date column from table 2

	if m["table1"] == "" || m["table2"] == "" || m["valCol1"] == "" || m["valCol2"] == "" || m["dateCol1"] == "" || m["dateCol2"] == "" {
		allNames = false
	}

	if spur {
		m["table3"] = RandomTableName() // get random 3rd table name
		guid3 := NameToGuid(m["table3"])
		columnNames3 := FetchTableCols(guid3)           // get all columns names in table 3
		m["valCol3"] = RandomAmountColumn(columnNames3) // get name of random numeric column from table 3
		m["dateCol3"] = RandomDateColumn(columnNames3)  // get name of random date column from table 3
		if m["table3"] == "" || m["valCol3"] == "" || m["dateCol3"] == "" {
			allNames = false
		}
	}

	return allNames
}

/**
 * @brief Bulk of the algorithm, take in map of column and table names and spit out correlation coefficient based on them
 */
func GetCoefP(m map[string]string) float64 {
	if len(m) == 0 {
		return 0.0
	}

	x := ExtractDateVal(m["table1"], m["dateCol1"], m["valCol1"]) // get the chosen random dates and amounts from table 1
	y := ExtractDateVal(m["table2"], m["dateCol2"], m["valCol2"]) // get the chosen random dates and amounts from table 2
	fromX, toX, rngX := DetermineRange(x)                         // get the date range for table 1
	fromY, toY, rngY := DetermineRange(y)                         // get the date range for table 2
	if rngX == 0 || rngY == 0 {
		return 0
	}

	// determine template range
	var bucketRange []FromTo

	if rngX <= rngY && (fromX == fromY && toX == toY || fromX.After(fromY) && toX.Before(toY)) { //// 1. X and Y ranges are equal or X range is within Y range
		bucketRange = CreateBuckets(fromX, toX, rngX)
	} else if rngY < rngX && fromY.After(fromX) && toY.Before(toX) { //////////////////////////////// 2. Y range is within X range
		bucketRange = CreateBuckets(fromY, toY, rngY)
	} else if fromX.Before(fromY) && toX.Before(fromY) || fromX.After(toY) && toX.After(toY) { ////// 3. ranges have no overlap
		return 0 /// pie charts
	} else if fromX.Before(fromY) { ///////////////////////////////////////////////////////////////// 4. ranges overlap between from Y and to X
		rngYX := dayNum(toX) - dayNum(fromY)
		bucketRange = CreateBuckets(fromY, toX, rngYX)
	} else { //////////////////////////////////////////////////////////////////////////////////////// 5. ranges overlap between from X and to Y
		rngXY := dayNum(toY) - dayNum(fromX)
		bucketRange = CreateBuckets(fromX, toY, rngXY)
	}

	var cf float64
	xBuckets := FillBuckets(x, bucketRange) // put table 1 values into buckets
	yBuckets := FillBuckets(y, bucketRange) // put table 2 values into buckets
	cf = Pearson(xBuckets, yBuckets)        // calculate coefficient of table 1 and table 2 values

	if cf == 0 {
		fmt.Println("\n\n\nxxxB - ", cf)
		fmt.Println("\n\n X", x)
		fmt.Println("\n\n Y", y)
		fmt.Println("\n\n DATE RANGE ", bucketRange)
		fmt.Println("\n\n X BUCKETS", xBuckets)
		fmt.Println("\n\n Y BUCKETS", yBuckets)
	}

	if cf != 0 {
		fmt.Println("\n\n\nxxxC - ", cf)
		fmt.Println("\n\n xxx - X", x)
		fmt.Println("\n\n xxx - Y", y)
		fmt.Println("\n\n xxx - DATE RANGE ", bucketRange)
		fmt.Println("\n\n xxx - X BUCKETS", xBuckets)
		fmt.Println("\n\n xxx - Y BUCKETS", yBuckets)
	}

	return cf
}

func GetCoefS(m map[string]string) float64 {
	if len(m) == 0 {
		return 0.0
	}

	x := ExtractDateVal(m["table2"], m["dateCol2"], m["valCol2"]) // get the chosen random dates and amounts from table 2
	y := ExtractDateVal(m["table3"], m["dateCol3"], m["valCol3"]) // get the chosen random dates and amounts from table 3
	z := ExtractDateVal(m["table1"], m["dateCol1"], m["valCol1"]) // get the chosen random dates and amounts from table 1
	fromX, toX, rngX := DetermineRange(x)                         // get the date range for table 2
	fromY, toY, rngY := DetermineRange(y)                         // get the date range for table 3
	fromZ, toZ, rngZ := DetermineRange(z)                         // get the date range for table 1

	if rngX == 0 || rngY == 0 || rngZ == 0 {
		return 0
	}

	fromOrig := fromX
	toOrig := toX
	rngOrig := rngX

	var bucketRange []FromTo

	if rngX <= rngY && (fromX == fromY && toX == toY || fromX.After(fromY) && toX.Before(toY)) { //// 1. X and Y ranges are equal or X range is within Y range
		bucketRange = CreateBuckets(fromX, toX, rngX)
	} else if rngY < rngX && fromY.After(fromX) && toY.Before(toX) { //////////////////////////////// 2. Y range is within X range
		bucketRange = CreateBuckets(fromY, toY, rngY)
		fromOrig = fromY
		toOrig = toY
		rngOrig = rngY
	} else if fromX.Before(fromY) && toX.Before(fromY) || fromX.After(toY) && toX.After(toY) { ////// 3. ranges have no overlap
		return 0 /// pie charts
	} else if fromX.Before(fromY) { ///////////////////////////////////////////////////////////////// 4. ranges overlap between from Y and to X
		rngYX := dayNum(toX) - dayNum(fromY)
		bucketRange = CreateBuckets(fromY, toX, rngYX)
		fromOrig = fromY
		rngOrig = rngYX
	} else { //////////////////////////////////////////////////////////////////////////////////////// 5. ranges overlap between from X and to Y
		rngXY := dayNum(toY) - dayNum(fromX)
		bucketRange = CreateBuckets(fromX, toY, rngXY)
		toOrig = toY
		rngOrig = rngXY
	}

	if rngOrig <= rngZ && (fromOrig == fromZ && toOrig == toZ || fromOrig.After(fromZ) && toOrig.Before(toZ)) {
		bucketRange = CreateBuckets(fromOrig, toOrig, rngOrig)
	} else if rngZ < rngOrig && fromZ.After(fromOrig) && toZ.Before(toOrig) {
		bucketRange = CreateBuckets(fromZ, toZ, rngZ)
	} else if fromOrig.Before(fromZ) && toOrig.Before(fromZ) || fromOrig.After(toZ) && toOrig.After(toZ) {
		return 0
	} else if fromOrig.Before(fromZ) {
		rngZOrig := dayNum(toOrig) - dayNum(fromZ)
		bucketRange = CreateBuckets(fromZ, toOrig, rngZOrig)
	} else {
		rngOrigZ := dayNum(toZ) - dayNum(fromOrig)
		bucketRange = CreateBuckets(fromOrig, toZ, rngOrigZ)
	}

	var cf float64
	xBuckets := FillBuckets(x, bucketRange)     // put table 2 values into buckets
	yBuckets := FillBuckets(y, bucketRange)     // put table 3 values into buckets
	zBuckets := FillBuckets(z, bucketRange)     // put table 1 values into buckets
	cf = Spurious(xBuckets, yBuckets, zBuckets) // calculate coefficient of table 1 and table 2 values

	if cf == 0 {
		fmt.Println("\n\n\nxxxB - ", cf)
		fmt.Println("\n\n X", x)
		fmt.Println("\n\n Y", y)
		fmt.Println("\n\n Y", z)
		fmt.Println("\n\n DATE RANGE ", bucketRange)
		fmt.Println("\n\n X BUCKETS", xBuckets)
		fmt.Println("\n\n Y BUCKETS", yBuckets)
		fmt.Println("\n\n Y BUCKETS", zBuckets)
	}

	if cf != 0 {
		fmt.Println("\n\n\nxxxC - ", cf)
		fmt.Println("\n\n X", x)
		fmt.Println("\n\n Y", y)
		fmt.Println("\n\n Y", z)
		fmt.Println("\n\n DATE RANGE ", bucketRange)
		fmt.Println("\n\n X BUCKETS", xBuckets)
		fmt.Println("\n\n Y BUCKETS", yBuckets)
		fmt.Println("\n\n Y BUCKETS", zBuckets)
	}

	return cf
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
 * @brief Extracts date column and amount column from specified table and returns slice of DateVal structs
 */
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

/**
 * @brief Returns the date range (from date, to date and the intervening difference between those dates in days) of an array of dates
 */
func DetermineRange(Dates []DateVal) (time.Time, time.Time, int) {
	lim := 1 // less dates than this gives nothing worth plotting
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
func CreateBuckets(fromDate time.Time, toDate time.Time, rng int) []FromTo {
	if rng == 0 {
		return nil
	}

	lim := 10
	max := 0

	if rng >= lim { /// no more than 10 buckets
		max = lim
	} else {
		max = rng
	}

	days, bNum := Steps(rng, max) // get days between dates and number of buckets
	result := make([]FromTo, bNum+1)
	date := fromDate // set starting date

	for i := 0; i < bNum; i++ {
		result[i].From = date                   // current date becomes from date
		result[i].To = date.AddDate(0, 0, days) // step amount to to date
		date = result[i].To
	}
	result[bNum].From = date
	result[bNum].To = toDate.AddDate(0, 0, 1)
	return result
}

/**
 * @brief Takes array of dates and amount values and specified discrete date ranges sums the values according to where the dates place them in that range
 */
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

/**
 * @brief Determine if date lies between 2 dates (from a inclusive up to b non inclusive)
 */
func (d DateVal) Between(from time.Time, to time.Time) bool {
	if d.Date == from || (d.Date.After(from) && d.Date.Before(to)) {
		return true
	}
	return false
}

/**
 * @brief Return date as number of days since 1900
 */
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

/**
 * @brief Return number of steps, number of buckets, remainder steps
 */
func Steps(a int, b int) (int, int) {
	stepNum := math.Ceil(float64(a) / float64(b))
	bucketNum := a / int(stepNum)
	return int(stepNum), bucketNum
}
