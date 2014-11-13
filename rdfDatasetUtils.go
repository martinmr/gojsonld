package gojsonld

import (
	"io/ioutil"
)

func ReadDatasetFromFile(path string) (*Dataset, error) {
	file, fileErr := ioutil.ReadFile(path)
	if fileErr != nil {
		return nil, fileErr
	}
	dataset, parseErr := parseDataset(file)
	return dataset, parseErr
}
