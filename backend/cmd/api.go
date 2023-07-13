package main

import (
	"api/internal/dto"
	"api/internal/web"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	db := readFromFile("resources/trails.json")

	webport := "8080"
	err := http.ListenAndServe(":"+webport, web.SetupRoutes(webport, &db))
	if err != nil {
		log.Fatal("Failed to start web service", err)
	}
}

func readFromFile(filename string) dto.TrailsDB {
	var db dto.TrailsDB
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("failed to open file for reading", err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&db)
	if err != nil {
		log.Fatal("failed to decode JSON data", err)
	}
	if db.GetTotalNumberOfTrails() == 0 {
		log.Fatal("no trails were loaded from json file")
	}

	log.Printf("Loaded %d number of trails from json\n", db.GetTotalNumberOfTrails())
	return db
}
