package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"os"
)

func GetDB() *sql.DB {
	fmt.Println("[database] Asked for MySQL connection")
	dbhost := "10.0.0.2:3306"
	if os.Getenv("database") != "" {
		dbhost = os.Getenv("database")
	}
	con, err := sql.Open("mysql", "root:@tcp("+dbhost+")/DataCon?allowAllFiles=true")
	con.Exec("SET NAMES UTF8")
	mysql.RegisterLocalFile("./")
	if err != nil {
		fmt.Println("[database] Unable to set up connection!")
	}
	con.Ping()
	return con
}
