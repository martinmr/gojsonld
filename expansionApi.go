package gojsonld

import (
	"sort"
	"strings"
)

func expand(activeContext *Context, activeProperty string,
	element interface{}) (interface{}, error) {
	//1)
	if element == nil {
		return nil, nil
	}
	// 2)
	if isScalar(element) {
		if activeProperty == "" || activeProperty == "@graph" {
			return nil, nil
		}
		return expandValue(activeContext, activeProperty, element)
	}
	// 3)
	if elementArray, isArray := element.([]interface{}); isArray {
		// 3.1)
		result := make([]interface{}, 0)
		for _, item := range elementArray {
			// 3.2.1)
			expandedItem, expandErr := expand(activeContext, activeProperty, item)
			//TODO verify handling error is done correctly
			if expandErr != nil {
				return nil, expandErr
			}
			// 3.2.2)
			expandedArray, isArray := expandedItem.([]interface{})
			_, isList := expandedItem.(map[string]interface{})["@list"]
			if (activeProperty == "@list" ||
				activeContext.getContainer(activeProperty) == "@list") &&
				(isArray || isList) {
				return nil, LIST_OF_LISTS
			}
			// 3.2.3)
			if isArray {
				for _, expandedItem := range expandedArray {
					result = append(result, expandedItem)
				}
			} else if expandedItem != nil {
				result = append(result, expandedItem)
			}
		}
		// 3.3)
		return result, nil
	}
	// 4)
	elementMap, isMap := element.(map[string]interface{})
	if !isMap {
		return nil, INVALID_INPUT
	}
	// 5)
	if context, containsContext := elementMap["@context"]; containsContext {
		processedContext, processErr := activeContext.parse(context, nil)
		if processErr != nil {
			//TODO check error handling is correct
			return nil, processErr
		}
		activeContext = processedContext
	}
	// 6)
	result := make(map[string]interface{}, 0)
	//7
	keys := make([]string, 0)
	for key := range elementMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := elementMap[key]
		var expandedValue interface{}
		// 7.1)
		if key == "@context" {
			continue
		}
		// 7.2)
		expandedProperty, expandErr := activeContext.expandIri(key, false, true,
			nil, nil)
		if expandErr != nil {
			//TODO check error handling
			return nil, expandErr
		}
		// 7.3)
		//TODO check "" equals nil in expandedProperty == ""
		if expandedProperty == "" || (!strings.Contains(expandedProperty, ":") &&
			!isKeyword(expandedProperty)) {
			continue
		}
		//7.4)
		if isKeyword(expandedProperty) {
			// 7.4.1)
			if activeProperty == "@reverse" {
				return nil, INVALID_REVERSE_PROPERTY_MAP
			}
			// 7.4.2)
			if _, containsProperty := result[expandedProperty]; containsProperty {
				return nil, COLLIDING_KEYWORDS
			}
			// 7.4.3)
			if valueString, isString := value.(string); !isString &&
				expandedProperty == "@id" {
				return nil, INVALID_ID_VALUE
			} else {
				//TODO check passing nil values to expandIri does not break
				// the code
				tmpExpandedValue, expandedValueErr := activeContext.expandIri(valueString,
					true, false, nil, nil)
				if expandedValueErr == nil {
					expandedValue = tmpExpandedValue
				} else {
					//TODO check error is being handled correctly
					return nil, expandedValueErr
				}
			}
			// 7.4.4)
			//TODO check logic is correct
			valueString, isString := value.(string)
			valueArray, isArray := value.([]interface{})
			valueMap, isMap := value.(map[string]interface{})
			if expandedProperty == "@type" {
				if isString {
					tmpExpandedValue, expandErr := activeContext.expandIri(valueString,
						true, true, nil, nil)
					if expandErr == nil {
						expandedValue = tmpExpandedValue
					} else {
						return nil, expandErr
					}
				} else if isArray {
					expandedArray := make([]string, 0)
					for _, item := range valueArray {
						if itemString, isItemString := item.(string); isItemString {
							tmpExpandedValue, expandErr := activeContext.
								expandIri(itemString, true, true, nil, nil)
							if expandErr == nil {
								expandedArray = append(expandedArray, tmpExpandedValue)
							} else {
								return nil, expandErr
							}
						} else {
							return nil, INVALID_TYPE_VALUE
						}
					}
					expandedValue = expandedArray
				} else if isMap {
					//TODO check if empty map check should be part of the spec
					if len(valueMap) != 0 {
						return nil, INVALID_TYPE_VALUE
					}
					expandedValue = value
				} else {
					return nil, INVALID_TYPE_VALUE
				}
			}
			// 7.4.5)
			if expandedProperty == "@graph" {
				tmpExpandedValue, expandErr := expand(activeContext, "@graph", value)
				if expandErr == nil {
					expandedValue = tmpExpandedValue
				} else {
					return nil, expandErr
				}
			}
			// 7.4.6)
			if expandedProperty == "@value" {
				if value != nil || !isScalar(value) {
					return nil, INVALID_VALUE_OBJECT_VALUE
				}
				expandedValue = value
				if value == nil {
					result["@value"] = nil
					continue
				}
			}
			// 7.4.7)
			if expandedProperty == "@language" {
				if !isString {
					return nil, INVALID_LANGUAGE_TAGGED_STRING
				}
				expandedValue = strings.ToLower(valueString)
			}
			// 7.4.8)
			if expandedProperty == "@index" {
				if !isString {
					return nil, INVALID_INDEX_VALUE
				}
				expandedValue = value
			}
			// 7.4.9)
			if expandedProperty == "@list" {
				// 7.4.9.1)
				//TODO check empty string works the same as null
				if activeProperty == "" || activeProperty == "@graph" {
					continue
				}
				// 7.4.9.2)
				tmpExpandedValue, expandErr := expand(activeContext, activeProperty,
					value)
				if expandErr == nil {
					expandedValue = tmpExpandedValue
				} else {
					return nil, expandErr
				}
				if _, isExpandedList := expandedValue.([]interface{}); isExpandedList {
					return nil, LIST_OF_LISTS
				}
			}
			// 7.4.10)
			if expandedProperty == "@set" {
				tmpExpandedValue, expandErr := expand(activeContext, activeProperty,
					value)
				if expandErr == nil {
					expandedValue = tmpExpandedValue
				} else {
					return nil, expandErr
				}
			}
			// 7.4.11)
			if expandedProperty == "@reverse" {
				if !isMap {
					return nil, INVALID_REVERSE_VALUE
				}
				// 7.4.11.1)
				tmpExpandedValue, expandErr := expand(activeContext, "@reverse",
					value)
				if expandErr == nil {
					expandedValue = tmpExpandedValue
				} else {
					return nil, expandErr
				}
				// 7.4.11.2)
				expandedValueMap, isExpandedMap := expandedValue.(map[string]interface{})
				if !isExpandedMap {
					//TODO check error handling
					return nil, UNKNOWN_ERROR
				}
				reverse, containsReserve := expandedValueMap["@reverse"]
				reverseMap, isReverseMap := reverse.(map[string]interface{})
				if containsReserve && isReverseMap {
					for property, item := range reverseMap {
						// 7.4.11.2.1)
						if _, containsProperty := result[property]; !containsProperty {
							result[property] = make([]interface{}, 0)
						}
						// 7.4.11.2.1)
						//TODO check if needs to handle lists differently
						resultArray := result[property].([]interface{})
						result[property] = append(resultArray, item)
					}
				}
				// 7.4.11.3)
				if containsReserve && len(expandedValueMap) > 1 {
					// 7.4.11.3.1)
					if _, containsReserve := result["@reverse"]; !containsReserve {
						result["@reverse"] = make(map[string]interface{})
					}
					// 7.4.11.3.2)
					// Naming the mapping of reverse in result to reverse result instead
					// of reverse map as in the spec because I am already using reverseMap
					// to hold the casting to a map of the variable reverse
					//TODO check that changes to reverseResultMap are reflected in reverseResult
					reverseResult := result["@reverse"]
					reverseResultMap, isResultMap := reverseResult.(map[string]interface{})
					if !isResultMap {
						//TODO check error handling
						return nil, UNKNOWN_ERROR
					}
					// 7.4.11.3.3)
					for property, items := range expandedValueMap {
						if property == "@reverse" {
							continue
						}
						// 7.4.11.3.3.1)
						itemsArray, isItemsArray := items.([]interface{})
						if !isItemsArray {
							//TODO check error handling
							return nil, UNKNOWN_ERROR
						}
						for _, item := range itemsArray {
							// 7.4.11.3.3.1.1)
							if isListObject(item) || isValueObject(item) {
								return nil, INVALID_REVERSE_PROPERTY_VALUE
							}
							// 7.4.11.3.3.1.2)
							_, containsProperty := reverseResultMap[property]
							if !containsProperty {
								reverseResultMap[property] = make([]interface{}, 0)
							}
							// 7.4.11.3.3.1.3)
							reverseArray := reverseResultMap[property].([]interface{})
							reverseResultMap[property] = append(reverseArray, item)
						}
					}
					//Reassign in case reverseResultMap is a copy of the original
					result["@reverse"] = reverseResultMap
				}
				// 7.4.11.4)
				continue
			}
			// 7.4.12)
			if expandedValue != nil {
				result[expandedProperty] = expandedValue
			}
			// 7.4.13)
			continue
		} else if _, isValueMap := value.(map[string]interface{}); isValueMap &&
			activeContext.getContainer(key) == "@language" {
			// 7.5)
			// 7.5.1)
			valueMap := value.(map[string]interface{})
			expandedValue = make([]interface{}, 0)
			// 7.5.2)
			keys := make([]string, 0)
			for key := range valueMap {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			for _, language := range keys {
				languageValue := valueMap[language]
				// 7.5.2.1)
				if _, isArray := languageValue.([]interface{}); !isArray {
					tmpArray := make([]interface{}, 0)
					tmpArray = append(tmpArray, languageValue)
					languageValue = tmpArray
				}
				languageArray := languageValue.([]interface{})
				// 7.5.2.2)
				for _, item := range languageArray {
					if _, isString := item.(string); !isString {
						return nil, INVALID_LANGUAGE_MAP_VALUE
					}
					newLanguageMap := make(map[string]interface{})
					newLanguageMap["@language"] = strings.ToLower(language)
					newLanguageMap["@value"] = item
					expandedValue = append(expandedValue.([]interface{}), newLanguageMap)
				}
			}
		} else if _, isValueMap := value.(map[string]interface{}); isValueMap &&
			activeContext.getContainer(key) == "@index" {
			// 7.6)
			// 7.1.6)
			valueMap := value.(map[string]interface{})
			expandedValue = make([]interface{}, 0)
			// 7.6.2)
			keys := make([]string, 0)
			for key := range valueMap {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			for _, index := range keys {
				indexValue := valueMap[index]
				// 7.6.2.1)
				if _, isArray := indexValue.([]interface{}); !isArray {
					tmpArray := make([]interface{}, 0)
					tmpArray = append(tmpArray, indexValue)
					indexValue = tmpArray
				}
				// 7.6.2.2)
				tmpIndexValue, expandErr := expand(activeContext, key, indexValue)
				if expandErr == nil {
					indexValue = tmpIndexValue
				} else {
					return nil, expandErr
				}
				// 7.6.2.3)
				indexArray := indexValue.([]interface{})
				for _, item := range indexArray {
					itemMap, isItemMap := item.(map[string]interface{})
					//TODO check error handling
					if !isItemMap {
						return nil, UNKNOWN_ERROR
					}
					if _, containsIndex := itemMap["@index"]; !containsIndex {
						itemMap["@index"] = index
					}
					expandedValue = append(expandedValue.([]interface{}), item)
				}
			}
		} else {
			// 7.7)
			tmpExpandedValue, tmpErr := expand(activeContext, key, value)
			if tmpErr == nil {
				expandedValue = tmpExpandedValue
			} else {
				return nil, expandErr
			}
		}
		// 7.8)
		if expandedValue == nil {
			continue
		}
		// 7.9)
		if !isListObject(expandedValue) && "@list" == activeContext.getContainer(key) {
			if _, isValueArray := expandedValue.([]interface{}); !isValueArray {
				tmpArray := make([]interface{}, 0)
				tmpArray = append(tmpArray, expandedValue)
				expandedValue = tmpArray
			}
			tmpMap := make(map[string]interface{}, 0)
			tmpMap["@list"] = expandedValue
			expandedValue = tmpMap
		} else if activeContext.isReverseProperty(key) {
			// 7.10)
			// 7.10.1)
			if _, containsReverse := result["@reverse"]; !containsReverse {
				result["@reverse"] = make(map[string]interface{})
			}
			// 7.10.2)
			reverseMap, isReverseMap := result["@reverse"].(map[string]interface{})
			if !isReverseMap {
				//TODO check error handling
				return nil, UNKNOWN_ERROR
			}
			// 7.10.3)
			if _, isExpandedArray := expandedValue.([]interface{}); !isExpandedArray {
				tmpArray := make([]interface{}, 0)
				tmpArray = append(tmpArray, expandedValue)
				expandedValue = tmpArray
			}
			// 7.10.4)
			expandedArray := expandedValue.([]interface{})
			for _, item := range expandedArray {
				// 7.10.4.1)
				if isValueObject(item) || isListObject(item) {
					return nil, INVALID_REVERSE_PROPERTY_VALUE
				}
				// 7.10.4.2)
				if _, containsProperty := reverseMap[expandedProperty]; !containsProperty {
					reverseMap[expandedProperty] = make([]interface{}, 0)
				}
				// 7.10.4.3)
				//TODO check if list needs to be handled different
				reverseMap[expandedProperty] = append(reverseMap[expandedProperty].([]interface{}),
					item)
			}
		} else if !activeContext.isReverseProperty(key) {
			// 7.11)
			// 7.11.1)
			if _, containsProperty := result[expandedProperty]; !containsProperty {
				result[expandedProperty] = make([]interface{}, 0)
			}
			// 7.11.2
			//TODO check if need to handle lists differently
			result[expandedProperty] = append(result[expandedProperty].([]interface{}),
				expandedValue)
		}
	}
	// 8)
	if value, containsValue := result["@value"]; containsValue {
		//8.1)
		if !isValidValueObject(value) {
			return nil, INVALID_VALUE_OBJECT
		}
		// 8.2)
		if value == nil {
			//TODO what if value is a string and equals ""
			//should we consider that as a nil value for the string type
			result = nil
		} else if _, isValueString := value.(string); !isValueString {
			// 8.3)
			if _, containsLanguage := result["@language"]; containsLanguage {
				return nil, INVALID_LANGUAGE_TAGGED_VALUE
			}
		} else if typeVal, containsType := result["@type"]; containsType {
			// 8.4)
			//TODO complete isIRI method
			if !isIRI(typeVal) {
				return nil, INVALID_TYPED_VALUE
			}
		}
	} else if typeVal, containsType := result["@type"]; containsType {
		// 9)
		if _, isTypeArray := typeVal.([]interface{}); !isTypeArray {
			tmpArray := make([]interface{}, 0)
			tmpArray = append(tmpArray, typeVal)
			result["@type"] = tmpArray
		}
	} else {
		//TODO make sure logic is still correct
		// 10)
		_, containsSet := result["@set"]
		_, containsList := result["@list"]
		if containsSet || containsList {
			// 10.1)
			maxLen := 0
			if _, containsIndex := result["@index"]; containsIndex {
				maxLen = 2
			} else {
				maxLen = 1
			}
			if len(result) > maxLen {
				return nil, INVALID_SET_OR_LIST_OBJECT
			}
			// 10.2)
			if containsSet {
				// TODO check comment's validity
				// result becomes an array here, thus the remaining checks
				// will never be true from here on
				// so simply return the value rather than have to make
				// result an object and cast it with every
				// other use in the function.
				return result["@set"], nil
			}
		}
	}
	// 11)
	if _, containsLanguage := result["@language"]; containsLanguage &&
		len(result) == 1 {
		result = nil
	}
	// 12)
	// TODO check that checking for "" instead of nil does not break the algorithm
	if activeProperty == "" || activeProperty == "@graph" {
		// 12.1)
		_, containsValue := result["@value"]
		_, containsList := result["@list"]
		_, containsID := result["@id"]
		//TODO check it's correct to test if result != nil
		if result != nil && (len(result) == 0 || containsList || containsValue) {
			result = nil
		} else if result != nil && len(result) == 1 && containsID {
			// 12.2)
			result = nil
		}
	}
	// 13)
	return result, nil
	//TODO figure out if something needs to be done with the below paragraph
	//If, after the above algorithm is run, the result is a JSON object
	//that contains only an @graph key, set the result to the value of @graph's
	//value. Otherwise, if the result is null, set it to an empty array.
	//Finally, if the result is not an array, then set the result to an array
	//containing only the result
}

func expandValue(activeContext *Context, activeProperty string,
	value interface{}) (interface{}, error) {
	// 1)
	result := make(map[string]interface{})
	termDefinitions := activeContext.termDefinitions
	typeValue, hasType := termDefinitions["@type"]
	if hasType && typeValue == "@id" {
		expandedValue, expandErr := activeContext.expandIri(value.(string),
			true, false, nil, nil)
		if expandErr == nil {
			result["@id"] = expandedValue
			return result, nil
		} else {
			return nil, expandErr
		}
	}
	// 2)
	if hasType && typeValue == "@vocab" {
		expandedValue, expandErr := activeContext.expandIri(value.(string), true, true, nil, nil)
		if expandErr == nil {
			// TODO make sure key is actually @id
			result["@id"] = expandedValue
			return result, nil
		} else {
			return nil, expandErr
		}
	}
	// 3)
	result["@value"] = value
	// 4)
	if hasType {
		result["@type"] = typeValue
	} else if _, isString := value.(string); isString {
		// 5.1)
		if language, hasLanguage := termDefinitions["@language"]; hasLanguage {
			_, isString := language.(string)
			//TODO check if we need to check for the empty string
			if isString {
				result["@language"] = language
			} else if defaultLanguage, hasDefaultLanguage := activeContext.table["language"]; hasDefaultLanguage {
				result["@language"] = defaultLanguage
			}
		}
	}
	// 6)
	return result, nil
}
