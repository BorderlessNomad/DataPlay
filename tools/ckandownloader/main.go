// This is a tool that will download EVERYTHING from data.gov provided you give
// it a list of all the items on the site, It targets .csv .xls, it does not really
// know what it is downloading so needs to rely on what the end webserver thinks it is
// If it does not know what the file is, it will name it .dunno
//
// After its done that it will download them all into files with the MD5 has of the
// URL it got it from. If you need to actually use this, I reccomend you keep the output
// of this program in the case that you need to pieace together what the tool actually
// downloaded, and where it came from on data.gov

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

func main() {
	runtime.GOMAXPROCS(3) // Hint at the run time that we would quite like more than one thread
	if len(os.Args) == 2 {
		listfile := string(os.Args[1])
		fileContents, e := ioutil.ReadFile(listfile)
		if e == nil {
			// split the file into lines.
			lines := strings.Split(string(fileContents), "\n")
			for i, line := range lines {
				fmt.Printf("!!Downloading dataset %d/%d \n", i, len(lines))
				ImportAllDatasets(line)
			}
		} else {
			fmt.Println("Failed to load that file.")
		}
	}

}
