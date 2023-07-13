package web

import (
	"api/internal/dto"
	"errors"
	"net/url"
	"strconv"
	"strings"
)

func getQueryStr(url *url.URL, query string, altquery string) (string, error) {
	for key, value := range url.Query() {
		if key == query || key == altquery {
			return strings.Join(value, ","), nil
		}
	}
	return "", errors.New("could not find key: " + query + "' or '" + altquery + "'")
}

func getQueryFloat(url *url.URL, query string, altquery string) (float64, error) {
	str, err := getQueryStr(url, query, altquery)
	if err != nil {
		return 0.0, err
	}
	number, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0.0, errors.New("could not parse number")
	}
	return number, nil
}

func getQueryCoordinate(url *url.URL, query string, altquery string) (dto.Coordinate, error) {
	str, err := getQueryStr(url, query, altquery)
	if err != nil {
		return dto.Coordinate{}, err
	}
	components, err := extractTwoNumberedComponents(str, ",")
	if err != nil {
		return dto.Coordinate{0, 0}, errors.New("could not parse coordinate, " + err.Error())
	}
	return dto.Coordinate{
		Lat: components[0],
		Lon: components[1],
	}, nil
}

func getQueryIntInterval(url *url.URL, query string, altquery string) ([2]int, error) {
	defaultInterval := [2]int{1, 50}
	str, err := getQueryStr(url, query, altquery)
	components, err := extractTwoNumberedComponents(str, "-")
	if err != nil {
		return defaultInterval, errors.New("could not parse interval" + err.Error())
	}
	return [2]int{int(components[0]), int(components[1])}, nil
}

func extractTwoNumberedComponents(txt, seperator string) ([2]float64, error) {
	response := [2]float64{}
	components := strings.Split(txt, seperator)
	if len(components) != 2 {
		return response, errors.New("incorrect number format")
	}

	for i := 0; i < 2; i++ {
		var err error
		response[i], err = strconv.ParseFloat(components[i], 64)
		if err != nil {
			return response, errors.New("could not parse numbers")
		}
	}
	return response, nil
}

func interpretTrailType(txt string, err error) (dto.Activity, error) {
	if err != nil {
		return dto.FOOT, err
	}
	switch txt {
	case "foot", "f":
		return dto.FOOT, nil
	case "bike", "b":
		return dto.BIKE, nil
	case "ski", "s":
		return dto.SKI, nil
	case "other", "o":
		return dto.OTHER, nil
	default:
		return dto.FOOT, errors.New("incorrect input for trail type")
	}
}
