package msql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// func GetSingleValue() string {

// }

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
