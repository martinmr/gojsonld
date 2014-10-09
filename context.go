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

// public Context parse(Object localContext) throws JsonLdError {
//     return this.parse(localContext, new ArrayList<String>());
// }

// /**
//  * IRI Expansion Algorithm
//  *
//  * http://json-ld.org/spec/latest/json-ld-api/#iri-expansion
//  *
//  * @param value
//  * @param relative
//  * @param vocab
//  * @param context
//  * @param defined
//  * @return
//  * @throws JsonLdError
//  */
// String expandIri(String value, boolean relative, boolean vocab, Map<String, Object> context,
//         Map<String, Boolean> defined) throws JsonLdError {
//     // 1)
//     if (value == null || JsonLdUtils.isKeyword(value)) {
//         return value;
//     }
//     // 2)
//     if (context != null && context.containsKey(value)
//             && !Boolean.TRUE.equals(defined.get(value))) {
//         this.createTermDefinition(context, value, defined);
//     }
//     // 3)
//     if (vocab && this.termDefinitions.containsKey(value)) {
//         final Map<String, Object> td = (LinkedHashMap<String, Object>) this.termDefinitions
//                 .get(value);
//         if (td != null) {
//             return (String) td.get("@id");
//         } else {
//             return null;
//         }
//     }
//     // 4)
//     final int colIndex = value.indexOf(":");
//     if (colIndex >= 0) {
//         // 4.1)
//         final String prefix = value.substring(0, colIndex);
//         final String suffix = value.substring(colIndex + 1);
//         // 4.2)
//         if ("_".equals(prefix) || suffix.startsWith("//")) {
//             return value;
//         }
//         // 4.3)
//         if (context != null && context.containsKey(prefix)
//                 && (!defined.containsKey(prefix) || defined.get(prefix) == false)) {
//             this.createTermDefinition(context, prefix, defined);
//         }
//         // 4.4)
//         if (this.termDefinitions.containsKey(prefix)) {
//             return (String) ((LinkedHashMap<String, Object>) this.termDefinitions.get(prefix))
//                     .get("@id") + suffix;
//         }
//         // 4.5)
//         return value;
//     }
//     // 5)
//     if (vocab && this.containsKey("@vocab")) {
//         return this.get("@vocab") + value;
//     }
//     // 6)
//     else if (relative) {
//         return JsonLdUrl.resolve((String) this.get("@base"), value);
//     } else if (context != null && JsonLdUtils.isRelativeIri(value)) {
//         throw new JsonLdError(Error.INVALID_IRI_MAPPING, "not an absolute IRI: " + value);
//     }
//     // 7)
//     return value;
// }

// /**
//  * IRI Compaction Algorithm
//  *
//  * http://json-ld.org/spec/latest/json-ld-api/#iri-compaction
//  *
//  * Compacts an IRI or keyword into a term or prefix if it can be. If the IRI
//  * has an associated value it may be passed.
//  *
//  * @param iri
//  *            the IRI to compact.
//  * @param value
//  *            the value to check or null.
//  * @param relativeTo
//  *            options for how to compact IRIs: vocab: true to split after
//  * @vocab, false not to.
//  * @param reverse
//  *            true if a reverse property is being compacted, false if not.
//  *
//  * @return the compacted term, prefix, keyword alias, or the original IRI.
//  */
// String compactIri(String iri, Object value, boolean relativeToVocab, boolean reverse) {
//     // 1)
//     if (iri == null) {
//         return null;
//     }

//     // 2)
//     if (relativeToVocab && getInverse().containsKey(iri)) {
//         // 2.1)
//         String defaultLanguage = (String) this.get("@language");
//         if (defaultLanguage == null) {
//             defaultLanguage = "@none";
//         }

//         // 2.2)
//         final List<String> containers = new ArrayList<String>();
//         // 2.3)
//         String typeLanguage = "@language";
//         String typeLanguageValue = "@null";

//         // 2.4)
//         if (value instanceof Map && ((Map<String, Object>) value).containsKey("@index")) {
//             containers.add("@index");
//         }

