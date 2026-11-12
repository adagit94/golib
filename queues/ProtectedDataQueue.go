package queues

// Operation intended to read/modify underlaying data in concurrency safe and sequenced way. Data d musn't be modified outside of function scope being executed and user code shall not make any assumptions about state of data after function being executed terminates. User code can clone data and use them in outer scope though.
type Op[D any] func(d *D)
type queue[D any] chan Op[D]

type protectedData[D any] struct {
	data D
	queue queue[D]
}

func (st *protectedData[D]) GetData() D {
	return st.data
}

func (st *protectedData[D]) Enqueue(op Op[D]) {
	st.queue <- op
}

func (st *protectedData[D]) Close() {
	close(st.queue)
}

func (st *protectedData[D]) loop() {
	for {
		f, closed := <- st.queue

		if closed {
			return
		}
		
		f(&st.data)
	}
}

type IProtectedDataQueue[D any] interface {
	GetData() D // Get data in current state. For more fine grained control in relation to predictability of data state in dependence to order of enqueued operations read data as part of enqueued operation itself. Underlaying data structure of returned data that isn't cloned - e.g. slice, map - musn't be changed by the user without explicit cloning in user code.
	Enqueue(op Op[D]) // Enqueue custom operation over underlaying data that's safe for concurrent usage and get's triggered in sequenced way. Users shall not make any other assumptions about safe concurrency, but the underlaying data state itself at the moment of operation execution - everything else have to be handled and guaranteed by the user.
	Close() // Close and drain underlaying queue of operations. Internal channel loop (running in it's own goroutine) will terminate afterwards. Any call to Enqueue after Close will panic.
}

// Create queue abstraction intended to read/modify underlaying data in concurrency safe and sequenced way. Receive loop that run's Enqueued operations is trigged automatically in it's own goroutine. 
func CreateProtectedDataQueue[D any](initialData D, queueSize int) IProtectedDataQueue[D] {
	st := &protectedData[D]{data: initialData, queue: make(chan Op[D], queueSize)}

	go st.loop()
	
	return st
}