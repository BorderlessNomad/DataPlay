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

	c := NewClient("d32943c0760c4eb7a25d40ad756877bb") // embedly API key MAYUR GMAIL
	// c := NewClient("73be1030c0ec4be1959485ac148ada91") // embedly API key GLYN GMAIL
	// c := NewClient("2ba4435681034ef6b92f729d527453e3") // embedly API key LEX
	// c := NewClient("8104b696aa0e471e8d58f83e4e4c39b1") // embedly API key GLYN@PLAYGEN
	options := Options{}
	c.Extract(urls, options, pos)
}
