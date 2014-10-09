package gojsonld

import (
	"sort"
	"strings"
)

func compact(activeContext *Context, activeProperty string,
	element interface{}, compactArrays bool) (interface{}, error) {
	// TODO check how to handle inverse
	// 1)
	if isScalar(element) {
		return element, nil
	}
	// 2)
	if elementArray, isArray := element.([]interface{}); isArray {
		// 2.1)
		result := make([]interface{}, 0)
		// 2.2)
		for _, item := range elementArray {
			// 2.2.1)
			compactedItem, compactErr := compact(activeContext,
				activeProperty, item, compactArrays)
			// 2.2.2)
			if compactErr == nil && compactedItem != nil {
				result = append(result, compactedItem)
			}
		}
		// 2.3)
		if compactArrays && len(result) == 1 &&
			activeContext.getContainer(activeProperty) != "" {
			return result[0], nil
		}
		// 2.4)
		return result, nil
	}
	// 3)
	elementMap, isMap := element.(map[string]interface{})
	if !isMap {
		//TODO handle error
		return nil, UNKNOWN_ERROR
	}
	// 4)
	_, hasID := elementMap["@id"]
	_, hasValue := elementMap["@value"]
	if hasID || hasValue {
		compactedValue := compactValue(activeContext, activeProperty,
			elementMap)
		if isScalar(compactedValue) {
			return compactedValue, nil
		}
	}
	// 5)
	insideReverse := "@reverse" == activeProperty
	// 6)
	result := make(map[string]interface{})
	// 7)
	keys := make([]string, 0)
	for key := range elementMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, expandedProperty := range keys {
		expandedValue := elementMap[expandedProperty]
		var compactedValue interface{}
		// 7.1)
		if expandedProperty == "@id" || expandedProperty == "@type" {
			// 7.1.1)
			if valueString, isString := expandedValue.(string); isString {
				//TODO handle default value of call to compactIri
				compactedValue = compactIri(activeContext, valueString, nil,
					expandedProperty == "@type", false)
			} else {
				// 7.1.2)
				// 7.1.2.1)
				compactedValue = make([]string, 0)
				expandedArray, isArray := expandedValue.([]string)
				if !isArray {
					//TODO handle error
					return nil, nil
				}
				// 7.1.2.2)
				for _, expandedType := range expandedArray {
					//TODO handle default reverse value
					tmpCompact := compactIri(activeContext, expandedType,
						nil, true, false)
					compactedValue = append(compactedValue.([]string),
						tmpCompact)
				}
				// 7.1.2.3)
				if len(compactedValue.([]string)) == 1 {
					compactedValue = compactedValue.([]string)[0]
				}
			}
			// 7.1.3)
			//TODO handle defaults
			alias := compactIri(activeContext, expandedProperty,
				nil, true, false)
			// 7.1.4)
			result[alias] = compactedValue
			continue
		}
		// 7.2)
		if expandedProperty == "@reverse" {
			var compactedValue interface{}
			//TODO handle default
			// 7.2.1)
			tmpCompacted, compactErr := compact(activeContext, "@reverse",
				expandedValue, compactArrays)
			if compactErr == nil {
				compactedValue = tmpCompacted
			} else {
				//TODO handle error
				return nil, UNKNOWN_ERROR
			}
			// 7.2.2)
			compactedMap, isCompactedMap := compactedValue.(map[string]interface{})
			if !isCompactedMap {
				//TODO handle error
				return nil, UNKNOWN_ERROR
			}
			keys := make([]string, 0)
			for key := range compactedMap {
				keys = append(keys, key)
			}
			for _, property := range keys {
				value := compactedMap[property]
				// 7.2.2.1)
				if activeContext.isReverseProperty(property) {
					// 7.2.2.1.1)
					if _, isValueArray := value.([]interface{}); isValueArray &&
						(activeContext.getContainer(property) == "@set" ||
							!compactArrays) {
						tmpArray := make([]interface{}, 0)
						tmpArray = append(tmpArray, value)
						//TODO check logic (java version is different)
						value = tmpArray
					}
					// 7.2.2.1.2)
					if _, hasProperty := result[property]; !hasProperty {
						result[property] = value
					} else {
						// 7.2.2.1.3)
						_, isResultArray := result[property].([]interface{})
						if !isResultArray {
							tmpArray := make([]interface{}, 0)
							tmpArray = append(tmpArray, result[property])
							result[property] = tmpArray
						}
						resultArray := result[property].([]interface{})
						valueArray, isValueArray := value.([]interface{})
						if isValueArray {
							for _, item := range valueArray {
								resultArray = append(resultArray, item)
							}
						} else {
							resultArray = append(resultArray, value)
						}
					}
					// 7.2.2.1.4)
					delete(compactedMap, property)
				}
			}
			// 7.2.3)
			if len(compactedMap) > 0 {
				// 7.2.3.1)
				alias := compactIri(activeContext, "@reverse",
					nil, true, false)
				// 7.2.3.2)
				result[alias] = compactedValue
			}
			// 7.2.4)
			continue
		}
		// 7.3)
		if "@index" == expandedProperty &&
			activeContext.getContainer(activeProperty) == "@index" {
			continue
		} else if "@index" == activeProperty || "@value" == activeProperty ||
			"@language" == activeProperty {
			// 7.4)
			// 7.4.1)
			//TODO handle defaults
			alias := compactIri(activeContext, expandedProperty,
				nil, true, false)
			// 7.4.2)
			result[alias] = expandedValue
			continue
		}
		// 7.5)
		expandedArray, isExpandedArray := expandedValue.([]interface{})
		if isExpandedArray && len(expandedArray) == 0 {
			// 7.5.1)
			itemActiveProperty := compactIri(activeContext, expandedProperty,
				expandedValue, true, insideReverse)
			// 7.5.2)
			activeValue, hasActiveProperty := result[itemActiveProperty]
			if !hasActiveProperty {
				result[itemActiveProperty] = make([]interface{}, 0)
			} else {
				_, isActiveArray := activeValue.([]interface{})
				if !isActiveArray {
					tmpArray := make([]interface{}, 0)
					tmpArray = append(tmpArray, result[itemActiveProperty])
					activeValue = tmpArray
				}
			}
		}
		// 7.6)
		//TODO handle err
		if !isExpandedArray {
			return nil, UNKNOWN_ERROR
		}
		for _, expandedItem := range expandedArray {
			// 7.6.1)
			itemActiveProperty := compactIri(activeContext, expandedProperty,
				expandedItem, true, insideReverse)
			// 7.6.2)
			container := activeContext.getContainer(itemActiveProperty)
			// 7.6.3)
			var compactedItem interface{}
			listValue, hasList := expandedItem.(map[string]interface{})["@list"]
			var passedElement interface{}
			if hasList {
				passedElement = listValue
			} else {
				passedElement = expandedItem
			}
			tmpCompact, compactErr := compact(activeContext,
				itemActiveProperty, passedElement, compactArrays)
			if compactErr == nil {
				compactedItem = tmpCompact
			} else {
				//TODO handle error
				return nil, compactErr
			}
			// 7.6.4)
			if isListObject(expandedItem) {
				// 7.6.4.1)
				_, isCompactedArray := compactedItem.([]interface{})
				if !isCompactedArray {
					tmpArray := make([]interface{}, 0)
					tmpArray = append(tmpArray, compactedItem)
					compactedItem = tmpArray
				}
				// 7.6.4.2)
				if container != "@list" {
					// 7.6.4.2.1)
					//TODO check default values. Java version is differrent
					// vocab = true
					listKey := compactIri(activeContext, "@list", nil, false, false)
					tmpMap := make(map[string]interface{}, 0)
					tmpMap[listKey] = compactedItem
					compactedItem = tmpMap
					// 7.6.4.2.2)
					index, hasIndex := expandedItem.(map[string]interface{})["@index"]
					if hasIndex {
						//TODO check defaults
						indexKey := compactIri(activeContext, "@index", nil,
							false, false)
						compactedItem.(map[string]interface{})[indexKey] = index
					}
				} else if _, hasProperty := result[itemActiveProperty]; hasProperty {
					// 7.6.4.3
					return nil, COMPACTION_TO_LIST_OF_LISTS
				}
			}
			// 7.6.5
			if "@index" == container || "@language" == container {
				// 7.6.5.1)
				_, hasActiveProperty := result[itemActiveProperty]
				if !hasActiveProperty {
					tmpMap := make(map[string]interface{}, 0)
					result[itemActiveProperty] = tmpMap
				}
				mapObject, isMap := result[itemActiveProperty].(map[string]interface{})
				if !isMap {
					//TODO handle error
					return nil, UNKNOWN_ERROR
				}
				// 7.6.5.2)
				compactedMap, isMap := compactedItem.(map[string]interface{})
				compactedValueKey, hasValueKey := compactedMap["@value"]
				if !isMap {
					//TODO handle error
					return nil, UNKNOWN_ERROR
				}
				if container == "@language" && hasValueKey {
					compactedItem = compactedValueKey
				}
				// 7.6.5.3)
				expandedItemMap, isMap := expandedItem.(map[string]interface{})
				mapKeyInterface, hasMapKey := expandedItemMap[container]
				mapKey, isString := mapKeyInterface.(string)
				if !isMap || !isString || !hasMapKey {
					//TODO handle error
					return nil, UNKNOWN_ERROR
				}
				// 7.6.5.4)
				if _, hasMapObjectKey := mapObject[mapKey]; !hasMapObjectKey {
					mapObject[mapKey] = compactedItem
				} else {
					_, isArray := mapObject[mapKey].([]interface{})
					if !isArray {
						tmpArray := make([]interface{}, 0)
						//TODO check logic. Java version appears to have bug
						//but make sure of that
						tmpArray = append(tmpArray, mapObject[mapKey])
						mapObject[mapKey] = tmpArray
					}
					mapObjectArray := mapObject[mapKey].([]interface{})
					mapObjectArray = append(mapObjectArray, compactedItem)
				}
			} else {
				// 7.6.6)
				// 7.6.6.1)
				//TODO check logic
				_, isCompactedArray := compactedItem.([]interface{})
				if (!compactArrays || "@set" == container || "@list" == container ||
					"@graph" == expandedProperty || "@list" == expandedProperty) &&
					!isCompactedArray {
					tmpArray := make([]interface{}, 0)
					tmpArray = append(tmpArray, compactedItem)
					compactedItem = tmpArray
				}
				// 7.6.6.2)
				activePropertyValue, hasActiveProperty := result[itemActiveProperty]
				if !hasActiveProperty {
					result[itemActiveProperty] = compactedItem
				} else {
					// 7.6.6.3)
					if _, isArray := activePropertyValue.([]interface{}); !isArray {
						tmpArray := make([]interface{}, 0)
						tmpArray = append(tmpArray, activePropertyValue)
						result[itemActiveProperty] = tmpArray
					}
					activePropertyArray := result[itemActiveProperty].([]interface{})
					compactedArray, isArray := compactedItem.([]interface{})
					if isArray {
						for _, item := range compactedArray {
							activePropertyArray = append(activePropertyArray, item)
						}
					} else {
						activePropertyArray = append(activePropertyArray, compactedItem)
					}
				}
			}
		}
	}
	// 8)
	return result, nil
}

