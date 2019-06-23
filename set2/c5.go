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

	adminInput := []byte(strings.Repeat("A", 10) + "admin" + strings.Repeat(string(byte(11)), 11))
	adminStr := post(adminInput, "pf")[16:32]
	mailInput := []byte("foo12@bar.com")
	mailStr := post(mailInput, "pf")[0:32]
	fmt.Println(string(post(append(mailStr, adminStr...), "val")))

}
