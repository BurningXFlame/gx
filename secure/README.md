# Secure

Provides utilities to improve application security.

- [Idle Timeout](#idle-timeout)
- [Secure TCP Server](#secure-tcp-server)

## Idle Timeout

```go
import "github.com/burningxflame/gx/secure/conns"

// If no data has been sent from conn by the time the timeout period elapses, conn.Read will return ErrIdleTimeout.
// Zero timeout means no timeout.
conn, err := conns.WithIdleTimeout(conn, timeout)
```

## Secure TCP Server

The Secure TCP Server has builtin abilities to defend against DDoS attacks, to defend against slow attacks, to close idle connections, etc.

```go
import (
  "github.com/burningxflame/gx/secure/tcp"
  "github.com/burningxflame/gx/sync/sem"
)

// Create a Secure TCP Server
srv := &tcp.Server{
  // The address to listen
  Addr: "host:port",
  // Connection handler is where you communicate with a client, i.e. receive/send data from/to a client.
  // Connection handler should return ASAP when ctx.Done channel is closed, which usually means an exit signal is sent.
  ConnHandler: func(ctx context.Context, conn net.Conn) error { ... },
  // If graceful shutdown takes longer than ShutdownTimeout, exit instantly.
  // Default to no timeout.
  ShutdownTimeout: time.Second*3,
  // Used to limit max number of concurrent connections
  // Default to no limit.
  ConnLimiter: sem.New(n),
  // If no data is sent from a connection in the specified duration, close the connection.
  // Default to no timeout.
  IdleTimeout: time.Minute,
  // Used for TLS handshake. If not provided, no TLS handshake.
  TlsConfig *tls.Config
  // If TLS handshake does not finish in the specified duration, close the connection.
  // Default to no timeout.
  TlsHandshakeTimeout: time.Second*5
  // If true, the Context argument of ConnHandler contains the identity of the TLS peer.
  // Call GetTlsPeer(ctx) to get peer identity.
  // And of course the peer should send a certificate, i.e. TlsConfig.ClientAuth should be RequireAnyClientCert or RequireAndVerifyClientCert.
  // Peer identity is a set of Common Name and SAN DNS Names of certificate holder, i.e. cert.Subject.CommonName and cert.DNSNames.
  CtxTlsPeer: false
  // If true, the Context argument of ConnHandler contains the connection id.
  // Call GetConnId(ctx) to get connection id.
  CtxConnId: false
  // Used to tag log messages
  Tag: "someTag",
  // A TagLogger used to log messages
  Log: ...,
}

// Start the Server
err := s.Serve(ctx)
```

**Connection Handler** is where you communicate with a client, i.e. receive/send data from/to a client. Connection handler should return ASAP when ctx.Done channel is closed, which usually means an exit signal is sent.

```go
func handleConn(ctx context.Context, conn net.Conn) error {
  ...

  // Return the identity of the TLS peer. See Server.CtxTlsPeer
  peer, ok := tcp.GetTlsPeer(ctx)
  ...

  // Return the connection id. See Server.CtxConnId
  connId, ok := tcp.GetConnId(ctx)
  ...
}
```
