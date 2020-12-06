package tools

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"strings"
	"time"
)

const passwordLength = 256

type password [passwordLength]byte

func init()  {
	// 设置随机种子，防止每次生成相同的随机数
	rand.Seed(time.Now().Unix())
}

func (password *password) String() string {
	return base64.StdEncoding.EncodeToString(password[:])
}

func ParsePassword(passwordString string) (*password, error) {
	bs, err := base64.StdEncoding.DecodeString(strings.TrimSpace(passwordString))
	if err != nil || len(bs) != passwordLength {
		return nil, errors.New("密码不合法，解析失败！！！")
	}
	pw := password{}
	copy(pw[:], bs)  // 复制数组切片
	return &pw, nil
}

// 生成密码，由 0-255 共256个数组成的排列，每位数值不重复且必须保证每个数值与其下标位不相等
func RandPassword() string {
	//	返回一个有n个元素的，[0,n)范围内整数的伪随机排列的切片
	arr := rand.Perm(passwordLength)
	pw := &password{}
	for i, v := range arr {
		pw[i] = byte(v)
		if i == v {
			return RandPassword()
		}
	}
	return pw.String()
}
