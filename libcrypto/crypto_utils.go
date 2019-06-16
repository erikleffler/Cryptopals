package libcrypto

import (
	"math/bits"
)

// Assumes b1 is longer than b2
func Xor(b1 []byte, b2 []byte) (xored []byte) {

	xored = make([]byte, len(b1))

	for i, _ := range b1 {
		xored[i] = b1[i] ^ b2[i % len(b2)]
	}

	return xored

}

// Assumes b1 is longer
func HammingDist(b1 []byte, b2[]byte) (dist int) {

	for i := 0; i < len(b1); i++ {
		dist += bits.OnesCount(uint(b1[i] ^ b2[i%len(b2)]))
	}
	return dist
}