func (activeContext *Context) getInverse() map[string]interface{} {
	if activeContext.inverse != nil {
		return activeContext.inverse
	}
	// 1)
	result := make(map[string]interface{})
	// 2)
	defaultLanguage := "@none"
	defaultContextLanguage, hasDefaultLanguage := activeContext.table["@language"]
	if hasDefaultLanguage {
		defaultLanguage = defaultContextLanguage.(string)
	}
	// 3)
	keys := make([]string, 0)
	for key := range activeContext.termDefinitions {
		keys = append(keys, key)
	}
	specialSortInverse(keys)
	for _, term := range keys {
		// 3.1
		definition := activeContext.termDefinitions[term].(map[string]interface{})
		if definition == nil {
			continue
		}
		// 3.2)
		container := "@none"
		if containerVal, hasContainer := definition["@container"]; hasContainer {
			container = containerVal.(string)
		}
		// 3,3
		iri := definition["@id"].(string)
		// 3.4)
		if _, hasIri := result[iri]; !hasIri {
			result[iri] = make(map[string]interface{}, 0)
		}
		// 3.5)
		containerMap := result[iri].(map[string]interface{})
		// 3.6)
		if _, hasContainer := containerMap[container]; !hasContainer {
			tmpMap := make(map[string]interface{}, 0)
			tmpMap["@language"] = make(map[string]interface{}, 0)
			tmpMap["@type"] = make(map[string]interface{}, 0)
		}
		// 3.7)
		typeLanguageMap := containerMap[container].(map[string]interface{})
		// 3.8)
		//TODO check equality is correct
		if definition["@reverse"] == true {
			// 3.8.1)
			typeMap := typeLanguageMap["@type"].(map[string]interface{})
			// 3.8.2)
			if _, hasReverse := typeMap["@reverse"]; !hasReverse {
				typeMap["@reverse"] = term
			}
			// 3.9)
		} else if _, hasType := definition["@type"]; hasType {
			// 3.9.1)
			typeMap := typeLanguageMap["@type"].(map[string]interface{})
			// 3.9.2)
			typeMapping := definition["@type"].(string)
			if _, hasTypeMapping := typeMap[typeMapping]; !hasTypeMapping {
				typeMap[typeMapping] = term
			}
			// 3.10)
		} else if _, hasLanguage := definition["@language"]; hasLanguage {
			// 3.10.1)
			languageMap := typeLanguageMap["@language"].(map[string]interface{})
			// 3.10.2)
			language, langErr := definition["@language"].(string)
			//TODO handle nil string
			if language == "" || !langErr {
				language = "@null"
			}
			// 3.10.3)
			if _, hasLanguage := languageMap[language]; !hasLanguage {
				languageMap[language] = term
			}
			// 3.11
		} else {
			// 3.11.1)
			languageMap := typeLanguageMap["@language"].(map[string]interface{})
			// 3.11.2)
			//TODO check why java version is different
			if _, hasLanguage := languageMap[defaultLanguage]; !hasLanguage {
				languageMap[defaultLanguage] = term
			}
			// 3.11.3)
			if _, hasNone := languageMap["@none"]; hasNone {
				languageMap["@none"] = term
			}
			// 3.11.4)
			typeMap := typeLanguageMap["@type"].(map[string]interface{})
			// 3.11.5)
			if _, hasNone := typeMap["@none"]; !hasNone {
				typeMap["@none"] = term
			}
		}
	}
	//4)
	//TODO check if we need to save result in activeContext.inverse
	activeContext.inverse = result
	return result
}

