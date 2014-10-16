package gojsonld

import (
	"strings"
)

func expand(activeContext *Context, activeProperty *string,
	element interface{}) (interface{}, error) {
	//1)
	if element == nil {
		return nil, nil
	}
	// 2)
	if isScalar(element) {
		if activeProperty == nil || *activeProperty == "@graph" {
			return nil, nil
		}
		return expandValue(activeContext, *activeProperty, element)
	}
	// 3)
	if elementArray, isArray := element.([]interface{}); isArray {
		// 3.1)
		result := make([]interface{}, 0)
		for _, item := range elementArray {
			// 3.2.1)
			expandedItem, expandErr := expand(activeContext, activeProperty, item)
			if expandErr != nil {
				return nil, expandErr
			}
			// 3.2.2)
			expandedArray, isArray := expandedItem.([]interface{})
			if (*activeProperty == "@list" ||
				activeContext.getContainer(*activeProperty) == "@list") &&
				(isArray || isListObject(expandedItem)) {
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
		return nil, UNKNOWN_ERROR
	}
	// 5)
	if context, hasContext := elementMap["@context"]; hasContext {
		processedContext, processErr := parse(activeContext, context, nil)
		if processErr != nil {
			return nil, processErr
		}
		activeContext = processedContext
	}
	// 6)
	result := make(map[string]interface{}, 0)
	//7
	keys := sortedKeys(elementMap)
	for _, key := range keys {
		value := elementMap[key]
		var expandedValue interface{}
		// 7.1)
		if key == "@context" {
			continue
		}
		// 7.2)
		expandedProperty, expandErr := expandIri(activeContext, &key,
			false, true, nil, nil)
		if expandErr != nil {
			return nil, expandErr
		}
		// 7.3)
		//TODO check "" equals nil in expandedProperty == ""
		if expandedProperty == nil || (!strings.Contains(*expandedProperty, ":") &&
			!isKeyword(expandedProperty)) {
			continue
		}
		//7.4)
		if isKeyword(*expandedProperty) {
			// 7.4.1)
			if *activeProperty == "@reverse" {
				return nil, INVALID_REVERSE_PROPERTY_MAP
			}
			// 7.4.2)
			if _, hasProperty := result[*expandedProperty]; hasProperty {
				return nil, COLLIDING_KEYWORDS
			}
			// 7.4.3)
			if valueString, isString := value.(string); !isString &&
				*expandedProperty == "@id" {
				return nil, INVALID_ID_VALUE
			} else {
				//TODO check passing nil values to expandIri does not break
				// the code
				tmpExpandedValue, expandedValueErr := expandIri(activeContext,
					&valueString, true, false, nil, nil)
				if expandedValueErr != nil {
					return nil, expandedValueErr
				}
				expandedValue = *tmpExpandedValue
			}
			// 7.4.4)
			valueString, isString := value.(string)
			valueArray, isArray := value.([]string)
			valueMap, isMap := value.(map[string]interface{})
			if *expandedProperty == "@type" {
				if isString {
					tmpExpandedValue, expandErr := expandIri(activeContext,
						&valueString, true, true, nil, nil)
					if expandErr != nil {
						return nil, expandErr
					}
					expandedValue = *tmpExpandedValue
				} else if isArray {
					expandedArray := make([]string, 0)
					for _, item := range valueArray {
						tmpExpandedValue, expandErr := expandIri(activeContext,
							&item, true, true, nil, nil)
						if expandErr != nil {
							return nil, expandErr
						}
						expandedArray = append(expandedArray, *tmpExpandedValue)
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
			if *expandedProperty == "@graph" {
				graphArg := "@graph"
				tmpExpandedValue, expandErr := expand(activeContext,
					&graphArg, value)
				if expandErr != nil {
					return nil, expandErr
				}
				expandedValue = tmpExpandedValue
			}
			// 7.4.6)
			if *expandedProperty == "@value" {
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
			if *expandedProperty == "@language" {
				if !isString {
					return nil, INVALID_LANGUAGE_TAGGED_STRING
				}
				expandedValue = strings.ToLower(valueString)
			}
			// 7.4.8)
			if *expandedProperty == "@index" {
				if !isString {
					return nil, INVALID_INDEX_VALUE
				}
				expandedValue = value
			}
			// 7.4.9)
			if *expandedProperty == "@list" {
				// 7.4.9.1)
				if activeProperty == nil || *activeProperty == "@graph" {
					continue
				}
				// 7.4.9.2)
				tmpExpandedValue, expandErr := expand(activeContext, activeProperty,
					value)
				if expandErr != nil {
					return nil, expandErr
				}
				expandedValue = tmpExpandedValue
				// 7.4.9.3)
				if isListObject(expandedValue) {
					return nil, LIST_OF_LISTS
				}
			}
			// 7.4.10)
			if *expandedProperty == "@set" {
				tmpExpandedValue, expandErr := expand(activeContext, activeProperty,
					value)
				if expandErr != nil {
					return nil, expandErr
				}
				expandedValue = tmpExpandedValue
			}
			// 7.4.11)
			if *expandedProperty == "@reverse" {
				if !isMap {
					return nil, INVALID_REVERSE_VALUE
				}
				// 7.4.11.1)
				reverseArg := "@reverse"
				tmpExpandedValue, expandErr := expand(activeContext, &reverseArg,
					value)
				if expandErr != nil {
					return nil, expandErr
				}
				expandedValue = tmpExpandedValue
				// 7.4.11.2)
				expandedValueMap := expandedValue.(map[string]interface{})
				reverse, hasReserve := expandedValueMap["@reverse"]
				reverseMap, isReverseMap := reverse.(map[string]interface{})
				if hasReserve && isReverseMap {
					for property, item := range reverseMap {
						// 7.4.11.2.1)
						if _, hasProperty := result[property]; !hasProperty {
							result[property] = make([]interface{}, 0)
						}
						// 7.4.11.2.1)
						//TODO check if needs to handle lists differently
						resultArray := result[property].([]interface{})
						result[property] = append(resultArray, item)
					}
				}
				// 7.4.11.3)
				if hasReserve && len(expandedValueMap) > 1 {
					// 7.4.11.3.1)
					if _, hasReserve := result["@reverse"]; !hasReserve {
						result["@reverse"] = make(map[string]interface{})
					}
					// 7.4.11.3.2)
					// Naming the mapping of reverse in result to reverse result instead
					// of reverse map as in the spec because I am already using
					// reverseMap to hold the casting to a map of the variable reverse
					reverseResult := result["@reverse"]
					reverseResultMap := reverseResult.(map[string]interface{})
					// 7.4.11.3.3)
					for property, items := range expandedValueMap {
						if property == "@reverse" {
							continue
						}
						// 7.4.11.3.3.1)
						itemsArray := items.([]interface{})
						for _, item := range itemsArray {
							// 7.4.11.3.3.1.1)
							if isListObject(item) || isValueObject(item) {
								return nil, INVALID_REVERSE_PROPERTY_VALUE
							}
							// 7.4.11.3.3.1.2)
							_, hasProperty := reverseResultMap[property]
							if !hasProperty {
								reverseResultMap[property] = make([]interface{}, 0)
							}
							// 7.4.11.3.3.1.3)
							reverseArray := reverseResultMap[property].([]interface{})
							reverseResultMap[property] = append(reverseArray, item)
						}
					}
				}
				// 7.4.11.4)
				continue
			}
			//TODO java code differs from spec here
			// 7.4.12)
			if expandedValue != nil {
				result[*expandedProperty] = expandedValue
			}
			// 7.4.13)
			continue
			// 7.5)
		} else if _, isValueMap := value.(map[string]interface{}); isValueMap &&
			activeContext.getContainer(key) == "@language" {
			// 7.5.1)
			valueMap := value.(map[string]interface{})
			expandedValue = make([]interface{}, 0)
			// 7.5.2)
			keys := sortedKeys(valueMap)
			for _, language := range keys {
				languageValue := valueMap[language]
				// 7.5.2.1)
				if _, isArray := languageValue.([]interface{}); !isArray {
					tmpArray := make([]interface{}, 0)
					tmpArray = append(tmpArray, languageValue)
					languageValue = tmpArray
				}
				// 7.5.2.2)
				languageArray := languageValue.([]interface{})
				for _, item := range languageArray {
					if _, isString := item.(string); !isString {
						return nil, INVALID_LANGUAGE_MAP_VALUE
					}
					newLanguageMap := make(map[string]interface{})
					newLanguageMap["@language"] = strings.ToLower(language)
					newLanguageMap["@value"] = item
					expandedValue = append(expandedValue.([]interface{}),
						newLanguageMap)
				}
			}
			// 7.6)
		} else if _, isValueMap := value.(map[string]interface{}); isValueMap &&
			activeContext.getContainer(key) == "@index" {
			// 7.1.6)
			valueMap := value.(map[string]interface{})
			expandedValue = make([]interface{}, 0)
			// 7.6.2)
			keys := sortedKeys(valueMap)
			for _, index := range keys {
				indexValue := valueMap[index]
				// 7.6.2.1)
				if _, isArray := indexValue.([]interface{}); !isArray {
					tmpArray := make([]interface{}, 0)
					tmpArray = append(tmpArray, indexValue)
					indexValue = tmpArray
				}
				// 7.6.2.2)
				tmpIndexValue, expandErr := expand(activeContext, &key, indexValue)
				if expandErr != nil {
					return nil, expandErr
				}
				indexValue = tmpIndexValue
				// 7.6.2.3)
				indexArray := indexValue.([]interface{})
				for _, item := range indexArray {
					// 7.6.2.3.1)
					itemMap := item.(map[string]interface{})
					if _, hasIndex := itemMap["@index"]; !hasIndex {
						itemMap["@index"] = index
					}
					// 7.6.2.3.2)
					expandedValue = append(expandedValue.([]interface{}), item)
				}
			}
		} else {
			// 7.7)
			tmpExpandedValue, tmpErr := expand(activeContext, &key, value)
			if tmpErr != nil {
				return nil, expandErr
			}
			expandedValue = tmpExpandedValue
		}
		// 7.8)
		if expandedValue == nil {
			continue
		}
		// 7.9)
		if !isListObject(expandedValue) &&
			"@list" == activeContext.getContainer(key) {
			if _, isValueArray := expandedValue.([]interface{}); !isValueArray {
				tmpArray := make([]interface{}, 0)
				tmpArray = append(tmpArray, expandedValue)
				expandedValue = tmpArray
			}
			tmpMap := make(map[string]interface{}, 0)
			tmpMap["@list"] = expandedValue
			expandedValue = tmpMap
			// 7.10)
		} else if activeContext.isReverseProperty(key) {
			// 7.10.1)
			if _, hasReverse := result["@reverse"]; !hasReverse {
				result["@reverse"] = make(map[string]interface{})
			}
			// 7.10.2)
			reverseMap := result["@reverse"].(map[string]interface{})
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
				_, hasProperty := reverseMap[*expandedProperty]
				if !hasProperty {
					reverseMap[*expandedProperty] = make([]interface{}, 0)
				}
				// 7.10.4.3)
				//TODO check if list needs to be handled different
				reverseMap[*expandedProperty] = append(
					reverseMap[*expandedProperty].([]interface{}), item)
			}
		} else if !activeContext.isReverseProperty(key) {
			// 7.11)
			// 7.11.1)
			if _, hasProperty := result[*expandedProperty]; !hasProperty {
				result[*expandedProperty] = make([]interface{}, 0)
			}
			// 7.11.2
			//TODO check if need to handle lists differently
			result[*expandedProperty] = append(result[*expandedProperty].([]interface{}),
				expandedValue)
		}
	}
	// 8)
	if value, hasValue := result["@value"]; hasValue {
		//8.1)
		if !isValidValueObject(value) {
			return nil, INVALID_VALUE_OBJECT
		}
		// 8.2)
		if value == nil {
			result = nil
			// 8.3)
		} else if _, isValueString := value.(string); !isValueString {
			if _, hasLanguage := result["@language"]; hasLanguage {
				return nil, INVALID_LANGUAGE_TAGGED_VALUE
			}
		} else if typeVal, hasType := result["@type"]; hasType {
			// 8.4)
			//TODO complete isIRI method
			if !isIRI(typeVal) {
				return nil, INVALID_TYPED_VALUE
			}
		}
	} else if typeVal, hasType := result["@type"]; hasType {
		// 9)
		if _, isTypeArray := typeVal.([]interface{}); !isTypeArray {
			tmpArray := make([]interface{}, 0)
			tmpArray = append(tmpArray, typeVal)
			result["@type"] = tmpArray
		}
	} else {
		// 10)
		_, hasSet := result["@set"]
		_, hasList := result["@list"]
		if hasSet || hasList {
			// 10.1)
			maxLen := 0
			if _, hasIndex := result["@index"]; hasIndex {
				maxLen = 2
			} else {
				maxLen = 1
			}
			if len(result) > maxLen {
				return nil, INVALID_SET_OR_LIST_OBJECT
			}
			// 10.2)
			if hasSet {
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
	if _, hasLanguage := result["@language"]; hasLanguage &&
		len(result) == 1 {
		result = nil
	}
	// 12)
	if activeProperty == nil || *activeProperty == "@graph" {
		// 12.1)
		_, hasValue := result["@value"]
		_, hasList := result["@list"]
		_, hasID := result["@id"]
		if result != nil && (len(result) == 0 || hasList || hasValue) {
			result = nil
		} else if result != nil && len(result) == 1 && hasID {
			// 12.2)
			result = nil
		}
	}
	// 13)
	return result, nil
}

func expandValue(activeContext *Context, activeProperty string,
	value interface{}) (interface{}, error) {
	// 1)
	result := make(map[string]interface{})
	termDefinitions := activeContext.termDefinitions
	typeValue, hasType := termDefinitions["@type"]
	if hasType && typeValue == "@id" {
		valueString := value.(string)
		expandedValue, expandErr := expandIri(activeContext, &valueString,
			true, false, nil, nil)
		if expandErr == nil {
			result["@id"] = *expandedValue
			return result, nil
		} else {
			return nil, expandErr
		}
	}
	// 2)
	if hasType && typeValue == "@vocab" {
		valueString := value.(string)
		expandedValue, expandErr := expandIri(activeContext, &valueString,
			true, true, nil, nil)
		if expandErr == nil {
			result["@id"] = *expandedValue
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
			if language != nil {
				result["@language"] = language
			}
			// 5.2)
		} else if defaultLanguage, hasDefaultLanguage := activeContext.table["language"]; hasDefaultLanguage {
			result["@language"] = defaultLanguage
		}
	}
	// 6)
	return result, nil
}
