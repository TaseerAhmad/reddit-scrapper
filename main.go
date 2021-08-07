package main

import (
	"fmt"
	"os"
	scrapper "reddit-scrapper/third_party/colly"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Not enough arguments passed")
		return
	}
	url := os.Args[3]
	var keepAlive bool
	var refreshRate int

	if os.Args[1] == "true" {
		keepAlive = true
	}

	rate, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Error inferring the value of refreshRate. Must be an integer. Passed value: ", os.Args[2])
		return
	}
	refreshRate = rate

	scrapper.Init("old.reddit.com")
	scrapper.Start(keepAlive, refreshRate, url)
}