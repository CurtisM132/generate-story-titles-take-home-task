package main

import (
	"fmt"
	"generate-story-titles/geospatial"
	parser "generate-story-titles/metadata-parser"
	reversegeocoder "generate-story-titles/reverse-geocoder"
	"io/ioutil"
	"log"
	"math"
	"path/filepath"
	"time"
)

const APIKey = ""

func main() {
	geocoder, err := reversegeocoder.NewReverseGeocoder(APIKey)
	if err != nil {
		log.Fatal(err)
	}

	// Get all CSV files in the directory
	files, err := ioutil.ReadDir("./test-data")
	if err != nil {
		log.Fatal(err)
	}

	storyTitles := []string{}

	// Parse and generate a story title for each CSV file in the directory
	for _, file := range files {
		path := "./test-data/" + file.Name()

		if !file.IsDir() && filepath.Ext(path) == ".csv" {
			records, err := parser.ParseCSVFile(path)
			if err != nil {
				log.Printf("failed to parse CSV file - %s", err)
			}

			if len(records) > 0 {
				storyTitle, err := generateStoryTitle(records, geocoder)
				if err != nil {
					continue
				}
				storyTitles = append(storyTitles, storyTitle)
			}
		}
	}

	for _, storyTitle := range storyTitles {
		fmt.Println(storyTitle)
	}
}

// generateStoryTitle Generates and returns a story title for a set of picture metadata records
func generateStoryTitle(records []*parser.Record, geocoder *reversegeocoder.ReverseGeocoder) (string, error) {
	storyTitle := generateTripStoryTitleComponent(records)

	locationComponent, err := generateLocationStoryTitleComponent(records, geocoder)
	if err != nil {
		return "", err
	}

	return storyTitle + locationComponent, nil
}

// generateTripStoryTitleComponent Generates and returns the trip aspect of the story title
func generateTripStoryTitleComponent(records []*parser.Record) string {
	tripLength, err := calculateTripLength(records)
	if err != nil {
		log.Printf("failed to calculate trip length - %s - falling back to default of %d", err, tripLength)
	}

	if tripLength <= 3 && records[len(records)/2].Datetime.Weekday() == time.Saturday {
		return "A weekend trip to "
	}

	// TODO: Check if that trip had any cultural significance (e.g., day of the dead)

	// TODO: Check if there was extreme weather during the trip

	return "A trip to "
}

// generateLocationStoryTitleComponent Generates and returns the location aspect of the story title
func generateLocationStoryTitleComponent(records []*parser.Record, geocoder *reversegeocoder.ReverseGeocoder) (
	string, error) {
	// Transform into a coordinates datasets
	coords := []geospatial.Coordinates{}
	for i := range records {
		coord := geospatial.Coordinates{
			Lat: records[i].Lat,
			Lon: records[i].Lon,
		}
		coords = append(coords, coord)
	}

	tripLocation, err := getTripLocation(coords, geocoder)
	if err != nil {
		return "", fmt.Errorf("failed to get trip location - %s", err)
	}

	locationSpread, err := geospatial.CalculateSpread(coords)
	if err != nil {
		log.Printf("failed to calculate trip location spread (in metres) - %s", err)
		return tripLocation.Locality, nil
	}

	// Based on how large the spread is, a different granularity of location is returned
	switch {
	case locationSpread > 15000:
		return tripLocation.Country, nil
	case locationSpread > 10000:
		return tripLocation.Locality, nil
	default:
		return tripLocation.Sublocality, nil
	}
}

// calculateTripLength Takes in a set of CSV records and returns the largest trip length (in days)
// that encompasses all records
func calculateTripLength(records []*parser.Record) (int, error) {
	if len(records) == 1 {
		return 1, nil
	}

	startDate := records[0].Datetime
	endDate := records[len(records)-1].Datetime

	return int(math.Abs(math.Ceil(endDate.Sub(startDate).Hours() / 24))), nil
}

// getTripLocation Takes a set of coordinates, calculates the midpoint of the dataset,
// then returns the human readable location of that midpoint (e.g., New York, Manhattan, United States)
func getTripLocation(coordDataset []geospatial.Coordinates, geocoder *reversegeocoder.ReverseGeocoder) (
	geospatial.Location, error) {
	// Get the mid point of the dataset so we can use that as the location to reverse geocode
	midPoint, err := geospatial.CalculateMidPoint(coordDataset)
	if err != nil {
		// If we fail for whatever reason then use the location from the middle of the dataset
		midPoint.Lat = coordDataset[len(coordDataset)/2].Lat
		midPoint.Lon = coordDataset[len(coordDataset)/2].Lon
		log.Printf("failed to calculate geographical midpoint of the trip - %s - falling back to %f, %f", err,
			midPoint.Lat, midPoint.Lon)
	}

	loc, err := geocoder.LocationFromLatLon(midPoint)
	if err != nil {
		return geospatial.Location{}, err
	}

	return loc, nil
}
