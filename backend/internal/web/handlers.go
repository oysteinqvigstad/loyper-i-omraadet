package web

import (
	"api/internal/dto"
	"encoding/json"
	"net/http"
)

func getNearByTrails(db *dto.TrailsDB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Unsupported HTTP method", http.StatusMethodNotAllowed)
			return
		}

		coords, err := getQueryCoordinate(r.URL, "coordinates", "c")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		typeOfTrail, err := interpretTrailType(getQueryStr(r.URL, "type", "t"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//quality, _ := getQueryStr(r.URL, "quality", "q")
		limitResultsInterval, _ := getQueryIntInterval(r.URL, "limit", "l")
		radius, err := getQueryFloat(r.URL, "radius", "r")
		if err != nil {
			radius = 8.0
		}

		results := db.GetNearbyTrails(typeOfTrail, coords, radius)
		results = limitResults(results, limitResultsInterval)

		// replace
		httpRespondJSON(w, results)

	})
}

func limitResults(trails dto.TrailsJSON, interval [2]int) dto.TrailsJSON {
	start, end := interval[0], interval[1]
	length := len(trails)
	if start < 1 {
		start = 1
	} else if start-1 > length {
		start = length
	}
	if end < start {
		end = start
	} else if end > length {
		end = length
	}
	return trails[start-1 : end]
}

func httpRespondJSON(w http.ResponseWriter, data any) {
	w.Header().Set("content-type", "application/json")
	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)
	if err != nil {
		http.Error(w, "could not encode JSON", http.StatusInternalServerError)
	}
}
