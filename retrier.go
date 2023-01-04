package twoface

import (
	"math"
	"time"

	"github.com/wrk-grp/errnie"
	"github.com/wrk-grp/spd"
)

/*
Retrier is an interface that can be implemented by any object that wants to
schedule itself onto a worker pool and be retried under certain conditions.
*/
type Retrier interface {
	Do(Job, *spd.Datagram) *spd.Datagram
}

func NewRetrier(retrierType Retrier) Retrier {
	return retrierType
}

/*
Fibonacci is a RetryStategy that retries a function n times with a Fibonacci
interval in seconds between retries.
*/
type Fibonacci struct {
	max int
	n   int
}

func NewFibonacci(max int) Retrier {
	return NewRetrier(&Fibonacci{
		max: max,
		n:   0,
	})
}

func (strategy *Fibonacci) Do(fn Job, dg *spd.Datagram) *spd.Datagram {
	errnie.Trace()

	// We have reached the maximum number of retries.
	// Bail.
	if strategy.n > strategy.max {
		return nil
	}

	// Error, retry.
	if dg = fn.Do(dg); dg == nil {
		// Backoff delay time by using Fibonacci sequence.
		strategy.n = int(
			math.Round((math.Pow(
				math.Phi, float64(strategy.n),
			) + math.Pow(
				math.Phi-1, float64(strategy.n),
			)) / math.Sqrt(5)),
		)

		// Wait for the next retry.
		time.Sleep(time.Duration(strategy.n) * time.Second)
		strategy.Do(fn, dg)
	}

	return dg
}
