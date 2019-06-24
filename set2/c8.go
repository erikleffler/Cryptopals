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
	req, err := http.NewRequest("POST", "http://localhost:8989/" + path, bytes.NewBuffer([]byte(dataStr)))
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


	input := []byte(strings.Repeat("A", 48 - (len(serverPrefix) % 16)))
	cipherText := post(input, "pf")
	xorRes := crypt.Xor([]byte(strings.Repeat("A", len(";admin=true;"))), []byte(";admin=true;"))
	blockIndex := len(serverPrefix) + len(input) - 32
	bitflippedSlice := crypt.Xor(cipherText[blockIndex:blockIndex + len(xorRes)], xorRes)
	bitflippedCipherText := append(cipherText[:blockIndex], bitflippedSlice...)
	bitflippedCipherText = append(bitflippedCipherText, cipherText[blockIndex + len(xorRes):]...)
	malRes := post(bitflippedCipherText, "val")
	fmt.Println(string(malRes))
}
