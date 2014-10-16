package gojsonld

type Context struct {
	table map[string]interface{}

	options         *Options
	termDefinitions map[string]interface{}
	inverse         map[string]interface{}
}

// public Context() {
//     this(new JsonLdOptions());
// }
// public Context(JsonLdOptions opts) {
//     super();
//     init(opts);
// }
// public Context(Map<String, Object> map, JsonLdOptions opts) {
//     super(map);
//     init(opts);
// }
// public Context(Map<String, Object> map) {
//     super(map);
//     init(new JsonLdOptions());
// }

// public Context(Object context, JsonLdOptions opts) {
//     // TODO: load remote context
//     super(context instanceof Map ? (Map<String, Object>) context : null);
//     init(opts);
// }

func (c *Context) init(options *Options) {
	c.options = options
	if len(options.base) > 0 {
		c.table["@base"] = options.base
	}
	c.termDefinitions = make(map[string]interface{})
}

// /**
//  * Return a map of potential RDF prefixes based on the JSON-LD Term
//  * Definitions in this context.
//  * <p>
//  * No guarantees of the prefixes are given, beyond that it will not contain
//  * ":".
//  *
//  * @param onlyCommonPrefixes
//  *            If <code>true</code>, the result will not include
//  *            "not so useful" prefixes, such as "term1":
//  *            "http://example.com/term1", e.g. all IRIs will end with "/" or
//  *            "#". If <code>false</code>, all potential prefixes are
//  *            returned.
//  *
//  * @return A map from prefix string to IRI string
//  */
// public Map<String, String> getPrefixes(boolean onlyCommonPrefixes) {
//     final Map<String, String> prefixes = new LinkedHashMap<String, String>();
//     for (final String term : termDefinitions.keySet()) {
//         if (term.contains(":")) {
//             continue;
//         }
//         final Map<String, Object> termDefinition = (Map<String, Object>) termDefinitions
//                 .get(term);
//         if (termDefinition == null) {
//             continue;
//         }
//         final String id = (String) termDefinition.get("@id");
//         if (id == null) {
//             continue;
//         }
//         if (term.startsWith("@") || id.startsWith("@")) {
//             continue;
//         }
//         if (!onlyCommonPrefixes || id.endsWith("/") || id.endsWith("#")) {
//             prefixes.put(term, id);
//         }
//     }
//     return prefixes;
// }

// String compactIri(String iri, boolean relativeToVocab) {
//     return compactIri(iri, null, relativeToVocab, false);
// }

// String compactIri(String iri) {
//     return compactIri(iri, null, false, false);
// }

func (c *Context) clone() *Context {
	var clonedContext *Context = new(Context)
	return clonedContext
}

func (c *Context) getContainer(property string) string {
	if "@graph" == property {
		return "@set"
	}
	if isKeyword(property) {
		return property
	}
	td := c.termDefinitions[property]
	if tdMap, ok := td.(map[string]interface{}); ok {
		return tdMap["@container"].(string)
	} else {
		return ""
	}
}

func (c *Context) isReverseProperty(property string) bool {
	td, isMap := c.termDefinitions[property].(map[string]interface{})
	if td == nil || !isMap {
		return false
	}
	reverse := td["@reverse"]
	reverseBool, isBool := reverse.(bool)
	return isBool && reverseBool
}

func (c *Context) getTypeMapping(property string) (string, bool) {
	td := c.termDefinitions[property]
	if tdMap, ok := td.(map[string]interface{}); ok {
		typeMapping, okMapping := tdMap["@type"].(string)
		return typeMapping, okMapping
	} else {
		return "", false
	}
}

func (c *Context) getLanguageMapping(property string) (string, bool) {
	td := c.termDefinitions[property]
	if tdMap, ok := td.(map[string]interface{}); ok {
		languageMapping, okMapping := tdMap["@language"].(string)
		return languageMapping, okMapping
	} else {
		return "", false
	}
}

func (c *Context) getTermDefinition(key string) (map[string]interface{}, bool) {
	termDefinition, ok := c.termDefinitions[key]
	if !ok {
		return nil, false
	}
	termDefinitionMap, okMap := termDefinition.(map[string]interface{})
	return termDefinitionMap, okMap
}

