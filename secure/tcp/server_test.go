package tcp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/burningxflame/gx/id/uuid"
	"github.com/burningxflame/gx/log/light"
	"github.com/burningxflame/gx/log/log"
	"github.com/burningxflame/gx/reliable/timeouts"
	"github.com/burningxflame/gx/sync/sem"
)

func TestServer(t *testing.T) {
	as := require.New(t)

	dur := time.Second / 10

	as.Nil(loadCerts())

	tcs := []testcase{
		{
			tag:     "Base",
			srv:     &Server{},
			cltSize: 3,
			cltOpt:  clientOpt{},
			cltOk:   3,
		},
		// ShutdownTimeout
		{
			tag:             "ShutdownTimeout/timeout",
			srv:             &Server{ShutdownTimeout: dur},
			cltSize:         3,
			cltOpt:          clientOpt{delayCloseConn: dur * 2},
			cltOk:           3,
			shutdownTimeout: true,
		},
		{
			tag:             "ShutdownTimeout/notTimeout",
			srv:             &Server{ShutdownTimeout: dur},
			cltSize:         3,
			cltOpt:          clientOpt{delayCloseConn: dur / 2},
			cltOk:           3,
			shutdownTimeout: false,
		},
		// ConnLimiter
		{
			tag:     "ConnLimiter/notExceeded",
			srv:     &Server{ConnLimiter: sem.New(3)},
			cltSize: 3,
			cltOpt:  clientOpt{timeout: dur, delayCloseConn: dur * 2},
			cltOk:   3,
		},
		{
			tag:     "ConnLimiter/exceeded",
			srv:     &Server{ConnLimiter: sem.New(3)},
			cltSize: 4,
			cltOpt:  clientOpt{timeout: dur, delayCloseConn: dur * 2},
			cltOk:   3,
		},
		{
			tag:     "ConnLimiter/released",
			srv:     &Server{ConnLimiter: sem.New(3)},
			cltSize: 9,
			cltOpt:  clientOpt{timeout: dur},
			cltOk:   9,
		},
		// IdleTimeout
		{
			tag:     "IdleTimeout/timeout",
			srv:     &Server{IdleTimeout: dur},
			cltSize: 3,
			cltOpt:  clientOpt{delay: dur * 2},
			cltOk:   0,
		},
		{
			tag:     "IdleTimeout/notTimeout",
			srv:     &Server{IdleTimeout: dur},
			cltSize: 3,
			cltOpt:  clientOpt{delay: dur / 2},
			cltOk:   3,
		},
		// TLS
		{
			tag: "TLS",
			srv: &Server{
				TlsConfig: &tls.Config{
					Certificates: srvCert,
					ClientCAs:    ca,
					ClientAuth:   tls.RequireAndVerifyClientCert,
				},
			},
			cltSize: 3,
			cltOpt: clientOpt{
				tlsConfig: &tls.Config{
					Certificates: cltCert,
					RootCAs:      ca,
					ServerName:   "DM-S",
				},
			},
			cltOk: 3,
		},
		{
			tag: "TLS/clientCertSignedByUnknownCA",
			srv: &Server{
				TlsConfig: &tls.Config{
					Certificates: srvCert,
					ClientCAs:    nil,
					ClientAuth:   tls.RequireAndVerifyClientCert,
				},
			},
			cltSize: 3,
			cltOpt: clientOpt{
				tlsConfig: &tls.Config{
					Certificates: cltCert,
					RootCAs:      ca,
					ServerName:   "DM-S",
				},
			},
			cltOk: 0,
		},
		{
			tag: "TLS/serverCertSignedByUnknownCA",
			srv: &Server{
				TlsConfig: &tls.Config{
					Certificates: srvCert,
					ClientCAs:    ca,
					ClientAuth:   tls.RequireAndVerifyClientCert,
				},
			},
			cltSize: 3,
			cltOpt: clientOpt{
				tlsConfig: &tls.Config{
					Certificates: cltCert,
					RootCAs:      nil,
					ServerName:   "DM-S",
				},
			},
			cltOk: 0,
		},
		{
			tag: "TLS/serverNameNotMatched",
			srv: &Server{
				TlsConfig: &tls.Config{
					Certificates: srvCert,
					ClientCAs:    ca,
					ClientAuth:   tls.RequireAndVerifyClientCert,
				},
			},
			cltSize: 3,
			cltOpt: clientOpt{
				tlsConfig: &tls.Config{
					Certificates: cltCert,
					RootCAs:      ca,
					ServerName:   "dummy",
				},
			},
			cltOk: 0,
		},
		// TlsHandshakeTimeout
		{
			tag: "TlsHandshakeTimeout/timeout",
			srv: &Server{
				TlsConfig: &tls.Config{
					Certificates: srvCert,
					ClientCAs:    ca,
					ClientAuth:   tls.RequireAndVerifyClientCert,
				},
				TlsHandshakeTimeout: dur,
			},
			cltSize: 3,
			cltOpt: clientOpt{
				tlsConfig: &tls.Config{
					Certificates: cltCert,
					RootCAs:      ca,
					ServerName:   "DM-S",
				},
				delay: dur,
			},
			cltOk: 0,
		},
		{
			tag: "TlsHandshakeTimeout/notTimeout",
			srv: &Server{
				TlsConfig: &tls.Config{
					Certificates: srvCert,
					ClientCAs:    ca,
					ClientAuth:   tls.RequireAndVerifyClientCert,
				},
				TlsHandshakeTimeout: dur,
			},
			cltSize: 3,
			cltOpt: clientOpt{
				tlsConfig: &tls.Config{
					Certificates: cltCert,
					RootCAs:      ca,
					ServerName:   "DM-S",
				},
			},
			cltOk: 3,
		},
		// CtxTlsPeer
		{
			tag: "CtxTlsPeer",
			srv: &Server{
				TlsConfig: &tls.Config{
					Certificates: srvCert,
					ClientCAs:    ca,
					ClientAuth:   tls.RequireAndVerifyClientCert,
				},
				CtxTlsPeer: true,
			},
			cltSize: 3,
			cltOpt: clientOpt{
				tlsConfig: &tls.Config{
					Certificates: cltCert,
					RootCAs:      ca,
					ServerName:   "DM-S",
				},
			},
			cltOk: 3,
		},
		// CtxConnId
		{
			tag:     "CtxConnId",
			srv:     &Server{CtxConnId: true},
			cltSize: 3,
			cltOpt:  clientOpt{},
			cltOk:   3,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.tag, func(t *testing.T) {
			test(t, tc)
		})
	}
}

