package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 2 {
		listfile := string(os.Args[1])
		fileContents, e := ioutil.ReadFile(listfile)
		if e == nil {
			// split the file into lines.
			lines := strings.Split(string(fileContents), "\n")
			for i, line := range lines {
				fmt.Printf("Downloading dataset %d/%d \n", i, len(lines))
				ImportAllDatasets(line, string(i))
			}
		} else {
			fmt.Println("Failed to load that file.")
		}
	}

}
