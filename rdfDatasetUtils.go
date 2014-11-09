package gojsonld

import (
	"io/ioutil"
)

func ReadDatasetFromFile(path string) ([]byte, error) {
	file, fileErr := ioutil.ReadFile(path)
	if fileErr != nil {
		return nil, fileErr
	}
	return file, nil
}
