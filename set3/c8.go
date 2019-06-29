package main

import(
	"time"
	"fmt"
	"math/rand"
	crypt "../libcrypto"
)

func uint32ToBytes(val uint32) (bytes []byte) {

	mask := uint32((1 << 8) - 1)
	bytes = make([]byte, 4)
	bytes[0] = byte(val & mask)
	bytes[1] = byte(val & (mask << 8))
	bytes[2] = byte(val & (mask << 16))
	bytes[3] = byte(val & (mask << 24))
	return bytes
}

func equal(a, b []byte) bool {
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

func mtEncrypt(clearText []byte, mt crypt.MT) (cipherText []byte) {

	var xorPad []byte
	var r uint32
	for i := 0; i < len(clearText) / 4; i++ {
		r, _ = mt.Rand()
		xorPad = append(xorPad, uint32ToBytes(r)...)
	}
	cipherText = crypt.Xor(clearText, xorPad)
	return cipherText
}

func main() {
	rand.Seed(time.Now().Unix())
	seed := uint32(rand.Intn(65535))
	fmt.Println("Original seed: ", seed)

	mt, err := crypt.NewMT(seed)
	if err != nil {
		fmt.Println(err)
		return
	}

	randLen := rand.Intn(35)
	randPre := make([]byte, randLen)
	rand.Read(randPre)
	clearText := append(randPre, []byte("AAAAAAAAAAAAAA")...)
	cipherText := mtEncrypt(clearText, mt)

	for i := uint32(0); i < 65535; i++ {
		mt, err = crypt.NewMT(i)
		if err != nil {
			fmt.Println(err)
			return
		}
		if equal(mtEncrypt(cipherText, mt)[len(cipherText)-14:], []byte("AAAAAAAAAAAAAA")) {
			fmt.Println("Found seed:", seed)
			return
		}
	}

}

