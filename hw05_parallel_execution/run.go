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

type producer struct {
	stopSignalCh <-chan struct{}
	ch           chan<- Task
	wg           *sync.WaitGroup
}

func (p *producer) produce(tasks []Task) {
	defer func() {
		close(p.ch)
		p.wg.Done()
	}()

	for _, task := range tasks {
		select {
		case <-p.stopSignalCh:
			return
		case p.ch <- task:
			continue
		}
	}
}

type worker struct {
	readCh       <-chan Task
	wg           *sync.WaitGroup
	errCounter   int32
	maxErrors    int32
	stopSignalCh chan<- struct{}
	stopped      int32
}

func (w *worker) process() {
	defer w.wg.Done()

	for task := range w.readCh {
		if atomic.LoadInt32(&w.stopped) == 1 {
			return
		}

		if err := task(); err == nil {
			continue
		}

		atomic.AddInt32(&w.errCounter, 1)

		if atomic.LoadInt32(&w.errCounter) >= w.maxErrors && w.maxErrors > 0 && atomic.CompareAndSwapInt32(&w.stopped, 0, 1) {
			w.stopSignalCh <- struct{}{}
		}
	}
}

func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrIncorrectNumberOfWorkers
	}

	taskCh := make(chan Task)
	stopSignalCh := make(chan struct{})
	wg := sync.WaitGroup{}

	producer := &producer{ch: taskCh, wg: &wg, stopSignalCh: stopSignalCh}
	wg.Add(1)
	go producer.produce(tasks)

	worker := &worker{readCh: taskCh, wg: &wg, maxErrors: int32(m), stopSignalCh: stopSignalCh}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker.process()
	}

	wg.Wait()

	if m > 0 && int(atomic.LoadInt32(&worker.errCounter)) >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
