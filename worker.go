package twoface

import (
	"time"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
Worker wraps a concurrent process that is able to process Job types
scheduled onto a Pool.
*/
type Worker struct {
	pool    chan chan Job
	jobs    chan Job
	latency time.Duration
}

func NewWorker(pool chan chan Job, jobs chan Job) *Worker {
	errnie.Trace()
	return &Worker{pool, jobs, 0 * time.Nanosecond}
}

func (worker *Worker) Read(p []byte) (n int, err error) {
	errnie.Trace()
	errnie.Debugs("not implemented")
	return
}

func (worker *Worker) Write(p []byte) (n int, err error) {
	errnie.Trace()
	t := time.Now()
	errnie.Debugs("not implemented")
	worker.latency = time.Duration(time.Since(t).Nanoseconds())
	return
}

func (worker *Worker) Close() error {
	errnie.Trace()
	errnie.Debugs("not implemented")
	return errnie.NewError(nil)
}
