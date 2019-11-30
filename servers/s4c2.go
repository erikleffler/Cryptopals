package main

import (
	"strings"
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
var nonce []byte

func embeddData(data []byte) []byte {

	data = []byte(strings.Replace(string(data), "=", "'='", -1))
	data = []byte(strings.Replace(string(data), ";", "';'", -1))

	return []byte("comment1=cooking%20MCs;userdata=" + string(data) + ";comment2=%20like%20a%20pound%20of%20bacon")
}

func httpHandlerPf(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "POST":

		var data Data
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)

		if err != nil {
			fmt.Println(err)
			return
		}

		clearText, err := crypt.B642Bytes(data.Secret)
		if err != nil {
			fmt.Println(err)
			return
		}
		clearText = embeddData(clearText)

		cipherText, err := crypt.Bytes2B64(crypt.CtrEncrypt(clearText, nonce, aesCipher))
		if err != nil {
			fmt.Println(err)
		}

		fmt.Fprintf(w, cipherText)
		return

	default:
		return
	}
}

func httpHandlerVal(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "POST":

		var data Data
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)

		if err != nil {
			fmt.Println(err)
			return
		}
		cipherText, err := crypt.B642Bytes(data.Secret)
		if err != nil {
			fmt.Println(err)
			return
		}
		clearText, err := crypt.Bytes2B64(crypt.CtrDecrypt(cipherText, nonce, aesCipher))
		if err != nil {
			fmt.Println(err)
		}

		fmt.Fprintf(w, clearText)
		return

	default:
		return
	}
}

func setupCipher() (cipher.Block, error) {
	randKey := make([]byte, 16)
	nonce = make([]byte, 16)
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Read(randKey)
	rand.Read(nonce)
	return aes.NewCipher(randKey)
}

func main() {

	var err error
	aesCipher, err = setupCipher()
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/pf", httpHandlerPf)
	http.HandleFunc("/val", httpHandlerVal)
	log.Fatal(http.ListenAndServe("127.0.0.1:8989", nil))
}
