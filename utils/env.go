package utils

import (
	"os"
	"strings"
)

// EnvMap maps all env. variables returned by os.Environ() to map data structure.
func EnvMap() map[string]string {
	args := os.Environ()
	vars := make(map[string]string, len(args))

	for _, arg := range args {
		s := strings.SplitN(arg, "=", 2)
		vars[s[0]] = s[1]
	}

	return vars
}
