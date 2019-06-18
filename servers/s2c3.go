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


func random_pad(cipher_text []byte) (padded_cipher_text []byte) {

	amount := 5 + rand.Intn(5)
	pad := make([]byte, amount)
	rand.Read(pad)

	padded_cipher_text = append(cipher_text, pad...)

	amount = 5 + rand.Intn(5)
	pad = make([]byte, amount)
	rand.Read(pad)

	return crypt.Pkcs7Pad(append(pad, padded_cipher_text...), 16)




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

		cipher_text, err := crypt.Bytes2Hex(random_encrypt(random_pad([]byte(clear_text.Secret))))
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
