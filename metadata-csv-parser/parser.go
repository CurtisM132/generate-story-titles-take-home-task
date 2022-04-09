package metadatacsvparser

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

type Record struct {
	Datetime string
	Lat      float64
	Lon      float64
}

func ParseCSVFile(filePath string) ([]*Record, error) {
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
		p.Datetime = row[0]
		if p.Lat, err = strconv.ParseFloat(row[1], 64); err != nil {
			continue
		}
		if p.Lon, err = strconv.ParseFloat(row[2], 64); err != nil {
			continue
		}
		records = append(records, p)
	}
}
