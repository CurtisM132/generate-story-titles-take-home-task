package metadataparser

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"
)

type Record struct {
	Datetime time.Time
	Lat      float64
	Lon      float64
}

func ParseCSVFile(filePath string) ([]*Record, error) {
	// The test data CSVs contain different time formats.
	// In a real system we would be able to sanitise the data that gets passed in so this wouldn't be needed.
	parseToDatetime := func(str string) (time.Time, error) {
		date, err := time.Parse("2006-01-02 15:04:05", str)
		if err != nil {
			date, err = time.Parse(time.RFC3339, str)
		}

		return date, err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvr := csv.NewReader(f)

	records := []*Record{}
	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return records, err
		}

		p := &Record{}
		if p.Datetime, err = parseToDatetime(row[0]); err != nil {
			continue
		}

		if p.Lat, err = strconv.ParseFloat(row[1], 64); err != nil {
			continue
		}
		if p.Lon, err = strconv.ParseFloat(row[2], 64); err != nil {
			continue
		}
		records = append(records, p)
	}
}
