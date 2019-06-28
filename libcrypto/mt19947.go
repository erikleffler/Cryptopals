package libcrypto

import (
	"fmt"
)


// Initialize algorithm constants
var w = uint32(32)
var n = uint32(624)
var m = uint32(397)
var a = uint32(2567483615)
var u = uint32(11)
var s = uint32(7)
var b = uint32(2636928640)
var t = uint32(15)
var c = uint32(4022730752)
var l = uint32(18)
var f = uint32(1812433253)
var lowMask = uint32((1 << uint32(31)) - 1)
var uppMask = uint32(1 << uint32(31))


type MT struct {
	State	[]uint32
	Seeded	bool
	Index	uint32
}

func (mt *MT) Rand() (uint32, error) {

	if !mt.Seeded {
		return 0, fmt.Errorf("Calling rand on unSeeded MT")
	}

	if mt.Index == n {
		err := mt.twist()
		if err != nil {
			return 0, err
		}
	}

	y := mt.State[mt.Index]
	y ^= (y >> u)
	y ^= ((y << s) & b)
	y ^= ((y << t) & c)
	y ^= y >> l

	mt.Index += 1

	return y, nil


}

func (mt *MT) twist() error {
	if mt.Index != n {
		return fmt.Errorf("Twisting before having retrieved all tempered values")
	}
	var x uint32
	var xA uint32
	for i := uint32(0); i < n; i++ {
		x = (mt.State[i] & uppMask) + (mt.State[(i + 1) % n] & lowMask)
		xA = x >> 1
		if x % 2 != 0 {
			xA ^= a
		}
		mt.State[i] = mt.State[(i + m) % n] ^ xA
	}
	mt.Index = 0
	return nil
}

func (mt *MT) seedMT(seed uint32) error {
	if mt.Seeded {
		return fmt.Errorf("MT Already Seeded")
	}
	mt.State = make([]uint32, n)
	mt.State[0] = seed
	for i := uint32(1); i < n; i++ {
		mt.State[i] = uint32(f * (mt.State[i-1] ^ (mt.State[i-1] >> (w-2))) + i)
	}
	mt.Index = n
	mt.Seeded = true
	return nil
}

func NewMT(seed uint32) (MT, error) {
	var mt MT
	err := mt.seedMT(seed)
	if err != nil {
		return mt, err
	}
	err = mt.twist()
	if err != nil {
		return mt, err
	}
	return mt, nil
}
