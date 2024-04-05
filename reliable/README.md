# Reliable

Provides utilities to improve application reliability.

- [Typical Use Case](#typical-use-case)
- [Process-Level Guardian](#process-level-guardian)
- [Goroutine-Level Guardian](#goroutine-level-guardian)
- [Auto-Reload on Config Changes](#auto-reload-on-config-changes)
- [Backoff](#backoff)
- [Readiness](#readiness)
- [Timeout Decorator](#timeout-decorator)

## Typical Use Case

For example, you develop a TCP service, and deploy the service to K8S.

**Case 1**: Container restarts because the service process exits abnormally.
**Problem**: It usually takes seconds ~ dozens of seconds before the container is up again.
**Better Solution**: Let the process-level guardian guard your service process. Even though your service process exits abnormally, the container keeps alive, because the guardian keeps alive. The guardian launches your service process instantly, and therefore decreases downtime to 0.

**Case 2**: Client connections are closed because server process exits.
**Problem**: Clients have to re-connect the server.
**Better Solution**: Let the goroutine-level guardian guard your service goroutine. Even though your service exits abnormally, the process keeps alive, because the guardian keeps alive. The guardian launches your service instantly, and therefore decreases downtime to 0. Furthermore, all established connections remain intact.

**Case 3**: Auto-restart service on config change.
**Problem**: Same as Case 1 and 2.
**Better Solution**: Let the auto-reloader watch the config file, and reload your service on config changes. It avoids container restart and process restart, and therefore decreases downtime to 0. Furthermore, all established connections remain intact.

## Process-Level Guardian

Supervisor starts and guards processes.

1. Install supervisor:

   ```sh
   # install location: $(go env GOPATH)/bin/.
   go install github.com/burningxflame/gx/reliable/supervisor/cmd/supervisor@latest
   ```

2. Run supervisor:

   ```sh
   supervisor conf.yaml
   ```

   [Sample config file](supervisor/cmd/supervisor/testdata/sample_conf.yaml).

## Goroutine-Level Guardian

Guard starts and guards a function.

```go
import (
  "github.com/burningxflame/gx/reliable/guard"
  "github.com/burningxflame/gx/reliable/backoff"
)

// Auto re-run a function until it succeeds (aka, returns nil error) or ctx.Done channel is closed.
// If AlsoRetryOnSuccess is true, auto re-run a function until ctx.Done channel is closed.
guard.WithGuard(ctx, guard.Conf{
    // The func to be guarded.
    // Fn should return ASAP when ctx.Done channel is closed, which usually means an exit signal is sent.
    Fn: func(ctx context.Context) error { ... },
    // Backoff strategy determines how long to wait between retries.
    Bf: backoff.Default(),
    // If true, re-run Fn even if it returns nil error.
    AlsoRetryOnSuccess: false,
    // Used to tag log messages
    Tag: "someTag",
    // A TagLogger used to log messages
    Log: ...,
})
```

Sample: Start and guard a service until ctx.Done channel is closed.

```go
go guard.WithGuard(ctx, guard.Conf{
    Fn: func(ctx context.Context) error {
      return serve(ctx)
    },
    Bf: backoff.Default(),
    AlsoRetryOnSuccess: true,
    Tag: "someService",
})
```

Sample: Retry a task until it succeeds or ctx.Done channel is closed.

```go
go guard.WithGuard(ctx, guard.Conf{
    Fn: func(ctx context.Context) error {
      return task(ctx)
    },
    Bf: backoff.Default(),
    Tag: "someTask",
})
```

## Auto-Reload on Config Changes

Auto-Reloader starts and re-run a function on config changes.

```go
import (
  "github.com/burningxflame/gx/reliable/autoreload"
  "github.com/burningxflame/gx/reliable/backoff"
)

// Watch a config file, and re-run a func on config changes.
// Will retry watching if the conf file is removed.
autoreload.WithAutoReload(ctx, Conf[C]{
    // Path of the config file
    Path: "/some/path",
    // Used to load config file on file write event. The returned C is the loaded config.
    // If C is the same as the last, it's ignored.
    Load: func(path string) (C, error),
    // Process is the func to be reloaded. C is the config loaded by Load.
    // Process should return ASAP when ctx.Done channel is closed.
    Process: func(ctx context.Context, c C),
    // Backoff strategy determines how long to wait between retries.
    Bf: backoff.Default(),
    // Used to tag log messages
    Tag: "someTag",
    // A TagLogger used to log messages
    Log: ...,
})
```

## Backoff

Backoff is usually used to determine how long to wait between retries.

```go
import "github.com/burningxflame/gx/reliable/backoff"

// Create a Backoff.
bf := backoff.New(backoff.Conf{
    // Min delay
    Min: time.Millisecond,
    // Max delay
    Max: time.Second * 30,
    // Unit of increment
    Unit: time.Second,
    // Strategy of increment. Linear or Exponent.
    Strategy: backoff.Exponent,
    // If a retry lasts longer than ResetAfter, the next delay will be reset to Min.
    ResetAfter: time.Second * 30,
})

// Return the next delay.
dur := bf.Next()
```

## Readiness

Readiness is a TCP Server for readiness check (aka, health check). Only for connectivity check. For security purpose, no sending data nor receiving data.

```go
import "github.com/burningxflame/gx/reliable/readiness"

// Create a TCP Server for readiness check
srv := &readiness.Server {
  // The address to listen
  Addr: "host:port",
  // Used to tag log messages. Default to "readiness".
  Tag: "readiness",
  // A TagLogger used to log messages
  Log: ...,
}

// Start the Server
err := srv.Serve(ctx)
```

You may guard it:

```go
import (
  "github.com/burningxflame/gx/reliable/readiness"
  "github.com/burningxflame/gx/reliable/guard"
  "github.com/burningxflame/gx/reliable/backoff"
)

go guard.WithGuard(ctx, guard.Conf{
    Fn: func(ctx context.Context) error {
      return srv.Serve(ctx)
    },
    Bf: backoff.Default(),
    AlsoRetryOnSuccess: true,
    Tag: "readiness",
})
```

## Timeout Decorator

A Timeout Decorator abandons the call to the wrapped function if the call does not finish in a specified duration.
There're 4 variants, each for a type of functions. All functions can be classfied into these four.

```go
// Timeout Decorator for functions with input and output parameters, i.e func(I) (O, error)
func WithTimeoutIO[I, O any](timeout time.Duration, fn func(I) (O, error)) func(I) (O, error)

// Timeout Decorator for functions with input parameters only, i.e func(I) error
func WithTimeoutI[I any](timeout time.Duration, fn func(I) error) func(I) error

// Timeout Decorator for functions with output parameters only, i.e func() (O, error)
func WithTimeoutO[O any](timeout time.Duration, fn func() (O, error)) func() (O, error)

// Timeout Decorator for functions with neither input nor output parameters, i.e func() error
func WithTimeout(timeout time.Duration, fn func() error) func() error
```

```go
import "github.com/burningxflame/gx/reliable/timeouts"

fn = timeouts.WithTimeout(timeout, fn)
```

**Samples**
[timeout_decorator](timeouts/timeout_test.go)
