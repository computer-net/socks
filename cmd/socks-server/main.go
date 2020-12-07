package main

import (
	"fmt"
	"log"
	"net"
	"socks-rocketeerli/cmd"
	"socks-rocketeerli/server"
)

var version = "master"

func main() {
	log.SetFlags(log.Lshortfile)
	// 默认配置
	config := &cmd.Config{}
	config.ReadConfig()
	config.SaveConfig()

	// 启动 server 端并监听
	rsServer, err := server.NewRsServer(config.Password, config.RemoteAddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(rsServer.Listen(func(listenAddr *net.TCPAddr) {
		log.Println(fmt.Sprintf(`
lightsocks-server:%s 启动成功，配置如下：
服务监听地址：
%s
密码：
%s`, version, listenAddr, config.Password))
	}))
}
