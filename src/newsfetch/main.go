package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func main() {

	e := DailyKimono() // generate the day's news urls

	if e == nil {
		fmt.Println("STARTING...")
		file, _ := os.Open("dailyurls.txt")
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

		c := NewClient(EmKey5)
		options := Options{}

		pos := 0

		for {
			e, p := c.Extract(urls, options, pos)
			pos = p
			if e == nil {
				file.Close()
				break
			}

			fmt.Println("RE-STARTING...", e.Error())
		}

	} else {
		fmt.Println(e.Error())
	}
}
