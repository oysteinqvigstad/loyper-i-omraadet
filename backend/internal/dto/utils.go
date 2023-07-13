package dto

import (
	"github.com/twpayne/go-polyline"
	"log"
	"math"
	"strconv"
	"strings"
)

func ParseXMLSegmentsToJSON(segments []SegmentXML, thresh1, thresh2, thresh3 float64) TrailsJSON {
	TrailDetailHigh := ConnectXMLSegments(CollapseAllXMLSegments(segments, thresh1), thresh1*4.0)
	TrailDetailMedium := ConnectXMLSegments(CollapseAllXMLSegments(segments, thresh2), thresh2*4.0)
	TrailDetailLow := ConnectXMLSegments(CollapseAllXMLSegments(segments, thresh3), thresh3*4.0)
	length := len(TrailDetailHigh)

	if length != len(TrailDetailMedium) || length != len(TrailDetailLow) {
		log.Fatal("Did not parse the same amount of trails for different qualities")
	}

	db := make(TrailsJSON, 0, length)
	for i := 0; i < length; i++ {
		db = append(db, TrailJSON{
			Name:                 TrailDetailHigh[i].Name,
			ID:                   TrailDetailHigh[i].ID,
			PolylineDetailHigh:   createPolylineString(TrailDetailHigh[i].Coordinate),
			PolylineDetailMedium: createPolylineString(TrailDetailMedium[i].Coordinate),
			PolyLineDetailLow:    createPolylineString(TrailDetailLow[i].Coordinate),
			Difficulty:           TrailDetailHigh[i].Difficulty,
			Subtype:              TrailDetailHigh[i].Subtype,
		})
	}

	return db
}

var counter = 0

func concatCoords(segments [][]Coordinate, thresholdDistance float64) []Coordinate {
	for len(segments) > 1 {
		minDistance := math.MaxFloat64
		minIndex := -1
		orientation := -1
		firstSegmentFirstPoint := segments[0][0]
		firstSegmentLastPoint := segments[0][len(segments[0])-1]

		for i := 1; i < len(segments); i++ {
			secondSegmentFirstPoint := segments[i][0]
			secondSegmentLastPoint := segments[i][len(segments[i])-1]

			comparePoints := [][2]Coordinate{
				{firstSegmentFirstPoint, secondSegmentFirstPoint},
				{firstSegmentFirstPoint, secondSegmentLastPoint},
				{firstSegmentLastPoint, secondSegmentFirstPoint},
				{firstSegmentLastPoint, secondSegmentLastPoint},
			}

			for j, points := range comparePoints {
				distance := points[0].DistanceInHaversine(points[1])
				if distance < minDistance {
					minDistance = distance
					minIndex = i
					orientation = j

				}
			}
		}

		if orientation == 1 || orientation == 3 {
			reverse(segments[minIndex])
		}
		if orientation == 0 || orientation == 1 {
			reverse(segments[0])
		}

		if minDistance > thresholdDistance*2 {
			counter += 1
			println(counter)
			segments[0] = append(segments[0], backtrack(segments[0], segments[minIndex][0])...)
		}

		segments[0] = append(segments[0], segments[minIndex]...)
		segments = append(segments[:minIndex], segments[minIndex+1:]...)
	}

	return segments[0]
}

func backtrack(initialSegment []Coordinate, destination Coordinate) []Coordinate {
	segment := copyReverse(initialSegment)
	minIndex := -1
	minDistance := math.MaxFloat64
	for i := 1; i < len(segment); i++ {
		distance := segment[i].DistanceInHaversine(destination)
		if distance < minDistance {
			minDistance = distance
			minIndex = i
		}
	}
	return segment[:minIndex+1]
}

func copyReverse[T any](original []T) (reversed []T) {
	reversed = make([]T, len(original))
	copy(reversed, original)
	reverse(reversed)
	return
}

func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func CollapseAllXMLSegments(segmentsXML []SegmentXML, minimumDistanceBetweenCoords float64) []SegmentCollapsedXML {
	s := make([]SegmentCollapsedXML, 0, len(segmentsXML))
	for _, segment := range segmentsXML {
		s = append(s, CollapseXMLSegment(segment, minimumDistanceBetweenCoords))
	}
	return s
}

