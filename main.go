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
	pages := os.Args[1]
	pagesToScrap, err := strconv.Atoi(pages)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	scrapper.Init()
	scrapper.Start(pagesToScrap, os.Args[2], os.Args[3])
}
