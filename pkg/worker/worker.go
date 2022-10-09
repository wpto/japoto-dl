package worker

import (
	"sync"
)

type doFunc = func() error

type BoundWorker struct {
	Input     chan doFunc
	semaphore chan struct{}
	done      chan struct{}
	err       error
	mtx       sync.Mutex
}

func NewBoundWorker(limit int) *BoundWorker {
	wk := &BoundWorker{
		Input:     make(chan doFunc),
		semaphore: make(chan struct{}, limit),
		done:      make(chan struct{}),
		err:       nil,
		mtx:       sync.Mutex{},
	}
	go wk.worker()
	return wk
}

func (wk *BoundWorker) Close() {
	close(wk.Input)
}

func (wk *BoundWorker) Done() <-chan struct{} {
	return wk.done
}

func (wk *BoundWorker) Wait() {
	<-wk.done
}

func (wk *BoundWorker) Err() (err error) {
	wk.mtx.Lock()
	defer wk.mtx.Unlock()
	return wk.err
}

func (wk *BoundWorker) worker() {
	defer close(wk.done)
	wg := sync.WaitGroup{}

	for do := range wk.Input {
		do := do
		wg.Add(1)

		go func() {
			defer wg.Done()

			select {
			case wk.semaphore <- struct{}{}:
			case <-wk.done:
				return
			}

			err := do()
			if err != nil {
				wk.mtx.Lock()
				wk.err = err
				close(wk.done)
				wk.mtx.Unlock()
			}

			select {
			case <-wk.semaphore:
			case <-wk.done:
				return
			}
		}()
	}

	wg.Wait()
}
