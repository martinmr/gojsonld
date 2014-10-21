package gojsonld

import (
	"encoding/json"
	//"fmt"
	"reflect"
	"testing"
)

const test_dir = "./test_files/"

func testExpand(inputFile string, outputFile string, t *testing.T) {
	inputJson, jsonErr := ReadJSONFromFile(test_dir + inputFile)
	if jsonErr != nil {
		t.Error("Could not open input file")
		return
	}
	outputJson, jsonErr := ReadJSONFromFile(test_dir + outputFile)
	if jsonErr != nil {
		t.Error("Could not open output file")
		return
	}

	api := API{}
	api.Options = &Options{
		Base:           "",
		CompactArrays:  true,
		ExpandContext:  nil,
		DocumentLoader: NewDocumentLoader(),
	}
	expandedJson, expandErr := api.Expand(inputJson)
	if expandErr != nil {
		t.Error("Expansion failed with error ", expandErr.Error())
		return
	}

	expandedString, _ := json.MarshalIndent(expandedJson, "", "    ")
	outputString, _ := json.MarshalIndent(outputJson, "", "    ")
	if !reflect.DeepEqual(expandedJson, outputJson) {
		t.Error("Expected:\n", string(outputString), "\nGot:\n",
			string(expandedString))
	}

}

func TestExpand0001(t *testing.T) {
	testExpand("expand-0001-in.jsonld", "expand-0001-out.jsonld", t)
}

func TestExpand0002(t *testing.T) {
	testExpand("expand-0002-in.jsonld", "expand-0002-out.jsonld", t)
}

func TestExpand0003(t *testing.T) {
	testExpand("expand-0003-in.jsonld", "expand-0003-out.jsonld", t)
}

func TestExpand0004(t *testing.T) {
	testExpand("expand-0004-in.jsonld", "expand-0004-out.jsonld", t)
}

func TestExpand0005(t *testing.T) {
	testExpand("expand-0005-in.jsonld", "expand-0005-out.jsonld", t)
}

func TestExpand0006(t *testing.T) {
	testExpand("expand-0006-in.jsonld", "expand-0006-out.jsonld", t)
}

func TestExpand0007(t *testing.T) {
	testExpand("expand-0007-in.jsonld", "expand-0007-out.jsonld", t)
}

func TestExpand0008(t *testing.T) {
	testExpand("expand-0008-in.jsonld", "expand-0008-out.jsonld", t)
}

func TestExpand0009(t *testing.T) {
	testExpand("expand-0009-in.jsonld", "expand-0009-out.jsonld", t)
}

func TestExpand0010(t *testing.T) {
	testExpand("expand-0010-in.jsonld", "expand-0010-out.jsonld", t)
}

func TestExpand0011(t *testing.T) {
	testExpand("expand-0011-in.jsonld", "expand-0011-out.jsonld", t)
}

func TestExpand0012(t *testing.T) {
	testExpand("expand-0012-in.jsonld", "expand-0012-out.jsonld", t)
}

func TestExpand0013(t *testing.T) {
	testExpand("expand-0013-in.jsonld", "expand-0013-out.jsonld", t)
}

func TestExpand0014(t *testing.T) {
	testExpand("expand-0014-in.jsonld", "expand-0014-out.jsonld", t)
}

func TestExpand0015(t *testing.T) {
	testExpand("expand-0015-in.jsonld", "expand-0015-out.jsonld", t)
}

func TestExpand0016(t *testing.T) {
	testExpand("expand-0016-in.jsonld", "expand-0016-out.jsonld", t)
}

func TestExpand0017(t *testing.T) {
	testExpand("expand-0017-in.jsonld", "expand-0017-out.jsonld", t)
}

func TestExpand0018(t *testing.T) {
	testExpand("expand-0018-in.jsonld", "expand-0018-out.jsonld", t)
}

func TestExpand0019(t *testing.T) {
	testExpand("expand-0019-in.jsonld", "expand-0019-out.jsonld", t)
}
