package main

import (
	msql "../../databasefuncs"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type ExecLogDef struct {
	origin        string
	link          string
	contents_hash string
	url_hash      string
}

func main() {
	database := msql.GetDB()
	defer database.Close()
	if len(os.Args) == 3 {
		execlogbytes, e := ioutil.ReadFile(os.Args[1])
		if e != nil {
			fmt.Println("Unable to read exec log :(")
			return
		}
		listoffilesbytes, e := ioutil.ReadFile(os.Args[2])
		if e != nil {
			fmt.Println("Unable to read list :(")
			return
		}
		execlog := strings.Split(string(execlogbytes), "\n")
		filelist := strings.Split(string(listoffilesbytes), "\n")
		fmt.Sprintf("%s", filelist)
		defList := make([]ExecLogDef, 0)
		for _, execline := range execlog {
			if strings.HasPrefix(execline, "http") {
				// An example log line is as follows
				// http://data.gov.uk/dataset/2011_census | http://wales.gov.uk/topics/statistics/headlines/population2013/2011-census-analysis-third-release-data-wales/?lang=en => ./data/df33ddd89e405949b67420215860bae2722c9a4f_d30548f56577eeb84f7923923b603caa2580d168
				spacesplit := strings.Split(execline, " ")
				origin := spacesplit[0]
				dlLink := spacesplit[2]
				filename := spacesplit[4]
				filteredname := strings.Replace(filename, "./data/", "", 1)
				splitname := strings.Split(filteredname, "_")
				contents_hash := splitname[0]
				urlhash := splitname[1]

				appender := ExecLogDef{
					origin:    origin,
					link:      dlLink,
					hashchunk: spacesplit[4],
				}
				defList = append(defList, appender)
			}
		}
		// Okay so now we have parsed the log files :toot:
		// Now we need to get the files out and import them into a database
		for _, file := range filelist {
			realfilename := strings.Split(file, ".")[0]
			if LookForItem(realfilename, defList) != -1 {
				// a
				filebytes := ioutil.ReadFile("./data/" + realfilename + ".csv")
				tablecode := ConstructTable(string(filebytes), false)
				_, e := database.Exec(tablecode)
				if e != nil {
					panic(e)
				}
				tableloader = tableloader + "LOAD DATA LOCAL INFILE '" + "./" + realfilename + "'\n"
				tableloader = tableloader + "INTO TABLE " + realfilename + "\n"
				tableloader = tableloader + "FIELDS TERMINATED BY ','" + "\n"
				tableloader = tableloader + "OPTIONALLY ENCLOSED BY '\"'" + "\n"
				tableloader = tableloader + "LINES TERMINATED BY '\\r\\n'\n"
				tableloader = tableloader + "IGNORE 1 LINES\n"
				tableloader = tableloader + "("
				for i := 0; i < colcount; i++ {
					tableloader = tableloader + "`Column " + fmt.Sprintf("%d", i) + "`,"
				}
				tableloader = tableloader + "`Column " + fmt.Sprintf("%d", colcount) + "`)"

				_, e = database.Exec(tableloader)

			}
		}

	} else {
		fmt.Println("$tool execlog listoffiles")
	}
}

func LookForItem(target string, list []ExecLogDef) int {
	for i, item := range list {
		if item.contents_hash == target {
			return i
		}
	}
	return -1
}

/*
	A Basic table will just make Column 1,  Column 2,  Column 3
	rather than a non basic one will make sensible ones
*/
func ConstructTable(csv string, basictable bool) string {
	rows := strings.Split(string(csv), "\n")
	var tablebuilder string
	if !basictable {
		colcount := len(strings.Split(rows[0], ","))
		tablebuilder = "CREATE TABLE `" + fmt.Sprintf("%x", hash.Sum(nil)) + "` ("
		for i := 0; i < colcount; i++ {
			tablebuilder = tablebuilder + "`Column " + fmt.Sprintf("%d", i) + "` TEXT NULL,"
		}
		tablebuilder = tablebuilder + "`Column " + fmt.Sprintf("%d", colcount) + "` TEXT NULL"
		tablebuilder = tablebuilder + ") COLLATE='latin1_swedish_ci' ENGINE=InnoDB;"
	}
	return tablebuilder
}
