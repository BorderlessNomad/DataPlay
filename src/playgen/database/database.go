package database

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"
)

type Database struct {
	gorm.DB
	SQL    *sql.DB
	User   string
	Pass   string
	Host   string
	Port   string
	Schema string
}

func (self *Database) Setup() {
	flag.StringVar(&self.User, "DBUser", "playgen", "The username to use while connecting to the postgresql DB")
	flag.StringVar(&self.Pass, "DBPasswd", "aDam3ntiUm", "The password to use while connecting to the postgresql DB")

	flag.StringVar(&self.Host, "DBHost", "10.0.0.2", "Where to connect to the postgresql DB")
	flag.StringVar(&self.Port, "DBPort", "5432", "Where to connect to the postgresql DB")

	flag.StringVar(&self.Schema, "DBDatabase", "dataplay", "The database name to use while connecting to the postgresql DB")
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
	self.DB, err = gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", self.User, self.Pass, self.Host, self.Port, self.Schema))
	if err != nil {
		panic(fmt.Sprintf("Error while connecting to Database: '%v'", err))
		return err
	}

	fmt.Println("[Database] Connected!", self.User, self.Pass, self.Host, self.Port, self.Schema)

	self.DB.DB().Exec("SET NAMES UTF8")
	self.DB.DB().SetMaxIdleConns(10)
	self.DB.DB().SetMaxOpenConns(100)
	self.DB.DB().Ping()

	/* Debug */
	self.DB.LogMode(true)
	// DB.SetLogger(gorm.Logger{revel.TRACE})

	return
}
