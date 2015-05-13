package database

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"
	"strconv"
)

type Database struct {
	gorm.DB
	SQL    *sql.DB
	User   string
	Pass   string
	Host   string
	Port   string
	Schema string
	Debug  bool
}

func (self *Database) Setup() {
	flag.StringVar(&self.User, "DBUser", "playgen", "The username to use while connecting to the postgresql DB")
	flag.StringVar(&self.Pass, "DBPasswd", "aDam3ntiUm", "The password to use while connecting to the postgresql DB")

	flag.StringVar(&self.Host, "DBHost", "109.231.124.33", "Where to connect to the postgresql DB")
	flag.StringVar(&self.Port, "DBPort", "9999", "Where to connect to the postgresql DB")

	flag.StringVar(&self.Schema, "DBDatabase", "dataplay", "The database name to use while connecting to the postgresql DB")
	flag.BoolVar(&self.Debug, "DBDebug", false, "Debug DB Queries")
}

func (self *Database) ParseEnvironment() {
	databaseHost := "109.231.124.33"
	// databasePort := "5432"
	databasePort := "9999"

	if os.Getenv("DP_DATABASE_HOST") != "" {
		databaseHost = os.Getenv("DP_DATABASE_HOST")
	}

	if os.Getenv("DP_DATABASE_PORT") != "" {
		databasePort = os.Getenv("DP_DATABASE_PORT")
	}

	self.Host = databaseHost
	self.Port = databasePort

	databaseDebug := false
	if os.Getenv("DP_DATABASE_DEBUG") != "" {
		databaseDebug, _ = strconv.ParseBool(os.Getenv("DP_DATABASE_DEBUG"))
	}

	self.Debug = databaseDebug
}

func (self *Database) Connect() (err error) {
	self.DB, err = gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", self.User, self.Pass, self.Host, self.Port, self.Schema))
	if err != nil {
		fmt.Println("[Database] Error while connecting: '%v'", err)
		return err
	}

	fmt.Println("[Database] Connected!", self.User, "@", self.Host, ":", self.Port, "/", self.Schema)

	maxIdleConns := 2048 // < 0 no idle connections are retained.
	maxOpenConns := 1024 // Unlimited  = < 0

	if os.Getenv("DP_DATABASE_MAXIDLECONNS") != "" {
		maxIdleConns, _ = strconv.Atoi(os.Getenv("DP_DATABASE_MAXIDLECONNS"))
	}

	if os.Getenv("DP_DATABASE_MAXOPENCONNS") != "" {
		maxOpenConns, _ = strconv.Atoi(os.Getenv("DP_DATABASE_MAXOPENCONNS"))
	}

	self.DB.DB().Exec("SET NAMES UTF8")
	self.DB.DB().SetMaxIdleConns(maxIdleConns)
	self.DB.DB().SetMaxOpenConns(maxOpenConns)
	self.DB.DB().Ping()

	/* Debug */
	self.DB.LogMode(self.Debug)

	return
}
