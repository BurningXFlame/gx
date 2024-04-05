# Runtime

- [Goroutine ID](#goroutine-id)
  - [Why](#why)
  - [Use](#use)
  - [Benchmark](#benchmark)
- [Goroutine Local Storage](#goroutine-local-storage)
  - [Why](#why-1)
  - [Key Design](#key-design)
  - [Use](#use-1)
  - [Benchmark](#benchmark-1)
- [G2G](#g2g)
  - [Use](#use-2)

## Goroutine ID

GID provides an extremely fast (1.2 ns) method to get the ID of the current goroutine.
More precisely speaking, it's the address of the current g (i.e. goroutine). Therefore, the ID makes sense only during the life cycle of the current goroutine.

It's implemented by using Go Assembly.
Supported CPU architectures: 386, amd64, arm, arm64, mips, mipsle, mips64, mips64le, ppc64, ppc64le, riscv64, s390x.

### Why

Normally you don't need to get the ID of the current goroutine, unless you're implementing a goroutine local storage.

### Use

```go
import "github.com/burningxflame/gx/runtime/gid"

// Return the ID of the current goroutine.
// More precisely speaking, it's the address of the current g (i.e. goroutine). Therefore, the ID makes sense only during the life cycle of the current goroutine.
id := gid.Gid()
```

### Benchmark

```txt
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz

BenchmarkGid-12    	903557202	         1.268 ns/op	       0 B/op	       0 allocs/op
```

## Goroutine Local Storage

GLS provides an extremely fast Goroutine Local Storage.

### Why

Normally you don't need goroutine local storage, because Context is better.
Unless you want to implement distributed tracing (or serverless, or other context-passing cases) at minimal cost for a legacy codebase where Context is not used widely or even not used at all.

### Key Design

Key points to design a goroutine local storage:

1. How to get the ID of the current goroutine? This is vital to relating goroutines and goroutine local storages. And this is vital to performance.
2. What to use as goroutine local storage? This is vital to performance.
3. How to detect goroutine creation? This is vital to creating goroutine local storage.
4. How to detect goroutine exit? This is vital to removing goroutine local storage instantly and therefore minimizing memory footprint.
5. Minimizing cost to apply goroutine local storage in a legacy codebase.

As follows are my choices:

1. Use Go Assembly to get the address of the current g (i.e. goroutine), and use it as the ID of the current goroutine. It's extremely fast. See [Goroutine ID](#goroutine-id).
2. Use std Value Context (context.WithValue) as goroutine local storage. In context-passing cases such as distributed tracing, serverless, etc, a goroutine local storage is copied mostly, read mostly, and written rarely.
Value Context is pointer, and therefore can be copied almost costlessly.
Value Context is immutable because it provides no API to modify itself. If you want to "modify" a Value Context, you can always create another based on the current one. (It feels like immutability in functional programming.) Therefore Value Context is concurrency-safe by nature.
In conclusion, Value Context is perfect for goroutine local storage.
3. Use a custom function `Go` instead of the keyword `go` to spawn a goroutine. Therefore we know exactly when a goroutine is created or exits.
4. See 3.
5. Use Go AST to create a code generator which scans a codebase and replaces all keyword `go` with function `Go`. See [G2G](#g2g)

### Use

```go
import "github.com/burningxflame/gx/runtime/gls"

// Associate the value with the key in the current goroutine local storage.
gls.Put(key, value)

// Return the value corresponding to the key in the current goroutine local storage.
value := gls.Get(key)

// Use function `Go` instead of the keyword `go` to spawn a goroutine. The spawned goroutine inherits local storage from its parent goroutine.
gls.Go(func(){
  ...

  // cv == value, because the spawned goroutine inherits local storage from its parent goroutine.
  cv := gls.Get(key)

  // Put another k-v pair.
  gls.Put(key2, value2)
})

```

### Benchmark

```txt
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz

BenchmarkGlsGet-12    	28428370	        37.74 ns/op	       8 B/op	       1 allocs/op

// Cost of GLS (create and drop GLS): 534-263=271 ns per goroutine
BenchmarkGlsCostA-12    	 2271580	       534.1 ns/op	     105 B/op	       5 allocs/op
BenchmarkGlsCostB-12    	 4524012	       263.5 ns/op	      16 B/op	       1 allocs/op
```

## G2G

G2G is a code generator which scans a codebase and replaces all keyword `go` with function `Go`.
It also takes care of the closure problem of for/range.
It's implemented by using Go AST.

### Use

```sh
go install github.com/burningxflame/gx/runtime/g2g@latest

g2g <path>
```
