package queues

type IChannelQueue[R any] IQueue[R]

// Create ChannelQueue of specified size with Closed channel initialized to be buffered with size of 1, so it's optional to use it as part of users code.
func CreateChannelQueue[R any](size int) *ChannelQueue[R] {
	qPtr := &ChannelQueue[R]{Queue: make(chan *Task[R], size), Closed: make(chan bool, 1)}
	return qPtr
}

// ChannelQueue is FIFO (usually buffered) queue. It iterates over tasks and triggers Work for them (or skips them - based on value in Task.Skip channel). It can be closed - afterwards, Enqueue will panic in case there would be attempt to enqueue the new task. True is send into Closed channel after queue get's drained.
type ChannelQueue[R any] struct {
	Queue  chan *Task[R]
	Closed chan bool
}

// Send task into (usually buffered) channel queue to schedule it for execution.
func (q *ChannelQueue[R]) Enqueue(t *Task[R]) {
	q.Queue <- t
}

// Triggers event loop like for cycle that awaits tasks to be received from the Queue and executes (triggers Work) or skips (Task.Skip channel) them. After Close, remaining tasks are drained from the queue and true is send to Close channel afterwards. Loop terminates.
//
// It's important to handle all errors gracefully inside of Task.Work to prevent shut down of whole loop.
func (q *ChannelQueue[R]) Loop() {
	for {
		t, ok := <-q.Queue

		if !ok {
			q.Closed <- true
			return
		}

		select {
		case <-t.Skip:
			continue

		default:
			t.Result <- t.Work()
		}
	}
}

// Close closes the queue and remaining tasks are drained from it. True is send to Close channel afterwards and event loop terminates.
func (q *ChannelQueue[R]) Close() {
	close(q.Queue)
}
