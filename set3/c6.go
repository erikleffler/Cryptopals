package main

import(
	"fmt"
	"time"
	"math/rand"
	crypt "../libcrypto"
)

func main() {
	now := time.Now().Unix()
	rand.Seed(now)
	seed := uint32(now + int64(rand.Intn(1000)))
	fmt.Println("Original seed: ", seed)

	mt, err := crypt.NewMT(seed)
	if err != nil {
		fmt.Println(err)
		return
	}
	origVal, err := mt.Rand()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Val: ", origVal)

	for i := int64(0); i < 1001; i++ {
		mt, err = crypt.NewMT(uint32(now + i))
		if err != nil {
			fmt.Println(err)
			return
		}
		val, err := mt.Rand()
		if err != nil {
			fmt.Println(err)
			return
		}
		if val == origVal {
			fmt.Println("Found seed:", seed)
			return
		}
	}

}
