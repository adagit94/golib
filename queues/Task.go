package queues

// Task represents unit of execution by scheduler.
// 
// Work - executes when triggered and sends result R into Result channel after execution finishes. It can be skipped via send to skip channel. It's important to handle all errors gracefully inside of Work to prevent shut down of whole scheduler for example.
//
// Result - channel from which result should be received. At least in case it's unbuffered so as to prevent deadlocks.
//
// Skip - task can be skipped via send to Skip channel. In such case Work doesn't get executed.
type Task[R any] struct {
	Work func() R
	Result chan R 
	Skip chan bool
}

// Task with buffered channels of size 1 for both Result and Skip channel to make receive from them optional.
func CreateTask[R any](work func() R) *Task[R] {
	t := &Task[R]{Work: work, Result: make(chan R, 1), Skip: make(chan bool, 1)}
	return t
}

type ISkippableTask interface {
	// Task can be skipped in case it wasn't started yet.
	Skip()
}