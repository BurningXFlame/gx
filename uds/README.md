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
s := &uh.Server{
  // http.Server in std lib
  Std: http.Server{
    Handler: someHandler,
    ...
  },
  // Used to tag log messages
  Tag:     "someTag",
  // The UDS address to listen
  UdsAddr: "/some/path",
  // File permission of the UdsAddr
  Perm:    0600,
  // If graceful shutdown takes longer than ShutdownTimeout, exit instantly.
  ShutdownTimeout: time.Second*3,
}

// Start an HTTP server over UDS
err := s.Serve(ctx)
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
