package gojsonld

import (
	"strings"
)

func (c *Context) parse(localContext interface{},
	remoteContexts []string) (*Context, error) {
	if remoteContexts == nil {
		remoteContexts = make([]string, 0)
	}
	// 1)
	result := c.clone()
	// 2)
	if _, ok := localContext.([]interface{}); !ok {
		var temp interface{} = localContext
		localContext = make([]interface{}, 0)
		localContext = append(localContext.([]interface{}), temp)
	}
	// 3)
	for _, context := range localContext.([]interface{}) {
		// 3.1)
		//TODO The base IRI of the active context is set
		//to the IRI of the currently being processed
		//document (which might be different from the currently
		//being processed context), if available; otherwise
		//to null. If set, the base option of a JSON-LD API
		//Implementation overrides the base IRI.
		if context == nil {
			newContext := new(Context)
			newContext.init(c.options)
			result = newContext
			continue
		}
		// 3.2)
		if _, isString := context.(string); isString {
			// 3.2.1)
			uri := result.table["@base"].(string)
			//TODO resolve uri
			// 3.2.2
			isRecursive := false
			for _, remoteContext := range remoteContexts {
				if remoteContext == uri {
					isRecursive = true
				}
			}
			if isRecursive {
				return nil, RECURSIVE_CONTEXT_INCLUSION
			}
			remoteContexts = append(remoteContexts, context.(string))
			// 3.2.3
			rd := c.options.documentLoader.loadDocument(uri)
			var remoteContext interface{} = rd.document
			//TODO is this check correct. Maybe needs to use the reflect package
			remoteContextMap, isMap := remoteContext.(map[string]interface{})
			_, containsContext := remoteContextMap["@context"]
			if !isMap {
				return nil, LOADING_REMOTE_CONTEXT_FAILED
			} else if !containsContext {
				return nil, INVALID_REMOTE_CONTEXT
			}
			context = remoteContextMap["@context"]
			// 3.2.4)
			recursiveResult, parseErr := result.parse(context, remoteContexts)
			if parseErr == nil {
				result = recursiveResult
			}
			//TODO handle case when parseErr is not nil
			// 3.2.5)
			continue
		}
		// 3.3)
		//TODO check JSON objects are already represented as maps
		contextMap, isMap := context.(map[string]interface{})
		if !isMap {
			return nil, INVALID_LOCAL_CONTEXT
		}
		// 3.4)
		// 3.4.1)
		baseValue, containsBase := contextMap["@base"]
		if len(remoteContexts) == 0 && containsBase {
			// 3.4.2)
			if baseValue == nil {
				delete(result.table, "@base")
			} else if stringValue, isString := baseValue.(string); isString {
				if isAbsoluteIri(stringValue) {
					// 3.4.3)
					result.table["@base"] = stringValue
				} else {
					//TODO check steps 3.4.4 and 3.4.5 are correctly implemented
					// 3.4.4)
					baseURI, isStringURI := result.table["@base"].(string)
					if !isStringURI || !isAbsoluteIri(baseURI) {
						//3.4.5)
						return nil, INVALID_BASE_IRI
					}
					//TODO
					//result.put("@base", JsonLdUrl.resolve(baseUri, (String) baseValue));
				}
			} else {
				return nil, INVALID_BASE_IRI
			}
		}
		// 3.5)
		// 3.5.1)
		vocabValue, containsVocab := contextMap["@vocab"]
		if containsVocab {
			// 3.5.2)
			if vocabValue == nil {
				delete(result.table, "@vocab")
			} else if stringValue, isString := vocabValue.(string); isString {
				// 3.5.3)
				if isAbsoluteIri(stringValue) || isBlankNodeIdentifier(stringValue) {
					result.table["@vocab"] = stringValue
				} else {
					return nil, INVALID_VOCAB_MAPPING
				}
			} else {
				return nil, INVALID_VOCAB_MAPPING
			}
		}
		// 3.6)
		// 3.6.1)
		languageValue, containsLanguage := contextMap["@language"]
		if containsLanguage {
			if languageValue == nil {
				delete(result.table, "@language")
			} else if stringValue, isString := languageValue.(string); isString {
				result.table["@language"] = strings.ToLower(stringValue)
			} else {
				return nil, INVALID_DEFAULT_LANGUAGE
			}
		}
		// 3.7
		defined := make(map[string]bool, 0)
		for key := range contextMap {
			if key == "@base" || key == "@vocab" || key == "@language" {
				continue
			}
			result.createTermDefinition(contextMap, key, defined)
		}
	}
	return result, nil
}

