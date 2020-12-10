package main

import (
	"bytes"
	"fmt"
	"io"
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
		log.Fatalln("建立 socks5 连接失败 !!!")
		log.Fatalln(err)
	}
	result := bytes.NewBuffer(nil)
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
	}
	fmt.Println(string(result.Bytes()))
}
