package util

import (
	"io/ioutil"
	"log"
	"time"
)

func LogToJson(posts []byte) {
	fiName := time.Now().String()
	err := ioutil.WriteFile(fiName, posts, 0755)
	if err != nil {
		log.Fatal(err)
		return
	}
}