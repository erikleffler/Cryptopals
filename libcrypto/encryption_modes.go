package libcrypto

import (
	"crypto/cipher"
)

func EcbDecrypt(cipher_text []byte, cipher cipher.Block) (clear_text []byte) {

	clear_text = make([]byte, len(cipher_text))

	for i := 0; i < len(cipher_text) / 16; i++ {

		cipher.Decrypt(clear_text[i*16:(i+1)*16], cipher_text[i*16:(i+1)*16])
	}
	return clear_text
}

func EcbEncrypt(clear_text []byte, cipher cipher.Block) (cipher_text []byte) {

	cipher_text = make([]byte, len(clear_text))

	for i := 0; i < len(clear_text) / 16; i++ {

		cipher.Encrypt(cipher_text[i*16:(i+1)*16], clear_text[i*16:(i+1)*16])
	}
	return cipher_text
}

func CbcDecrypt(cipher_text []byte, iv []byte, cipher cipher.Block) (clear_text []byte) {


	pre_xored_block := make([]byte, 16)

	cipher.Decrypt(pre_xored_block, cipher_text[0:16])
	clear_text = append(clear_text, Xor(pre_xored_block, iv)...)

	for i := 1; i < len(cipher_text) / 16; i++ {

		cipher.Decrypt(pre_xored_block, cipher_text[i*16:(i+1)*16])
		clear_text = append(clear_text, Xor(pre_xored_block, cipher_text[(i-1)*16:i*16])...)
	}

	return clear_text
}

func CbcEncrypt(clear_text []byte, iv []byte, cipher cipher.Block) (cipher_text []byte) {

	cipher_text = make([]byte, len(clear_text))

	xored_block := Xor(clear_text[0:16], iv)

	i := 0
	for ; i < (len(clear_text) - 1) / 16; i++ {

		cipher.Encrypt(cipher_text[i*16:(i+1)*16], xored_block)
		xored_block = Xor(clear_text[(i+1)*16:(i+2)*16], cipher_text[i*16:(i+1)*16])
	}

	cipher.Encrypt(cipher_text[i*16:(i+1)*16], xored_block)
	return cipher_text
}
