package main

import (
	"fmt"
	"strings"
	"math/rand"
	"crypto/aes"
	"time"
	crypt "../libcrypto"
)

func main() {

	key := make([]byte, 16)
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Read(key)

	// Attacker
	cipher, err := aes.NewCipher(key)

	if err != nil {
		fmt.Println(err)
		return
	}

	p1 := []byte(strings.Repeat("A", 16))
	p2 := []byte(strings.Repeat("B", 16))
	p3 := []byte(strings.Repeat("C", 16))
	clearText := append(p1, p2...)
	clearText = append(clearText, p3...)

	// Server
	cipherText := crypt.CbcEncrypt(clearText, key, cipher)

	//Attacker
	c1 := cipherText[:16]
	modCipherText := append(c1, make([]byte, 16)...)
	modCipherText = append(modCipherText, c1...)

	//Server
	modClearText := crypt.CbcDecrypt(modCipherText, key, cipher)

	// Attacker
	retKey := crypt.Xor(modClearText[:16], modClearText[32:])
	fmt.Println(crypt.Bytes2Hex(key))
	fmt.Println(crypt.Bytes2Hex(retKey))
}
