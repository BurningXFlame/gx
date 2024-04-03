# ID

- [NanoID](#nanoid)
  - [Use](#use)
  - [Benchmark](#benchmark)
- [UUID](#uuid)
  - [Use](#use-1)
  - [Benchmark](#benchmark-1)

## NanoID

A Go implementation of [NanoID](https://github.com/ai/nanoid/blob/main/README.md).
NanoID can be expressed as regex `[A-Za-z0-9_-]{21}`. NanoID is case-sensative.

In contrast to UUID v4,

- similar collision probability.
- shorter: NanoID uses a larger alphabet (64 vs 16), and therefore shorter size (21 vs 36) under similar collision probability.
- safer: NanoID uses a cryptographically secure random generator.

### Use

```go
import "github.com/burningxflame/gx/id/nanoid"

// generate a NanoID
id, err := nanoid.New()


// Create a generator which generates NanoID in batches, for better performance.
// The batchSize specifies the number of NanoIDs to generate in each batch.
gen := nanoid.Generator(batchSize)
// generate a NanoID by calling generator
id, err := gen()


// Concurrency-safe version of Generator
gen := nanoid.ConGenerator(batchSize)
id, err := gen()
```

### Benchmark

```txt
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz

BenchmarkNew-12        904230       1270 ns/op       48 B/op        2 allocs/op

// batchSize: 10 ~ 10000
BenchmarkGenerator/10-12           6925892        172.1 ns/op       24 B/op        1 allocs/op
BenchmarkGenerator/100-12          7849317        154.5 ns/op       24 B/op        1 allocs/op
BenchmarkGenerator/1000-12         8230513        144.5 ns/op       24 B/op        1 allocs/op
BenchmarkGenerator/10000-12        8430886        141.2 ns/op       24 B/op        1 allocs/op

// batchSize: 10 ~ 10000
BenchmarkConGenerator/10-12 	  4647276	       254.7 ns/op	      24 B/op	       1 allocs/op
BenchmarkConGenerator/100-12         	 5425093	       217.4 ns/op	      24 B/op	       1 allocs/op
BenchmarkConGenerator/1000-12        	 5392015	       219.8 ns/op	      24 B/op	       1 allocs/op
BenchmarkConGenerator/10000-12       	 5167356	       233.1 ns/op	      24 B/op	       1 allocs/op
```

## UUID

UUID v4 and secure UUID v4.
UUID v4 uses a pseudo-random generator.
Secure UUID v4 uses a cryptographically secure random generator.

### Use

```go
import "github.com/burningxflame/gx/id/uuid"

// Generate a UUID v4
id := uuid.New()


// Generate a secure UUID v4
id, err := uuid.NewSec()


// Create a generator which generates secure UUID v4 in batches, for better performance.
// The batchSize specifies the number of UUIDs to generate in each batch.
gen := uuid.SecGenerator(batchSize)
// generate a UUID by calling generator
id, err := gen()


// Concurrency-safe version of SecGenerator
gen := uuid.SecConGenerator(batchSize)
id, err := gen()
```

### Benchmark

```txt
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz

BenchmarkNew-12                 13437386         87.45 ns/op       64 B/op        2 allocs/op

BenchmarkNewSec-12                755618       1400 ns/op       80 B/op        3 allocs/op

BenchmarkSecGenerator/10-12      5504864        221.1 ns/op       64 B/op        2 allocs/op
BenchmarkSecGenerator/100-12     6727554        170.5 ns/op       64 B/op        2 allocs/op
BenchmarkSecGenerator/1000-12    7639945        145.9 ns/op       64 B/op        2 allocs/op
BenchmarkSecGenerator/10000-12   7874086        145.8 ns/op       64 B/op        2 allocs/op

BenchmarkSecConGenerator/10-12   4009354        269.7 ns/op       64 B/op        2 allocs/op
BenchmarkSecConGenerator/100-12           5280829        223.8 ns/op       64 B/op        2 allocs/op
BenchmarkSecConGenerator/1000-12          5449786        218.1 ns/op       64 B/op        2 allocs/op
BenchmarkSecConGenerator/10000-12         4941277        224.3 ns/op       64 B/op        2 allocs/op
```
