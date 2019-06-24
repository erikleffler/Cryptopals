package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"bytes"
	crypt "../libcrypto"
)

var iv []byte

func post(data []byte, path string) (cipherText []byte) {
	dataStr, err := crypt.Bytes2B64(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	b64Iv, err := crypt.Bytes2B64(iv)
	if err != nil {
		fmt.Println(err)
		return
	}
	dataStr = "{\"Secret\" : \"" + dataStr + "\", \"Iv\": \"" + b64Iv + "\"}"
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
	cipherText, err = crypt.B642Bytes(string(body))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return cipherText
}

func padOracleAtk(cipherText []byte, iv []byte) (clearText []byte) {

	clearText = make([]byte, len(cipherText))
	cipherText = append(iv, cipherText...)

	blockCpy := make([]byte, 32)
	var response string

	for blockIndex := 0; blockIndex < (len(clearText) / 16); blockIndex++ {

		copy(blockCpy, cipherText[blockIndex*16:(blockIndex+2)*16])

		for byteIndex := 15; byteIndex >= 0; byteIndex-- {

			for byteVal := 0; byteVal < 256; byteVal++ {

				blockCpy[byteIndex] ^= byte(byteVal)
				response = string(post(blockCpy, "val"))

				if response == "Valid padding" {

					if byteVal == 0 { //already padded, need to figure out amount
						for padAmount := 1; padAmount <= 16; padAmount++ {
							blockCpy[16-padAmount] ^= 177 //some number > block_size
							response = string(post(blockCpy, "val"))
							blockCpy[16-padAmount] ^= 177 // ^
							if response == "Valid padding" {
								byteIndex = 17 - padAmount
								break
							} else if padAmount == 16 {
								return clearText //full block padded
							}
						}
					} else {

						clearText[blockIndex*16 + byteIndex] = byte((16 - byteIndex) ^ byteVal)
						fmt.Println(string(clearText))
					}

					for i := byteIndex; i < 16; i++ {
						blockCpy[i] ^= byte((16 - byteIndex) ^ (17 - byteIndex))
					}
					break

				} else {
					blockCpy[byteIndex] ^= byte(byteVal)
					if byteVal == 255 {
						fmt.Println("failed", blockIndex)
					}
				}
			}
		}
	}
	return clearText
}

func main() {

	cipherTextAndIv := post([]byte{}, "cip")
	cipherText := cipherTextAndIv[:len(cipherTextAndIv) - 16]
	iv = cipherTextAndIv[len(cipherTextAndIv) - 16:]
	fmt.Println(iv)
	padOracleAtk(cipherText, iv)
}
