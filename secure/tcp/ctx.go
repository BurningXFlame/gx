package tcp

import (
	"context"
	"crypto/tls"

	"github.com/burningxflame/gx/ds/set"
	"github.com/burningxflame/gx/id/uuid"
)

func connId() string {
	return uuid.New()
}

var ctxKeyConnId int

func withConnId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, &ctxKeyConnId, id)
}

// Return the connection id. See Server.CtxConnId
func GetConnId(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(&ctxKeyConnId).(string)
	return id, ok
}

var ctxKeyPeer int

func withTlsPeer(ctx context.Context, conn *tls.Conn) context.Context {
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) < 1 {
		return ctx
	}

	cert := certs[0]
	if cert == nil {
		return ctx
	}

	s := set.New[string]()
	s.Add(cert.Subject.CommonName)
	for _, v := range cert.DNSNames {
		s.Add(v)
	}

	if s.Len() == 0 {
		return ctx
	}

	l := make([]string, 0, s.Len())
	s.ForEach(func(v string) {
		l = append(l, v)
	})

	return context.WithValue(ctx, &ctxKeyPeer, l)
}

// Return the identity of the TLS peer. See Server.CtxTlsPeer
func GetTlsPeer(ctx context.Context) ([]string, bool) {
	id, ok := ctx.Value(&ctxKeyPeer).([]string)
	return id, ok
}
