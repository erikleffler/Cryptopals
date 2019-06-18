package main

import (
	"fmt"
	"net/http"
	"math/rand"
	"crypto/aes"
	"log"
	"encoding/json"


	crypt "../libcrypto"
)

type data struct {
	Secret string
}

func random_encrypt(clear_text []byte) []byte {

	rand_key := make([]byte, 16)
	rand.Read(rand_key)

	cipher, err := aes.NewCipher(rand_key)
	if err != nil {
		fmt.Println(err)
	}

	if rand.Intn(2) >= 1 {
		return crypt.EcbEncrypt(clear_text, cipher)
	}

	iv := make([]byte, 16)
	rand.Read(iv)

	return crypt.CbcEncrypt(clear_text, iv, cipher)
}

func http_handler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "POST":

		var clear_text data
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&clear_text)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(clear_text)

		cipher_text, err := crypt.Bytes2B64(random_encrypt(crypt.Pkcs7Pad([]byte(clear_text.Secret), 16)))
		if err != nil {
			fmt.Println(err)
		}

		fmt.Fprintf(w, cipher_text)

	default:
		return
	}
}

func main() {
	http.HandleFunc("/", http_handler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8989", nil))
}
