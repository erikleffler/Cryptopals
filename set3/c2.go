package main

import(
	"fmt"
	"crypto/aes"
	crypt "../libcrypto"
)

func main() {
	cipher, err := aes.NewCipher([]byte("YELLOW SUBMARINE"))
	if err != nil {
		fmt.Println(err)
		return
	}
	nonce := make([]byte, 8)
	cipherText, err := crypt.B642Bytes("L77na/nrFsKvynd6HzOoG7GHTLXsTVu9qvY/2syLXzhPweyyMTJULu/6/kXX0KSvoOLSFQ==")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(crypt.CtrDecrypt(cipherText, nonce, cipher)))

}
