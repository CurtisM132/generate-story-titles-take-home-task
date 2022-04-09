package main

import (
	"fmt"
	reversegeocoder "generate-story-titles/reverse-geocoder"
	"log"
)

const APIKey = ""

func main() {
	geocoder, err := reversegeocoder.NewReverseGeocoder(APIKey)
	if err != nil {
		log.Fatal(err)
	}

	coords := reversegeocoder.Coordinates{
		Lat: 40.728808,
		Lon: -73.996106,
	}

	loc, err := geocoder.LocationFromLatLon(coords)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(loc)
}
