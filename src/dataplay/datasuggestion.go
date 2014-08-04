package main

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

type DateAmt struct {
	Dates   []string
	Amounts []float64
}

/**
 * @brief Gets (or generates if one does not exist) a JSON string containing the details of the correlation between a random numeric column of the
 * passed table name and a random numeric column of a randomly selected table from the database
 * @details
 *
 * @param string
 * @return
 */
func GetCorrelation(table1 string) string {
	if table1 == "" {
		return ""
	}

	c := Correlation{}
	table2 := RandomTableGUID()               // gives random 2nd table
	cols1 := FetchTableCols(table1)           // gets all columns in table 1
	cols2 := FetchTableCols(table2)           // gets all columns in table 2
	column1 := SelectRandomValidColumn(cols1) //gives name of random valid (ie: numeric) column from table 1
	column2 := SelectRandomValidColumn(cols2) //gives name of random valid (ie: numeric) column from table 2

	var coef []float64 // check if correlation already exists for this pairing first
	err := DB.Model(&c).Where("tbl1 = ?", table1).Where("col1 = ?", column1).Where("tbl2 = ?", table2).Where("col2 = ?", column2).Where("method = ?", "Pearson").Pluck("coef", &coef).Error
	check(err)

	if coef == nil {
		x := DateAmt{
			Dates:   ExtractColumnWithExpression("date", table1, cols1), // extract date column from named table, passing in all columns in table
			Amounts: ExtractDataColumn(table1, column1),                 // extract  column of data from named column in named table
		}

		y := DateAmt{
			Dates:   ExtractColumnWithExpression("date", table2, cols2), // extract date column from named table, passing in all columns in table
			Amounts: ExtractDataColumn(table2, column2),                 // extract  column of data from named column in named table
		}

		DataClean(&x, &y)                   //
		cf := Pearson(x.Amounts, y.Amounts) // calculate coefficient

		correlation := Correlation{
			Tbl1:   table1,
			Col1:   column1,
			Tbl2:   table2,
			Col2:   column2,
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

	var result []string //now query again and result now exists!
	err = DB.Model(&c).Where("tbl1 = ?", table1).Where("col1 = ?", column1).Where("tbl2 = ?", table2).Where("col2 = ?", column2).Where("method = ?", "Pearson").Pluck("json", &result).Error
	check(err)
	return result[0]
}

/**
 * @brief Takes a bunch of column names and types and returns a random column of a numeric type
 * @details
 *
 * @param ColType slice of column names
 * @return column name
 */
func SelectRandomValidColumn(cols []ColType) string {
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
 * @brief Returns a random table name from the database schema
 * @details
 * @return
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
 * @brief Clean the data ready for Pearsons
 * @details Modify two Date/Amount pair columns to be same length based on discrete date values
 *
 * @param DateAmt [description]
 * @return [description]
 */
func DataClean(x *DateAmt, y *DateAmt) {
	for i, _ := range x.Dates {
		FixDate(&x.Dates[i])
	}
	for i, _ := range y.Dates {
		FixDate(&y.Dates[i])
	}
}

func FixDate(d *string) {
	//to fix date
}

/**
 * @brief searches through table columns for one containing expression or part of expression and returns values in that column
 * @details
 *
 * @param string Expression (part of column name) being searched for, e.g. "date"
 * @param string Guid of table
 * @param ColType names of columns in table
 * @return column of date values
 */
func ExtractColumnWithExpression(expression string, guid string, cols []ColType) []string {
	var result []string
	if expression == "" || guid == "" {
		return result
	}

	for _, v := range cols {
		dated, _ := regexp.MatchString(expression, strings.ToLower(v.Name))
		if dated == true {
			DB.Table(guid).Pluck(v.Name, &result)
			return result
		}
	}
	return result
}
