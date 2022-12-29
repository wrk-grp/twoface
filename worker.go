package twoface

import (
	"io"
	"time"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
Worker wraps a concurrent process that is able to process Job types
scheduled onto a Pool.
*/
type Worker struct {
	pool    *Pool
	jobs    chan Job
	latency time.Duration
}

func NewWorker(pool *Pool) *Worker {
	errnie.Trace()
	return &Worker{pool, make(chan Job), 0 * time.Nanosecond}
}

func (worker *Worker) Read(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented")
	return
}

func (worker *Worker) Write(p []byte) (n int, err error) {
	errnie.Trace()

	go func() {
		t := time.Now()

		for {
			select {
			case <-worker.pool.ctx.Done():
				if errnie.Handles(worker.Close()) != nil {
					return
				}
			default:
				worker.pool.workers <- worker.jobs
				job := <-worker.jobs
				job.Do()
			}
		}

		worker.latency = time.Duration(time.Since(t).Nanoseconds())
	}()

	return len(p), io.EOF
}

func (worker *Worker) Close() error {
	errnie.Trace()
	errnie.Debugs("not implemented")
	return errnie.NewError(nil)
}
