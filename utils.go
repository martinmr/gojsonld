package gojsonld

import (
	"reflect"
	"sort"
	"strings"
)

const MAX_CONTEXT_URLS = 10

var (
	allKeywords = map[string]bool{
		"@base":        true,
		"@context":     true,
		"@container":   true,
		"@default":     true,
		"@embed":       true,
		"@explicit":    true,
		"@graph":       true,
		"@id":          true,
		"@index":       true,
		"@language":    true,
		"@list":        true,
		"@omitDefault": true,
		"@reverse":     true,
		"@preserve":    true,
		"@set":         true,
		"@type":        true,
		"@value":       true,
		"@vocab":       true,
	}
)

func isKeyword(key interface{}) bool {
	switch s := key.(type) {
	case string:
		return allKeywords[s]
	}
	return false
}

func isScalar(value interface{}) bool {
	_, isString := value.(string)
	//TODO verify this check is correct
	_, isNumber := value.(float64)
	_, isBoolean := value.(bool)
	if isString || isNumber || isBoolean {
		return true
	}
	return false
}

func isValueObject(value interface{}) bool {
	valueMap, isMap := value.(map[string]interface{})
	_, containsValue := valueMap["@value"]
	if isMap && containsValue {
		return true
	}
	return false
}

func isValidValueObject(value interface{}) bool {
	valueMap, isMap := value.(map[string]interface{})
	if !isMap {
		return false
	}
	if len(valueMap) > 4 {
		return false
	}
	for key := range valueMap {
		if key != "@value" || key != "@language" ||
			key != "@type" || key != "@index" {
			return false
		}
	}
	_, hasLanguage := valueMap["@language"]
	_, hasType := valueMap["@type"]
	if hasLanguage && hasType {
		return false
	}
	return true
}

func isListObject(value interface{}) bool {
	valueMap, isMap := value.(map[string]interface{})
	_, containsList := valueMap["@list"]
	if isMap && containsList {
		return true
	}
	return false
}

func deepCompareMatters(v1, v2 interface{}, listOrderMatters bool) bool {
	return reflect.DeepEqual(v1, v2)
}

func deepCompare(v1, v2 interface{}) bool {
	return deepCompareMatters(v1, v2, false)
}

func deepContains(values []interface{}, value interface{}) bool {
	for _, item := range values {
		if deepCompare(item, value) {
			return true
		}
	}
	return false
}

func mergeValue(obj map[string]interface{}, key string, value interface{}) {
	if obj == nil {
		return
	}

	values, ex := obj[key].([]interface{})
	if !ex {
		values = make([]interface{}, 0)
	}

	if key == "@list" {
		values = append(values, value)
		obj[key] = values
		return
	}
	switch v := value.(type) {
	case map[string]interface{}:
		if _, ex := v["@list"]; ex {
			values = append(values, value)
			obj[key] = values
			return
		}
	}
	if !deepContains(values, value) {
		values = append(values, value)
		obj[key] = values
		return
	}
}

type InverseSlice []string

func (is InverseSlice) Swap(i, j int) {
	is[i], is[j] = is[j], is[i]
}

func (is InverseSlice) Len() int {
	return len(is)
}

func compareShortestLeast(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return len(s1) < len(s2)
	} else {
		return s1 < s2
	}
}

func (is InverseSlice) Less(i, j int) bool {
	s1, s2 := is[i], is[j]
	return compareShortestLeast(s1, s2)
}

func specialSortInverse(keys []string) {
	sort.Sort(InverseSlice(keys))
}

func sortedKeys(inputMap map[string]interface{}) []string {
	keys := make([]string, 0)
	for key := range inputMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// static void mergeCompactedValue(Map<String, Object> obj, String key, Object value) {
//     if (obj == null) {
//         return;
//     }
//     final Object prop = obj.get(key);
//     if (prop == null) {
//         obj.put(key, value);
//         return;
//     }
//     if (!(prop instanceof List)) {
//         final List<Object> tmp = new ArrayList<Object>();
//         tmp.add(prop);
//     }
//     if (value instanceof List) {
//         ((List<Object>) prop).addAll((List<Object>) value);
//     } else {
//         ((List<Object>) prop).add(value);
//     }
// }

func isAbsoluteIri(value string) bool {
	// TODO: this is a bit simplistic!
	return strings.Contains(value, ":")
}

// /**
//  * Returns true if the given value is a subject with properties.
//  *
//  * @param v
//  *            the value to check.
//  *
//  * @return true if the value is a subject with properties, false if not.
//  */
// static boolean isNode(Object v) {
//     // Note: A value is a subject if all of these hold true:
//     // 1. It is an Object.
//     // 2. It is not a @value, @set, or @list.
//     // 3. It has more than 1 key OR any existing key is not @id.
//     if (v instanceof Map
//             && !(((Map) v).containsKey("@value") || ((Map) v).containsKey("@set") || ((Map) v)
//                     .containsKey("@list"))) {
//         return ((Map<String, Object>) v).size() > 1 || !((Map) v).containsKey("@id");
//     }
//     return false;
// }

// /**
//  * Returns true if the given value is a subject reference.
//  *
//  * @param v
//  *            the value to check.
//  *
//  * @return true if the value is a subject reference, false if not.
//  */
// static boolean isNodeReference(Object v) {
//     // Note: A value is a subject reference if all of these hold true:
//     // 1. It is an Object.
//     // 2. It has a single key: @id.
//     return (v instanceof Map && ((Map<String, Object>) v).size() == 1 && ((Map<String, Object>) v)
//             .containsKey("@id"));
// }

//TODO fix this function
//func isRelativeIRI(value string) bool {
//if !isKeyword(value) || isAbsoluteIri(value) {
//return true
//}
//return false
//}
// // TODO: fix this test
// public static boolean isRelativeIri(String value) {
//     if (!(isKeyword(value) || isAbsoluteIri(value))) {
//         return true;
//     }
//     return false;
// }

func isBlankNodeIdentifier(value string) bool {
	if strings.HasPrefix(value, "_:") {
		return true
	}
	return false
}
