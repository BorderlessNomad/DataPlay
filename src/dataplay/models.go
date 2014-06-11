package main

type User struct {
	Uid      int `primaryKey:"yes"`
	Email    string
	Password string
}

func (u User) TableName() string {
	return "priv_users"
}

type Tracking struct {
	Id   int `primaryKey:"yes"`
	User string
	Guid string
}

func (t Tracking) TableName() string {
	return "priv_tracking"
}

type Index struct {
	Guid    int `primaryKey:"yes"`
	Name    string
	Title   string
	Notes   string
	CkanUrl string
	Owner   int
}

func (i Index) TableName() string {
	return "index"
}

type OnlineData struct {
	Guid        int `primaryKey:"yes"`
	Datasetguid string
	Tablename   string
	Defaults    string
}

func (o OnlineData) TableName() string {
	return "priv_onlinedata"
}
