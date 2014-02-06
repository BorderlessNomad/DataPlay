package main

import (
	msql "../../databasefuncs"
	"database/sql"
	"fmt"
	"github.com/cheggaaa/pb"
	"github.com/skelterjohn/go.matrix" // daa59528eefd43623a4c8e36373a86f9eef870a2
)

var degree = 2

func GetPolyResults(xGiven []float64, yGiven []float64) []float64 {
	m := len(yGiven)
	if m != len(xGiven) {
		return []float64{0, 0, 0} // Send it back, There is nothing sane here.
	}
	n := degree + 1
	y := matrix.MakeDenseMatrix(yGiven, m, 1)
	x := matrix.Zeros(m, n)
	for i := 0; i < m; i++ {
		ip := float64(1)
		for j := 0; j < n; j++ {
			x.Set(i, j, ip)
			ip *= xGiven[i]
		}
	}

	q, r := x.QR()
	qty, err := q.Transpose().Times(y)
	if err != nil {
		fmt.Println(err)
		return []float64{0, 0, 0}
	}
	c := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		c[i] = qty.Get(i, 0)
		for j := i + 1; j < n; j++ {
			c[i] -= c[j] * r.Get(i, j)
		}
		c[i] /= r.Get(i, i)
	}
	fmt.Println(c)
	return c
}

type ScanJob struct {
	TableName string
	X         string
	Y         string
}

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
	fmt.Println("Building Job List...")
	jobs := make([]ScanJob, 0)
	for _, v := range TableScanTargets {

		var CreateSQL string
		database.QueryRow(fmt.Sprintf("SHOW CREATE TABLE `DataCon`.`%s`", v)).Scan(&v, &CreateSQL)

		Bits := ParseCreateTableSQL(CreateSQL)
		for _, bit := range Bits {
			if bit.Sqltype == "int" || bit.Sqltype == "float" {
				for _, bit2 := range Bits {
					if bit.Sqltype == "int" || bit.Sqltype == "float" {
						newJob := ScanJob{
							TableName: v,
							X:         bit.Name,
							Y:         bit2.Name,
						}
						jobs = append(jobs, newJob)
					}
				}
			}
		}
	}
	fmt.Printf("Preparing to do %d jobs", len(jobs))
	bar := pb.StartNew(len(TableScanTargets))
	for _, job := range jobs {
		DoPoly(job, database)
		bar.Increment()
	}
}

func DoPoly(job ScanJob, db *sql.DB) {
	q, _ := db.Query(fmt.Sprintf("SELECT `%s`,`%s` FROM `%s`", job.X, job.Y, job.TableName))
	for q.Next() {
		var f1, f2 float64
		q.Scan(&f1, &f2)
		fmt.Println(f1, f2)
	}
}