func compactIri(activeContext *Context, iri string,
	value interface{}, vocab bool, reverse bool) string {
	// 1)
	if iri == "" {
		return ""
	}
	// 2)
	//TODO set activeContext.inverse
	activeContext.getInverse()
	if _, hasIri := activeContext.inverse[iri]; hasIri && vocab == true {
		// 2.1)
		defaultLanguage := "@none"
		if language, hasLanguage := activeContext.table["@language"]; hasLanguage {
			defaultLanguage = language.(string)
		}
		// 2.2)
		containers := make([]string, 0)
		// 2.3)
		typeLanguage := "@language"
		typeLanguageValue := "@none"
		// 2.4)
		_, hasIndex := value.(map[string]interface{})["@index"]
		if hasIndex {
			containers = append(containers, "@index")
		}
		// 2.5)
		if reverse {
			typeLanguage = "@type"
			typeLanguageValue = "@reverse"
			containers = append(containers, "@set")
			// 2.6)
		} else if isListObject(value) {
			// 2.6.1)
			if !hasIndex {
				containers = append(containers, "@list")
			}
			// 2.6.2)
			list := value.(map[string]interface{})["@list"].([]interface{})
			// 2.6.3)
			commonType := ""
			commonLanguage := ""
			if len(list) == 0 {
				commonLanguage = defaultLanguage
			}
			// 2.6.4)
			for _, item := range list {
				// 2.6.4.1)
				itemLanguage := "@none"
				itemType := "@none"
				// 2.6.4.2)
				itemMap, _ := item.(map[string]interface{})
				if isValueObject(item) {
					language, hasLanguage := itemMap["@language"]
					// 2.6.4.2.1)
					if hasLanguage {
						itemLanguage = language.(string)
						// 2.6.4.2.2)
					} else if typeVal, hasType := itemMap["@type"]; hasType {
						itemType = typeVal.(string)
						// 2.6.4.2.3)
					} else {
						itemLanguage = "@null"
					}
					// 2.6.4.3)
				} else {
					itemType = "@id"
				}
				// 2.4.6.4)
				if commonLanguage == "" {
					commonLanguage = itemLanguage
					// 2.6.4.5)
				} else if itemLanguage != commonLanguage && isValueObject(item) {
					commonLanguage = "@none"
				}
				// 2.6.4.6)
				if commonType == "" {
					commonType = itemType
					// 2.6.4.7)
				} else if commonType != itemType {
					commonType = "@none"
				}
				// 2.6.4.8)
				if commonLanguage == "@none" && commonLanguage == "@none" {
					break
				}
			}
			// 2.6.5)
			if commonLanguage == "" {
				commonLanguage = "@none"
			}
			// 2.6.6)
			if commonType == "" {
				commonType = "@none"
			}
			// 2.6.7)
			if commonType != "@none" {
				typeLanguage = "@type"
				typeLanguageValue = commonType
				// 2.6.8)
			} else {
				typeLanguageValue = commonLanguage
			}
			// 2.7)
		} else {
			// 2.7.1)
			if isValueObject(value) {
				// 2.7.1.1)
				valueMap := value.(map[string]interface{})
				language, hasLanguage := valueMap["@language"]
				_, hasIndex := valueMap["@index"]
				typeVal, hasType := valueMap["@type"]
				if hasLanguage && !hasIndex {
					typeLanguageValue = language.(string)
					containers = append(containers, "@language")
					// 2.7.1.2)
				} else if hasType {
					typeLanguageValue = typeVal.(string)
					typeLanguage = "@type"
				}
				// 2.7.2)
			} else {
				typeLanguage = "@type"
				typeLanguageValue = "@id"
			}
			// 2.7.3)
			containers = append(containers, "@set")
		}
		// 2.8)
		containers = append(containers, "@none")
		// 2.9)
		if typeLanguageValue == "" {
			typeLanguageValue = "@null"
		}
		// 2.10)
		preferredValues := make([]string, 0)
		// 2.11)
		if typeLanguageValue == "@reverse" {
			preferredValues = append(preferredValues, "@reverse")
		}
		// 2.12)
		id, hasID := value.(map[string]interface{})["@id"]
		if (typeLanguageValue == "@id" || typeLanguageValue == "@reverse") &&
			hasID {
			// 2.12.1)
			result := compactIri(activeContext, id.(string), nil, true, true)
			tdResult, hasDefinition := activeContext.termDefinitions[result]
			resultIri, hasIri := tdResult.(map[string]interface{})["@id"]
			valueIri := value.(map[string]interface{})["@id"]
			if hasDefinition && hasIri && valueIri.(string) == resultIri.(string) {
				preferredValues = append(preferredValues, "@vocab")
				preferredValues = append(preferredValues, "@id")
				preferredValues = append(preferredValues, "@none")
				// 2.12.2)
			} else {
				preferredValues = append(preferredValues, "@id")
				preferredValues = append(preferredValues, "@vocab")
				preferredValues = append(preferredValues, "@none")
			}
			// 2.13)
		} else {
			preferredValues = append(preferredValues, typeLanguageValue)
			preferredValues = append(preferredValues, "@none")
		}
		// 2.14)
		term := selectTerm(activeContext, iri, containers, typeLanguage,
			preferredValues)
		// 2.15)
		if term != "" {
			return term
		}
	}
	// 3)
	vocabContext, hasVocab := activeContext.table["@vocab"]
	if vocab && hasVocab {
		// 3.1
		if strings.Index(iri, vocabContext.(string)) == 0 && vocabContext != iri {
			suffix := iri[len(vocabContext.(string)):]
			if _, hasSuffix := activeContext.termDefinitions[suffix]; !hasSuffix {
				return suffix
			}
		}
	}
	// 4)
	compactIri := ""
	// 5)
	for term := range activeContext.termDefinitions {
		termDefinition := activeContext.termDefinitions[term].(map[string]interface{})
		// 5.1)
		if strings.Contains(term, ":") {
			continue
		}
		// 5.2)
		if termDefinition == nil || iri == termDefinition["@id"] ||
			!strings.HasPrefix(iri, termDefinition["@id"].(string)) {
			continue
		}
		// 5.3)
		candidate := term + ":" + iri[len(termDefinition["@id"].(string)):]
		// 5.4)
		//TODO check logic
		//TODO check what to do when checking for value == nil and value is a string
		candidateIsShorter := compareShortestLeast(candidate, compactIri)
		tdCandidate, hasCandidate := activeContext.termDefinitions[candidate]
		candidateIri, _ := (tdCandidate.(map[string]interface{}))["@id"].(string)
		if (compactIri == "" || candidateIsShorter) && (!hasCandidate ||
			(iri == candidateIri && value == nil)) {
			compactIri = candidate
		}
	}
	// 6)
	if compactIri != "" {
		return compactIri
	}
	// 7)
	if !vocab {
		//TODO remove base from iri
	}
	//8
	return iri
}

