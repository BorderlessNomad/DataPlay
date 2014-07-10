package main

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"math/rand"
	"time"
)

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
	table2 := RandomTableGUID()
	cols1 := FetchTableCols(table1)
	cols2 := FetchTableCols(table2)
	column1 := SelectRandomValidColumn(cols1)
	column2 := SelectRandomValidColumn(cols2)

	var coef []float64
	err := DB.Model(&c).Where("tbl1 = ?", table1).Where("col1 = ?", column1).Where("tbl2 = ?", table2).Where("col2 = ?", column2).Where("method = ?", "Pearson").Pluck("coef", &coef).Error
	check(err)

	if coef == nil {
		x := ExtractData(table1, column1)
		y := ExtractData(table2, column2)
		c := Pearson(x, y)

		correlation := Correlation{
			Tbl1:   table1,
			Col1:   column1,
			Tbl2:   table2,
			Col2:   column2,
			Tbl3:   "Null",
			Col3:   "Null",
			Method: "Pearson",
			Coef:   c,
		}

		jv, _ := json.Marshal(correlation)
		correlation.Json = string(jv)

		err := DB.Save(&correlation).Error
		check(err)
	}

	var result []string
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
		return "Null"
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
