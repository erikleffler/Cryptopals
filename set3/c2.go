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
	
	nonce = make([]byte, 8)
	cipherText = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")

	fmt.Println(string(crypt.CtrDecrypt(crypt.CtrEncrypt(cipherText, nonce, cipher), nonce, cipher)))

}
