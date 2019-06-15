package main

import (
	"bufio"
	"fmt"
	"os"
	enc "../libcrypto"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	bytes, err := enc.Hex2Bytes(text[:len(text)-1])
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println(enc.Bytes2B64(bytes))
}
