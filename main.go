package main

import (
	"fmt"
	"generate-story-titles/geospatial"
	csv "generate-story-titles/metadata-csv-parser"
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

	// Get all csv files in the directory
	files, err := ioutil.ReadDir("./test-data")
	if err != nil {
		log.Fatal(err)
	}

	storyTitles := []string{}
	for _, file := range files {
		path := "./test-data/" + file.Name()

		if !file.IsDir() && filepath.Ext(path) == ".csv" {
			records, err := csv.ParseCSVFile(path)
			if err != nil {
				log.Printf("failed to parse CSV file - %s", err)
			}

			if len(records) > 0 {
				storyTitle, err := generateStoryTitle(records, geocoder)
				if err != nil {
					log.Printf("failed to generate story title for dataset - %s", err)
					continue
				}
				storyTitles = append(storyTitles, storyTitle)
			}
		}
	}

	// for _, storyTitle := range storyTitles {
	// 	fmt.Println(storyTitle)
	// }
}

func generateStoryTitle(records []*csv.Record, geocoder *reversegeocoder.ReverseGeocoder) (string, error) {
	tripLength, err := calculateTripLength(records)
	if err != nil {
		log.Printf("failed to calculate trip length - %s - falling back to default", err)
	}

	coords := []geospatial.Coordinates{}
	for i := range records {
		coord := geospatial.Coordinates{
			Lat: records[i].Lat,
			Lon: records[i].Lon,
		}
		coords = append(coords, coord)
	}

	locationSpread, err := geospatial.CalculateSpread(coords)
	if err != nil {
		log.Printf("failed to calculate trip spread (in metres) - %s", err)
	}

	tripLocation, err := getTripLocation(coords, geocoder)
	if err != nil {
		log.Printf("failed to get trip location - %s", err)
	}

	fmt.Println(tripLength, locationSpread, tripLocation)

	return "blah", nil
}

// calculateTripLength Takes in a set of CSV records and returns the largest trip length (in days)
// that encompasses all records
func calculateTripLength(records []*csv.Record) (int, error) {
	timeBetweenInDays := func(layout string) (int, error) {
		startDate, err := time.Parse(layout, records[0].Datetime)
		if err != nil {
			return 0, fmt.Errorf("failed to parse time - %s", err)
		}

		endDate, err := time.Parse(layout, records[len(records)-1].Datetime)
		if err != nil {
			return 0, fmt.Errorf("failed to parse time - %s", err)
		}

		return int(math.Abs(math.Ceil(endDate.Sub(startDate).Hours() / 24))), nil
	}

	if len(records) == 1 {
		return 1, nil
	}

	// The test data CSVs contain different time formats.
	// In a real system we would be able to sanitise the data that gets passed in so this wouldn't be needed.
	tripLength, err := timeBetweenInDays("2006-01-02 15:04:05")
	if err != nil {
		tripLength, err = timeBetweenInDays(time.RFC3339)
	}

	return tripLength, err
}

// getTripLocation Takes a set of coordinates, calculates the midpoint of the dataset,
// then returns the human readable location of that midpoint (e.g., New York, Manhattan, United States)
func getTripLocation(coordDataset []geospatial.Coordinates, geocoder *reversegeocoder.ReverseGeocoder) (
	geospatial.Location, error) {
	midPoint, err := geospatial.CalculateMidPoint(coordDataset)
	if err != nil {
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
