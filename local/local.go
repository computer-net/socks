package local

import (
	"log"
	"net"
	"socks-rocketeerli/tools"
)

type RsLocal struct {
	Cipher *tools.Cipher
	ListenAddr *net.TCPAddr
	RemoteAddr *net.TCPAddr
}

/**
本地端的职责：
1. 监听来自用户本机浏览器的代理请求
2. 转发前，加密数据
3. 转发 socket 数据到墙外代理服务端
4. 接收服务端返回的数据，并转发给用户的浏览器
*/
// 根据密码+本地地址+远程地址，新建一个本地端
func NewRsLocal(password string, listenAddr, remoteAddr string) (*RsLocal, error) {
	// 解析字符串，生成密码
	pw, err := tools.ParsePassword(password)
	if err != nil {
		return nil, err
	}
	//	解析本机地址和远程地址
	lsAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	rmAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		return nil, err
	}
	return &RsLocal{
		Cipher: tools.NewCipher(pw),
		ListenAddr: lsAddr,
		RemoteAddr: rmAddr,
	}, nil
}

func (local *RsLocal)handleConn(userConn *tools.SecureTCPConn) {
	defer userConn.Close()
	proxyServer, err := tools.DialEncryptedTCP(local.RemoteAddr, local.Cipher)
	if err != nil {
		log.Println(err)
		return
	}
	defer proxyServer.Close()
	go func() {
		err := proxyServer.DecodeCopy(userConn)
		if err != nil {
			// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
			userConn.Close()
			proxyServer.Close()
		}
	}()
	// 从 localUser 发送数据发送到 proxyServer，这里因为处在翻墙阶段出现网络错误的概率更大
	userConn.EncodeCopy(proxyServer)
}

func (local *RsLocal) Listen(didListen func(listenAddr *net.TCPAddr)) error {
	return tools.ListenEncryptedTCP(local.ListenAddr, local.Cipher, local.handleConn, didListen)
}
