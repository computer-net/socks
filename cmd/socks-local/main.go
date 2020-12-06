package main

import (
	"log"
	"strconv"
)
import "socks-rocketeerli/tools"

func main() {
	pwBase64 := tools.RandPassword()
	log.Println(pwBase64)
	pwByte, err := tools.ParsePassword(pwBase64)
	if err!=nil {
		log.Println("...Error")
	}
	for i, v := range pwByte {
		log.Println(strconv.Itoa(i) + "\t" + strconv.Itoa(int(v)))
	}
}