package main

import(
	"fmt"
	crypt "../libcrypto"
)

func main() {
	seed := uint32(13)
	mt1, err := crypt.NewMT(seed)
	mt2, err := crypt.NewMT(seed)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i:= 0; i <= 625; i++ {
		fmt.Println(mt1.Rand())
		fmt.Println(mt2.Rand())
	}
}
