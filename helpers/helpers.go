package helpers

import (
	"sync"

	"github.com/pkg/errors"
)

func EachLimit(works <-chan func() error, finished chan<- error, limit int) {
	done := make(chan struct{})
	defer close(done)
	errc := make(chan error)

	var wg sync.WaitGroup
	wg.Add(limit)
	for i := 0; i < limit; i++ {
		go func() {
			digester(done, works, errc)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	err := <-errc
	if err != nil {
		finished <- errors.Wrap(err, "eachlimit")
	} else {
		finished <- nil
	}
}
func digester(done <-chan struct{}, fns <-chan func() error, errc chan<- error) {
	for fn := range fns {
		err := fn()
		select {
		case <-done:
			return
		default:
		}
		if err != nil {
			errc <- err
			return
		}
	}
}
