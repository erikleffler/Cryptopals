package main

import(
	"bufio"
	"fmt"
	"os"
	"crypto/aes"

	crypt "../libcrypto"
)

func main() {

	var text string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text += scanner.Text()
	}

	raw_text, _ := crypt.B642Bytes(text)
	iv := make([]byte, 16)
	cipher, _ := aes.NewCipher([]byte("YELLOW SUBMARINE"))
	clear_text := crypt.CbcDecrypt(raw_text, iv, cipher)
	fmt.Println(string(clear_text))

	fmt.Println(crypt.Bytes2B64(crypt.CbcEncrypt(clear_text, iv, cipher)))

}
