package libcrypto

import (
	"crypto/cipher"
	"encoding/binary"
)

func EcbDecrypt(cipherText []byte, cipher cipher.Block) (clearText []byte) {

	clearText = make([]byte, len(cipherText))

	for i := 0; i < len(cipherText) / 16; i++ {

		cipher.Decrypt(clearText[i*16:(i+1)*16], cipherText[i*16:(i+1)*16])
	}
	return clearText
}

func EcbEncrypt(clearText []byte, cipher cipher.Block) (cipherText []byte) {

	cipherText = make([]byte, len(clearText))

	for i := 0; i < len(clearText) / 16; i++ {

		cipher.Encrypt(cipherText[i*16:(i+1)*16], clearText[i*16:(i+1)*16])
	}
	return cipherText
}

func CbcDecrypt(cipherText []byte, iv []byte, cipher cipher.Block) (clearText []byte) {


	preXoredBlock := make([]byte, 16)

	cipher.Decrypt(preXoredBlock, cipherText[0:16])
	clearText = append(clearText, Xor(preXoredBlock, iv)...)

	for i := 1; i < len(cipherText) / 16; i++ {

		cipher.Decrypt(preXoredBlock, cipherText[i*16:(i+1)*16])
		clearText = append(clearText, Xor(preXoredBlock, cipherText[(i-1)*16:i*16])...)
	}

	return clearText
}

func CbcEncrypt(clearText []byte, iv []byte, cipher cipher.Block) (cipherText []byte) {

	cipherText = make([]byte, len(clearText))

	xoredBlock := Xor(clearText[0:16], iv)

	i := 0
	for ; i < (len(clearText) - 1) / 16; i++ {

		cipher.Encrypt(cipherText[i*16:(i+1)*16], xoredBlock)
		xoredBlock = Xor(clearText[(i+1)*16:(i+2)*16], cipherText[i*16:(i+1)*16])
	}

	cipher.Encrypt(cipherText[i*16:(i+1)*16], xoredBlock)
	return cipherText
}

func CtrDecrypt(cipherText []byte, nonce []byte, cipher cipher.Block) (clearText []byte) {

	xorPad := make([]byte, 16)
	encBlock := append(nonce, make([]byte, 8)...)

	var ctr uint64
	for ctr = 0; int(ctr) < len(cipherText) / 16; ctr++ {
		binary.LittleEndian.PutUint64(encBlock[8:], ctr)
		cipher.Encrypt(xorPad, encBlock)
		clearText = append(clearText, Xor(cipherText[ctr*16:(ctr+1)*16], xorPad)...)
	}

	if len(clearText) % 16 != 0{
		binary.LittleEndian.PutUint64(encBlock[8:], ctr)
		cipher.Encrypt(xorPad, encBlock)
		clearText = append(clearText, Xor(cipherText[ctr*16:], xorPad)...)
	}

	return clearText[:len(cipherText)]
}
func CtrEncrypt(clearText []byte, nonce []byte, cipher cipher.Block) (cipherText []byte) {

	xorPad := make([]byte, 16)
	encBlock := append(nonce, make([]byte, 8)...)

	var ctr uint64
	for ctr = 0; int(ctr) < len(clearText) / 16; ctr++ {
		binary.LittleEndian.PutUint64(encBlock[8:], ctr)
		cipher.Encrypt(xorPad, encBlock)
		cipherText = append(cipherText, Xor(clearText[ctr*16:(ctr+1)*16], xorPad)...)
	}

	if len(clearText) % 16 != 0{
		binary.LittleEndian.PutUint64(encBlock[8:], ctr)
		cipher.Encrypt(xorPad, encBlock)
		cipherText = append(cipherText, Xor(clearText[ctr*16:], xorPad)...)
	}

	return cipherText[:len(clearText)]
}
