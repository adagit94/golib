package writers

import (
	u "github.com/adagit94/golib/utils"
	"os"
	"sync"
)

type pathSlot struct {
	mu     sync.Mutex
	action func()
	runs   bool
}

func (st *pathSlot) run() {
	st.mu.Lock()

	if action := st.action; action != nil {
		st.action = nil
		st.mu.Unlock()
		action()
		st.run()
	} else {
		st.runs = false
		st.mu.Unlock()
	}
}

type pathsSlots map[string]*pathSlot

type singleSlotFilesWriter[D BuffData] struct {
	mu         sync.Mutex
	pathsSlots pathsSlots
	saveFile   SaveFile[D, error]
}

func (st *singleSlotFilesWriter[D]) Set(path string, data D) {
	st.mu.Lock()

	pathSlotPtr, found := st.pathsSlots[path]

	if !found {
		pathSlotPtr = &pathSlot{}
		st.pathsSlots[path] = pathSlotPtr
	}

	st.mu.Unlock()
	pathSlotPtr.mu.Lock()

	pathSlotPtr.action = func() {
		st.saveFile(path, data)
	}

	if !pathSlotPtr.runs {
		pathSlotPtr.runs = true
		go pathSlotPtr.run()
	}

	pathSlotPtr.mu.Unlock()
}

type ISingleSlotFilesWriter[D BuffData] interface {
	// Sets a write task to the ready slot of passed path (directory structure get's created in case it doesn't exists). Data will be scheduled to be saved in discrete goroutine in case no write is in progress for a path, or only last attempt to save the data will be taken into account in case there is running save operation at the moment for a path or preceding task wasn't taken from the slot yet.
	Set(path string, data D)
}

// Single Slot File Writer (SSFW) is intended fo use cases when it's favourable to save only last scheduled data under specific path. It cannot be guaranteed that data being [ISingleSlotFilesWriter.Set] for a future save will be actually saved since there can be other, following calls for a same path in meanwhile, and as such, this abstraction is provided for use cases that doesn't require to wait for a save to complete and has such guarantee of data being actually saved.
func CreateSingleSlotFileWriter[D BuffData](saveFile SaveFile[D, error]) ISingleSlotFilesWriter[D] {
	writerPtr := &singleSlotFilesWriter[D]{pathsSlots: make(pathsSlots), saveFile: saveFile}
	return writerPtr
}

// Single Slot String File Writer (SSSFW) is intended fo use cases when it's favourable to save only last scheduled data under specific path. It cannot be guaranteed that data being [ISingleSlotFilesWriter.Set] for a future save will be actually saved since there can be other, following calls for a same path in meanwhile, and as such, this abstraction is provided for use cases that doesn't require to wait for a save to complete and has such guarantee of data being actually saved. Perms are set for both, directories, and destination file alike.
func CreateSingleSlotStringFileWriter(perms os.FileMode) ISingleSlotFilesWriter[string] {
	writerPtr := CreateSingleSlotFileWriter(func(path string, data string) error {
		return u.WriteFileString(path, perms, data)
	})

	return writerPtr
}

// Single Slot Bytes File Writer (SSBFW) is intended fo use cases when it's favourable to save only last scheduled data under specific path. It cannot be guaranteed that data being [ISingleSlotFilesWriter.Set] for a future save will be actually saved since there can be other, following calls for a same path in meanwhile, and as such, this abstraction is provided for use cases that doesn't require to wait for a save to complete and has such guarantee of data being actually saved. Perms are set for both, directories, and destination file alike.
func CreateSingleSlotBytesFileWriter(perms os.FileMode) ISingleSlotFilesWriter[[]byte] {
	writerPtr := CreateSingleSlotFileWriter(func(path string, data []byte) error {
		return u.WriteFileBytes(path, perms, data)
	})

	return writerPtr
}
