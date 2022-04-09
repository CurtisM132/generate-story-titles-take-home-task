package reversegeocoder

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"generate-story-titles/geospatial"
)

type ReverseGeocoder struct {
	APIKey string
}

// const APIURL = "https://maps.googleapis.com/maps/api/geocode/json?latlng=%f,%f&result_type=sublocality|administrative_area_level_1|country&key=%s"
const APIURL = "https://maps.googleapis.com/maps/api/geocode/json?latlng=%f,%f&key=%s"

func NewReverseGeocoder(apiKey string) (*ReverseGeocoder, error) {
	if apiKey == "" {
		return nil, errors.New("no API key supplied")
	}

	return &ReverseGeocoder{
		APIKey: apiKey,
	}, nil
}

func (r *ReverseGeocoder) LocationFromLatLon(coords geospatial.Coordinates) (geospatial.Location, error) {
	url := fmt.Sprintf(APIURL, coords.Lat, coords.Lon, r.APIKey)

	fmt.Println(url)

	resp, err := r.sendHTTPRequest(http.MethodGet, url)
	if err != nil {
		return geospatial.Location{}, err
	}

	loc, err := r.parseLocationResponse(resp)
	if err != nil {
		return geospatial.Location{}, err
	}

	return loc, nil
}

func (r *ReverseGeocoder) sendHTTPRequest(reqType string, url string) (string, error) {
	switch reqType {
	case http.MethodGet:
		resp, err := http.Get(url)
		if err != nil {
			return "", fmt.Errorf("failed to GET %s - %s", url, err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to parse response body from GET to %s - %s", url, err)
		}

		return string(body), nil
	default:
		return "", fmt.Errorf("unknown request type - %s", reqType)
	}
}
