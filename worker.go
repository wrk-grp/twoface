package twoface

import (
	"bytes"
	"io"
	"time"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
Worker wraps a concurrent process that is able to process Job types
scheduled onto a Pool.
*/
type Worker struct {
	ctx     *Context
	pool    *Pool
	jobs    chan Job
	buffer  *bytes.Buffer
	latency time.Duration
}

func NewWorker(pool *Pool) *Worker {
	errnie.Trace()

	return &Worker{
		NewContext(),
		pool,
		make(chan Job),
		bytes.NewBuffer([]byte{}),
		0 * time.Nanosecond,
	}
}

func (worker *Worker) Read(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented")
	return
}

func (worker *Worker) Write(p []byte) (n int, err error) {
	errnie.Trace()

	go func() {
		for {
			select {
			case <-worker.pool.ctx.Done():
				if errnie.Handles(worker.Close()) != nil {
					return
				}
			case <-worker.ctx.root.Done():
				return
			default:
				worker.pool.workers <- worker.jobs
				job := <-worker.jobs

				t := time.Now()
				job.Do()
				worker.latency = time.Duration(
					time.Since(t).Nanoseconds(),
				)
			}
		}

	}()

	return len(p), io.EOF
}

func (worker *Worker) Close() error {
	errnie.Trace()
	worker.ctx.cancel()
	return errnie.NewError(nil)
}
