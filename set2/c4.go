package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"bytes"
	"strings"
	crypt "../libcrypto"
)

func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func post(data []byte) (cipher_text []byte) {
	dataStr, err := crypt.Bytes2B64(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	dataStr = "{\"Secret\" : \"" + dataStr + "\"}"
	req, err := http.NewRequest("POST", "http://localhost:8989", bytes.NewBuffer([]byte(dataStr)))
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

func bruteByte(brutedBytes []byte) (bruted byte, err error) {
	blockIndex := len(brutedBytes) / 16
	prefix := strings.Repeat("A", (15 - (len(brutedBytes) % 16)))
	cipherBlock := post([]byte(prefix))

	cipherBlock = cipherBlock[blockIndex*16:(blockIndex+1)*16]
	input := append([]byte(prefix), brutedBytes...)
	input = input[len(input) - 15:]

	for bruted := 0; bruted <= 255; bruted++ {
		clearTextGuess := append(input, byte(bruted))
		returnedBlock := post(clearTextGuess)
		returnedBlock = returnedBlock[:16]
		if Equal(cipherBlock, returnedBlock) {
			return byte(bruted), nil
		}
	}

	return byte(0), fmt.Errorf("No byte matched")
}

func main() {

	var brutedBytes []byte
	for i := 0; i < len(post([]byte{})); i++ {
		b, err := bruteByte(brutedBytes)
		if err != nil {
			fmt.Println(err)
		} else {
			brutedBytes = append(brutedBytes, b)
			fmt.Println(string(brutedBytes))
		}
	}

}
