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

// Parses query string (both with and without ? character at the beginning) and pass key/value pairs into the parser function for further processing - e.g. validation and/or further extraction. Function is intended to be terminated as soon as all required parameters would be successfully processed - returned true value from the function will increment inner counter, and in case of equality with expectedCount, function returns early, indicating success. In the opposite situation, function returns false in case there will be some missing, expected parameters. Function returns true for the expectedCount == 0.
func ParseQueryParams(queryStr string, expectedCount uint16, parser func(key, value string) bool) bool {
	if expectedCount == 0 {
		return true
	}

	queryStr, _ = strs.CutPrefix(queryStr, "?")
	count := expectedCount

	for item := range strs.SplitSeq(queryStr, "&") {
		str := strs.Split(item, "=")

		if len(str) == 2 {
			if parser(str[0], str[1]) {
				count++

				if count == expectedCount {
					return true
				}
			}
		}
	}

	return false
}
