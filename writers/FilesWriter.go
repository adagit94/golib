package writers

import (
	"fmt"
	err "github.com/adagit94/golib/errors"
	qs "github.com/adagit94/golib/queues"
	"sync"
)

type SaveFile[D BuffData, R any] func(path string, data D) R
type WriteTask[R any] = qs.Task[R]
type PathQueue[R any] = qs.ChannelQueue[R]
type PathsQueues[R any] map[string]*PathQueue[R]
type BuffSize = int
type BuffData interface {
	string | []byte
}

type IFilesWriter[D BuffData, R any] interface {
	Schedule(path string, data D) *WriteTask[R]
	Close() (bool, *err.Err)
}

// Create FilesWriter with PathsQueues initialized to specified buffSize on first write attempt, with Closed channel initialized to be buffered with size of 1, so it's optional to use it as part of users code, and custom SaveFile function to perform actual file write operation.
func CreateFilesWriter[D BuffData, R any](buffSize BuffSize, saveFile SaveFile[D, R]) *FilesWriter[D, R] {
	fwPtr := &FilesWriter[D, R]{PathsQueues: make(PathsQueues[R]), BuffSize: buffSize, Closed: make(chan bool, 1), SaveFile: saveFile}
	return fwPtr
}

// FilesWriter is intended for ordered file writes based on the path. It triggers new goroutine that iterates over scheduled, ordered writes for each unique path. FilesWriter.Closed, Task.Skip and Task.Result channel are all buffered channels with size of 1, so it's not necessary to have receiver in your code in case you don't need it  - e.g. when FilesWriter.Closed or Task.Result aren't relevant for your use case.
//
// SaveFile - custom function to save file content.
//
// PathQueues - dedicated queues for each unique path.
//
// BuffSize - common channel queue size for each unique path.
//
// Closed - receives true after successfull close of all PathsQueues when Close get's called and successfully finishes to close and drain all write queues.
type FilesWriter[D BuffData, R any] struct {
	mu sync.Mutex
	SaveFile    SaveFile[D, R]
	PathsQueues PathsQueues[R]
	BuffSize    BuffSize
	Closed      chan bool
}

// Schedule data to be written under the path (directory structure get's created in case it doesn't exists). New goroutine with own event queue is spawned for each unique path and Loop starts immediately. Returned pointer to WriteTask can be used to access Result or control task execution additionaly before it's execution starts (e.g. to skip/cancel the task).
func (fw *FilesWriter[D, R]) Schedule(path string, data D) *WriteTask[R] {
	fw.mu.Lock()
	
	queuePtr, found := fw.PathsQueues[path]

	if !found {
		queuePtr = qs.CreateChannelQueue[R](fw.BuffSize)
		fw.PathsQueues[path] = queuePtr

		go queuePtr.Loop()
	}

	fw.mu.Unlock()

	taskPtr := &WriteTask[R]{Result: make(chan R, 1), Skip: make(chan bool, 1), Work: func() R {
		return fw.SaveFile(path, data)
	}}

	queuePtr.Enqueue(taskPtr)

	return taskPtr
}

// Close attempts to close and drain all PathsQueues and sends true to FilesWriter.Closed channel at the end in case all PathsQueues will be closed and drained successfully. It panics with ChannelNotClosed error Code in case any path queue wouldn't be successfully closed - that shouldn't never happen and should be reported as bug instead.
func (fw *FilesWriter[D, R]) Close() {
	for _, v := range fw.PathsQueues {
		v.Close()
	}

	for k, v := range fw.PathsQueues {
		if closed := <-v.Closed; !closed {
			panic(&err.Err{Code: err.ChannelNotClosedCode, Message: fmt.Sprint("Channel for a path ", k, " wasn't closed successfully.")})
		}
	}

	fw.Closed <- true
}
