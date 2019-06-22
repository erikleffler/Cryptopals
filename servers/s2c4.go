package main

import (
	"fmt"
	"net/http"
	"math/rand"
	"crypto/aes"
	"crypto/cipher"
	"log"
	"encoding/json"
	"time"

	crypt "../libcrypto"
)

type Data struct {
	Secret string
}

var aesCipher cipher.Block
var secret_string, _ = crypt.B642Bytes("Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK")


func http_handler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "POST":

		var data Data
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)

		if err != nil {
			fmt.Println(err)
			return
		}
		clear_text, err := crypt.B642Bytes(data.Secret)
		if err != nil {
			fmt.Println(err)
			return
		}
		clear_text = append(clear_text, secret_string...)

		cipher_text, err := crypt.Bytes2B64(crypt.EcbEncrypt(crypt.Pkcs7Pad(clear_text, 16), aesCipher))
		if err != nil {
			fmt.Println(err)
		}

		fmt.Fprintf(w, cipher_text)
		return

	default:
		return
	}
}

func setup_cipher() (cipher.Block, error) {
	rand_key := make([]byte, 16)
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Read(rand_key)
	return aes.NewCipher(rand_key)
}

func main() {

	var err error
	aesCipher, err = setup_cipher()
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/", http_handler)
	log.Fatal(http.ListenAndServe("127.0.0.1:8989", nil))
}
