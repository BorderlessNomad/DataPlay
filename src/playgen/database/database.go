package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

type Database struct {
	DB     *sql.DB
	User   string
	Pass   string
	Host   string
	Port   int
	Schema string
}

func (self *Database) Setup(username string, password string, host string, port int, schema string) {
	self.User = username
	self.Pass = password
	self.Host = host
	self.Port = port
	self.Schema = schema

	// flag.StringVar(&self.User, "DBUser", "playgen", "The username to use while connecting to the postgresql DB")
	// flag.StringVar(&self.Pass, "DBPasswd", "aDam3ntiUm", "The password to use while connecting to the postgresql DB")
	// flag.StringVar(&self.Schema, "DBDatabase", "dataplay", "The database name to use while connecting to the postgresql DB")
	// flag.StringVar(&self.Host, "DBHost", "10.0.0.2:5432", "Where to connect to the postgresql DB")
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

	self.DB, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", self.User, self.Pass, self.Host, self.Port, self.Schema))
	if err != nil {
		// logger.Printf("Unable to set up database connection: %s\n", err)
		return
	}
	self.DB.Exec("SET NAMES UTF8")
	self.DB.Ping()
	return
}
