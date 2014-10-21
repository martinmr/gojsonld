package gojsonld

import (
	"strconv"
	"strings"
)

func flatten(element interface{}, activeContext *Context) (interface{}, error) {
	// 1)
	nodeMap := make(map[string]interface{}, 0)
	nodeMap["@default"] = make(map[string]interface{}, 0)
	// 2)
	var idGenerator = BlankNodeIdGenerator{}
	idGenerator.counter = 0
	idGenerator.identifierMap = make(map[string]string, 0)
	defaultArg := "@default"
	err := generateNodeMap(element, nodeMap, &defaultArg, nil,
		nil, nil, &idGenerator)
	if err != nil {
		return nil, err
	}
	// 3)
	defaultGraph := nodeMap["@default"].(map[string]interface{})
	// 4)
	for graphName := range nodeMap {
		if graphName == "@default" {
			continue
		}
		graph := nodeMap[graphName].(map[string]interface{})
		// 4.1)
		if _, hasGraph := defaultGraph[graphName]; !hasGraph {
			tmpMap := make(map[string]interface{}, 0)
			tmpMap["@id"] = graphName
			defaultGraph[graphName] = graphName
		}
		// 4.2)
		entry := defaultGraph[graphName].(map[string]interface{})
		// 4.3)
		//TODO check spec (java comment states this step should only
		//be done if it does not already exists
		entry["@graph"] = make([]interface{}, 0)
		// 4.4)
		keys := sortedKeys(graph)
		for _, id := range keys {
			node := graph[id].(map[string]interface{})
			if _, hasId := node["@id"]; !(hasId && len(node) == 1) {
				entry["@graph"] = append(entry["@graph"].([]interface{}), node)
			}
		}
	}
	// 5)
	flattened := make([]interface{}, 0)
	// 6)
	keys := sortedKeys(defaultGraph)
	for _, id := range keys {
		node := defaultGraph[id].(map[string]interface{})
		if _, hasId := defaultGraph["@id"]; !(hasId && len(node) == 1) {
			flattened = append(flattened, node)
		}
	}
	// 7)
	if activeContext == nil {
		return flattened, nil
	}
	// 8)
	//TODO make sure "" is ok to send as default value
	compacted, compactErr := compact(activeContext, "", flattened,
		activeContext.options.CompactArrays)
	if compactErr != nil {
		return nil, compactErr
	}
	if _, isArray := compacted.([]interface{}); !isArray {
		tmpArray := make([]interface{}, 0)
		tmpArray = append(tmpArray, compacted)
		compacted = tmpArray
	}
	graphArg := "@graph"
	alias, compactErr := compactIri(activeContext, &graphArg, nil, false, false)
	if compactErr != nil {
		return nil, compactErr
	}
	returnValue, serializeErr := activeContext.serialize()
	if serializeErr != nil {
		return nil, serializeErr
	}
	returnValue[*alias] = compacted
	return returnValue, nil
}

