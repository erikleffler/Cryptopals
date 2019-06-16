package main

import (
	"fmt"
	crypt "../libcrypto"
)

func main() {
	fmt.Println(crypt.Pkcs7Pad([]byte("Yellow Submarine"), 20))
}
