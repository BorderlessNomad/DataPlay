package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println("--- Start ---")

	fmt.Println("Open embedly.key")
	// key, kErr := ioutil.ReadFile("embedly.txt")
	// check(kErr)

	// newsClient := NewClient(string(key))
	newsClient := NewClient("d32943c0760c4eb7a25d40ad756877bb")

	fmt.Println("Generate the day's news urls")
	e := DailyKimono() // generate the day's news urls

	if e == nil {
		fmt.Println("Open dailyurls.txt")

		file, fErr := os.Open("dailyurls.txt")
		check(fErr)
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

		options := Options{}

		pos := 0

		for {
			e, p := newsClient.Extract(urls, options, pos)
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
