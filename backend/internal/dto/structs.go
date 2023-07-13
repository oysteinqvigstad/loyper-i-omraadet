package dto

import (
	"fmt"
	"math"
	"sort"
)

type Coordinate struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
}

type RoutesXML struct {
	FootSegments  []SegmentXML `xml:"featureMember>Fotrute"`
	BikeSegments  []SegmentXML `xml:"featureMember>Sykkelrute"`
	SkiSegments   []SegmentXML `xml:"featureMember>Skiløype"`
	OtherSegments []SegmentXML `xml:"featureMember>AnnenRute"`
}

type TrailXML struct {
	Name       string `xml:"rutenavn"`
	ID         string `xml:"rutenummer"`
	Difficulty string `xml:"gradering"`
	Subtype    string `xml:"spesialAnnenrutetype"` // only applicable for 'OtherTrails'
	Coordinate []Coordinate
}

type TrailsXML []TrailXML

type SegmentXML struct {
	PosList     string    `xml:"senterlinje>LineString>posList"`
	FootTrails  TrailsXML `xml:"fotruteInfo>FotruteInfo"`
	SkiTrails   TrailsXML `xml:"skiløypeInfo>SkiløypeInfo"`
	BikeTrails  TrailsXML `xml:"sykkelruteInfo>SykkelruteInfo"`
	OtherTrails TrailsXML `xml:"annenRuteInfo>AnnenRuteInfo"`
	Coordinates []Coordinate
}

type SegmentCollapsedXML struct {
	Coordinates []Coordinate
	Trails      TrailsXML
}

type ParsedTrailsXML struct {
	FootTrails  TrailsXML
	BikeTrails  TrailsXML
	SkiTrails   TrailsXML
	OtherTrails TrailsXML
}

type TrailJSON struct {
	Name                 string     `json:"name"`
	ID                   string     `json:"id"`
	Location             Coordinate `json:"location"`
	PolylineDetailHigh   string     `json:"polyline_detail_high"`
	PolylineDetailMedium string     `json:"polyline_detail_medium"`
	PolyLineDetailLow    string     `json:"polyline_detail_low"`
	Difficulty           string     `json:"difficulty"`
	Subtype              string     `json:"subtype,omitempty"`
	DistanceFromUser     float64    `json:",omitempty"`
}

type TrailsJSON []TrailJSON

type TrailsDB struct {
	FootTrails  TrailsJSON `json:"foot_trails"`
	BikeTrails  TrailsJSON `json:"bike_trails"`
	SkiTrails   TrailsJSON `json:"ski_trails"`
	OtherTrails TrailsJSON `json:"other_trails"`
}

type Activity int

const (
	FOOT Activity = iota
	BIKE
	SKI
	OTHER
)

func (trails TrailsXML) spliceSegments(segments []SegmentCollapsedXML, thresholdDistance float64) {
	for i, trail := range trails {
		coordinates := make([][]Coordinate, 0)
		for _, segment := range segments {
			if segment.Trails.isTrailPresent(trail) && len(segment.Coordinates) > 1 {
				coordinates = append(coordinates, segment.Coordinates)
			}
		}
		if len(coordinates) > 0 {
			trails[i].Coordinate = concatCoords(coordinates, thresholdDistance)
		}
		trails[i].Coordinate = makeNewCoordinatesList(trails[i].Coordinate, thresholdDistance/4.0)
	}
}

func (trails TrailsXML) isTrailPresent(trail TrailXML) bool {
	for _, val := range trails {
		if val.Name == trail.Name && val.ID == trail.ID {
			return true
		}
	}
	return false
}

func (p *Coordinate) DistanceInHaversine(p2 Coordinate) float64 {
	const R = 6371 // Radius of the Earth in kilometers
	lat1 := p.Lat * math.Pi / 180
	lat2 := p2.Lat * math.Pi / 180
	dlat := (p2.Lat - p.Lat) * math.Pi / 180
	dlng := (p2.Lon - p.Lon) * math.Pi / 180

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func (p *Coordinate) approximateWithinDistance(p2 Coordinate, thresholdInKM float64) bool {
	kmPerDegree := 111.2
	heuristicFactor := 1.41
	deltaLatitude := math.Abs(p.Lat - p2.Lat)
	deltaLongitude := math.Abs(p.Lon - p2.Lon)
	manhattanDistance := (deltaLatitude + deltaLongitude) * kmPerDegree
	return manhattanDistance/heuristicFactor < thresholdInKM
}

func (db *TrailsDB) PrintStatistics() {
	fmt.Printf("Foot  trails:\t%4d\n", len(db.FootTrails))
	fmt.Printf("Bike  trails:\t%4d\n", len(db.BikeTrails))
	fmt.Printf("Ski   trails:\t%4d\n", len(db.SkiTrails))
	fmt.Printf("Other trails:\t%4d\n", len(db.OtherTrails))

	subtypes := make(map[string]int)
	for _, trail := range db.OtherTrails {
		subtypes[trail.Subtype] += 1
	}

	for subtypeName, count := range subtypes {
		fmt.Printf("  -%-9s:\t%4d\n", subtypeName, count)
	}

}

func (db *TrailsDB) GetNearbyTrails(activity Activity, coord Coordinate, distanceThreshold float64) TrailsJSON {
	var nearbyTrails TrailsJSON
	allTrails := db.getTrailsByActivityType(activity)
	for _, trail := range allTrails {
		if trail.Location.approximateWithinDistance(coord, distanceThreshold) {
			userDistance := trail.Location.DistanceInHaversine(coord)
			if userDistance < distanceThreshold {
				trail.DistanceFromUser = userDistance
				nearbyTrails = append(nearbyTrails, trail)
			}
		}
	}
	sort.Slice(nearbyTrails, func(i, j int) bool {
		return nearbyTrails[i].DistanceFromUser < nearbyTrails[j].DistanceFromUser
	})
	return nearbyTrails

}

func (db *TrailsDB) getTrailsByActivityType(activity Activity) TrailsJSON {
	var trails TrailsJSON
	switch activity {
	case FOOT:
		trails = db.FootTrails
	case BIKE:
		trails = db.BikeTrails
	case SKI:
		trails = db.SkiTrails
	case OTHER:
		trails = db.OtherTrails
	}
	return trails
}

func (db *TrailsDB) GetTotalNumberOfTrails() int {
	return len(db.FootTrails) + len(db.BikeTrails) + len(db.SkiTrails) + len(db.OtherTrails)
}
