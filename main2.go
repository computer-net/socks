package main

import (
	"fmt"
	"log"
	"socks-rocketeerli/socks5"
)

func main() {
	//	新建一个 `socks5`请求，测试是否能够访问成功
	socks, err := socks5.NewSocks5Dialer("http://127.0.0.1:7474")
	if err != nil {
		log.Fatalln("创建 socks5 对象失败 !!!")
	}
	conn, err := socks.Dial("tcp", "www.baidu.com:80")
	if err != nil {
		log.Fatal(err)
	}
	n, err := conn.Write([]byte("GET / HTTP/1.1\r\nHost: www.baidu.com\r\n\r\n"))
	log.Println(n, err)
	var b [64 * 1024]byte
	n, err = conn.Read(b[:])
	if err != nil {
		log.Fatal(err)
	}
	//TODO: 处理Content-Length
	fmt.Println(string(b[:n]))
}
