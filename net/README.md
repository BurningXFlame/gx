# net

- [Connection Pool](#connection-pool)
- [SOCKS5](#socks5)

## Connection Pool

[ConnPool](connpool/pool.go) is a concurrency-safe connection pool.

```go
import "github.com/burningxflame/gx/net/connpool"

// Create a connection pool
pool, err := connpool.New(connpool.Conf{
  // Initial number of connections in the pool
  Init: 8,
  // Max number of connections in the pool
  Cap:  16,
  // Used to create a connection
  New:  func() (net.Conn, error) {...},
  // Used to check whether a connection is still connected. Ping returns nil if still connected.
  Ping: func(conn net.Conn) error {...},
  // Timeout of New and Ping. Abandon the call to New or Ping after the timeout elapses.
  Timeout time.Duration
})

// Get a connection from the pool.
// If the connection is not connected any more, will drop it and get another, until a connected connection is found.
// If no connection in the pool is connected, will create a new one.
conn, err := pool.Get()

// Put a connection back to the pool.
// Will close and drop the connection if the pool is full.
// Nil connection is ignored.
pool.Put(conn)

// Close all connections in the pool
pool.Close()
```

## SOCKS5

[SOCKS](socks/client.go) is a client-side implementation of the SOCKS5 proxy protocol.

```go
import (
    "net"
    "github.com/burningxflame/gx/net/socks"
)

// Connect to SOCKS5 proxy
conn, err := net.Dial("tcp", socksAddr)
if err != nil {
    return err
}

// Client-side handshake of SOCKS5 proxy protocol.
// The destAddr is the destination address to connect to through SOCKS5 proxy.
err = socks.ClientHandshake(conn, destAddr)
```
