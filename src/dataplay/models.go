package main

import (
	"time"
)

type appError struct {
	Error   error
	Message string
	Code    int
}

type Activity struct {
	ActivityName string
	Date         time.Time
	Uid          int
}

func (a Activity) TableName() string {
	return "priv_activities"
}

type Correlation struct {
	Tbl1    string
	Col1    string
	Tbl2    string
	Col2    string
	Tbl3    string
	Col3    string
	Method  string
	Coef    float64
	Json    []byte
	Id      int `primaryKey:"yes"`
	Abscoef float64
}

func (c Correlation) TableName() string {
	return "priv_correlated"
}

type Index struct {
	Guid    string
	Name    string
	Title   string
	Notes   string
	CkanUrl string
	Owner   int
}

func (i Index) TableName() string {
	return "index"
}

type Observation struct {
	Text           string
	PatternId      int
	DiscoveredBy   int
	Coordinates    string
	Rating         float64
	Valid          int
	Invalid        int
	ObservationId  int `primaryKey:"yes"`
	DateDiscovered time.Time
}

func (ob Observation) TableName() string {
	return "priv_observations"
}

type OnlineData struct {
	Guid        string
	Datasetguid string
	Tablename   string
	Defaults    string
}

func (od OnlineData) TableName() string {
	return "priv_onlinedata"
}

type StatsCheck struct {
	Id     int `primaryKey:"yes"`
	Table  string
	X      string
	Y      string
	P1     int
	P2     int
	P3     int
	Xstart int
	Xend   int
}

func (sc StatsCheck) TableName() string {
	return "priv_statcheck"
}

type StringSearch struct {
	Tablename string
	X         string
	Value     string
	Count     int
}

func (ss StringSearch) TableName() string {
	return "priv_stringsearch"
}

type SearchTerm struct {
	Term  string `primaryKey:"yes"`
	Count int
}

func (st SearchTerm) TableName() string {
	return "priv_searchterms"
}

type Tracking struct {
	Id      int `primaryKey:"yes"`
	User    int
	Guid    string
	Info    string
	Created time.Time
}

func (t Tracking) TableName() string {
	return "priv_tracking"
}

type TrackingInfo struct {
	Id   int `primaryKey:"yes"`
	Info []byte
}

func (ti TrackingInfo) TableName() string {
	return "priv_tracking"
}

type TableSchema struct {
	ColumnName string
	DataType   string
}

func (ts TableSchema) TableName() string {
	return "information_schema.columns"
}

type User struct {
	Uid        int `primaryKey:"yes"`
	Email      string
	Password   string
	Reputation int
	// ProfilePic string
}

func (u User) TableName() string {
	return "priv_users"
}

type Validated struct {
	PatternId      int `primaryKey:"yes"`
	DiscoveredBy   int
	DateDiscovered time.Time
	Correlated     bool
	Rating         float64
	Valid          int
	Invalid        int
	DataHash       string
	IdentifierHash string
	Json           []byte
	OriginId       int
}

func (v Validated) TableName() string {
	return "priv_validated"
}

type Validation struct {
	PatternId   int
	ValidatedBy int
}

func (vn Validation) TableName() string {
	return "priv_validations"
}