//         // 2.5)
//         if (reverse) {
//             typeLanguage = "@type";
//             typeLanguageValue = "@reverse";
//             containers.add("@set");
//         }
//         // 2.6)
//         else if (value instanceof Map && ((Map<String, Object>) value).containsKey("@list")) {
//             // 2.6.1)
//             if (!((Map<String, Object>) value).containsKey("@index")) {
//                 containers.add("@list");
//             }
//             // 2.6.2)
//             final List<Object> list = (List<Object>) ((Map<String, Object>) value).get("@list");
//             // 2.6.3)
//             String commonLanguage = (list.size() == 0) ? defaultLanguage : null;
//             String commonType = null;
//             // 2.6.4)
//             for (final Object item : list) {
//                 // 2.6.4.1)
//                 String itemLanguage = "@none";
//                 String itemType = "@none";
//                 // 2.6.4.2)
//                 if (JsonLdUtils.isValue(item)) {
//                     // 2.6.4.2.1)
//                     if (((Map<String, Object>) item).containsKey("@language")) {
//                         itemLanguage = (String) ((Map<String, Object>) item).get("@language");
//                     }
//                     // 2.6.4.2.2)
//                     else if (((Map<String, Object>) item).containsKey("@type")) {
//                         itemType = (String) ((Map<String, Object>) item).get("@type");
//                     }
//                     // 2.6.4.2.3)
//                     else {
//                         itemLanguage = "@null";
//                     }
//                 }
//                 // 2.6.4.3)
//                 else {
//                     itemType = "@id";
//                 }
//                 // 2.6.4.4)
//                 if (commonLanguage == null) {
//                     commonLanguage = itemLanguage;
//                 }
//                 // 2.6.4.5)
//                 else if (!commonLanguage.equals(itemLanguage) && JsonLdUtils.isValue(item)) {
//                     commonLanguage = "@none";
//                 }
//                 // 2.6.4.6)
//                 if (commonType == null) {
//                     commonType = itemType;
//                 }
//                 // 2.6.4.7)
//                 else if (!commonType.equals(itemType)) {
//                     commonType = "@none";
//                 }
//                 // 2.6.4.8)
//                 if ("@none".equals(commonLanguage) && "@none".equals(commonType)) {
//                     break;
//                 }
//             }
//             // 2.6.5)
//             commonLanguage = (commonLanguage != null) ? commonLanguage : "@none";
//             // 2.6.6)
//             commonType = (commonType != null) ? commonType : "@none";
//             // 2.6.7)
//             if (!"@none".equals(commonType)) {
//                 typeLanguage = "@type";
//                 typeLanguageValue = commonType;
//             }
//             // 2.6.8)
//             else {
//                 typeLanguageValue = commonLanguage;
//             }
//         }
//         // 2.7)
//         else {
//             // 2.7.1)
//             if (value instanceof Map && ((Map<String, Object>) value).containsKey("@value")) {
//                 // 2.7.1.1)
//                 if (((Map<String, Object>) value).containsKey("@language")
//                         && !((Map<String, Object>) value).containsKey("@index")) {
//                     containers.add("@language");
//                     typeLanguageValue = (String) ((Map<String, Object>) value).get("@language");
//                 }
//                 // 2.7.1.2)
//                 else if (((Map<String, Object>) value).containsKey("@type")) {
//                     typeLanguage = "@type";
//                     typeLanguageValue = (String) ((Map<String, Object>) value).get("@type");
//                 }
//             }
//             // 2.7.2)
//             else {
//                 typeLanguage = "@type";
//                 typeLanguageValue = "@id";
//             }
//             // 2.7.3)
//             containers.add("@set");
//         }

//         // 2.8)
//         containers.add("@none");
//         // 2.9)
//         if (typeLanguageValue == null) {
//             typeLanguageValue = "@null";
//         }
//         // 2.10)
//         final List<String> preferredValues = new ArrayList<String>();
//         // 2.11)
//         if ("@reverse".equals(typeLanguageValue)) {
//             preferredValues.add("@reverse");
//         }
//         // 2.12)
//         if (("@reverse".equals(typeLanguageValue) || "@id".equals(typeLanguageValue))
//                 && (value instanceof Map) && ((Map<String, Object>) value).containsKey("@id")) {
//             // 2.12.1)
//             final String result = this.compactIri(
//                     (String) ((Map<String, Object>) value).get("@id"), null, true, true);
//             if (termDefinitions.containsKey(result)
//                     && ((Map<String, Object>) termDefinitions.get(result)).containsKey("@id")
//                     && ((Map<String, Object>) value).get("@id").equals(
//                             ((Map<String, Object>) termDefinitions.get(result)).get("@id"))) {
//                 preferredValues.add("@vocab");
//                 preferredValues.add("@id");
//             }
//             // 2.12.2)
//             else {
//                 preferredValues.add("@id");
//                 preferredValues.add("@vocab");
//             }
//         }
//         // 2.13)
//         else {
//             preferredValues.add(typeLanguageValue);
//         }
//         preferredValues.add("@none");

//         // 2.14)
//         final String term = selectTerm(iri, containers, typeLanguage, preferredValues);
//         // 2.15)
//         if (term != null) {
//             return term;
//         }
//     }

//     // 3)
//     if (relativeToVocab && this.containsKey("@vocab")) {
//         // determine if vocab is a prefix of the iri
//         final String vocab = (String) this.get("@vocab");
//         // 3.1)
//         if (iri.indexOf(vocab) == 0 && !iri.equals(vocab)) {
//             // use suffix as relative iri if it is not a term in the
//             // active context
//             final String suffix = iri.substring(vocab.length());
//             if (!termDefinitions.containsKey(suffix)) {
//                 return suffix;
//             }
//         }
//     }
// 4)
//     String compactIRI = null;
//     // 5)
//     for (final String term : termDefinitions.keySet()) {
//         final Map<String, Object> termDefinition = (Map<String, Object>) termDefinitions
//                 .get(term);
//         // 5.1)
//         if (term.contains(":")) {
//             continue;
//         }
//         // 5.2)
//         if (termDefinition == null || iri.equals(termDefinition.get("@id"))
//                 || !iri.startsWith((String) termDefinition.get("@id"))) {
//             continue;
//         }

