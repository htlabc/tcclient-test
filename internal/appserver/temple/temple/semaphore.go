package temple

import "golang.org/x/sync/semaphore"

func newSemaphore(workers int) *semaphore.Weighted {
	return semaphore.NewWeighted(int64(workers))
}
