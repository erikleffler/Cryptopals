package main

import(
	"fmt"
	"time"
	"math/rand"
	crypt "../libcrypto"
)

// MT19937 variables
var u = uint32(11)
var s = uint32(7)
var b = uint32(2636928640)
var t = uint32(15)
var c = uint32(4022730752)
var l = uint32(18)

func reverseTemper(y uint32) uint32 {

	/*
		Original algorithm: 

		y := mt.State[mt.Index]	(1)
		y ^= (y >> u)			(2)
		y ^= ((y << s) & b)		(3)
		y ^= ((y << t) & c)		(4)
		y ^= y >> l				(5)

		In order to reverse step 5, we perform the identical operation. And get

		y_res = (orig_y ^ orig_y >> l) ^ (orig_y >> l ^ orig_y >> 2l),
		which since 2l > 32 => y >> 2l = 0 and y_res = orig_y

		to reverse step 4 we reproduce the logic above, however, we also need to be left
		wtih a y << x*t such that x*t > 32. After reproducing step 4 we are left with:

		y_res = y_orig ^ ((y_org << t) & c)^ ((y_orig << t) & c) ^ (((y_orig << 2t) & (c << t)) & c)
		
		Notice the two middle terms cancel out each other, however the remaining last term
		can only be cancled out by repeating the operation twice again, I.e a total of 1 + 2 times.
		This reasoning leads us to deduce that we need to redo the operations a: 
			
			sum([2**n for n in range(0, ceil(32/bit_shift_val)-1)])	

		amount of times where bit_shift_val in this case is u, s, t or l from above.
	*/

	y ^= y >> l

	// 1 + 2 times
	y ^= (y << t) & c
	y ^= (y << t) & c
	y ^= (y << t) & c

	// 1 + 2 + 4 times
	y ^= (y << s) & b
	y ^= (y << s) & b
	y ^= (y << s) & b
	y ^= (y << s) & b
	y ^= (y << s) & b
	y ^= (y << s) & b
	y ^= (y << s) & b

	// 1 + 2 times
	y ^= y >> u
	y ^= y >> u
	y ^= y >> u

	return y
}

func cloneMT(obs []uint32) crypt.MT {
	clonedState := make([]uint32, 624)
	for i := 0; i < 624; i++ {
		clonedState[i] = reverseTemper(obs[i])
	}
	mt := crypt.MT{clonedState, true, 624}
	return mt
}

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

	obs := make([]uint32, 624)
	for i := 0; i < 624; i++ {
		obs[i], err = mt.Rand()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	clonedMT := cloneMT(obs)
	var cv, ov uint32
	for i := 0; i < 10; i++ {
		cv, err = clonedMT.Rand()
		if err != nil {
			fmt.Println(err)
			return
		}
		ov,err  = mt.Rand()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Cloned:	", cv, "	Orig:	", ov)
	}
}