func (c *Context) createTermDefinition(localContext map[string]interface{},
	term string, defined map[string]bool) error {
	// 1)
	if definedValue, isDefined := defined[term]; isDefined {
		if definedValue {
			return nil
		}
		return CYCLIC_IRI_MAPPING
	}
	// 2)
	defined[term] = false
	// 3)
	if isKeyword(term) {
		return KEYWORD_REDEFINITION
	}
	// 4)
	delete(c.termDefinitions, term)
	// 5)
	var value interface{} = localContext[term]
	// 6)
	valueMap, isMap := value.(map[string]interface{})
	if value == nil || (isMap && valueMap["@id"] == nil) {
		c.termDefinitions[term] = nil
		defined[term] = true
		return nil
	}
	// 7)
	if _, isString := value.(string); isString {
		tempMap := make(map[string]interface{}, 0)
		tempMap["@id"] = value
		value = tempMap
	}
	// 8)
	if !isMap {
		return INVALID_TERM_DEFINITION
	}
	//9)
	definition := make(map[string]interface{}, 0)
	// 10)
	if typeVal, containsType := valueMap["@type"]; containsType {
		// 10.1)
		typeString, isString := typeVal.(string)
		if !isString {
			return INVALID_TYPE_MAPPING
		}
		// 10.2)
		//TODO handle typeErr
		expandedType, typeErr := c.expandIri(typeString, false,
			true, localContext, defined)
		if typeErr == nil {
			typeString = expandedType
		}
		//TODO handle error returned by expandIRI
		if typeString != "@id" || typeString != "@vocab" || !isAbsoluteIri(typeString) {
			return INVALID_TYPE_MAPPING
		}
		// 10.3)
		definition["@type"] = typeString
	}
	// 11)
	if reverse, containsReverse := valueMap["@reverse"]; containsReverse {
		// 11.1)
		if _, containsID := valueMap["@id"]; containsID {
			return INVALID_REVERSE_PROPERTY
		}
		reverseString, isString := reverse.(string)
		// 11.2)
		if !isString {
			return INVALID_IRI_MAPPING
		}
		// 11.3)
		expandedReverse, reverseErr := c.expandIri(reverseString, false,
			true, localContext, defined)
		if reverseErr == nil {
			reverseString = expandedReverse
		}
		if !isAbsoluteIri(reverseString) || !isBlankNodeIdentifier(reverseString) {
			return INVALID_IRI_MAPPING
		}
		definition["@id"] = reverseString
		// 11.4)
		if container, containsContainer := valueMap["@container"]; containsContainer {
			if container != "@set" || container != "@index" || container != nil {
				return INVALID_REVERSE_PROPERTY
			}
			definition["@container"] = container
		}
		definition["@reverse"] = true
		c.termDefinitions[term] = definition
		defined[term] = true
		return nil
	}
	//12)
	definition["@reverse"] = false
	// 13)
	id, containsID := valueMap["@id"]
	if containsID && id == term {
		idString, isString := id.(string)
		// 13.1)
		if !isString {
			return INVALID_IRI_MAPPING
		}
		// 13.2)
		expandedID, idErr := c.expandIri(idString, false, true, localContext, defined)
		if idErr == nil {
			idString = expandedID
		}
		if isKeyword(idString) || isBlankNodeIdentifier(idString) ||
			isAbsoluteIri(idString) {
			if idString == "@context" {
				return INVALID_KEYWORD_ALIAS
			}
			definition["@id"] = idString
		} else {
			return INVALID_IRI_MAPPING
		}
	} else if strings.Contains(term, ":") {
		// 14)
		colIndex := strings.Index(term, ":")
		prefix := term[:colIndex]
		suffix := term[colIndex+1:]
		// 14.1)
		if _, containsPrefix := localContext[prefix]; containsPrefix {
			c.createTermDefinition(localContext, prefix, defined)
		}
		// 14.2)
		prefixVal, containsPrefix := c.termDefinitions[prefix]
		prefixMap, _ := prefixVal.(map[string]interface{})
		if containsPrefix {
			definition["@id"] = prefixMap["@id"].(string) + suffix
		} else {
			//14/3
			definition["@id"] = term
		}
	} else if vocab, containsVocab := c.table["@vocab"]; containsVocab {
		// 15)
		definition["@id"] = vocab.(string) + term
	} else {
		return INVALID_IRI_MAPPING
	}
	// 16)
	// 16.1)
	if container, containsContainer := valueMap["@container"]; containsContainer {
		// 16.2)
		if container != "@list" || container != "@set" || container != "@index" ||
			container != "@language" {
			return INVALID_CONTAINER_MAPPING
		}
		// 16.3)
		definition["@container"] = container
	}
	// 17)
	language, containsLanguage := valueMap["@language"]
	_, containsType := valueMap["@type"]
	if containsLanguage && !containsType {
		// 17.1)
		languageString, isString := language.(string)
		if !isString || language != nil {
			return INVALID_LANGUAGE_MAPPING
		}
		// 17.2)
		if isString {
			definition["@language"] = strings.ToLower(languageString)
		} else {
			definition["@language"] = nil
		}
	}
	//18
	c.termDefinitions[term] = definition
	defined[term] = true
	return nil
}

