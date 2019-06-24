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
	Secret	string
	Iv		string
}

var aesCipher cipher.Block
var randIv []byte

func randClearText() ([]byte, error){
	clearTextArray := []string {
		"MDAwMDAwTm93IHRoYXQgdGhlIHBhcnR5IGlzIGp1bXBpbmc=",
		"MDAwMDAxV2l0aCB0aGUgYmFzcyBraWNrZWQgaW4gYW5kIHRoZSBWZWdhJ3MgYXJlIHB1bXBpbic=",
		"MDAwMDAyUXVpY2sgdG8gdGhlIHBvaW50LCB0byB0aGUgcG9pbnQsIG5vIGZha2luZw==",
		"MDAwMDAzQ29va2luZyBNQydzIGxpa2UgYSBwb3VuZCBvZiBiYWNvbg==",
		"MDAwMDA0QnVybmluZyAnZW0sIGlmIHlvdSBhaW4ndCBxdWljayBhbmQgbmltYmxl",
		"MDAwMDA1SSBnbyBjcmF6eSB3aGVuIEkgaGVhciBhIGN5bWJhbA==",
		"MDAwMDA2QW5kIGEgaGlnaCBoYXQgd2l0aCBhIHNvdXBlZCB1cCB0ZW1wbw==",
		"MDAwMDA3SSdtIG9uIGEgcm9sbCwgaXQncyB0aW1lIHRvIGdvIHNvbG8=",
		"MDAwMDA4b2xsaW4nIGluIG15IGZpdmUgcG9pbnQgb2g=",
		"MDAwMDA5aXRoIG15IHJhZy10b3AgZG93biBzbyBteSBoYWlyIGNhbiBibG93",
	}
	return crypt.B642Bytes(clearTextArray[rand.Intn(len(clearTextArray))])
}

func httpHandlerCip(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "POST":

		clearText, err := randClearText()
		if err != nil {
			fmt.Println(err)
			return
		}

		cipherTextAndIv, err := crypt.Bytes2B64(append(crypt.CbcEncrypt(crypt.Pkcs7Pad(clearText, 16),randIv, aesCipher), randIv...))
		if err != nil {
			fmt.Println(err)
		}

		fmt.Fprintf(w, cipherTextAndIv)
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
		iv, err := crypt.B642Bytes(data.Iv)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = crypt.Pkcs7Unpad(crypt.CbcDecrypt(cipherText, iv, aesCipher), 16)
		if err != nil {
			msg, _ := crypt.Bytes2B64([]byte(err.Error()))
			fmt.Fprintf(w, msg)
		} else {
			msg, _ := crypt.Bytes2B64([]byte("Valid padding"))
			fmt.Fprintf(w, msg)
		}

		return

	default:
		return
	}
}

func setupCipher() (cipher.Block, error) {
	randKey := make([]byte, 16)
	randIv = make([]byte, 16)
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Read(randKey)
	rand.Read(randIv)
	return aes.NewCipher(randKey)
}

func main() {

	var err error
	aesCipher, err = setupCipher()
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/cip", httpHandlerCip)
	http.HandleFunc("/val", httpHandlerVal)
	log.Fatal(http.ListenAndServe("127.0.0.1:8989", nil))
}
