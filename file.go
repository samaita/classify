package classify

import (
	"bufio"
	"encoding/csv"
	"os"
)

func readFileToArrayString(filename string) ([]string, error) {
	var result []string

	csvFile, err := os.Open(filename)
	if err != nil {
		return result, err
	}
	r := csv.NewReader(bufio.NewReader(csvFile))

	doc, err := r.ReadAll()
	if err != nil {
		return result, err
	}

	for _, row := range doc {
		if len(row) > 0 {
			result = append(result, row[0])
		}
	}

	return result, nil
}

func readFileToMapString(filename string) (map[string]string, error) {
	result := make(map[string]string)

	csvFile, err := os.Open(filename)
	if err != nil {
		return result, err
	}
	r := csv.NewReader(bufio.NewReader(csvFile))

	doc, err := r.ReadAll()
	if err != nil {
		return result, err
	}

	for _, row := range doc {
		if len(row) > 1 {
			result[row[0]] = row[1]
		}
	}

	return result, nil
}
