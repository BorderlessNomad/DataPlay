package database

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type Database struct {
	DB     *sql.DB
	User   string
	Pass   string
	Host   string
	Schema string
}

func (self *Database) SetupFlags() {
	flag.StringVar(&self.User, "DBUser", "root", "The username to use while connecting to the mysql DB")
	flag.StringVar(&self.Pass, "DBPasswd", "", "The password to use while connecting to the mysql DB")
	flag.StringVar(&self.Schema, "DBDatabase", "DataCon", "The database name to use while connecting to the mysql DB")
	flag.StringVar(&self.Host, "DBHost", "10.0.0.2:3306", "Where to connect to the mysql DB")
}

func (self *Database) ParseEnvironment() {
	env := os.Getenv("DATABASE")
	if env == "" {
		// backwards compat
		env = os.Getenv("database")
	}
	if env != "" {
		self.Host = env
	}
}

func (self *Database) Connect() (err error) {
	// debuglogger.Println("GetDB was called")
	if self.Host == "" {
		return fmt.Errorf("No Database host inputted")
	}

	self.DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", self.User, self.Pass, self.Host, self.Schema))
	if err != nil {
		// logger.Printf("Unable to set up database connection: %s\n", err)
		return
	}
	self.DB.Exec("SET NAMES UTF8")
	self.DB.Ping()
	return
}
