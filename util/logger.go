package util

import (
	"encoding/json"
	"fmt"
	"log"

	"os"
	"reddit-scrapper/models"
)

func LogToJson(posts []models.Post, fileName string) {
	jsonData, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	fmt.Println("[INFO] Attempting to write scrapped data to", fileName)

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("Unable to save scrapped data. Error: ", err.Error())
		return
	}
	defer file.Close()

	file.Write(jsonData)

	fmt.Println("[SUCCESS] Write to", fileName, "completed")
}
