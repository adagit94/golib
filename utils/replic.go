package utils

import (
	"sync"
	t "github.com/adagit94/golib/types"
)

// Replicate is intended for replication uses cases where same function is triggered repeatedly with potentially variable argument (like FS or network URI). Every func. call for respective arguments runs in separate goroutine, so whole replication process can be performed in parallel. In case of failure op. can be repeated based on maxRetries. It returns slice of Results with same order as arguments passed. Non-nil errors should be returned from failed calls to attempt retry when maxRetries > 0 and report error correctly in Result. 
func Replicate[Arg any, Res any](maxRetries uint8, op func(arg Arg) (Res, error), args ...Arg) []*t.Result[Res] {
	var wg sync.WaitGroup
	resultsChan := make(chan *t.Result[Res], len(args))

	for _, arg := range args {
		wg.Go(func() {
			res, err := Retried(func() (Res, error) {
				return op(arg)
			}, maxRetries)

			resultsChan <- &t.Result[Res]{Result: res, Error: err}
		})
	}

	wg.Wait()
	close(resultsChan)
	
	return ChanToSlice(resultsChan)
}
