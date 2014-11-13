package gojsonld

import (
	"fmt"
	"testing"
)

func TestDataset0001(t *testing.T) {
	dataset, parseErr := ReadDatasetFromFile(
		test_dir + "fromRdf-0001-in.nq")
	if !isNil(parseErr) {
		t.Error(parseErr.Error())
		return
	}
	fmt.Println(dataset.Graphs["@default"])
}
