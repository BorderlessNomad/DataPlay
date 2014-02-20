package main

import (
	msql "../../databasefuncs"
	"database/sql"
	"fmt"
	"github.com/cheggaaa/pb" // 66139f61bba9938c8f87e64bea6a8a47f40fdc32
	"regexp"
	"strings"
)

type ScanJob struct {
	TableName string
	X         string
}

func main() {
	database := msql.GetDB()
	database.Ping()

	q, e := database.Query("SELECT `TableName` FROM priv_onlinedata")
	if e != nil {
		panic(e)
	}
	TableScanTargets := make([]string, 0)
	for q.Next() {
		TTS := ""
		q.Scan(&TTS)
		TableScanTargets = append(TableScanTargets, TTS)
	}

	fmt.Println("Building Job List...")
	jobs := MakeJobs(database, TableScanTargets)
	fmt.Printf("Preparing to do %d jobs", len(jobs))
	bar := pb.StartNew(len(jobs))
	for _, job := range jobs {
		IndexTable(job, database)
		database.Exec(fmt.Sprintf("ALTER TABLE `%s` COMMENT='!Indexed!';", job.TableName))
		bar.Increment()
	}
}

func IndexTable(job ScanJob, db *sql.DB) {
	q, _ := db.Query(fmt.Sprintf("SELECT `%s` FROM `%s`", job.X, job.TableName))
	checkingdict := make(map[string]int)
	InsertQ, e := db.Prepare("INSERT INTO `priv_stringsearch` (`tablename`, `x`, `value`, `count`) VALUES (?, ?, ?, ?);")
	if e != nil {
		panic(e)
	}
	// Count up all the vars in this col
	for q.Next() {
		var strout string
		q.Scan(&strout)
		checkingdict[strout]++
	}
	// Now spit them into INSERT's on the table
	for k, v := range checkingdict {
		InsertQ.Exec(job.TableName, job.X, k, v)
	}
}

func MakeJobs(database *sql.DB, TableScanTargets []string) (jobs []ScanJob) {
	jobs = make([]ScanJob, 0)
	for _, v := range TableScanTargets {
		var CreateSQL string
		database.QueryRow(fmt.Sprintf("SHOW CREATE TABLE `DataCon`.`%s`", v)).Scan(&v, &CreateSQL)

		Bits := ParseCreateTableSQL(CreateSQL)
		for _, bit := range Bits {
			if bit.Sqltype == "varchar" {
				newJob := ScanJob{
					TableName: v,
					X:         bit.Name,
				}
				jobs = append(jobs, newJob)
			}
		}
	}
	return jobs
}

type IdentifyResponce struct {
	Cols    []ColType
	Request string
}

type ColType struct {
	Name    string
	Sqltype string
}

func FetchTableCols(guid string, database *sql.DB) (output []ColType) {
	if guid == "" {
		return output
	}

	var tablename string
	database.QueryRow("SELECT TableName FROM `priv_onlinedata` WHERE GUID = ? LIMIT 1", guid).Scan(&tablename)
	if tablename == "" {
		return output
	}

	var createcode string
	database.QueryRow("SHOW CREATE TABLE "+tablename).Scan(&tablename, &createcode)
	if createcode == "" {
		return output
	}
	results := ParseCreateTableSQL(createcode)
	return results
}

func BuildREArrayForCreateTable(input string) []string {
	re := ".*?(`.*?`).*?((?:[a-z][a-z]+))" // http://i.imgur.com/dkbyB.jpg
	// This regex looks for things that look like
	// `colname` INT,

	var sqlRE = regexp.MustCompile(re)
	results := sqlRE.FindStringSubmatch(input)
	return results
}

func ParseCreateTableSQL(input string) []ColType {
	returnerr := make([]ColType, 0) // Setup the array that I will be append()ing to.
	SQLLines := strings.Split(input, "\n")
	// The mysql server gives you the SQL create code formatted. So I exploit this by
	// using it to split the system up by \n

	for c, line := range SQLLines {
		if c != 0 && strings.HasPrefix(strings.TrimSpace(line), "`") { // Clipping off the create part since its useless for me.
			results := BuildREArrayForCreateTable(line)
			if len(results) == 3 {
				// We expect there to be 3 matches from the Regex, if not then we probs don't
				// have what we want
				DeQuoted := strings.Replace(results[1], "`", "", -1)
				NewCol := ColType{
					Name:    DeQuoted,
					Sqltype: results[2],
				}
				returnerr = append(returnerr, NewCol)
			}
		}
	}
	return returnerr
}

type SuggestionResponce struct {
	Request string
}

func CheckIfColExists(createcode string, targettable string) bool {

	SQLLines := strings.Split(createcode, "\n")

	for c, line := range SQLLines {
		if c != 0 { // Clipping off the create part since its useless for me.
			results := BuildREArrayForCreateTable(line)
			if len(results) == 3 {
				if results[1] == "`"+targettable+"`" {
					return true
				}
			}
		}
	}
	return false
}
