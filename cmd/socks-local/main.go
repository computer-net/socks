package main

import (
	"fmt"
	"log"
	"net"
	"socks-rocketeerli/cmd"
	"socks-rocketeerli/local"
)

const (
	DefaultListenAddr = ":7448"
)

var version = "master"

func main() {
	// 输出文件名+行号格式
	log.SetFlags(log.Lshortfile)
	// 默认配置
	config := &cmd.Config{
		ListenAddr: DefaultListenAddr,
	}
	config.ReadConfig()
	config.SaveConfig()
	// 启动 local 端并监听
	rsLocal, err := local.NewRsLocal(config.Password, config.ListenAddr, config.RemoteAddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(rsLocal.Listen(func(listenAddr *net.TCPAddr) {
		log.Println(fmt.Sprintf(`
lightsocks-local:%s 启动成功，配置如下：
本地监听地址：
%s
远程服务地址：
%s
密码：
%s`, version, listenAddr, config.RemoteAddr, config.Password))
	}))
}