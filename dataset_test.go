package gojsonld

import (
	"fmt"
	"testing"
)

func TestDataset0001(t *testing.T) {
	datasetBytes, readErr := ReadDatasetFromFile(
		test_dir + "fromRdf-0001-in.nq")
	if !isNil(readErr) {
		t.Error(readErr.Error())
		return
	}
	dataset, parseErr := parseDataset(datasetBytes)
	if !isNil(parseErr) {
		t.Error(parseErr.Error())
		return
	}
	fmt.Println(dataset.Graphs["@default"])
}
