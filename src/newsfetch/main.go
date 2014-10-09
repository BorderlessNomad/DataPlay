package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func main() {
	Start(0)
}

func Start(pos int) {
	fmt.Println("STARTING...")
	file, _ := os.Open("urls5000.csv") // url file
	defer file.Close()
	reader := csv.NewReader(file)
	urls := make([]string, 0)

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}

		urls = append(urls, record[0])
	}

	c := NewClient(EmKey1)
	options := Options{}
	c.Extract(urls, options, pos)
}
