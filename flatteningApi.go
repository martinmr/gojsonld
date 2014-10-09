package gojsonld

import (
	"strings"
)

func flatten(element interface{}, context interface{}) interface{} {
	// 1)
	nodeMap := make(map[string]interface{}, 0)
	nodeMap["@default"] = make(map[string]interface{}, 0)
	// 2)
	//TODO make sure element is expanded
	//TODO check "" works instead of nil
	generateNodeMap(element, nodeMap, "@default", nil, "", nil)
	// 3)
	defaultGraph := nodeMap["@default"].(map[string]interface{})
	delete(nodeMap, "@default")
	// 4)
	for graphName := range nodeMap {
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
	if context == nil {
		return flattened
	}
	// 8)
	//TODO figure out how to pass options
	activeContext := Context{}
	//TODO check correct value of second argument to parse
	activeContext.parse(context, nil)
	compacted, _ := compact(&activeContext, "", flattened, false)
	//TODO java version is different than spec
	//TODO figure out what to do in step 8
	//for now returns compacted
	return compacted
}

func generateNodeMap(element interface{}, nodeMap map[string]interface{},
	activeGraph string, activeSubject interface{}, activeProperty string,
	list map[string]interface{}) error {
	// 1)
	if _, isArray := element.([]interface{}); isArray {
		// 1.1)
		err := generateNodeMap(element, nodeMap, activeGraph, activeSubject,
			activeProperty, list)
		if err != nil {
			return err
		}
	}
	// 2)
	elementMap := element.(map[string]interface{})
	// TODO check if needs to add map if activeGraph does not exist in nodeMap
	if _, hasGraph := nodeMap[activeGraph]; !hasGraph {
		nodeMap[activeGraph] = make(map[string]interface{}, 0)
	}
	graph := nodeMap[activeGraph].(map[string]interface{})
	var node map[string]interface{}
	//TODO handle activeSubject being a string
	if activeSubject == nil || activeSubject.(string) == "" {
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
				//TODO replace item by new blank node identifier
				newTypes = append(newTypes, item)
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
			//TODO merge
			if _, hasProperty := node[activeProperty]; !hasProperty {
				tmpArray := make([]interface{}, 0)
				tmpArray = append(tmpArray, element)
				node[activeProperty] = tmpArray
			}
			// 4.2)
		} else {
			//TODO merge
			list["@list"] = append(list["@list"].([]interface{}), element)
		}
		// 5)
	} else if _, hasList := elementMap["@list"]; hasList {
		// 5.1)
		result := make(map[string]interface{}, 0)
		result["@list"] = make([]interface{}, 0)
		// 5.2)
		generateNodeMap(elementMap["@list"], nodeMap, activeGraph, activeSubject,
			activeProperty, result)
		// 5.3)
		//TODO merge
		node[activeProperty] = append(node[activeGraph].([]interface{}), result)
		// 6)
	} else {
		//6.1)
		var id string
		if _, hasId := elementMap["@id"]; hasId {
			if strings.HasPrefix(id, "_:") {
				//TODO generate blank node identifier
				id = elementMap["@id"].(string)
			} else {
				id = elementMap["@id"].(string)
			}
			delete(elementMap, "@id")
			// 6.2)
		} else {
			//TODO generate identifier
			id = ""
		}
		// 6.3)
		if _, hasId := graph[id]; !hasId {
			tmpMap := make(map[string]interface{}, 0)
			tmpMap["@id"] = id
			graph[id] = tmpMap
		}
		// 6.4)
		//TODO line asked by the spec but breaks various tests in java version
		node := graph[id].(map[string]interface{})
		// 6.5)
		if _, isMap := activeSubject.(map[string]interface{}); isMap {
			// 6.5.1)
			if _, hasProperty := node[activeProperty]; !hasProperty {
				tmpArray := make([]interface{}, 0)
				tmpArray = append(tmpArray, activeSubject)
				node[activeProperty] = tmpArray
			} else {
				//TODO merge
			}
			// 6.6)
		} else if activeProperty != "" {
			// 6.6.1)
			reference := make(map[string]interface{}, 0)
			reference["@id"] = id
			// 6.6.2)
			if list == nil {
				// 6.6.2.1)
				if _, hasProperty := node[activeProperty]; !hasProperty {
					tmpArray := make([]interface{}, 0)
					tmpArray = append(tmpArray, reference)
					node[activeProperty] = tmpArray
					// 6.6.2.2)
				} else {
					//TODO merge
				}
				// 6.6.3)
			} else {
				//TODO merge
				//TODO code differs from spec
			}
		}
		//TODO code differs from spec. see below
		// TODO: SPEC this is removed in the spec now, but it's still needed
		// (see 6.4)
		node = graph[id].(map[string]interface{})
		//6.7)
		if _, hasType := elementMap["@type"]; hasType {
			for _, typeVal := range elementMap["@type"].([]interface{}) {
				//TODO merge
			}
		}
		// 6.8)
		if indexVal, hasIndex := elementMap["@index"]; hasIndex {
			//TODO deep compare for now just set the key-value pair
			// without checking
			node["@index"] = indexVal
			//remove if deep compare returns false
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
						referencedNode, property, nil)
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
			generateNodeMap(graphVal, nodeMap, id, nil, "", nil)
			delete(elementMap, "@graph")
		}
		// 6.11
		keys := sortedKeys(elementMap)
		for _, property := range keys {
			value := elementMap[property]
			// 6.11.1)
			if strings.HasPrefix(property, "_:") {
				//TODO set property to a blank node identifier
			}
			// 6.11.2)
			if _, hasProperty := node[property]; !hasProperty {
				tmpArray := make([]interface{}, 0)
				node[property] = tmpArray
			}
			// 6.11.3)
			generateNodeMap(value, nodeMap, activeGraph, id, property, nil)
		}
	}
	return nil
}

//Blank node identifier
//TODO figure out how to hold state of the blank node identifier