func (c *Context) expandIri(value string, relative bool, vocab bool,
	localContext map[string]interface{}, defined map[string]bool) (string, error) {
	//1)
	//Using "" as nil value for strings (Go doesn't support strings of nil value)
	if isKeyword(value) || value == "" {
		return value, nil
	}
	//2)
	//TODO figure out what to do when value not in defined
	//for now we take the if branch if value not in defined
	//same thing as in step 4.3
	_, containsValue := localContext[value]
	definedValue, inDefined := defined[value]
	if containsValue && (!inDefined || definedValue == false) {
		createErr := c.createTermDefinition(localContext, value, defined)
		if createErr != nil {
			//TODO handle error
		}
	}
	// 3)
	if td, hasTermDefinition := c.termDefinitions[value]; vocab && hasTermDefinition {
		tdMap, isMap := td.(map[string]interface{})
		idString, isString := tdMap["@id"].(string)
		if isMap && isString {
			return idString, nil
		}
		//TODO check returning correct error
		return "", INVALID_TERM_DEFINITION
	}
	// 4)
	if colIndex := strings.Index(value, ":"); colIndex >= 0 {
		// 4.1)
		prefix := value[:colIndex]
		suffix := value[colIndex+1:]
		// 4.2)
		if prefix == "_" || strings.HasPrefix(suffix, "//") {
			return value, nil
		}
		// 4.3)
		_, containsPrefix := localContext[prefix]
		definedPrefix, inDefined := defined[prefix]
		if containsPrefix && (!inDefined || definedPrefix == false) {
			createErr := c.createTermDefinition(localContext, prefix, defined)
			if createErr != nil {
				//TODO handle error
			}
		}
		// 4.4)
		if td, hasTermDefinition := c.termDefinitions[prefix]; hasTermDefinition {
			tdMap, isMap := td.(map[string]interface{})
			id, isString := tdMap["@id"].(string)
			if isMap && isString {
				return id + suffix, nil
			}
			//TODO check error handling
			return "", INVALID_TERM_DEFINITION
		}
		// 4.5
		return value, nil
	}
	// 5)
	if vocabMapping, containsVocab := c.table["@vocab"]; vocab && containsVocab {
		vocabString, isString := vocabMapping.(string)
		if isString {
			return vocabString + value, nil
		}
		//TODO check error handling
		return "", INVALID_VOCAB_MAPPING
	} else if relative {
		// 6)
		//TODO set value to the result of resolving value agains base IRI
	}
	// 7)
	return value, nil
}
