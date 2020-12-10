package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"socks-rocketeerli/socks5"
)

func main() {
	//	新建一个 `socks5`请求，测试是否能够访问成功
	socks, err := socks5.NewSocks5Dialer("http://127.0.0.1:7474")
	if err != nil {
		log.Fatalln("创建 socks5 对象失败 !!!")
	}
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
}