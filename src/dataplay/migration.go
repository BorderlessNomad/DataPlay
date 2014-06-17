package main

import (
	"fmt"
	// "reflect"
	"regexp"
	// "strconv"
	"strings"
	"time"
)

func MigrateColumns() {
	schema := []struct {
		TableName  string
		ColumnName string
		DataType   string
	}{}
	err1 := DB.Table("information_schema.columns").Select("table_name, column_name, data_type").Where("table_catalog = ?", "dataplay").Where("table_schema = ?", "public").Where("data_type = ?", "character varying").Where("table_name != ?", "proc").Find(&schema).Error

	check(err1)

	r, _ := regexp.Compile(`(amount|price)+`)
	float, _ := regexp.Compile(`^[0-9]*\.[0-9]+$`)
	integer, _ := regexp.Compile(`^\d+$`)
	alphabet, _ := regexp.Compile(`.*[a-zA-Z]+.*`)
	empty, _ := regexp.Compile(`^\s*$`)

	excludeColumns := map[string]bool{
		"date":                  true,
		"table":                 true,
		"time spent_(hh:mm:ss)": true,
	}

	// fmt.Println("\"Table\",\"Columns\",\"Money\",\"Float\",\"Integer\",\"Date\",\"String\"")
	for _, info := range schema {
		/**
		 * 	Apply TRIM, LOWER, CLEAN
		 *
		 * 	DataTypes
		 * 		Integer -> bigint
		 * 		Money -> numeric(100, 2)
		 * 		Float -> numeric(100, 10)
		 * 		Date (ISO) e.g. 1999-01-08		 *
		 * 		String -> DO NOTHING
		 */

		var isMoney, hasFloat, hasInteger, hasDate, isString bool = false, false, false, false, false

		info.TableName = strings.ToLower(info.TableName)
		info.ColumnName = strings.ToLower(info.ColumnName)

		if excludeColumns[info.ColumnName] { // :P
			continue
		}

		values := make([]string, 0)
		err3 := DB.Table(info.TableName).Pluck(info.ColumnName, &values).Error

		check(err3)

		if r.MatchString(info.ColumnName) {
			isMoney = true
		}

		for _, data := range values {
			if !isString && (alphabet.MatchString(data) || empty.MatchString(data)) {
				isString = true
			}

			if !hasFloat && float.MatchString(data) {
				hasFloat = true
			}

			if !hasInteger && integer.MatchString(data) {
				hasInteger = true
			}

			if !hasDate {
				_, errd := time.Parse("2006-01-02", data)
				if errd == nil {
					hasDate = true
				}
				_, errd = time.Parse("2006/01/02", data)
				if errd == nil {
					hasDate = true
				}
				_, errd = time.Parse("02-01-2006", data)
				if errd == nil {
					hasDate = true
				}
				_, errd = time.Parse("02/01/2006", data)
				if errd == nil {
					hasDate = true
				}
			}

		}

		if isMoney && !hasFloat && !hasInteger && !hasDate && !isString {
			fmt.Println("Money:", info.TableName, info.ColumnName)
			AlterTableToMoney(info.TableName, info.ColumnName)
		} else if !isMoney && hasFloat && !hasInteger && !hasDate && !isString {
			// fmt.Println("Float:", info.TableName, info.ColumnName)
		} else if !isMoney && !hasFloat && hasInteger && !hasDate && !isString {
			// fmt.Println("Integer:", info.TableName, info.ColumnName)
		} else if !isMoney && !hasFloat && !hasInteger && hasDate && !isString {
			// fmt.Println("Date:", info.TableName, info.ColumnName)
		} else if !isMoney && !hasFloat && !hasInteger && !hasDate {
			// DO NOTHING
		} else {
			// fmt.Printf("\"%s\",\"%s\",\"%t\",\"%t\",\"%t\",\"%t\",\"%t\"\n", info.TableName, info.ColumnName, isMoney, hasFloat, hasInteger, hasDate, isString)
		}
	}
}

func AlterTableToMoney(table string, column string) {
	// DB.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE numeric(100, 2) USING ((replace(replace(replace(replace(replace(trim(%s), ',', ''), '?', ''), ' ', ''), '<', ''), '>', ''))::numeric(100, 4));", table, column, column))
	DB.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE numeric(100, 2) USING ((regexp_replace(%s, '[^\\d-]+'))::numeric(100, 4));", table, column, column))
}
