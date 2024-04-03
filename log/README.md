# Logging

- [Logging Facade](#logging-facade)
  - [For Logger Implementers](#for-logger-implementers)
    - [Implement interface log.Logger](#implement-interface-loglogger)
    - [Register Self as the Global Logger](#register-self-as-the-global-logger)
  - [For Logger Users](#for-logger-users)
    - [Register a Logger](#register-a-logger)
    - [Use Logging Facade](#use-logging-facade)
- [Light Logger](#light-logger)
  - [Use](#use)
  - [Performance](#performance)
- [Concurrent Buffer Writer](#concurrent-buffer-writer)
- [Auto-Flusher](#auto-flusher)
- [Log Rotator](#log-rotator)
- [Extension](#extension)

## Logging Facade

[Log](log/log.go) is a logging facade with leveled logging, tagged logging.

### For Logger Implementers

#### Implement interface [log.Logger](log/log.go)

```go
// A concrete logger should implement interface Logger.
type Logger interface {
    // Print a log message
    Printf(format string, v ...any)
    // Close the logger. Flush buffer, close files, etc.
    Close() error
}
```

#### Register Self as the Global Logger

Call `log.Set` to register self as the global logger, in two methods:

**Method 1:**
Call `log.Set` in `init`, and therefore registered when the package is imported.

```go
package dummy_log_x

import "github.com/burningxflame/gx/log/log"

func init() {
    logger := ...
    log.Set(logger, log.LevelInfo)
}
```

**Method 2:**
Provide public Init functions, and call `log.Set` in those Init functions.

```go
package dummy_log_y

import "github.com/burningxflame/gx/log/log"

func InitRuntimeLog(...) {
    ...
    log.Set(...)
}

func InitTestLog() {
    ...
    log.Set(...)
}
```

### For Logger Users

#### Register a Logger

Corresponding to the 2 methods in above section.

**Method 1:**

```go
package main

// register a logger on package import
import (
    _ "dummy_log_x"
    "github.com/burningxflame/gx/log/log"
)

func main(){
    // Close the logger. Flush buffer, close files, etc.
    // Must be called before process exit.
    defer log.Close()
    ...
}
```

**Method 2:**

```go
package main

import (
    "dummy_log_y"
    "github.com/burningxflame/gx/log/log"
)

func main(){
    // explicitly register a logger
    dummy_log_y.InitRuntimeLog(...)
    // Close the logger. Flush buffer, close files, etc.
    // Must be called before process exit.
    defer log.Close()
    ...
}
```

#### Use Logging Facade

```go
import "github.com/burningxflame/gx/log/log"

// leveled logging
log.Error(...)
log.Warn(...)
log.Info(...)
log.Debug(...)
log.Trace(...)

// tagged logging
// every log message will be prefixed with the tag, e.g., INFO  [tag1] some info
tl := log.WithTag("tag1")
tl.Error(...)
tl.Warn(...)
tl.Info(...)
tl.Debug(...)
tl.Trace(...)

// tag chain
// every log message will be prefixed with tag chain, e.g., INFO  [tag1] [tag2] some info
tl2 := tl.WithTag("tag2") // equivalent to log.WithTag("tag1").WithTag("tag2")
tl2.Error(...)
tl2.Warn(...)
tl2.Info(...)
tl2.Debug(...)
tl2.Trace(...)
```

## Light Logger

[Light](light/light.go) is an all-in-one logger. Light Logger = Go Std Log Wrapper + Concurrent Buffer Writer (with Auto Flusher) + Log Rotator.

### Use

**Use in Production:**

```go
import "github.com/burningxflame/gx/log/light"

err := light.Init(light.Conf{
    Level: log.LevelInfo, // Log level. Default to LevelError.
    // Log format flag. Refer to go std log. Default to LstdFlags |   Lmicroseconds | Lmsgprefix.
    Format: ...
    BufSize: 1<<20, // Buffer Size in bytes. Default to 1M.
    FlushInterval: time.Second*5, // Auto-flush interval. Default to 5s.
    Rc: light.RotateConf{ // Log-rotating config
        FilePath: ..., // Fullpath of log file
        // Max byte size of a log file. If a file exceeds this size, the file   will be rotated. Default to 10MB.
        FileSize: 10<<20,
        // Max number of old log files. Older files will be removed.
        NBak: 2,
        // Permission of log file. Default to 0600.
        Perm: 0600,
        // If true, rotated log files will not be compressed. Otherwise,   rotated   log files will be compressed with gzip.
        NoCompress: false,
        // If ture, rotated log files will be renamed based on UTC time. Local     time otherwise.
        Utc: false,
    },
})
```

**Use in Test:**

```go
import "github.com/burningxflame/gx/log/light"

// For test purpose, create a simple logger, and register it as the global logger.
// All log messages will be written to stdout, and log leve is debug.
err := light.InitTestLog()
```

### Performance

```txt
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz

Benchmark/short-12    	 2644009	       454.4 ns/op	      82 B/op	       1 allocs/op
Benchmark/long-12     	 1719120	       627.6 ns/op	     357 B/op	       1 allocs/op
```

A short log message is like:

```txt
2023/12/23 20:51:34.523690 INFO  ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-
```

converted into disk IO speed:
$ 10^9 / 454.4 * 98 /2^{20} = 206 MB/s $

A long log message is like:

```txt
2023/12/23 20:51:34.523690 INFO  ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-
```

converted into disk IO speed:
$ 10^9 / 627.6 * 354 /2^{20} = 538 MB/s $

## Concurrent Buffer Writer

[Conbuf](conbuf/buf.go) is a concurrency-safe buffer writer.

```go
import "github.com/burningxflame/gx/log/conbuf"

// Wrap a Writer and create a concurrent Buffer Writer
bw := conbuf.NewWriter(w, bufSize)
```

**Performance**:

```txt
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz

BenchmarkBufWriter-12    	12018157	        96.62 ns/op	       0 B/op	       0 allocs/op
BenchmarkNoBuf-12        	  286538	      3540 ns/op	       0 B/op	       0 allocs/op
```

## Auto-Flusher

[AutoFlusher](conbuf/flush.go) wraps a WriteFlusher and return another WriteFlusher which auto flushes the wrapped WriteFlusher at certain intervals.
AutoFlusher can be useful when you want to check fresh log messages, but the buffer is not full and therefore not flushed yet.

```go
import "github.com/burningxflame/gx/log/conbuf"

// Wrap a WriteFlusher and return another WriteFlusher which auto flushes the wrapped WriteFlusher at certain intervals until ctx.Done channel is closed.
fw := conbuf.WithAutoFlush(
    ctx,
    w, // the WriteFlusher to be wrapped
    interval, // the flush interval
    // Specify a buffer chan if you want to receive background flush errors if any. Leave it nil otherwise.
    chErr,
)
```

## Log Rotator

[Log Rotator](rotate/rotate.go) provides abilities such as

- rotating log files
- compressing rotated files
- removing old files
- re-create log files if deleted from outside

```go
import "github.com/burningxflame/gx/log/rotate"

// Create a log rotator. The returned wc is a WriteCloser.
wc, err := rotate.New(rotate.Conf{
    FilePath: ..., // Fullpath of log file
    // Max byte size of a log file. If a file exceeds this size, the file will   be rotated. Default to 10MB.
    FileSize: 10<<20,
    // Max number of old log files. Older files will be removed.
    NBak: 2,
    // Permission of log file. Default to 0600.
    Perm: 0600,
    // If true, rotated log files will not be compressed. Otherwise, rotated   log files will be compressed with gzip.
    NoCompress: false,
    // If ture, rotated log files will be renamed based on UTC time. Local time   otherwise.
    Utc: false,
})
...

// Close the log rotator. Close files, finish handling backups, etc.
// Must be called before process exit.
wc.Close()
```

## Extension

[Light](#light-logger) is an all-in-one logger. However you may perfer another.

You can

- wrap the logger you preferred and implement the [log.Logger](log/log.go), and then use the [Logging Facade](#logging-facade) for consistent user experience.
- enhance the logger you preferred by combining the logger with [Concurrent Buffer Writer](#concurrent-buffer-writer), [Auto-Flusher](#auto-flusher), and/or [Log Rotator](#log-rotator).
