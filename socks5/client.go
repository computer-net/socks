package socks5

import (
	"errors"
	"io"
	"log"
	"net"
	"strconv"
)

// NewSocks5Dialer returns a socks5 proxy dialer.
func NewSocks5Dialer(s string) (*Socks5, error) {
	return NewSocks5(s)
}

func (s *Socks5) Addr() string {
	return s.addr
}

func (s *Socks5) Dial(network, addr string) (net.Conn, error) {
	switch network {
	case "tcp", "tcp6", "tcp4":
	default:
		return nil, errors.New("[socks5]: no support for connection type " + network)
	}
	c, err := net.Dial(network, s.addr)
	if err != nil {
		log.Printf("[socks5]: dial to %s error: %s", s.addr, err)
		return nil, err
	}
	if err := s.connect(c, addr); err != nil {
		c.Close()
		log.Printf("cnask %s", err)
		return nil, err
	}
	return c, nil
}

func (s *Socks5) connect(conn net.Conn, target string) error {
	// 获取主机名和端口号
	host, portStr, err := net.SplitHostPort(target)
	if err != nil {
		return err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return errors.New("proxy: failed to parse port number: " + portStr)
	}
	if port < 1 || port > 0xffff {
		return errors.New("proxy: port number out of range: " + portStr)
	}
	// 构造建立连接请求
	// the size here is just an estimate
	buf := make([]byte, 0, 6+len(host))
	/**
	客户端向服务端连接连接，客户端发送的数据包如下：
		   The localConn connects to the dstServer, and sends a ver
		   identifier/method selection message:
			          +----+----------+----------+
			          |VER | NMETHODS | METHODS  |
			          +----+----------+----------+
			          | 1  |    1     | 1 to 255 |
			          +----+----------+----------+
		   The VER field is set to X'05' for this ver of the protocol.  The
		   NMETHODS field contains the number of method identifier octets that
		   appear in the METHODS field.
	其中各个字段的含义如下：
	-VER：代表 SOCKS 的版本，SOCKS5 默认为0x05，其固定长度为1个字节；
	-NMETHODS：表示第三个字段METHODS的长度，它的长度也是1个字节；
	-METHODS：表示客户端支持的验证方式，可以有多种，他的长度是1-255个字节。
	*/
	// 第一个字段VER代表Socks的版本，Socks5默认为0x05，其固定长度为1个字节
	buf = append(buf, Version)
	// 这里仅支持两种验证方式：不需要验证和用户名密码验证。
	if len(s.user) > 0 && len(s.user) < 256 && len(s.password) < 256 {
		buf = append(buf, 2 /* num auth methods */, 0, 2)
	} else {
		buf = append(buf, 1 /* num auth methods */, 0)
	}
	// 发送建立连接的请求
	if _, err := conn.Write(buf); err != nil {
		return errors.New("proxy: failed to write greeting to SOCKS5 proxy at " + s.addr + ": " + err.Error())
	}
	/**
	服务端发送来的响应信息格式如下：
	   The dstServer selects from one of the methods given in METHODS, and
	   sends a METHOD selection message:
		          +----+--------+
		          |VER | METHOD |
		          +----+--------+
		          | 1  |   1    |
		          +----+--------+
	*/
	if _, err := io.ReadFull(conn, buf[:2]); err != nil {
		return errors.New("proxy: failed to read greeting from SOCKS5 proxy at " + s.addr + ": " + err.Error())
	}
	if buf[0] != Version {
		return errors.New("proxy: SOCKS5 proxy at " + s.addr + " has unexpected version " + strconv.Itoa(int(buf[0])))
	}
	if buf[1] == 0xff { // 验证方式为 0xff 时，表示 NO ACCEPTABLE METHODS（都不支持，没法连接了）
		return errors.New("proxy: SOCKS5 proxy at " + s.addr + " requires authentication")
	}
	// 如果是用户名和密码验证方式，需要验证用户名和密码
	if buf[1] == 2 {
		buf = buf[:0]
		buf = append(buf, 1 /* password protocol version */)
		buf = append(buf, uint8(len(s.user)))
		buf = append(buf, s.user...)
		buf = append(buf, uint8(len(s.password)))
		buf = append(buf, s.password...)

		if _, err := conn.Write(buf); err != nil {
			return errors.New("proxy: failed to write authentication request to SOCKS5 proxy at " + s.addr + ": " + err.Error())
		}

		if _, err := io.ReadFull(conn, buf[:2]); err != nil {
			return errors.New("proxy: failed to read authentication reply from SOCKS5 proxy at " + s.addr + ": " + err.Error())
		}

		if buf[1] != 0 {
			return errors.New("proxy: SOCKS5 proxy at " + s.addr + " rejected username/password")
		}
	}
	//	客户端发起的连接由服务端验证通过后，客户端下一步应该告诉真正目标服务的地址给服务器，服务器得到地址后再去请求真正的目标服务。
	/**
	客户端需要发送的包格式如下：
		  +----+-----+-------+------+----------+----------+
		  |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
		  +----+-----+-------+------+----------+----------+
		  | 1  |  1  | X'00' |  1   | Variable |    2     |
		  +----+-----+-------+------+----------+----------+
	VER：代表 SOCKS 协议的版本，SOCKS 默认为0x05，其值长度为1个字节；
	CMD：代表客户端请求的类型，值长度也是1个字节，有三种类型；其中 0x01 代表 CONNECT；
	RSV：保留字，值长度为1个字节；
	ATYP：代表请求的远程服务器地址类型：
		IPV4： address: 0x01；
		DOMAINNAME: 0x03；
		IPV6： address: 0x04；
	DST.ADDR：代表远程服务器的地址；
	DST.PORT：代表远程服务器的端口.
	*/
	buf = buf[:0]
	buf = append(buf, Version, 1, 0 /* reserved */)
	// 分析请求的远程服务器地址类型
	if ip := net.ParseIP(host); ip != nil {  // 是否可以解析成功
		if ip4 := ip.To4(); ip4 != nil {
			buf = append(buf, 1)
			ip = ip4
		} else {
			buf = append(buf, 4)
		}
		buf = append(buf, ip...)
	} else {
		if len(host) > 255 {
			return errors.New("proxy: destination hostname too long: " + host)
		}
		buf = append(buf, 3)
		buf = append(buf, byte(len(host)))
		buf = append(buf, host...)
	}
	buf = append(buf, byte(port>>8), byte(port))
	if _, err := conn.Write(buf); err != nil {
		return errors.New("proxy: failed to write connect request to SOCKS5 proxy at " + s.addr + ": " + err.Error())
	}
	/**
	服务端响应的信息格式：
	  +----+-----+-------+------+----------+----------+
	  |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	  +----+-----+-------+------+----------+----------+
	  | 1  |  1  | X'00' |  1   | Variable |    2     |
	  +----+-----+-------+------+----------+----------+
	*/
	if _, err := io.ReadFull(conn, buf[:4]); err != nil {
		return errors.New("proxy: failed to read connect reply from SOCKS5 proxy at " + s.addr + ": " + err.Error())
	}
	// 判断出错的类型
	failure := "unknown error"
	if int(buf[1]) < len(Errors) {
		failure = Errors[buf[1]].Error()
	}
	if len(failure) > 0 {
		return errors.New("proxy: SOCKS5 proxy at " + s.addr + " failed to connect: " + failure)
	}
	// 判断请求的类型
	bytesToDiscard := 0
	switch buf[3] {
	case 1:
		bytesToDiscard = net.IPv4len
	case 4:
		bytesToDiscard = net.IPv6len
	case 3:
		_, err := io.ReadFull(conn, buf[:1])
		if err != nil {
			return errors.New("proxy: failed to read domain length from SOCKS5 proxy at " + s.addr + ": " + err.Error())
		}
		bytesToDiscard = int(buf[0])
	default:
		return errors.New("proxy: got unknown address type " + strconv.Itoa(int(buf[3])) + " from SOCKS5 proxy at " + s.addr)
	}
	// 根据类型分配地址
	if cap(buf) < bytesToDiscard {
		buf = make([]byte, bytesToDiscard)
	} else {
		buf = buf[:bytesToDiscard]
	}
	// 读取响应的地址
	if _, err := io.ReadFull(conn, buf); err != nil {
		return errors.New("proxy: failed to read address from SOCKS5 proxy at " + s.addr + ": " + err.Error())
	}
	// 读取最后的端口号
	if _, err := io.ReadFull(conn, buf[:2]); err != nil {
		return errors.New("proxy: failed to read port from SOCKS5 proxy at " + s.addr + ": " + err.Error())
	}
	return nil
}
