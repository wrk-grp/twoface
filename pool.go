package twoface

import (
	"github.com/theapemachine/wrkspc/tweaker"
	"github.com/wrk-grp/errnie"
)

/*
Pool is a set of Worker types, each running their own (pre-warmed) goroutine.
Any object that implements the Job interface is able to schedule work on the
worker pool, which keeps the amount of goroutines in check, while still being
able to benefit from high concurrency in all kinds of scenarios.
*/
type Pool struct {
	ctx     *Context
	workers chan chan Job
	jobs    chan Job
}

/*
NewPool instantiates a worker pool with bound size of maxWorkers, taking in a
Context type to be able to cleanly cancel all of the sub processes it starts.
*/
func NewPool(ctx *Context) *Pool {
	errnie.Trace()
	jbuf := tweaker.GetInt("twoface.pool.job.buffer")
	errnie.Informs("job buffer", jbuf)

	return &Pool{
		ctx:     ctx,
		workers: make(chan chan Job),
		jobs:    make(chan Job, jbuf),
	}
}

/*
Do is the entry point for new jobs that want to be scheduled onto the worker pool.
*/
func (pool *Pool) Do(jobType Job) {
	errnie.Trace()

	// Send the job to the job channel.
	pool.jobs <- NewJob(jobType)
}

/*
Run the workers, after creating and assigning them to the pool.
*/
func (pool *Pool) Run() *Pool {
	errnie.Trace()

	// Start the auto-scaler to control the pool size dynamically.
	NewScaler(pool)

	// Start the job scheduling process.
	go pool.dispatch()
	return pool
}

func (pool *Pool) dispatch() {
	errnie.Trace()

	// Make sure that we cleanly close the channels if our dispatcher
	// returns for whatever reason.
	defer close(pool.jobs)
	defer close(pool.workers)

	for {
		select {
		case job := <-pool.jobs:
			errnie.Debugs("pickup job")
			// A new job was received from the jobs queue, get the first available
			// worker from the pool once ready.
			jobChannel := <-pool.workers
			errnie.Debugs("pickup worker")
			// Then send the job to the worker for processing.
			jobChannel <- job
			errnie.Debugs("do job")
		case <-pool.ctx.Done():
			errnie.Informs("pool going down")
			return
		}
	}
}
