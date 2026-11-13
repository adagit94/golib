package utils

import (
	strs "strings"
)

// Parses query string (both with and without ? character at the beginning) and returns key/values pairs in form of a map. Multiple values separated by , will be included in slice without any trimming, so it's best not use any whitespaces that shouldn't be part of value itself.
func ParseQueryStr(queryStr string) map[string][]string {
	queryStr, _ = strs.CutPrefix(queryStr, "?")
	queryParams := make(map[string][]string)

	for item := range strs.SplitSeq(queryStr, "&") {
		str := strs.Split(item, "=")

		if len(str) == 2 {
			queryParams[str[0]] = ParseQueryParamValues(str[1])
		}
	}

	return queryParams
}

// Parse values separated by , character - they  will be included in slice without any trimming, so it's best not use any whitespaces that shouldn't be part of value itself.
func ParseQueryParamValues(values string) []string {
	return strs.Split(values, ",")
}

// Parses query string (both with and without ? character at the beginning) and pass key/value pairs into parser function for further processing - e.g. validation and/or further extraction.
func ParseQueryParams(queryStr string, parser func(key, value string)) {
	queryStr, _ = strs.CutPrefix(queryStr, "?")
	
	for item := range strs.SplitSeq(queryStr, "&") {
		str := strs.Split(item, "=")

		if len(str) == 2 {
			parser(str[0], str[1])
		}
	}
}
