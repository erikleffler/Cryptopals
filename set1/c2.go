package main

import (
	"fmt"
	crypt "../libcrypto"
)

func main() {

	str1, err := crypt.Hex2Bytes("1c0111001f010100061a024b53535009181c")
	if err != nil{
		fmt.Println(err)
		return
	}
	str2, err := crypt.Hex2Bytes("686974207468652062756c6c277320657965")
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println(crypt.Bytes2Hex(crypt.Xor(str1, str2)))

}
