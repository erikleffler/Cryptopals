package main

import (
	"fmt"
)

func main() {
	a := []rune("asdasdasd")
	for i, r := range a {
		fmt.Printf("i%d r %c\n", i, r)
	}
}
