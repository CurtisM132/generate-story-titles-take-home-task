package reversegeocoder

import (
	"encoding/json"
	"fmt"
	"generate-story-titles/geospatial"
)

type response struct {
	Results []result
	Status  string
}

type result struct {
	AddressComponents []address `json:"address_components"`
}

type address struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

// parseLocationResponseBody Takes the body of a reverse geocode response then extracts and returns the pertinent information
func (r *ReverseGeocoder) parseLocationResponseBody(body string) (geospatial.Location, error) {
	resp := response{}

	err := json.Unmarshal([]byte(body), &resp)
	if err != nil {
		return geospatial.Location{}, fmt.Errorf("failed to parse request body - %s", err)
	}

	// Validate that the request actually succeeded
	if resp.Status != "OK" {
		return geospatial.Location{}, fmt.Errorf("request failed with status - %s", resp.Status)
	}

	// Extract the important parts of the address
	location := geospatial.Location{}
	for i := range resp.Results {
		for p := range resp.Results[i].AddressComponents {
			for _, addType := range resp.Results[i].AddressComponents[p].Types {
				if addType == "locality" {
					location.Locality = resp.Results[i].AddressComponents[p].LongName
					break
				} else if addType == "sublocality" || addType == "administrative_area_level_1" {
					location.Sublocality = resp.Results[i].AddressComponents[p].LongName
					break
				} else if addType == "country" {
					location.Country = resp.Results[i].AddressComponents[p].LongName
					break
				}
			}
		}
	}

	return location, nil
}
