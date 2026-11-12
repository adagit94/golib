package queues

// Generic interface that's meant to provide common syntactic and implicit semantic ground for general purpose scheduler (potentially with more or less different/variable implementation semantics across use cases).  
//
// Enqueue - intended for task scheduling. Agnostic of specific implementation semantics - e.g. channel, fixed size array or dynamically sized slice with heap allocations can be used.
//
// Loop - event loop like trigger intended for iteration over tasks/event queue as part of a scheduler.
//
// Close - mean't to close the channel queue or any other type of queue and prevent future submissions of tasks to the scheduler.
type IQueue[R any] interface {
	Enqueue(t *Task[R])
	Loop()
	Close()
}