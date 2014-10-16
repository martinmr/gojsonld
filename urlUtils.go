package gojsonld

import (
	"net/url"
	"path/filepath"
)

func resolve(base, ref string) (string, error) {
	baseUrl, baseErr := url.Parse(base)
	refUrl, refErr := url.Parse(ref)
	if baseErr != nil {
		return "", baseErr
	}
	if refErr != nil {
		return "", refErr
	}
	resolvedUrl := baseUrl.ResolveReference(refUrl)
	return resolvedUrl.String(), nil
}

func removeBase(base, ref string) (string, error) {
	baseUrl, baseErr := url.Parse(base)
	refUrl, refErr := url.Parse(ref)
	if baseErr != nil {
		return "", baseErr
	}
	if refErr != nil {
		return "", refErr
	}
	if baseUrl.Host != refUrl.Host || baseUrl.Scheme != refUrl.Scheme {
		//TODO handle error
		return "", UNKNOWN_ERROR
	}
	rel, err := filepath.Rel(baseUrl.Path, refUrl.Path)
	if err != nil {
		return "", err
	}
	refUrl.Host = ""
	refUrl.Scheme = ""
	refUrl.Path = rel
	return refUrl.String(), nil
}