//         // 5.3)
//         final String candidate = term + ":"
//                 + iri.substring(((String) termDefinition.get("@id")).length());
//         // 5.4)
//         if ((compactIRI == null || compareShortestLeast(candidate, compactIRI) < 0)
//                 && (!termDefinitions.containsKey(candidate) || (iri
//                         .equals(((Map<String, Object>) termDefinitions.get(candidate))
//                                 .get("@id")) && value == null))) {
//             compactIRI = candidate;
//         }

//     }

//     // 6)
//     if (compactIRI != null) {
//         return compactIRI;
//     }

//     // 7)
//     if (!relativeToVocab) {
//         return JsonLdUrl.removeBase(this.get("@base"), iri);
//     }

//     // 8)
//     return iri;
// }

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

func (c *Context) createInverse() {
	return
}

// /**
//  * Term Selection
//  *
//  * http://json-ld.org/spec/latest/json-ld-api/#term-selection
//  *
//  * This algorithm, invoked via the IRI Compaction algorithm, makes use of an
//  * active context's inverse context to find the term that is best used to
//  * compact an IRI. Other information about a value associated with the IRI
//  * is given, including which container mappings and which type mapping or
//  * language mapping would be best used to express the value.
//  *
//  * @return the selected term.
//  */
// private String selectTerm(String iri, List<String> containers, String typeLanguage,
//         List<String> preferredValues) {
//     final Map<String, Object> inv = getInverse();
//     // 1)
//     final Map<String, Object> containerMap = (Map<String, Object>) inv.get(iri);
//     // 2)
//     for (final String container : containers) {
//         // 2.1)
//         if (!containerMap.containsKey(container)) {
//             continue;
//         }
//         // 2.2)
//         final Map<String, Object> typeLanguageMap = (Map<String, Object>) containerMap
//                 .get(container);
//         // 2.3)
//         final Map<String, Object> valueMap = (Map<String, Object>) typeLanguageMap
//                 .get(typeLanguage);
//         // 2.4 )
//         for (final String item : preferredValues) {
//             // 2.4.1
//             if (!valueMap.containsKey(item)) {
//                 continue;
//             }
//             // 2.4.2
//             return (String) valueMap.get(item);
//         }
//     }
//     // 3)
//     return null;
// }

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
	return reverse != nil && isBool && reverseBool
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

// public Map<String, Object> serialize() {
//     final Map<String, Object> ctx = new LinkedHashMap<String, Object>();
//     if (this.get("@base") != null && !this.get("@base").equals(options.getBase())) {
//         ctx.put("@base", this.get("@base"));
//     }
//     if (this.get("@language") != null) {
//         ctx.put("@language", this.get("@language"));
//     }
//     if (this.get("@vocab") != null) {
//         ctx.put("@vocab", this.get("@vocab"));
//     }
//     for (final String term : termDefinitions.keySet()) {
//         final Map<String, Object> definition = (Map<String, Object>) termDefinitions.get(term);
//         if (definition.get("@language") == null
//                 && definition.get("@container") == null
//                 && definition.get("@type") == null
//                 && (definition.get("@reverse") == null || Boolean.FALSE.equals(definition
//                         .get("@reverse")))) {
//             final String cid = this.compactIri((String) definition.get("@id"));
//             ctx.put(term, term.equals(cid) ? definition.get("@id") : cid);
//         } else {
//             final Map<String, Object> defn = new LinkedHashMap<String, Object>();
//             final String cid = this.compactIri((String) definition.get("@id"));
//             final Boolean reverseProperty = Boolean.TRUE.equals(definition.get("@reverse"));
//             if (!(term.equals(cid) && !reverseProperty)) {
//                 defn.put(reverseProperty ? "@reverse" : "@id", cid);
//             }
//             final String typeMapping = (String) definition.get("@type");
//             if (typeMapping != null) {
//                 defn.put("@type", JsonLdUtils.isKeyword(typeMapping) ? typeMapping
//                         : compactIri(typeMapping, true));
//             }
//             if (definition.get("@container") != null) {
//                 defn.put("@container", definition.get("@container"));
//             }
//             final Object lang = definition.get("@language");
//             if (definition.get("@language") != null) {
//                 defn.put("@language", Boolean.FALSE.equals(lang) ? null : lang);
//             }
//             ctx.put(term, defn);
//         }
//     }

//     final Map<String, Object> rval = new LinkedHashMap<String, Object>();
//     if (!(ctx == null || ctx.isEmpty())) {
//         rval.put("@context", ctx);
//     }
//     return rval;
// }
