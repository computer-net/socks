package cmd

import (
	"io"
	"log"
	"math/rand"
	"net"
	"reflect"
	"socks-rocketeerli/local"
	"socks-rocketeerli/server"
	"socks-rocketeerli/tools"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/proxy"
)

const (
	MaxPackSize               = 1024 * 1024 * 5 // 5Mb
	EchoServerAddr            = "127.0.0.1:3453"
	RSocksProxyLocalAddr  = "127.0.0.1:8448"
	RSocksProxyServerAddr = "127.0.0.1:8449"
)

var (
	rsocksDialer proxy.Dialer
)

func init() {
	log.SetFlags(log.Lshortfile)
	go runEchoServer()
	go runRSocksProxyServer()
	// 初始化代理socksDialer
	var err error
	// 等它们启动好
	time.Sleep(time.Second)
	rsocksDialer, err = proxy.SOCKS5("tcp", RSocksProxyLocalAddr, nil, proxy.Direct)
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

func runRSocksProxyServer() {
	password := tools.RandPassword()
	serverS, err := local.NewRsLocal(password, RSocksProxyLocalAddr, RSocksProxyServerAddr)
	if err != nil {
		log.Fatalln(err)
	}
	localS, err := server.NewRsServer(password, RSocksProxyServerAddr)
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		log.Fatalln(serverS.Listen(func(listenAddr *net.TCPAddr) {
			log.Println(listenAddr)
		}))
	}()
	log.Fatalln(localS.Listen(func(listenAddr *net.TCPAddr) {
		log.Println(listenAddr)
	}))
}

// 发生一次连接测试经过代理后的数据传输的正确性
// packSize 代表这个连接发生数据的大小
func testConnect(packSize int) {
	// 随机生产 MaxPackSize byte的[]byte
	data := make([]byte, packSize)
	_, err := rand.Read(data)

	// 连接
	conn, err := rsocksDialer.Dial("tcp", EchoServerAddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// 写
	go func() {
		conn.Write(data)
	}()

	// 读
	buf := make([]byte, len(data))
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		log.Fatalln(err)
	}
	if !reflect.DeepEqual(data, buf) {
		log.Fatalln("通过 Rsocks 代理传输得到的数据前后不一致")
	} else {
		log.Println("数据一致性验证通过")
	}
}

func TestLightsocks(t *testing.T) {
	testConnect(rand.Intn(MaxPackSize))
}

// 获取并发发送 data 到 echo server 并且收到全部返回 所花费到时间
func benchmarkRSocks(concurrenceCount int) {
	wg := sync.WaitGroup{}
	wg.Add(concurrenceCount)
	for i := 0; i < concurrenceCount; i++ {
		go func() {
			testConnect(rand.Intn(MaxPackSize))
			wg.Done()
		}()
	}
	wg.Wait()
}

// 获取 发送 data 到 echo server 并且收到全部返回 所花费到时间
func BenchmarkRSocks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		benchmarkRSocks(10)
		b.StopTimer()
	}
}
