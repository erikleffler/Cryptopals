package main

import (
	"fmt"
	"os"
	"bufio"

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

func detect_block_repetition(cipher_text []byte, block_size int) bool{

	for i := 0; i < len(cipher_text) / block_size; i++ {
		for j := i+1; j < len(cipher_text) / block_size; j++ {
			if Equal(cipher_text[j*block_size:(j+1)*block_size], cipher_text[i*block_size:(i+1)*block_size]) {
				return true
			}
		}
	}
	return false

}

func main() {
	file, err := os.Open("./8.txt")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		raw_text, err := crypt.Hex2Bytes(scanner.Text())
		if err != nil {
			fmt.Println(err)
			return
		}
		if detect_block_repetition(raw_text, 16) {
			fmt.Println(crypt.Bytes2Hex(raw_text))
			fmt.Println("======================================================")
		}
	}

}
