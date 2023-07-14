package main

import (
	"api/internal/dto"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

func main() {
	// minimum distance between points
	highQuality := 0.020
	mediumQuality := 0.100
	lowQuality := 0.250

	jsonTrails := parse("resources/turer.xml", highQuality, mediumQuality, lowQuality)
	filteredTrails := cleanupAndReformatDB(jsonTrails, 10, 300)
	writeToFile(filteredTrails)
}

func parse(filename string, thresh1, thresh2, thresh3 float64) dto.TrailsDB {
	xmlFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()
	var routes dto.RoutesXML
	decoder := xml.NewDecoder(xmlFile)
	err = decoder.Decode(&routes)
	if err != nil {
		log.Fatal(err)
	}

	return dto.TrailsDB{
		FootTrails:  dto.ParseXMLSegmentsToJSON(routes.FootSegments, thresh1, thresh2, thresh3),
		BikeTrails:  dto.ParseXMLSegmentsToJSON(routes.BikeSegments, thresh1, thresh2, thresh3),
		SkiTrails:   dto.ParseXMLSegmentsToJSON(routes.SkiSegments, thresh1, thresh2, thresh3),
		OtherTrails: dto.ParseXMLSegmentsToJSON(routes.OtherSegments, thresh1, thresh2, thresh3),
	}
}

func cleanupAndReformatDB(db dto.TrailsDB, lowerlengthLimit, upperLengthLimit int) dto.TrailsDB {
	return dto.TrailsDB{
		FootTrails:  cleanupAndReformatTrails(db.FootTrails, lowerlengthLimit, upperLengthLimit, "F"),
		BikeTrails:  cleanupAndReformatTrails(db.BikeTrails, lowerlengthLimit, upperLengthLimit, "B"),
		SkiTrails:   cleanupAndReformatTrails(db.SkiTrails, lowerlengthLimit, upperLengthLimit, "S"),
		OtherTrails: cleanupAndReformatTrails(db.OtherTrails, lowerlengthLimit, upperLengthLimit, "O"),
	}

}

func cleanupAndReformatTrails(trails dto.TrailsJSON, lowerlengthLimit, upperLengthLimit int, prefix string) dto.TrailsJSON {
	filteredTrails := make(dto.TrailsJSON, 0, len(trails))
	for i, trail := range trails {
		if len(trail.PolyLineDetailLow) <= upperLengthLimit && len(trail.PolyLineDetailLow) >= lowerlengthLimit {
			filteredTrails = append(filteredTrails, dto.TrailJSON{
				Name:                 trail.Name,
				ID:                   prefix + dto.IntToBase62(i),
				Location:             dto.EstimateTrailLocation(trail.PolylineDetailMedium),
				PolylineDetailHigh:   trail.PolylineDetailHigh,
				PolylineDetailMedium: trail.PolylineDetailMedium,
				PolyLineDetailLow:    trail.PolyLineDetailLow,
				Difficulty:           dto.ReformatDifficulty(trail.Difficulty),
				Subtype:              dto.ReformatSubtype(trail.Subtype),
			})
		}
	}
	return filteredTrails
}

func writeToFile(db dto.TrailsDB) {
	data, err := json.MarshalIndent(db, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("resources/trails.json", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
