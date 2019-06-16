package main

import (
	"fmt"
	crypt "../libcrypto"
)

func main() {

	text := []byte("Burning 'em, if you ain't quick and nimble I go crazy when I hear a cymbal")
	key := []byte("ICE")
	fmt.Println(crypt.Bytes2Hex(crypt.Xor(text, key)))

}
