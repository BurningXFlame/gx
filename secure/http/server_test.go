package http

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
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

	as.Nil(loadCerts())

	const dur = time.Second / 10

	tcs := []testcase{
		{
			tag:     "Base",
			srv:     &Server{},
			srvOpt:  serverOpt{},
			cltSize: 3,
			cltOpt:  clientOpt{},
			cltOk:   3,
		},
		// Limiter
		{
			tag:     "Limiter/notExceeded",
			srv:     &Server{Limiter: sem.New(3)},
			srvOpt:  serverOpt{delay: dur},
			cltSize: 3,
			cltOpt:  clientOpt{timeout: dur * 3 / 2},
			cltOk:   3,
		},
		{
			tag:     "Limiter/exceeded",
			srv:     &Server{Limiter: sem.New(3)},
			srvOpt:  serverOpt{delay: dur},
			cltSize: 4,
			cltOpt:  clientOpt{timeout: dur * 3 / 2},
			cltOk:   3,
		},
		{
			tag:     "Limiter/released",
			srv:     &Server{Limiter: sem.New(3)},
			srvOpt:  serverOpt{},
			cltSize: 9,
			cltOpt:  clientOpt{timeout: dur * 3 / 2},
			cltOk:   9,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.tag, func(t *testing.T) {
			test(t, tc)
		})
	}
}

type testcase struct {
	tag     string
	srv     *Server
	srvOpt  serverOpt
	cltSize int
	cltOpt  clientOpt
	cltOk   int
}

type serverOpt struct {
	delay time.Duration
}

type clientOpt struct {
	timeout time.Duration
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
		tc.srv.Std.Addr = addr
		tc.srv.Std.Handler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			time.Sleep(tc.srvOpt.delay)
			handleEcho(rw, r)
		})
		tc.srv.Std.TLSConfig = &tls.Config{
			Certificates: srvCert,
			ClientCAs:    ca,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}
		tc.srv.Serve(ctx)
	}()
	time.Sleep(time.Millisecond * 10)

	var wg sync.WaitGroup
	chClt := make(chan error, tc.cltSize)

	for i := 0; i < tc.cltSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			chClt <- req(addr, tc.cltOpt)
		}()
	}

	wg.Wait()
	close(chClt)

	cancel()

	okCnt := 0
	for err := range chClt {
		if err == nil {
			okCnt++
		} else {
			log.Debug("%v", err)
		}
	}
	as.Equal(tc.cltOk, okCnt)
}

func randAddr() string {
	port := 1024 + rand.Intn(math.MaxInt8)
	return fmt.Sprintf("127.0.0.1:%v", port)
}

func handleEcho(w http.ResponseWriter, r *http.Request) {
	io.Copy(w, r.Body)
}

func req(addr string, tc clientOpt) error {
	return timeouts.WithTimeout(tc.timeout, func() error {
		return _req(addr, tc)
	})()
}

func _req(addr string, _ clientOpt) error {
	c := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: cltCert,
				RootCAs:      ca,
				ServerName:   "DM-S",
			},
		},
		Timeout: time.Second * 3,
	}

	msg := uuid.New()
	resp, err := c.Post("https://"+addr+"/echo", "raw", strings.NewReader(msg))
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if msg != string(body) {
		log.Debug("msg: %s, body: %s", msg, body)
		return errDummy
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
