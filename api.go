package gojsonld

import (
	"strings"
)

type API struct {
	Options *Options
}

func (api *API) Expand(input interface{}, context interface{}) ([]interface{}, error) {
	// 1)
	//TODO implement API with promises
	// 2)
	//TODO handle remote context
	inputString, isString := input.(string)
	if isString && strings.Contains(inputString, ":") {
		remoteDocument, remoteErr := api.Options.documentLoader.
			loadDocument(inputString)
		if remoteErr != nil {
			return nil, LOADING_DOCUMENT_FAILED
		}
		if api.Options.base == "" {
			api.Options.base = inputString
		}
		input = remoteDocument.document
	}
	// 3)
	activeContext := Context{}
	activeContext.init(api.Options)
	// 4)
	if api.Options.expandContext != nil {
		var expandContext interface{}
		expandContext = api.Options.expandContext
		mapContext, hasContext := expandContext.(map[string]interface{})["@context"]
		if hasContext {
			expandContext = mapContext
		}
		emptyArray := make([]string, 0)
		tmpContext, parseErr := parse(&activeContext, expandContext, emptyArray)
		if parseErr != nil {
			return nil, parseErr
		}
		activeContext = *tmpContext
	}
	// 5)
	//TODO load remote context
	// 6)
	expanded, expandErr := expand(&activeContext, nil, input)
	if expandErr != nil {
		return nil, expandErr
	}
	// Final step of expansion algorithm
	expandedMap, isMap := expanded.(map[string]interface{})
	graphVal, hasGraph := expandedMap["@graph"]
	if isMap && hasGraph && len(expandedMap) == 1 {
		expanded = graphVal
	} else if expanded == nil {
		expanded = make([]interface{}, 0)
	}
	if _, isArray := expanded.([]interface{}); isArray {
		tmpArray := make([]interface{}, 0)
		tmpArray = append(tmpArray, expand)
		expanded = tmpArray
	}
	// 7)
	return expanded.([]interface{}), nil
}