func selectTerm(activeContext *Context, iri string, containers []string,
	typeLanguage string, preferredValues []string) string {
	inverse := activeContext.getInverse()
	// 1)
	containerMap := inverse[iri].(map[string]interface{})
	// 2)
	for _, container := range containers {
		// 2.1)
		if _, hasContainer := containerMap[container]; !hasContainer {
			continue
		}
		// 2.2)
		typeLanguageMap := containerMap[container].(map[string]interface{})
		// 2.3)
		valueMap := typeLanguageMap[typeLanguage].(map[string]interface{})
		// 2.4)
		for _, item := range preferredValues {
			// 2.4.1)
			if _, hasItem := valueMap[item]; !hasItem {
				continue
				// 2.4.2)
			} else {
				return valueMap[item].(string)
			}
		}
	}
	// 3)
	return ""
}

//Value Compaction Algorithm
//http://json-ld.org/spec/latest/json-ld-api/#value-compaction
func compactValue(activeContext *Context, activeProperty string,
	value map[string]interface{}) interface{} {
	// 1)
	numberMembers := len(value)
	// 2)
	if _, ok := value["@index"]; ok &&
		"@index" == activeContext.getContainer(activeProperty) {
		numberMembers--
	}
	// 3)
	if numberMembers > 2 {
		return value
	}
	// 4)
	typeMapping, _ := activeContext.getTypeMapping(activeProperty)
	languageMapping, _ := activeContext.getLanguageMapping(activeProperty)
	if _, ok := value["@id"]; ok {
		// 4.1)
		if numberMembers == 1 && "@id" == typeMapping {
			//TODO compact IRI
			return ""
		}
		// 4.2)
		if numberMembers == 1 && "@vocab" == typeMapping {
			//TODO compact IRI
			return ""
		}
		return value
	}
	// 5)
	var valueValue interface{} = value["@value"]
	if typeKey, ok := value["@type"]; ok && typeKey.(string) == typeMapping {
		return valueValue
	}
	// 6)
	language, _ := value["@language"]
	languageString, validString := language.(string)
	if validString &&
		(languageString == languageMapping ||
			languageString == activeContext.table["@language"]) {
		return valueValue
	}
	// 7)
	_, isValueValueString := valueValue.(string)
	_, hasDefaultLanguage := activeContext.table["@language"]
	termDefinition, _ := activeContext.getTermDefinition(activeProperty)
	activePropertyLanguage, containsLanguage := termDefinition["@language"]
	if numberMembers == 1 && (!isValueValueString || !hasDefaultLanguage ||
		(containsLanguage && activePropertyLanguage == "")) {
		//TODO check that checking for "" is correct
		return valueValue
	}
	// 8)
	return value
}