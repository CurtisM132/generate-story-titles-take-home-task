package geospatial

import (
	"errors"
	"math"
)

const earthsRadius float64 = 6378100 // Earth radius in Metres

func ToRadians(deg float64) float64 {
	return float64(deg) * (math.Pi / 180.0)
}

func ToDegrees(rad float64) float64 {
	return float64(rad) * (180.0 / math.Pi)
}

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance returns the distance (in meters) between two points of a given longitude and latitude
func Distance(firstCoord, secondCoord Coordinates) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2 float64
	la1 = firstCoord.Lat * math.Pi / 180
	lo1 = firstCoord.Lon * math.Pi / 180
	la2 = secondCoord.Lat * math.Pi / 180
	lo2 = secondCoord.Lon * math.Pi / 180

	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * earthsRadius * math.Asin(math.Sqrt(h))
}

func CalculateUpperAndLowerBounds(coordDataset []Coordinates) (smallestCoord, largestCoord Coordinates, err error) {
	if len(coordDataset) == 0 {
		err = errors.New("no coordinates in dataset")
		return
	}

	smallestCoord.Lat = coordDataset[0].Lat
	smallestCoord.Lon = coordDataset[0].Lon

	largestCoord.Lat = coordDataset[0].Lat
	largestCoord.Lon = coordDataset[0].Lon

	for i := 1; i < len(coordDataset); i++ {
		if coordDataset[i].Lat < smallestCoord.Lat {
			smallestCoord.Lat = coordDataset[i].Lat
		} else if coordDataset[i].Lon < smallestCoord.Lon {
			smallestCoord.Lon = coordDataset[i].Lon
		}

		if coordDataset[i].Lat > largestCoord.Lat {
			largestCoord.Lat = coordDataset[i].Lat
		} else if coordDataset[i].Lon > largestCoord.Lon {
			largestCoord.Lon = coordDataset[i].Lon
		}
	}

	return
}

func CalculateMidPoint(coordDataset []Coordinates) (Coordinates, error) {
	if len(coordDataset) == 0 {
		return Coordinates{}, errors.New("no coordinates in dataset")
	}

	sCoord, lCoord, err := CalculateUpperAndLowerBounds(coordDataset)
	if err != nil {
		return Coordinates{}, err
	}

	dLon := ToRadians(lCoord.Lon - sCoord.Lon)

	//convert to radians
	lat1 := ToRadians(sCoord.Lat)
	lat2 := ToRadians(lCoord.Lat)
	lon1 := ToRadians(sCoord.Lon)

	Bx := math.Cos(lat2) * math.Cos(dLon)
	By := math.Cos(lat2) * math.Sin(dLon)
	lat3 := math.Atan2(math.Sin(lat1)+math.Sin(lat2), math.Sqrt((math.Cos(lat1)+Bx)*(math.Cos(lat1)+Bx)+By*By))
	lon3 := lon1 + math.Atan2(By, math.Cos(lat1)+Bx)

	return Coordinates{Lat: ToDegrees(lat3), Lon: ToDegrees(lon3)}, nil
}

func CalculateSpread(coordDataset []Coordinates) (float64, error) {
	if len(coordDataset) == 0 {
		return 0, errors.New("no coordinates in dataset")
	}

	sCoord, lCoord, err := CalculateUpperAndLowerBounds(coordDataset)
	if err != nil {
		return 0, err
	}

	return Distance(sCoord, lCoord), nil
}
