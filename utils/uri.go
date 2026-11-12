package utils

import (
	strs "strings"
)

// Parses query string (without ? character at the beginning) and returns key/values pairs in form of map. Multiple values separated by , will be included in slice without any trimming, so it's best not use any whitespaces that shouldn't part of value itself.
func ParseQueryStr(queryStr string) map[string][]string {
	queryParams := make(map[string][]string)

	for item := range strs.SplitSeq(queryStr, "&") {
		str := strs.Split(item, "=")

		if len(str) == 2 {
			queryParams[str[0]] = strs.Split(str[1], ",")
		}
	}

	return queryParams
}
