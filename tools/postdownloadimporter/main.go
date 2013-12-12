package main

import (
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
					origin:        origin,
					link:          dlLink,
					contents_hash: contents_hash,
					url_hash:      urlhash,
				}
				defList = append(defList, appender)
			}
		}

	} else {
		fmt.Println("$tool execlog listoffiles")
	}
}
