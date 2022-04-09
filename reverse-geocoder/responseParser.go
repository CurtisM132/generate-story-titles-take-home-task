package reversegeocoder

import (
	"encoding/json"
	"fmt"
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

func (r *ReverseGeocoder) parseLocationResponse(resp string) (Location, error) {
	m := response{}

	err := json.Unmarshal([]byte(resp), &m)
	if err != nil {
		return Location{}, fmt.Errorf("failed to parse request body - %s", err)
	}

	// Validate that the request actually succeeded
	if m.Status != "OK" {
		return Location{}, fmt.Errorf("failed to parse request body - %s", err)
	}

	// Extract the important parts of the address
	location := Location{}
	for i := range m.Results {
		for p := range m.Results[i].AddressComponents {
			for _, addType := range m.Results[i].AddressComponents[p].Types {
				if addType == "locality" {
					location.Locality = m.Results[i].AddressComponents[p].LongName
					break
				} else if addType == "sublocality" {
					location.Sublocality = m.Results[i].AddressComponents[p].LongName
					break
				} else if addType == "country" {
					location.Country = m.Results[i].AddressComponents[p].LongName
					break
				}
			}
		}
	}

	return location, nil
}
