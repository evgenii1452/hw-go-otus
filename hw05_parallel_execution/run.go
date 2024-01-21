package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded      = errors.New("errors limit exceeded")
	ErrIncorrectNumberOfWorkers = errors.New("number of workers must be greater than 0")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrIncorrectNumberOfWorkers
	}

	var errCounter int32
	taskChan := make(chan Task, len(tasks))

	produce(taskChan, tasks)

	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go process(taskChan, &errCounter, m, &wg)
	}

	wg.Wait()

	if m > 0 && int(atomic.LoadInt32(&errCounter)) >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func process(ch <-chan Task, errCounter *int32, maxErrors int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		task, ok := <-ch
		if !ok || int(atomic.LoadInt32(errCounter)) >= maxErrors && maxErrors > 0 {
			return
		}

		if err := task(); err == nil {
			continue
		}

		atomic.AddInt32(errCounter, 1)
	}
}

func produce(ch chan<- Task, tasks []Task) {
	defer close(ch)

	for _, task := range tasks {
		ch <- task
	}
}
