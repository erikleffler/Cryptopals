package main

import (
	"fmt"
	crypt "../libcrypto"
)

func main() {
	text, err := crypt.Pkcs7Unpad([]byte("ICE ICE BABY\x04\x04\x04\x04"), 16)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(text)
	}
	text, err = crypt.Pkcs7Unpad([]byte("ICE ICE BABY\x05\x05\x05\x05"), 16)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(text)
	}
	text, err = crypt.Pkcs7Unpad([]byte("ICE ICE BABY\x01\x02\x03\x04"), 16)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(text)
	}
}
