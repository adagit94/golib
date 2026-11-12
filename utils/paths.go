package utils

import (
	fp "path/filepath"
)

// Function removes final segment from the path in case it includes file extension. Otherwise, original path is returned.
func RemoveFileFromPath(path string) string {
	ext := fp.Ext(path)

	if ext != "" {
		path = fp.Dir(path)
	}

	return path
}
