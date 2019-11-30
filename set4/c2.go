package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"bytes"
	"strings"
	crypt "../libcrypto"
)

func post(data []byte, path string) (cipher_text []byte) {
	dataStr, err := crypt.Bytes2B64(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	dataStr = "{\"Secret\" : \"" + dataStr + "\"}"
	req, err := http.NewRequest("POST", "http://127.0.0.1:8989/" + path, bytes.NewBuffer([]byte(dataStr)))
	if err != nil {
		fmt.Println(err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	cipher_text, err = crypt.B642Bytes(string(body))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return cipher_text
}

func main() {

	serverPrefix := "comment1=cooking%20MCs;userdata="

	payload := []byte("a;admin=true")

	input := []byte(strings.Repeat("A", len(payload)))
	cipherText := post(input, "pf")
	xorRes := crypt.Xor(input, payload)

	bitflippedSlice := crypt.Xor(cipherText[len(serverPrefix):len(payload) + len(serverPrefix)], xorRes)
	bitflippedCipherText := append(cipherText[:len(serverPrefix)], bitflippedSlice...)
	bitflippedCipherText = append(bitflippedCipherText, cipherText[len(serverPrefix) + len(xorRes):]...)
	malRes := post(bitflippedCipherText, "val")
	fmt.Println(string(malRes))
}
