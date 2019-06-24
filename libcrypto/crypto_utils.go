package libcrypto

import (
	"math/bits"
	"fmt"
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

func Pkcs7Pad(bytes []byte, block_size int) []byte{

	amount := block_size - (len(bytes) % block_size)
	padding := make([]byte, amount)
	for i := range padding {
		padding[i] = byte(amount)
	}
	return append(bytes, padding...)
}

func Pkcs7Unpad(bytes []byte, block_size int) (unpadded []byte, err error) {
	lastByte := bytes[len(bytes) - 1]
	if int(lastByte) <= block_size && int(lastByte) > 0 {
		padding := bytes[len(bytes) - int(lastByte):]
		for _, v := range padding {
			if v != lastByte {
				return unpadded, fmt.Errorf("Invalid padding")
			}
		}
		return bytes[:len(bytes) - int(lastByte)], nil
	}
	return unpadded, fmt.Errorf("Recieved byte array without padding")
}
