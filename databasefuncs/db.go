package msql

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"os"
)

// err = db.QueryRow("select name from users where id = ?", 1).Scan(&name)

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
	dbhost := "localhost:3306"
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
