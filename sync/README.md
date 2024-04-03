# sync

- [Semaphore](#semaphore)
  - [Use](#use)
  - [Benchmark](#benchmark)
- [Keyed-Semaphores](#keyed-semaphores)

## Semaphore

[Semaphore](sem/sem.go) is commonly used for limiting max concurrency, e.g. max concurrent connections.

### Use

```go
import "github.com/burningxflame/gx/sync/sem"

// Create a semaphore.
// The ca specifies the capacity of the semaphore.
s := sem.New(ca)

// Acquire a permit from the semaphore.
// If none is available, block until one is available or ctx.Done channel is closed.
err := s.Acquire(ctx)

// Try to acquire a permit from the semaphore.
// Return true if available, false otherwise.
ok := s.TryAcquire()

// Release a permit to the semaphore.
s.Release()

// Return the number of available permits.
n := s.Available()
```

### Benchmark

```txt
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz

// Acquire + Release
BenchmarkAcquireRelease-12     22697646         51.96 ns/op        0 B/op        0 allocs/op
```

## Keyed-Semaphores

Commonly used for limiting max concurrency per key, e.g. max concurrent connections per client.

```go
import "github.com/burningxflame/gx/sync/sem"

// Create Keyed-Semaphores.
// The ca specifies the capacity of every semaphore.
// If the number of semaphores exceeds sizeHint, will try to shrink, i.e. remove semaphores who have no permits taken.
ks := sem.NewKSem[string](ca, sizeHint)

// Get the semaphore of a key, create if not exist.
s := ks.Get(key)

// ... use the semaphore

// Get the semaphore of another key, create if not exist.
s2 := ks.Get(key2)

// ... use the semaphore
```
