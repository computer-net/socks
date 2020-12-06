package cmd

import (
	"io"
	"log"
	"net"
	"time"
)

const (
	MaxPackSize               = 1024 * 1024 * 5 // 5Mb
	EchoServerAddr            = "127.0.0.1:3453"
	LightSocksProxyLocalAddr  = "127.0.0.1:8448"
	LightSocksProxyServerAddr = "127.0.0.1:8449"
)

var (
	lightsocksDialer proxy.Dialer
)

func init() {
	log.SetFlags(log.Lshortfile)
	go runEchoServer()
	go runLightsocksProxyServer()
	// 初始化代理socksDialer
	var err error
	// 等它们启动好
	time.Sleep(time.Second)
	lightsocksDialer, err = proxy.SOCKS5("tcp", LightSocksProxyLocalAddr, nil, proxy.Direct)
	if err != nil {
		log.Fatalln(err)
	}
}

// 启动echo server
func runEchoServer() {
	listener, err := net.Listen("tcp", EchoServerAddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		log.Println("echoServer connect Accept")
		go func() {
			defer func() {
				conn.Close()
				log.Println("echoServer connect Close")
			}()
			io.Copy(conn, conn)
		}()
	}
}
