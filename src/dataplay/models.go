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
	ActivityId    int `gorm:"primary_key:yes"`
	Type          string
	Created       time.Time
	Uid           int
	DiscoveredId  int
	ObservationId int
}

func (a Activity) TableName() string {
	return "priv_activity"
}

type Correlation struct {
	Tbl1          string
	Col1          string
	Tbl2          string
	Col2          string
	Tbl3          string
	Col3          string
	Method        string
	Coef          float64
	Json          []byte
	CorrelationId int `gorm:"primary_key:yes"`
	Abscoef       float64
}

func (c Correlation) TableName() string {
	return "priv_correlation"
}

type Dictionary struct {
	Term      string
	Frequency int
	Id        int `gorm:"primary_key:yes"`
}

func (di Dictionary) TableName() string {
	return "priv_dictionary"
}

type Departments struct {
	Key  string
	Dept string
	Id   int `gorm:"primary_key:yes"`
}

func (dp Departments) TableName() string {
	return "priv_departments"
}

type Discovered struct {
	DiscoveredId  int `gorm:"primary_key:yes"`
	Uid           int
	Created       time.Time
	Rating        float64
	Credited      int
	Discredited   int
	Json          []byte
	CorrelationId int
	RelationId    string
}

func (di Discovered) TableName() string {
	return "priv_discovered"
}

type Events struct {
	Key   string
	Event string
	Id    int `gorm:"primary_key:yes"`
}

func (e Events) TableName() string {
	return "priv_events"
}

type Index struct {
	Guid        string
	Name        string
	Title       string
	Notes       string
	CkanUrl     string
	Owner       int
	PrimaryDate string
}

func (i Index) TableName() string {
	return "index"
}

type Observation struct {
	Comment       string
	DiscoveredId  int
	Uid           int
	Rating        float64
	Credited      int
	Discredited   int
	ObservationId int `gorm:"primary_key:yes"`
	Created       time.Time
	X             string
	Y             string
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
	Id     int `gorm:"primary_key:yes"`
	Table  string
	X      string
	Y      string
	P1     int
	P2     int
	P3     int
	Xstart int
	Xend   int
}

type Regions struct {
	Key    string
	Region string
	Id     int `gorm:"primary_key:yes"`
}

func (r Regions) TableName() string {
	return "priv_regions"
}

type SearchTerm struct {
	Id    int `gorm:"primary_key:yes"`
	Term  string
	Count int
}

func (st SearchTerm) TableName() string {
	return "priv_searchterms"
}

type Social struct {
	Uid           int
	Network       string
	FullName      string
	FirstName     string
	LastName      string
	NetworkUserId string
	Sid           int
}

func (so Social) TableName() string {
	return "priv_social"
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

type Tracking struct {
	Id      int `gorm:"primary_key:yes"`
	User    int
	Guid    string
	Info    string
	Created time.Time
}

func (t Tracking) TableName() string {
	return "priv_tracking"
}

type TrackingInfo struct {
	Id   int `gorm:"primary_key:yes"`
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
	Uid        int `gorm:"primary_key:yes"`
	Email      string
	Password   string
	Reputation int
	Avatar     string
	Username   string
	Usertype   int
}

func (u User) TableName() string {
	return "priv_users"
}

type UserTokens struct {
	Id      int `gorm:"primary_key:yes"`
	Uid     int
	Hash    string
	Used    bool
	Created time.Time
}

func (ut UserTokens) TableName() string {
	return "priv_user_tokens"
}

type Credit struct {
	DiscoveredId  int
	Uid           int
	CreditId      int `gorm:"primary_key:yes"`
	Created       time.Time
	ObservationId int
	Credflag      bool
}

func (vn Credit) TableName() string {
	return "priv_credits"
}
