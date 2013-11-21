package msql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func GetSingleNumberValue(rows *sql.Rows) (out int, e error) {
	rows.Next()
	var outputnumber int
	ee := rows.Scan(&outputnumber)
	if e == nil {
		return outputnumber, ee
	}
	return 0, ee
}

func GetSingleStringValue(rows *sql.Rows) (out string, e error) {
	rows.Next()
	var outputstr string
	ee := rows.Scan(&outputstr)
	if e == nil {
		return outputstr, ee
	}
	return "", ee
}

func GetQueryMonstrosity(*sql.Rows) string {
	// This one gives you back a csv file of the results
	return ""
}

func GetDB() *sql.DB {
	fmt.Println("[database] Asked for MySQL connection")
	con, err := sql.Open("mysql", "root:@tcp(10.0.0.2:3306)/DataCon?charset=utf8")
	con.Exec("SET NAMES UTF8")

	if err != nil {
		fmt.Println("[database] Unable to set up connection!")
	}
	con.Ping()
	return con
}
