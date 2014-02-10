package main

import (
	msql "../../databasefuncs"
	"database/sql"
	"fmt"
	"github.com/cheggaaa/pb" // 66139f61bba9938c8f87e64bea6a8a47f40fdc32
)

func main() {
	database := msql.GetDB()
	database.Ping()

	q, e := database.Query("SELECT `TableName` FROM priv_onlinedata")
	if e != nil {
		panic(":(")
	}
	TableScanTargets := make([]string, 0)
	for q.Next() {
		TTS := ""
		q.Scan(&TTS)
		TableScanTargets = append(TableScanTargets, TTS)
	}
}
