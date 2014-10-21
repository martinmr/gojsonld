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

func TestExpand0020(t *testing.T) {
	testExpand("expand-0020-in.jsonld", "expand-0020-out.jsonld", t)
}

func TestExpand0021(t *testing.T) {
	testExpand("expand-0021-in.jsonld", "expand-0021-out.jsonld", t)
}

func TestExpand0022(t *testing.T) {
	testExpand("expand-0022-in.jsonld", "expand-0022-out.jsonld", t)
}

func TestExpand0023(t *testing.T) {
	testExpand("expand-0023-in.jsonld", "expand-0023-out.jsonld", t)
}

func TestExpand0024(t *testing.T) {
	testExpand("expand-0024-in.jsonld", "expand-0024-out.jsonld", t)
}

func TestExpand0025(t *testing.T) {
	testExpand("expand-0025-in.jsonld", "expand-0025-out.jsonld", t)
}

func TestExpand0026(t *testing.T) {
	testExpand("expand-0026-in.jsonld", "expand-0026-out.jsonld", t)
}

func TestExpand0027(t *testing.T) {
	testExpand("expand-0027-in.jsonld", "expand-0027-out.jsonld", t)
}

func TestExpand0028(t *testing.T) {
	testExpand("expand-0028-in.jsonld", "expand-0028-out.jsonld", t)
}

func TestExpand0029(t *testing.T) {
	testExpand("expand-0029-in.jsonld", "expand-0029-out.jsonld", t)
}

func TestExpand0030(t *testing.T) {
	testExpand("expand-0030-in.jsonld", "expand-0030-out.jsonld", t)
}

func TestExpand0031(t *testing.T) {
	testExpand("expand-0031-in.jsonld", "expand-0031-out.jsonld", t)
}

func TestExpand0032(t *testing.T) {
	testExpand("expand-0032-in.jsonld", "expand-0032-out.jsonld", t)
}

func TestExpand0033(t *testing.T) {
	testExpand("expand-0033-in.jsonld", "expand-0033-out.jsonld", t)
}

func TestExpand0034(t *testing.T) {
	testExpand("expand-0034-in.jsonld", "expand-0034-out.jsonld", t)
}

func TestExpand0035(t *testing.T) {
	testExpand("expand-0035-in.jsonld", "expand-0035-out.jsonld", t)
}

func TestExpand0036(t *testing.T) {
	testExpand("expand-0036-in.jsonld", "expand-0036-out.jsonld", t)
}

func TestExpand0037(t *testing.T) {
	testExpand("expand-0037-in.jsonld", "expand-0037-out.jsonld", t)
}

func TestExpand0038(t *testing.T) {
	testExpand("expand-0038-in.jsonld", "expand-0038-out.jsonld", t)
}

func TestExpand0039(t *testing.T) {
	testExpand("expand-0039-in.jsonld", "expand-0039-out.jsonld", t)
}

func TestExpand0040(t *testing.T) {
	testExpand("expand-0040-in.jsonld", "expand-0040-out.jsonld", t)
}

func TestExpand0041(t *testing.T) {
	testExpand("expand-0041-in.jsonld", "expand-0041-out.jsonld", t)
}

func TestExpand0042(t *testing.T) {
	testExpand("expand-0042-in.jsonld", "expand-0042-out.jsonld", t)
}

func TestExpand0043(t *testing.T) {
	testExpand("expand-0043-in.jsonld", "expand-0043-out.jsonld", t)
}

func TestExpand0045(t *testing.T) {
	testExpand("expand-0045-in.jsonld", "expand-0045-out.jsonld", t)
}

func TestExpand0046(t *testing.T) {
	testExpand("expand-0046-in.jsonld", "expand-0046-out.jsonld", t)
}

func TestExpand0047(t *testing.T) {
	testExpand("expand-0047-in.jsonld", "expand-0047-out.jsonld", t)
}

func TestExpand0048(t *testing.T) {
	testExpand("expand-0048-in.jsonld", "expand-0048-out.jsonld", t)
}

func TestExpand0049(t *testing.T) {
	testExpand("expand-0049-in.jsonld", "expand-0049-out.jsonld", t)
}

func TestExpand0050(t *testing.T) {
	testExpand("expand-0050-in.jsonld", "expand-0050-out.jsonld", t)
}

func TestExpand0051(t *testing.T) {
	testExpand("expand-0051-in.jsonld", "expand-0051-out.jsonld", t)
}

func TestExpand0052(t *testing.T) {
	testExpand("expand-0052-in.jsonld", "expand-0052-out.jsonld", t)
}

func TestExpand0053(t *testing.T) {
	testExpand("expand-0053-in.jsonld", "expand-0053-out.jsonld", t)
}

func TestExpand0054(t *testing.T) {
	testExpand("expand-0054-in.jsonld", "expand-0054-out.jsonld", t)
}