func generateNodeMap(element interface{}, nodeMap map[string]interface{},
	activeGraph *string, activeSubject interface{}, activeProperty *string,
	list map[string]interface{}, idGenerator *BlankNodeIdGenerator) error {
	// 1)
	if elementArray, isArray := element.([]interface{}); isArray {
		// 1.1)
		for _, item := range elementArray {
			err := generateNodeMap(item, nodeMap, activeGraph, activeSubject,
				activeProperty, list, idGenerator)
			if err != nil {
				return err
			}
		}
	}
	// 2)
	elementMap := element.(map[string]interface{})
	// TODO check if needs to add map if activeGraph does not exist in nodeMap
	if _, hasGraph := nodeMap[*activeGraph]; !hasGraph {
		nodeMap[*activeGraph] = make(map[string]interface{}, 0)
	}
	graph := nodeMap[*activeGraph].(map[string]interface{})
	var node map[string]interface{}
	if activeSubject == nil {
		node = nil
	} else {
		node = graph[activeSubject.(string)].(map[string]interface{})
	}
	// 3)
	if _, hasType := elementMap["@type"]; hasType {
		var oldTypes []string
		newTypes := make([]string, 0)
		_, isArray := elementMap["@type"].([]string)
		if isArray {
			oldTypes = elementMap["@type"].([]string)
		} else {
			oldTypes = make([]string, 0)
			oldTypes = append(oldTypes, elementMap["@type"].(string))
		}
		for _, item := range oldTypes {
			if strings.HasPrefix(item, "_:") {
				newTypes = append(newTypes,
					idGenerator.generateBlankNodeIdentifier(&item))
			} else {
				newTypes = append(newTypes, item)
			}
		}
		if isArray {
			elementMap["@type"] = newTypes
		} else {
			elementMap["@type"] = newTypes[0]
		}
	}
	// 4)
	if _, hasValue := elementMap["@value"]; hasValue {
		// 4.1)
		if list == nil {
			mergeValue(node, *activeProperty, elementMap)
			// 4.2)
		} else {
			mergeValue(list, "@list", elementMap)
		}
		// 5)
	} else if _, hasList := elementMap["@list"]; hasList {
		// 5.1)
		result := make(map[string]interface{}, 0)
		result["@list"] = make([]interface{}, 0)
		// 5.2)
		generateNodeMap(elementMap["@list"], nodeMap, activeGraph, activeSubject,
			activeProperty, result, idGenerator)
		// 5.3)
		mergeValue(node, *activeProperty, result)
		// 6)
	} else {
		//6.1)
		var id string
		if elementID, hasId := elementMap["@id"]; hasId {
			if strings.HasPrefix(id, "_:") {
				elementIDString := elementID.(string)
				id = idGenerator.generateBlankNodeIdentifier(&elementIDString)
			} else {
				id = elementMap["@id"].(string)
			}
			delete(elementMap, "@id")
			// 6.2)
		} else {
			id = idGenerator.generateBlankNodeIdentifier(nil)
		}
		// 6.3)
		if _, hasId := graph[id]; !hasId {
			tmpMap := make(map[string]interface{}, 0)
			tmpMap["@id"] = id
			graph[id] = tmpMap
		}
		// 6.4)
		//TODO line asked by the spec but breaks various tests in java version
		// node := graph[id].(map[string]interface{})
		// 6.5)
		if _, isMap := activeSubject.(map[string]interface{}); isMap {
			mergeValue(graph["@id"].(map[string]interface{}), *activeProperty,
				activeSubject)
			// 6.6)
		} else if activeProperty != nil {
			// 6.6.1)
			reference := make(map[string]interface{}, 0)
			reference["@id"] = id
			// 6.6.2)
			if list == nil {
				mergeValue(node, *activeProperty, reference)
				// 6.6.3)
			} else {
				//TODO code differs from spec. For now following Java code
				mergeValue(list, "@list", reference)
			}
		}
		//TODO code differs from spec. see below
		// TODO: SPEC this is removed in the spec now, but it's still needed
		// (see 6.4)
		node = graph[id].(map[string]interface{})
		//6.7)
		if _, hasType := elementMap["@type"]; hasType {
			for _, typeVal := range elementMap["@type"].([]interface{}) {
				mergeValue(node, "@type", typeVal)
			}
			delete(elementMap, "@type")
		}
		// 6.8)
		if indexVal, hasIndex := elementMap["@index"]; hasIndex {
			if nodeIndexVal, hasNodeIndex := node["@index"]; hasNodeIndex {
				if !deepCompare(nodeIndexVal, indexVal) {
					return CONFLICTING_INDEXES
				}
			}
			//TODO java version differs. Check spec is correct
			node["@index"] = indexVal
			delete(elementMap, "@index")
		}
		// 6.9)
		if _, hasReverse := elementMap["@reverse"]; hasReverse {
			//6.9.1)
			referencedNode := make(map[string]interface{}, 0)
			referencedNode["@id"] = id
			// 6.9.2)
			reverseMap := elementMap["@reverse"].(map[string]interface{})
			// 6.9.3)
			for property := range reverseMap {
				values := reverseMap[property].([]interface{})
				// 6.9.3.1)
				for _, value := range values {
					// 6.9.3.1.1)
					err := generateNodeMap(value, nodeMap, activeGraph,
						referencedNode, &property, nil, idGenerator)
					if err != nil {
						return err
					}
				}
			}
			//6.9.4)
			delete(elementMap, "@reverse")
		}
		// 6.10)
		if graphVal, hasGraph := elementMap["@graph"]; hasGraph {
			err := generateNodeMap(graphVal, nodeMap, &id, nil, nil, nil,
				idGenerator)
			if err != nil {
				return err
			}
			delete(elementMap, "@graph")
		}
		// 6.11
		keys := sortedKeys(elementMap)
		for _, property := range keys {
			value := elementMap[property]
			// 6.11.1)
			if strings.HasPrefix(property, "_:") {
				property = idGenerator.
					generateBlankNodeIdentifier(&property)
			}
			// 6.11.2)
			if _, hasProperty := node[property]; !hasProperty {
				tmpArray := make([]interface{}, 0)
				node[property] = tmpArray
			}
			// 6.11.3)
			generateNodeMap(value, nodeMap, activeGraph, id, &property,
				nil, idGenerator)
		}
	}
	return nil
}

type BlankNodeIdGenerator struct {
	counter       int64
	identifierMap map[string]string
}

func (g *BlankNodeIdGenerator) generateBlankNodeIdentifier(identifier *string) string {
	// 1)
	if identifier != nil {
		if id, hasId := g.identifierMap[*identifier]; hasId {
			return id
		}
	}
	// 2)
	newId := "_:b" + strconv.FormatInt(g.counter, 10)
	// 3)
	g.counter += 1
	// 4)
	if identifier != nil {
		g.identifierMap[*identifier] = newId
	}
	// 5)
	return newId
}
