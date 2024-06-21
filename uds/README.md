# UDS - Unix Domain Socket

- [HTTP Server over UDS](#http-server-over-uds)
- [HTTP Client over UDS](#http-client-over-uds)

## HTTP Server over UDS

```go
import (
  "net/http"
  uh "github.com/burningxflame/gx/uds/http"
)

// Create an HTTP server over UDS
srv := &uh.Server{
  // http.Server in std lib
  Std: http.Server{
    Handler: someHandler,
    ...
  },
  // The UDS address to listen
  UdsAddr: "/some/path",
  // File permission of the UdsAddr
  Perm: 0600,
  // Used to limit max number of concurrent requests.
  // Default to no limit.
  Limiter: sem.New(n),
  // If graceful shutdown takes longer than ShutdownTimeout, exit instantly.
  ShutdownTimeout: time.Second*3,
  // Used to tag log messages
  Tag: "someTag",
  // A TagLogger used to log messages
  Log: ...,
}

// Start an HTTP server over UDS
err := srv.Serve(ctx)
```

## HTTP Client over UDS

```go
import (
  "net/http"
  uh "github.com/burningxflame/gx/uds/http"
)

// Create an HTTP client over UDS.
// The return value c is *http.Client.
c := uh.NewClient(udsAddr)
// Use the client
c.Timeout = time.Second * 3
resp, err := c.Post(...)
```
