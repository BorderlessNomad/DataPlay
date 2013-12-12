package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

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

	} else {
		fmt.Println("$tool execlog listoffiles")
	}
}
