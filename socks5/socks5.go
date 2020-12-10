package socks5

import (
	"errors"
	"log"
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


type Socks5 struct {
	addr     string
	user     string
	password string
}

func NewSocks5(s string) (*Socks5, error) {
	u, err := url.Parse(s)
	if err != nil {
		log.Printf("Parse url err: %s", err)
	}
	addr := u.Host
	user := u.User.Username()
	pass, _ := u.User.Password()
	return &Socks5{
		addr:     addr,
		user:     user,
		password: pass,
	}, nil
}
