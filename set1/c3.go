package main

import (
	"fmt"
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

func break_xor_ceasar(cipher []byte) (clear []byte) {

	var best_score float64
	var current_score float64
	var attempt []byte

	for i := 0; i <= 255; i++{
		attempt = crypt.Xor(cipher, []byte{byte(i)})
		current_score = score(attempt)
		if current_score >= best_score {
			clear = attempt
			best_score = current_score
		}
	}
	return clear

}

func main() {

	cipher, err := crypt.Hex2Bytes("1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(break_xor_ceasar(cipher)))


}
