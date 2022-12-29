# twoface

Higher abstractions of concurrency "primitives" to make working
with concurrent and/or distributed systems easier.

## context

A wrapper around native go `context` adding more boilerplate for
ergonomics.

## job

A primitive interface that allows scheduling onto pools.

Comes with a built in `RetrieableJob` implementation type that uses
retries with gradual backoff.

## pool

A goroutine worker pool that is combined with `scaler`.

## scaler

Auto scaling and load balancing for goroutine pools.

## worker

A wrapper that exposes worker goroutines as `io.ReadWriteCloser` types.
