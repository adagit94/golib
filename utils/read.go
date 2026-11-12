package utils

import (
	"io/fs"
	"os"
	fp "path/filepath"
)

// Read all files under directory path as bytes and save them into the map with key corresponding to the path beginning from dirPath passed. In case recurse is true, files of all nested directories are included in flat fashion with complete path to destination file as the map key. In case error occurs during transition of directory, nil is returned for a map and vice versa.
func ReadDirFilesFlat(dirPath string, recurse bool) (map[string][]byte, error) {
	m := make(map[string][]byte)

	if err := fs.WalkDir(os.DirFS(dirPath), ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			if !recurse && path != "." {
				return fs.SkipDir
			}

			return nil
		}

		p := fp.Join(dirPath, path)

		if fileContent, err := os.ReadFile(p); err != nil {
			return err
		} else {
			m[p] = fileContent
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return m, nil
}
