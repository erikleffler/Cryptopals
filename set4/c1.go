package main

import(
	"fmt"
	"math/rand"
	"os"
	"bufio"
	"time"

	"crypto/aes"
	"crypto/cipher"
	crypt "../libcrypto"
)

func editCipherText(cipherText []byte,
	nonce []byte,
	offset int,
	newText []byte,
	cipher cipher.Block) (newCipherText []byte) {

	// First get the number of blocks for the xor pad
	numBlocks := (len(newText) + (offset % 16)) / 16
	// If overflow, one more block
	if (len(newText) + (offset % 16)) % 16 != 0 {
		numBlocks += 1
	}

	// Create the xor pad
	var xorPad []byte
	blockXorPad := make([]byte, 16)

	encBlock := append(nonce, make([]byte, 16)...)
	ctrBase := offset / 16
	for ctrInc := 0; ctrInc < numBlocks; ctrInc++ {
		encBlock[8] = byte(ctrBase+ctrInc)
		cipher.Encrypt(blockXorPad, encBlock)
		xorPad = append(xorPad, blockXorPad...)
	}

	// Trim the start
	xorPad = xorPad[offset%16:]
	// Trim the end
	xorPad = xorPad[:len(newText)]
	//construct newCipherText
	newCipherText = make([]byte, len(cipherText))
	copy(newCipherText, cipherText)

	newCipherText = append(newCipherText[:offset], crypt.Xor(newText, xorPad)...)
	newCipherText = append(newCipherText, cipherText[offset + len(newText):]...)

	return newCipherText
}

// It takes the nonce and cipher as a paramter but only to pass to edit.
func recoverClearText(cipherText []byte, nonce []byte, cipher cipher.Block) (recoveredClearText []byte){

	recoveredClearText = make([]byte, len(cipherText))
	zeros := make([]byte, len(cipherText))
	for i, _ := range(zeros) {
		zeros[i] = 0
	}
	encZeros := editCipherText(cipherText, nonce, 0, zeros, cipher)
	return crypt.Xor(cipherText, encZeros)
}

func equals(b1 []byte, b2 []byte) bool {
	for i, _ := range(b1) {
		if b1[i] != b2[i] {
			return false
		}
	}
	return true
}

func main() {
	file, err := os.Open("./1.txt")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var firstCipherText []byte
	for scanner.Scan() {
		var bytes []byte
		bytes, err := crypt.B642Bytes(scanner.Text())
		if err != nil {
			fmt.Println(err)
			return
		}
		firstCipherText = append(firstCipherText, bytes...)
	}
	key := []byte("YELLOW SUBMARINE")
	cipher, err := aes.NewCipher(key)

	if err != nil {
		fmt.Println(err)
		return
	}

	clearText := make([]byte, len(firstCipherText))
	for i := 0; i < len(firstCipherText) / 16; i++ {
		cipher.Decrypt(clearText[i*16:(i+1)*16], firstCipherText[i*16:(i+1)*16])
	}

	key = make([]byte, 16)
	nonce := make([]byte, 16)

	rand.Seed(time.Now().UTC().UnixNano())
	rand.Read(key)
	rand.Read(nonce)

	cipher, err = aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println(err)
		return
	}
	cipherText := crypt.CtrEncrypt(clearText, nonce, cipher)
	recoveredClearText := recoverClearText(cipherText[:], nonce, cipher)
	fmt.Println(equals(recoveredClearText, clearText))
	fmt.Println(string(recoveredClearText))

}

