/*
GX (https://github.com/BurningXFlame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package socks

import (
	"fmt"
	"io"
	"net"
)

const (
	ver uint8 = 5

	authNone          uint8 = 0
	authNotAcceptable uint8 = 0xff

	cmdConnect uint8 = 1

	rsv uint8 = 0

	repOk uint8 = 0
)

// Client-side handshake of SOCKS5 proxy protocol.
// The destAddr is the destination address to connect to through SOCKS5 proxy.
func ClientHandshake(conn net.Conn, destAddr string) error {
	buf := *getBuf()
	defer putBuf(&buf)

	dest, err := parseAddr(destAddr)
	if err != nil {
		return err
	}

	err = auth(conn, buf)
	if err != nil {
		return err
	}

	err = sendAddr(conn, buf, dest)
	if err != nil {
		return err
	}

	return recvResp(conn, buf)
}

func auth(conn net.Conn, buf []byte) error {
	// send protocol version and auth mothods
	buf[0], buf[1], buf[2] = ver, 1, authNone
	_, err := conn.Write(buf[:3])
	if err != nil {
		return err
	}

	// recv protocol version and auth mothod
	_, err = io.ReadFull(conn, buf[:2])
	if err != nil {
		return err
	}

	verServer, authServer := buf[0], buf[1]

	if verServer != ver {
		return fmt.Errorf("version not match, client: %v, server: %v", ver, verServer)
	}

	if authServer != authNone {
		return fmt.Errorf("auth method not match, client: %v, server: %v", authNone, authServer)
	}

	return nil
}

func sendAddr(conn net.Conn, buf []byte, dest addr) error {
	buf[0], buf[1], buf[2] = ver, cmdConnect, rsv
	_, err := conn.Write(buf[:3])
	if err != nil {
		return err
	}

	return dest.send(conn, buf)
}

func recvResp(conn net.Conn, buf []byte) error {
	// recv rep code
	_, err := io.ReadFull(conn, buf[:3])
	if err != nil {
		return err
	}

	rep := buf[1]
	if rep != repOk {
		return RepErr{rep}
	}

	// discard what's left
	return discardAddr(conn, buf)
}

type RepErr struct {
	Rep byte
}

func (e RepErr) Error() string {
	return fmt.Sprintf("rep 0x%x", e.Rep)
}
