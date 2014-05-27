package main

import "testing"
import "time"

func TestMain(t *testing.T) {
	go func() {
		main()
	}()

	time.Sleep(time.Second * 120)
}
