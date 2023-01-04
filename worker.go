package twoface

import (
	"bytes"
	"io"
	"time"

	"github.com/wrk-grp/errnie"
	"github.com/wrk-grp/spd"
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
				errnie.Warns("pool closing")
				if errnie.Handles(worker.Close()) != nil {
					return
				}
			case <-worker.ctx.root.Done():
				errnie.Warns("worker done")
				return
			default:
				worker.pool.workers <- worker.jobs
				job := <-worker.jobs

				t := time.Now()
				dg := &spd.Empty
				dg.Decode(p)

				_ = job.Do(dg)

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
