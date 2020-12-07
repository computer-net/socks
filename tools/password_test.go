package tools

import (
	"log"
	"strconv"
	"testing"
)

func TestGenPassword(t *testing.T) {
	pwBase64 := RandPassword()
	log.Println(pwBase64)
	pwByte, err := ParsePassword(pwBase64)
	if err!=nil {
		log.Println("...Error")
	}
	for i, v := range pwByte {
		if i == int(v) {
			log.Println("generate password failed!!!")
			log.Println(strconv.Itoa(i) + "\t" + strconv.Itoa(int(v)))
			return
		}
	}
	log.Println("generate password successfully!!!")
}