func CollapseXMLSegment(segmentXML SegmentXML, minimumDistanceBetweenCoords float64) SegmentCollapsedXML {
	segment := SegmentCollapsedXML{
		Coordinates: convertPosListToCoords(segmentXML.PosList, minimumDistanceBetweenCoords),
		Trails:      segmentXML.FootTrails,
	}
	// combining all slices, it is safe because it will always be of a single type,
	// this just makes it easier to process later
	segment.Trails = append(segment.Trails, segmentXML.BikeTrails...)
	segment.Trails = append(segment.Trails, segmentXML.SkiTrails...)
	segment.Trails = append(segment.Trails, segmentXML.OtherTrails...)
	return segment
}

func convertPosListToCoords(posList string, minimumDistanceBetweenCoords float64) []Coordinate {
	splitWords := strings.Split(posList, " ")
	coords := make([]Coordinate, 0, len(splitWords)/2)
	for i := 0; i < len(splitWords); i += 2 {
		coords = append(coords, stringToCoordinate(splitWords[i], splitWords[i+1]))
	}
	return makeNewCoordinatesList(coords, minimumDistanceBetweenCoords)
}

func stringToCoordinate(lat, lon string) Coordinate {
	x, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		log.Fatal(err)
	}
	y, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		log.Fatal(err)
	}
	return Coordinate{x, y}
}
func isbeyondMinimumDistance(coords []Coordinate, p Coordinate, threshold float64) bool {
	return len(coords) == 0 || coords[len(coords)-1].DistanceInHaversine(p) > threshold
}

func makeNewCoordinatesList(coords []Coordinate, minimumDistance float64) []Coordinate {
	if len(coords) < 4 {
		return coords
	}
	newCoordinates := make([]Coordinate, 0, len(coords))
	for _, p := range coords {
		if isbeyondMinimumDistance(newCoordinates, p, minimumDistance) {
			newCoordinates = append(newCoordinates, p)
		}
	}
	if coords[len(coords)-1].DistanceInHaversine(newCoordinates[len(newCoordinates)-1]) > minimumDistance/2.0 {
		newCoordinates = append(newCoordinates, coords[len(coords)-1])
	}
	return newCoordinates
}

func createPolylineString(coords []Coordinate) string {
	array := make([][]float64, 0, len(coords)/2)
	for _, c := range coords {
		array = append(array, []float64{c.Lat, c.Lon})
	}
	return string(polyline.EncodeCoords(array))
}

func ConnectXMLSegments(segments []SegmentCollapsedXML, thresholdDistance float64) TrailsXML {
	trails := getAllUniqueTrails(segments)
	trails.spliceSegments(segments, thresholdDistance)
	return trails
}

func getAllUniqueTrails(segments []SegmentCollapsedXML) TrailsXML {
	trails := make(TrailsXML, 0, len(segments))
	for _, segment := range segments {
		for _, trail := range segment.Trails {
			if !trails.isTrailPresent(trail) && trail.Name != "" && trail.Name != "Ukjent" {
				trails = append(trails, trail)
			}
		}
	}
	return trails
}

func ReformatDifficulty(code string) string {
	switch code {
	case "G":
		return "Enkel"
	case "B":
		return "Middels"
	case "R":
		return "Krevende"
	case "S":
		return "Ekspert"
	default:
		return ""
	}
}

func ReformatSubtype(code string) string {
	switch code {
	case "1":
		return "Padlerute"
	case "2":
		return "Riderute"
	case "3":
		return "Trugerute"
	case "4":
		return "Klatring"
	default:
		return "Annet"
	}
}

func IntToBase62(num int) string {
	const base = 62
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	if num == 0 {
		return string(charset[0])
	}
	var res string
	for num > 0 {
		res = string(charset[num%base]) + res
		num = num / base
	}
	return res
}

func EstimateTrailLocation(encodedPolyline string) Coordinate {
	coords, _, err := polyline.DecodeCoords([]byte(encodedPolyline))
	if err != nil {
		log.Fatal(err)
	}
	firstLat, firstLon := coords[0][0], coords[0][1]
	lastLat, lastLon := coords[len(coords)-1][0], coords[len(coords)-1][1]

	return Coordinate{
		Lat: (firstLat + lastLat) / 2.0,
		Lon: (firstLon + lastLon) / 2.0,
	}
}
