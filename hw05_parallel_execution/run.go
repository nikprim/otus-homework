package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) (returnErr error) {
	var errorCount int32
	tasksCh := make(chan Task)

	if len(tasks) < n {
		n = len(tasks)
	}

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			for task := range tasksCh {
				if err := task(); err != nil && m > 0 {
					atomic.AddInt32(&errorCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if m > 0 && int(atomic.LoadInt32(&errorCount)) >= m {
			returnErr = ErrErrorsLimitExceeded
			break
		}

		tasksCh <- task
	}

	close(tasksCh)
	wg.Wait()

	return returnErr
}
