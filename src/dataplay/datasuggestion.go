package main

import (
	// "encoding/json"
	//"fmt"
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
	From   time.Time
	To time.Time
}

/**
 * @brief Gets (or generates if one does not exist) a JSON string containing the details of the correlation between a random numeric column of the
 * passed table name and a random numeric column of a randomly selected table from the database
 */
func GetCorrelation(table1 string) string {
	if table1 == "" {
		return ""
	}
	c := Correlation{}

	////ALGORITHM////////

	table2 := RandomTableGUID()                 // gives random 2nd table name
	columnNames1 := FetchTableCols(table1)      // gets all column names in table 1
	columnNames2 := FetchTableCols(table2)      // gets all columns names in table 2
	amtCol1 := RandomAmountColumn(columnNames1) // gives name of random numeric column from table 1
	amtCol2 := RandomAmountColumn(columnNames2) // gives name of random numeric column from table 2
	dateCol1 := RandomDateColumn(columnNames1)  // gives name of random date column from table 1
	dateCol2 := RandomDateColumn(columnNames2)  // gives name of random date column from table 2

	var coef []float64 // check if correlation already exists for this pairing first @TODO: Add date col!!!
	err := DB.Model(&c).Where("tbl1 = ?", table1).Where("col1 = ?", amtCol1).Where("tbl2 = ?", table2).Where("col2 = ?", amtCol2).Where("method = ?", "Pearson").Pluck("coef", &coef).Error
	check(err)

	if coef == nil {
		var cf float64
		x := ExtractDateAmt(table1, dateCol1, amtCol1) //get the chosen random dates and amounts from table 1
		y := ExtractDateAmt(table2, dateCol2, amtCol2) //get the chosen random dates and amounts from table 2
		fromDateX, toDateX, rngX := DetermineRange(x)  // get the date range for table 1
		fromDateY, toDateY, rngY := DetermineRange(y)  // get the date range for table 2

		//choose whichever range is smaller to be the template range (as long as it has overlap with the other range), if there's no overlap return 0
		if rngX == rngY || (rngX < rngY && (fromDateX.After(fromDateY) && fromDateX.Before(fromDateY))) || (fromDateX.After(fromDateY) || fromDateX.Before(fromDateY)) {
			///use range X as template
			cf = Pearson(x.Amount, y.Amount) // calculate coefficient
		} else if rngY < rngX && (fromDateX.After(fromDateY) && fromDateX.Before(fromDateY)) || (fromDateX.After(fromDateY) || fromDateX.Before(fromDateY)) {
			/// use range Y as template
			cf = Pearson(x.????, y.???) // calculate coefficient
		} else {
			cf = 0
		}

		correlation := Correlation{ //data to be saved to correlation table
			Tbl1:   table1,
			Col1:   amtCol1,
			Tbl2:   table2,
			Col2:   amtCol2,
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
	err = DB.Model(&c).Where("tbl1 = ?", table1).Where("col1 = ?", amtCol1).Where("tbl2 = ?", table2).Where("col2 = ?", amtCol2).Where("method = ?", "Pearson").Pluck("json", &result).Error
	check(err)
	return result[0]
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
func RandomTableGUID() string {
	var guid []string
	err := DB.Table("index").Order("random()").Limit(1).Pluck("guid", &guid).Error
	if err != nil && err != gorm.RecordNotFound {
		panic(err)
	}
	return guid[0]
}

/**
 * @brief Extracts date column and amount column from specified table and returns slice of DateAmt structs
 */
func ExtractDateAmt(guid string, dateCol string, amtCol string) []DateAmt {
	var result []DateAmt
	var dates []time.Time
	var amounts []float64

	if guid == "" || dateCol == "" || amtCol == "" {
		return result
	}

	DB.Table(guid).Pluck(dateCol, &dates)
	DB.Table(guid).Pluck(amtCol, &amounts)
	result = make([]DateAmt, len(amounts))

	for i, v := range amounts {
		result[i].Amount = v
	}
	for j, w := range dates {
		result[j].Date = w
	}
	return result
}

/**
 * @brief Returns the date range
 */
func DetermineRange(Dates []DateAmt) (time.Time, time.Time, int) {
	lim := 6
	var fromDate time.Time
	var toDate time.Time

	if len(Dates) <= lim {
		return fromDate, toDate, 0
	}

	dVal := 0
	high, low := 0, 1000000

	for _, v := range Dates {
		dVal = v.Date.Year()*365 + int(v.Date.Month())*30 + v.Date.Day()
		if dVal > high {
			high = dVal
			toDate = v.Date
		}
		if dVal < low {
			low = dVal
			fromDate = v.Date
		}
	}
	rng := (toDate.Year()*365 + int(toDate.Month())*30 + toDate.Day()) - (fromDate.Year()*365 + int(fromDate.Month())*30 + fromDate.Day())
	return fromDate, toDate, rng
}

func CreateBucket(Dates []DateAmt, rng int) []FromTo {
	buckets := 0
	if rng >= 10 {
		buckets = 10
	} else {
		buckets = rng
	}

	step := rng / buckets

	for i := 0; i < buckets; i++ {

	}


}