// public Object expandValue(String activeProperty, Object value) throws JsonLdError {
//     final Map<String, Object> rval = new LinkedHashMap<String, Object>();
//     final Map<String, Object> td = getTermDefinition(activeProperty);
//     // 1)
//     if (td != null && "@id".equals(td.get("@type"))) {
//         // TODO: i'm pretty sure value should be a string if the @type is
//         // @id
//         rval.put("@id", expandIri(value.toString(), true, false, null, null));
//         return rval;
//     }
//     // 2)
//     if (td != null && "@vocab".equals(td.get("@type"))) {
//         // TODO: same as above
//         rval.put("@id", expandIri(value.toString(), true, true, null, null));
//         return rval;
//     }
//     // 3)
//     rval.put("@value", value);
//     // 4)
//     if (td != null && td.containsKey("@type")) {
//         rval.put("@type", td.get("@type"));
//     }
//     // 5)
//     else if (value instanceof String) {
//         // 5.1)
//         if (td != null && td.containsKey("@language")) {
//             final String lang = (String) td.get("@language");
//             if (lang != null) {
//                 rval.put("@language", lang);
//             }
//         }
//         // 5.2)
//         else if (this.get("@language") != null) {
//             rval.put("@language", this.get("@language"));
//         }
//     }
//     return rval;
// }

// public Object getContextValue(String activeProperty, String string) throws JsonLdError {
//     throw new JsonLdError(Error.NOT_IMPLEMENTED,
//             "getContextValue is only used by old code so far and thus isn't implemented");
// }

func (c *Context) serialize() (map[string]interface{}, error) {
	context := make(map[string]interface{})
	if base, hasBase := c.table["@base"]; hasBase && base != c.options.base {
		context["@base"] = base
	}
	if language, hasLanguage := c.table["@language"]; hasLanguage {
		context["@language"] = language
	}
	if vocab, hasVocab := c.table["@vocab"]; hasVocab {
		context["@vocab"] = vocab
	}
	for term := range c.termDefinitions {
		definition := c.termDefinitions[term].(map[string]interface{})
		_, hasType := definition["@type"]
		container, hasContainer := definition["@container"]
		language, hasLanguage := definition["@language"]
		reverse, hasReverse := definition["@reverse"]
		reverseBool, isBool := reverse.(bool)
		if !hasLanguage &&
			!hasContainer &&
			!hasType &&
			(!hasReverse || (isBool && reverseBool == false)) {
			id := definition["@id"].(string)
			compactID, compactErr := compactIri(c, &id, nil, false, false)
			if compactErr != nil {
				return nil, compactErr
			}
			if term == *compactID {
				context[term] = id
			} else {
				context[term] = *compactID
			}
		} else {
			defn := make(map[string]interface{}, 0)
			id := definition["@id"].(string)
			compactID, compactErr := compactIri(c, &id, nil, false, false)
			if compactErr != nil {
				return nil, compactErr
			}
			var reverseProperty bool
			if isBool && reverseBool == true {
				reverseProperty = true
			} else {
				reverseProperty = false
			}
			if !(term == *compactID && !reverseProperty) {
				if reverseProperty {
					defn["@reverse"] = *compactID
				} else {
					defn["@id"] = *compactID
				}
			}
			typeMapping, hasTypeMapping := definition["@type"].(string)
			if hasTypeMapping {
				if isKeyword(typeMapping) {
					defn["@type"] = typeMapping
				} else {
					compactType, compactErr := compactIri(c, &typeMapping,
						nil, true, false)
					if compactErr != nil {
						return nil, compactErr
					}
					defn["@type"] = *compactType
				}
			}
			if hasContainer {
				defn["@container"] = container
			}
			if hasLanguage {
				languageBool, isBool := language.(bool)
				if isBool && languageBool == false {
					defn["@language"] = nil
				} else {
					defn["@language"] = language
				}
			}
			context[term] = defn
		}
	}
	returnValue := make(map[string]interface{}, 0)
	if !(context == nil || len(context) == 0) {
		returnValue["@context"] = context
	}
	return returnValue, nil
}
