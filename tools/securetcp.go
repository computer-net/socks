package tools

import "io"

const bufSize =  1024

// 加密传输的 TCP Socket
type SecureTCPConn struct {
	io.ReadWriteCloser
	Cipher *Cipher
}

// 对来自输入流中的数据进行加密
func (secureSocket *SecureTCPConn)DecodeRead(bs []byte) (n int, err error) {
	n, err = secureSocket.Read(bs)
	if err != nil {
		return
	}
	secureSocket.Cipher.Decode(bs[:n])
	return
}

// 将bs中的数据加密后写入输出流
func (secureSocket *SecureTCPConn)EncodeWrite(bs []byte) (int, error) {
	secureSocket.Cipher.Encode(bs)
	return secureSocket.Write(bs)
}

// 从 src 中不断读取数据，加密后，写到 dst 中
func (secureSocket *SecureTCPConn) EncodeCopy(dst io.ReadWriteCloser) error {
	buf := make([]byte, bufSize)
	for {
		// 从 src 中读取原始数据
		readNum, readErr := secureSocket.Read(buf)
		if readErr != nil {
			if readErr != io.EOF {
				return readErr
			} else {
				return nil
			}
		}
		if readNum > 0 {  // 如果读取到数据，就对其进行加密，写入到 dst 中
			writeNum, writeErr := (&SecureTCPConn{
				ReadWriteCloser: dst,
				Cipher: secureSocket.Cipher,
			}).EncodeWrite(buf[:readNum])
			if writeErr != nil {
				return writeErr
			}
			// 如果读取的数据和加密后的数据长度不等，代表出现错误
			if readNum != writeNum {
				return io.ErrShortWrite
			}
		}
	}
}

// 从 src 中持续读取加密后的数据，解密后，写入到 dst 中
func (secureSocket *SecureTCPConn) DecodeCopy (dst io.Writer) error {
	buf := make([]byte, bufSize)
	for {
		// 读取加密后的数据
		readNum, errRead := secureSocket.DecodeRead(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			} else {
				return nil
			}
		}
		if readNum > 0 {
			writeNum, writeErr := dst.Write(buf[:readNum])
			if writeErr != nil {
				return writeErr
			}
			// 如果读取的数据和解密后的数据长度不等，代表出现错误
			if readNum != writeNum {
				return io.ErrShortWrite
			}
		}
	}
}
