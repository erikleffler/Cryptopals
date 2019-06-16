package main

import (
	"fmt"
	"os"
	"bufio"

	crypt "../libcrypto"
)

var freq = map[string]float64{
	"a": 0.08167,
	"b": 0.01492,
	"c": 0.02782,
	"d": 0.04253,
	"e": 0.12702,
	"f": 0.02228,
	"g": 0.02015,
	"h": 0.06094,
	"i": 0.06966,
	"j": 0.00153,
	"k": 0.00772,
	"l": 0.04025,
	"m": 0.02406,
	"n": 0.06749,
	"o": 0.07507,
	"p": 0.01929,
	"q": 0.00095,
	"r": 0.05987,
	"s": 0.06327,
	"t": 0.09056,
	"u": 0.02758,
	"v": 0.00978,
	"w": 0.02360,
	"x": 0.00150,
	"y": 0.01974,
	"z": 0.00074,
	" ": 0.1,
}

func score(text []byte) (score float64) {

	for _, v := range text {
		score += freq[string(v)]
	}
	return score
}

func break_xor_ceasar(cipher []byte) (key byte) {

	var current_score float64
	var best_score float64
	var attempt []byte

	for i := 0; i <= 255; i++{
		attempt = crypt.Xor(cipher, []byte{byte(i)})
		current_score = score(attempt)
		if current_score >= best_score {
			key = byte(i)
			best_score = current_score
		}
	}
	return

}
func score_keysize(cipher []byte, keysize int) (score float64) {


	for i := 0; i < len(cipher) / keysize - 2; i++ {
		score += float64(crypt.HammingDist(cipher[i*keysize:(i+1)*keysize], cipher[(i+1)*keysize:(i+2)*keysize]))
	}

	return score / float64(len(cipher) - keysize)
}

func find_keysize(cipher []byte, min int, max int) (keysize int) {

	score := 10.0
	best_score := 10.0
	for i := min; i <= max; i++ {
		score = score_keysize(cipher, i)
		if score <= best_score {
			best_score = score
			keysize = i
		}
	}
	return keysize
}

func transpose_cipher(cipher []byte, keysize int) [][]byte {

	cipher_blocks := make([][]byte, keysize)

	for i, v := range cipher {
		cipher_blocks[i%keysize] = append(cipher_blocks[i%keysize], v)
	}
	return cipher_blocks
}

func main() {
	file, err := os.Open("./6.txt")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	cipher := ""
	for scanner.Scan() {
		cipher += scanner.Text()
	}
	raw_cipher, err := crypt.B642Bytes(cipher)
	if err != nil {
		fmt.Println(err)
		return
	}



	keysize := find_keysize(raw_cipher, 2, 40)
	key := make([]byte, keysize)
	cipher_blocks := transpose_cipher(raw_cipher, keysize)

	for i, v := range cipher_blocks {
		key[i] = break_xor_ceasar(v)
	}
	fmt.Println(string(crypt.Xor(raw_cipher, key)))

}
