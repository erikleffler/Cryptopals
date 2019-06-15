package libcrypto

// Assumes b1 is longer than b2
func Xor(b1 []byte, b2 []byte) (xored []byte) {

	xored = make([]byte, len(b1))

	for i, _ := range b1 {
		xored[i] = b1[i] ^ b2[i % len(b2)]
	}

	return xored

}
