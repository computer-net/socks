package tools

type Cipher struct {
	//	编码时使用的密码
	encodePassword *password
	//	解码时使用的密码
	decodePassword *password
}

// Cipher 加密数据的方法
func (cipher *Cipher)Encode (bs []byte) {
	for i, v := range bs{
		bs[i] = cipher.encodePassword[v]
	}
}

// Cipher 解密数据的方法
func (cipher *Cipher)Decode(bs []byte) {
	for i, v := range bs{
		bs[i] = cipher.decodePassword[v]
	}
}

func NewCipher(encodePassword *password) *Cipher {
	decodePassword := &password{}
	for i, v := range encodePassword {
		decodePassword[v] = byte(i)
	}
	return &Cipher{
		encodePassword: encodePassword,
		decodePassword: decodePassword,
	}
}
