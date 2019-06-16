package main

import (
	"fmt"
	"os"
	"bufio"
	"crypto/aes"

	crypt "../libcrypto"
)

func aes_ecb_decrypt(cipher_text []byte, key []byte) (clear_text []byte) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return
	}

	clear_text = make([]byte, len(cipher_text))
	for i := 0; i < len(cipher_text) / 16; i++ {
		cipher.Decrypt(clear_text[i*16:(i+1)*16], cipher_text[i*16:(i+1)*16])
	}
	return clear_text
}

func main () {
	file, err := os.Open("./7.txt")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	cipher := ""
	for scanner.Scan() {
		cipher += scanner.Text()
	}
	raw_cipher, err := crypt.B642Bytes(cipher)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(aes_ecb_decrypt(raw_cipher, []byte("YELLOW SUBMARINE"))))

}
