package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"socks-rocketeerli/socks5"
)

var useTransport bool

func init() {
	flag.BoolVar(&useTransport, "t", true, "whether use transport or not")
	flag.Parse()
}

func main() {
	//	新建一个 `socks5`请求，测试是否能够访问成功
	socks, err := socks5.NewSocks5Dialer("http://127.0.0.1:7474")
	if err != nil {
		log.Fatalln("创建 socks5 对象失败 !!!")
	}
	/**
	两种方式利用 TCP Dial 访问：
	1. 将 TCP 的请求方法封装到 Transport 中，利用封装后的 http 客户端进行请求
	2. 在建立连接成功后，直接手动实现 http 请求，将 HTTP 的头部信息写到连接中
	*/
	if useTransport {
		client := http.Client{
			Transport: &http.Transport{
				Dial: socks.Dial,
			},
		}
		resp, err := client.Get("https://baidu.com")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))
		fmt.Println(useTransport)
	} else {
		conn, err := socks.Dial("tcp", "www.baidu.com:80")
		if err != nil {
			log.Fatalln("创建 socks5 连接失败 !!!")
			log.Fatal(err)
		}
		n, err := conn.Write([]byte("GET / HTTP/1.1\r\nHost: www.baidu.com\r\n\r\n"))
		log.Println(n, err)
		var b [4 * 64*1024]byte
		n, err = conn.Read(b[:])
		if err != nil {
			log.Fatal(err)
		}
		//TODO: 处理 Content-Length
		fmt.Println(string(b[:n]))
		fmt.Println(useTransport)
	}
}
