/*
GX (github.com/burningxflame/gx).
Copyright © 2022-2024 BurningXFlame. All rights reserved.

Dual-licensed: AGPLv3/Commercial.
Read the LICENSE file for details.
*/

package socks

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"strconv"
)

const (
	atypIpv4   byte = 1
	atypIpv6   byte = 4
	atypDomain byte = 3

	ipv4Size      = 4
	ipv6Size      = 16
	domainMaxSize = math.MaxUint8

	portSize = 2
)

var (
	errDomainOversize = errors.New("domain name oversize")
	errInvalidAddr    = errors.New("invalid address")
	ErrInvalidAtyp    = errors.New("invalid atyp")
)

type Addr struct {
	host host
	port uint16
}

var addr0 Addr

func parseAddr(dest string) (Addr, error) {
	hostStr, portStr, err := net.SplitHostPort(dest)
	if err != nil {
		return addr0, fmt.Errorf("invalid addr: %v", err)
	}

	host, err := toHost(hostStr)
	if err != nil {
		return addr0, err
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return addr0, fmt.Errorf("invalid port: %v", err)
	}

	return Addr{
		host: host,
		port: uint16(port),
	}, nil
}

func (a *Addr) send(conn net.Conn, buf []byte) error {
	err := a.host.send(conn, buf)
	if err != nil {
		return err
	}

	return binary.Write(conn, binary.BigEndian, a.port)
}

func (a *Addr) Recv(conn net.Conn, buf []byte) error {
	_, err := io.ReadFull(conn, buf[:1])
	if err != nil {
		return err
	}

	atyp := buf[0]
	switch atyp {
	case atypIpv4:
		a.host = &ipv4{}
	case atypIpv6:
		a.host = &ipv6{}
	case atypDomain:
		a.host = &domain{}
	default:
		return ErrInvalidAtyp
	}

	err = a.host.recv(conn, buf)
	if err != nil {
		return err
	}

	return binary.Read(conn, binary.BigEndian, &a.port)
}

func (a *Addr) String() string {
	return net.JoinHostPort(a.host.String(), strconv.Itoa(int(a.port)))
}

type host interface {
	String() string
	send(net.Conn, []byte) error
	recv(net.Conn, []byte) error
}

func toHost(s string) (host, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		rs, err := toDomain(s)
		if err != nil {
			return nil, err
		}

		return rs, nil
	}

	if rs := toIpv4(ip); rs != nil {
		return rs, nil
	}

	if rs := toIpv6(ip); rs != nil {
		return rs, nil
	}

	return nil, errInvalidAddr
}

type ipv4 struct {
	v [ipv4Size]byte
}

func toIpv4(v net.IP) *ipv4 {
	_v := v.To4()
	if _v == nil {
		return nil
	}

	var rs ipv4
	copy(rs.v[:], _v)
	return &rs
}

func (a *ipv4) String() string {
	return net.IP(a.v[:]).String()
}

func (a *ipv4) send(conn net.Conn, buf []byte) error {
	buf[0] = atypIpv4
	copy(buf[1:], a.v[:])

	_, err := conn.Write(buf[:1+ipv4Size])
	return err
}

func (a *ipv4) recv(conn net.Conn, _ []byte) error {
	_, err := io.ReadFull(conn, a.v[:])
	return err
}

type ipv6 struct {
	v [ipv6Size]byte
}

func toIpv6(v net.IP) *ipv6 {
	_v := v.To16()
	if _v == nil {
		return nil
	}

	var rs ipv6
	copy(rs.v[:], _v)
	return &rs
}

func (a *ipv6) String() string {
	return net.IP(a.v[:]).String()
}

func (a *ipv6) send(conn net.Conn, buf []byte) error {
	buf[0] = atypIpv6
	copy(buf[1:], a.v[:])

	_, err := conn.Write(buf[:1+ipv6Size])
	return err
}

func (a *ipv6) recv(conn net.Conn, _ []byte) error {
	_, err := io.ReadFull(conn, a.v[:])
	return err
}

type domain struct {
	v string
}

func toDomain(s string) (*domain, error) {
	if len(s) > domainMaxSize {
		return nil, errDomainOversize
	}

	return &domain{s}, nil
}

func (a *domain) String() string {
	return a.v
}

func (a *domain) send(conn net.Conn, buf []byte) error {
	size := len(a.v)
	if size > domainMaxSize {
		return errDomainOversize
	}

	buf[0] = atypDomain
	buf[1] = byte(size)
	copy(buf[2:], a.v)

	_, err := conn.Write(buf[:2+size])
	return err
}

func (a *domain) recv(conn net.Conn, buf []byte) error {
	_, err := io.ReadFull(conn, buf[:1])
	if err != nil {
		return err
	}

	size := buf[0]
	_, err = io.ReadFull(conn, buf[:size])
	if err != nil {
		return err
	}

	a.v = string(buf[:size])
	return nil
}

func discardAddr(conn net.Conn, buf []byte) error {
	_, err := io.ReadFull(conn, buf[:1])
	if err != nil {
		return err
	}
	atyp := buf[0]

	var hostSize int
	switch atyp {
	case atypIpv4:
		hostSize = ipv4Size
	case atypIpv6:
		hostSize = ipv6Size
	case atypDomain:
		_, err := io.ReadFull(conn, buf[:1])
		if err != nil {
			return err
		}
		hostSize = int(buf[0])
	default:
		return ErrInvalidAtyp
	}

	_, err = io.ReadFull(conn, buf[:hostSize+portSize])
	return err
}
