package reversegeocoder

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"generate-story-titles/geospatial"
)

type ReverseGeocoder struct {
	apiKey string
}

const APIURL = "https://maps.googleapis.com/maps/api/geocode/json?latlng=%f,%f&key=%s"

func NewReverseGeocoder(apiKey string) (*ReverseGeocoder, error) {
	if apiKey == "" {
		return nil, errors.New("no API key supplied")
	}

	return &ReverseGeocoder{
		apiKey: apiKey,
	}, nil
}

// LocationFromLatLon Calls the Google API to reverse geocode a coordinate pair and returns a struct
// containing the pertinent information
func (r *ReverseGeocoder) LocationFromLatLon(coords geospatial.Coordinates) (geospatial.Location, error) {
	url := fmt.Sprintf(APIURL, coords.Lat, coords.Lon, r.apiKey)

	resp, err := r.sendHTTPGETRequest(url)
	if err != nil {
		return geospatial.Location{}, err
	}

	loc, err := r.parseLocationResponseBody(resp)
	if err != nil {
		return geospatial.Location{}, err
	}

	return loc, nil
}

func (r *ReverseGeocoder) sendHTTPGETRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to GET %s - %s", url, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse response body from GET to %s - %s", url, err)
	}

	return string(body), nil
}
