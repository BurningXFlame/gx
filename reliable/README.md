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
**Better Solution**: Let the process-level guardian guard your service process. Even though your service process exits abnormally, the container keeps alive, because the guardian keeps alive. The guardian launches your service process instantly, and therefore decreases downtime to nearly 0.

**Case 2**: Client connections are closed because server process exits.
**Problem**: Clients have to re-connect the server.
**Better Solution**: Let the goroutine-level guardian guard your service goroutine. Even though your service exits abnormally, the process keeps alive, because the guardian keeps alive. The guardian launches your service instantly, and therefore decreases downtime to nearly 0. Furthermore, all established connections remain intact.

**Case 3**: Auto-restart service on config change.
**Problem**: Same as Case 1 and 2.
**Better Solution**: Let the auto-reloader watch the config file, and reload your service on config changes. It avoids container restart and process restart, and therefore decreases downtime to nearly 0. Furthermore, all established connections remain intact.

## Process-Level Guardian

[Supervisor](supervisor/cmd/supervisor/supervisor.go) starts and guards processes.

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

[Guard](guard/guard.go) starts and guards a function.

```go
import (
  "github.com/burningxflame/gx/reliable/guard"
  "github.com/burningxflame/gx/reliable/backoff"
)

// Auto re-run a function until it succeeds (aka, returns nil error) or ctx.Done channel is closed.
// If AlsoRetryOnSuccess is true, auto re-run a function until ctx.Done channel is closed.
guard.WithGuard(ctx, guard.Conf{
    Tag: "someTag", // Used to tag log messages
    // The func to be guarded. Fn must be canceled when ctx.Done channel is     closed.
    Fn: func(ctx context.Context) error { ... },
    // Backoff strategy determines how long to wait between retries.
    Bf: backoff.Default(),
    // If true, re-run Fn even if it returns nil error.
    AlsoRetryOnSuccess: false,
})
```

Sample: Start and guard a service until ctx.Done channel is closed.

```go
go guard.WithGuard(ctx, guard.Conf{
    Tag: "someService",
    Fn: func(ctx context.Context) error {
      return serve(ctx)
    },
    Bf: backoff.Default(),
    AlsoRetryOnSuccess: true,
})
```

Sample: Retry a task until it succeeds or ctx.Done channel is closed.

```go
go guard.WithGuard(ctx, guard.Conf{
    Tag: "someTask",
    Fn: func(ctx context.Context) error {
      return task(ctx)
    },
    Bf: backoff.Default(),
})
```

## Auto-Reload on Config Changes

[Auto-Reloader](autoreload/auto_reload.go) starts and re-run a function on config changes.

```go
import (
  "github.com/burningxflame/gx/reliable/autoreload"
  "github.com/burningxflame/gx/reliable/backoff"
)

// Watch a config file, and re-run a func on config changes.
// Will retry watching if the conf file is removed.
autoreload.WithAutoReload(ctx, Conf[C]{
    Tag:     "someTag", // Used to tag log messages
    Path:    "/some/path", // Path of the config file
    // Used to load config file on file write event. The returned C is the    loaded config. If C is the same as the last, it's ignored.
    Load:    func(path string) (C, error),
    // Process is the func to be reloaded. C is the config loaded by Load.    Process must be canceled when ctx.Done channel is closed.
    Process: func(context.Context, C),
    // Backoff strategy determines how long to wait between retries.
    Bf:      backoff.Default(),
})
```

## Backoff

[Backoff](backoff/backoff.go) is usually used to determine how long to wait between retries.

```go
import "github.com/burningxflame/gx/reliable/backoff"

// Create a Backoff.
bf := backoff.New(backoff.Conf{
    Min        time.Duration // Min delay
    Max        time.Duration // Max delay
    Unit       time.Duration // Unit of increment
    Strategy   Strategy      // Strategy of increment. Linear or Exponent.
    // If a retry lasts longer than ResetAfter, the next delay will be reset to Min.
    ResetAfter time.Duration
})

// Return the next delay.
dur := bf.Next()
```

## Readiness

[Readiness](readiness/readiness.go) is a TCP server for readiness check (aka, health check). Only for connectivity check. For security purpose, no sending data nor receiving data.

```go
import "github.com/burningxflame/gx/reliable/readiness"

// Start a TCP server at the specific address for readiness check.
readiness.Serve(ctx, "ip:port")
```

You may guard it:

```go
import (
  "github.com/burningxflame/gx/reliable/readiness"
  "github.com/burningxflame/gx/reliable/guard"
  "github.com/burningxflame/gx/reliable/backoff"
)

go guard.WithGuard(ctx, guard.Conf{
    Tag: "readiness",
    Fn: func(ctx context.Context) error {
      return readiness.Serve(ctx, "ip:port")
    },
    Bf: backoff.Default(),
    AlsoRetryOnSuccess: true,
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
import	"github.com/burningxflame/gx/reliable/timeout"

fn = timeout.WithTimeout(duration, fn)
```

**Samples**
[timeout_decorator](timeout/timeout_test.go)