type testcase struct {
	tag             string
	srv             *Server
	cltSize         int
	cltOpt          clientOpt
	cltOk           int
	shutdownTimeout bool
}

type clientOpt struct {
	delay          time.Duration
	timeout        time.Duration
	delayCloseConn time.Duration
	tlsConfig      *tls.Config
}

func test(t *testing.T, tc testcase) {
	light.InitTestLog()
	as := require.New(t)

	log.Info("test case: %v", tc.tag)

	addr := randAddr()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chSrv := make(chan error, 1)
	go func() {
		defer close(chSrv)

		tc.srv.Tag = "echo"
		tc.srv.Addr = addr
		tc.srv.ConnHandler = func(ctx context.Context, conn net.Conn) error {
			return handlConn(ctx, conn, tc)
		}
		chSrv <- tc.srv.Serve(ctx)
	}()
	time.Sleep(time.Millisecond * 10)

	var wg sync.WaitGroup
	chClt := make(chan error, tc.cltSize)

	for i := 0; i < tc.cltSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			chClt <- client(addr, tc.cltOpt)
		}()
	}

	wg.Wait()
	close(chClt)

	cancel()
	err := <-chSrv
	if tc.shutdownTimeout {
		as.ErrorIs(err, timeouts.ErrTimeout)
	}

	okCnt := 0
	for err := range chClt {
		if err == nil {
			okCnt++
		}
	}
	as.Equal(tc.cltOk, okCnt)
}

func randAddr() string {
	port := 1024 + rand.Intn(math.MaxInt8)
	return fmt.Sprintf("127.0.0.1:%v", port)
}

func handlConn(ctx context.Context, conn net.Conn, tc testcase) error {
	_, ok := GetConnId(ctx)
	if tc.srv.CtxConnId != ok {
		return errDummy
	}

	_, ok = GetTlsPeer(ctx)
	if tc.srv.CtxTlsPeer != ok {
		return errDummy
	}

	_, err := io.Copy(conn, conn)
	return err
}

func client(addr string, tc clientOpt) error {
	return timeouts.WithTimeout(tc.timeout, func() error {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return err
		}
		defer func() {
			if tc.delayCloseConn <= 0 {
				conn.Close()
				return
			}

			go func() {
				time.Sleep(tc.delayCloseConn)
				conn.Close()
			}()
		}()

		if tc.delay > 0 {
			time.Sleep(tc.delay)
		}

		if tc.tlsConfig != nil {
			conn = tls.Client(conn, tc.tlsConfig)
		}

		err = comm(conn)
		if err != nil {
			return err
		}

		return nil
	})()
}

func comm(conn net.Conn) error {
	for i := 0; i < 3; i++ {
		msg := uuid.New()
		buf := make([]byte, len(msg))

		_, err := conn.Write([]byte(msg))
		if err != nil {
			return err
		}

		_, err = io.ReadFull(conn, buf)
		if err != nil {
			return err
		}

		if msg != string(buf) {
			return errDummy
		}
	}

	return nil
}

var errDummy = errors.New("dummy")

var (
	ca      *x509.CertPool
	srvCert []tls.Certificate
	cltCert []tls.Certificate
)

func loadCerts() error {
	ca = x509.NewCertPool()
	certPem, err := os.ReadFile("testdata/ca.crt")
	if err != nil {
		return err
	}
	ok := ca.AppendCertsFromPEM(certPem)
	if !ok {
		return errDummy
	}

	cert, err := tls.LoadX509KeyPair("testdata/server.crt", "testdata/server.key")
	if err != nil {
		return err
	}
	srvCert = []tls.Certificate{cert}

	cert, err = tls.LoadX509KeyPair("testdata/client.crt", "testdata/client.key")
	if err != nil {
		return err
	}
	cltCert = []tls.Certificate{cert}

	return nil
}
