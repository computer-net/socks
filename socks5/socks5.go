package socks5

import (
	"errors"
	"log"
	"net"
	"net/url"
)

// Version is socks5 version number.
const Version = 5

// Errors are socks5 errors
var Errors = []error{
	errors.New(""),
	errors.New("general failure"),
	errors.New("connection forbidden"),
	errors.New("network unreachable"),
	errors.New("host unreachable"),
	errors.New("connection refused"),
	errors.New("TTL expired"),
	errors.New("command not supported"),
	errors.New("address type not supported"),
	errors.New("socks5UDPAssociate"),
}

type Dialer interface {
	// Addr is the dialer's addr
	Addr() string

	// Dial connects to the given address
	Dial(network, addr string) (c net.Conn, err error)
}

type Proxy interface {
	// Dial connects to the given address via the proxy.
	Dial(network, addr string) (c net.Conn, dialer Dialer, err error)
	// Get the dialer by dstAddr.
	NextDialer(dstAddr string) Dialer
}

type Socks5 struct {
	dialer   Dialer
	proxy    Proxy
	addr     string
	user     string
	password string
}

func NewSocks5(s string, d Dialer, p Proxy) (*Socks5, error) {
	u, err := url.Parse(s)
	if err != nil {
		log.Printf("Parse url err: %s", err)
	}
	addr := u.Host
	user := u.User.Username()
	pass, _ := u.User.Password()
	return &Socks5{
		dialer:   d,
		proxy:    p,
		addr:     addr,
		user:     user,
		password: pass,
	}, nil
}
