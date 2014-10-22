package gojsonld

import (
	"strings"
)

func Expand(input interface{}, options *Options) ([]interface{}, error) {
	// 1)
	//TODO implement API with promises
	// 2)
	//TODO handle remote context
	inputString, isString := input.(string)
	if isString && strings.Contains(inputString, ":") {
		remoteDocument, remoteErr := options.DocumentLoader.
			loadDocument(inputString)
		if remoteErr != nil {
			return nil, LOADING_DOCUMENT_FAILED
		}
		if options.Base == "" {
			options.Base = inputString
		}
		input = remoteDocument.document
	}
	// 3)
	activeContext := Context{}
	activeContext.init(options)
	// 4)
	if options.ExpandContext != nil {
		var expandContext interface{}
		mapContext,
			hasContext := options.ExpandContext.(map[string]interface{})["@context"]
		if hasContext {
			expandContext = mapContext
		}
		emptyArray := make([]string, 0)
		tmpContext, parseErr := parse(&activeContext, expandContext, emptyArray)
		if !isNil(parseErr) {
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
	} else if isNil(expanded) {
		expanded = make([]interface{}, 0)
	}
	if _, isArray := expanded.([]interface{}); !isArray {
		tmpArray := make([]interface{}, 0)
		tmpArray = append(tmpArray, expanded)
		expanded = tmpArray
	}
	// 7)
	return expanded.([]interface{}), nil
}

func Compact(input interface{}, context interface{},
	options *Options) (map[string]interface{}, error) {
	// 1)
	// TODO use promises
	// 2)
	expanded, expandErr := Expand(input, options)
	if !isNil(expandErr) {
		return nil, expandErr
	}
	//7)
	contextMap, isMap := context.(map[string]interface{})
	contextValue, hasContext := contextMap["@context"]
	if isMap && hasContext {
		context = contextValue
	}
	activeContext := Context{}
	activeContext.init(options)
	emptyArray := make([]string, 0)
	tmpContext, parseErr := parse(&activeContext, context, emptyArray)
	if !isNil(parseErr) {
		return nil, parseErr
	}
	activeContext = *tmpContext
	//8)
	//TODO check passing "" works
	compacted, compactErr := compact(&activeContext, "", expanded,
		options.CompactArrays)
	if !isNil(compactErr) {
		return nil, compactErr
	}
	//final step of Compaction algorithm
	compactedArray, isArray := compacted.([]interface{})
	if isArray {
		if len(compactedArray) == 0 {
			compacted = make(map[string]interface{}, 0)
		} else {
			graphArg := "@graph"
			iri, compactErr := compactIri(&activeContext, &graphArg, nil,
				true, false)
			if !isNil(compactErr) {
				return nil, compactErr
			}
			tmpMap := make(map[string]interface{}, 0)
			tmpMap[*iri] = compacted
			compacted = tmpMap
		}
	}
	if !isNil(compacted) && !isNil(context) {
		contextMap, isMap := context.(map[string]interface{})
		contextArray, isArray := context.([]interface{})
		if (isMap && len(contextMap) > 0) || (isArray && len(contextArray) > 0) {
			compacted.(map[string]interface{})["@context"] = context
		}
	}
	//9
	return compacted.(map[string]interface{}), nil
}
