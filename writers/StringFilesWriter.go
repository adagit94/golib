package writers

import (
	u "github.com/adagit94/golib/utils"
	"os"
)

// Create FilesWriter with PathsQueues initialized to specified buffSize on first write attempt, with Closed channel initialized to be buffered with size of 1, so it's optional to use it as part of users code, and fs.WriteFileString as SaveFile function to perform actual file write operation.
func CreateStringFilesWriter(buffSize BuffSize, filePerms os.FileMode) *FilesWriter[string, error] {
	fwPtr := CreateFilesWriter(buffSize, func(path string, data string) error {
		return u.WriteFileString(path, filePerms, data)
	})

	return fwPtr
}
