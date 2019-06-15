package libcrypto

import (
	"fmt"
)

func Hex2Bytes(hex string) (bytes []byte, err error) {

	if len(hex) % 2 != 0 {
		return nil, fmt.Errorf("hex2Bytes recieved string with an odd number of letters")
	}

	bytes = make([]byte, len(hex) / 2)
	var byte_val uint

	for i, v := range []rune(hex) {

		c := uint(v)

		switch ; {
			case (c >= 48) && (c <= 57):
				c -= 48
			case (c >= 97) && (c <= 122):
				c -= 87
			case (c >= 65) && (c <= 90):
				c -= 55
			default:
				return nil, fmt.Errorf("hex2Bytes recieved string with a non hex character: %c", v)
		}


		if i % 2 == 0 {
			byte_val = 16 * c 
		} else {
			byte_val += c
			bytes[i / 2] = byte(byte_val)
		}
	}

	return bytes, nil
}

func Bytes2Hex(bytes []byte) (hex string, err error) {

	var val uint8
	for _, v := range bytes {

		val = uint8(v >> 4)
		hex += hex_table(val)

		val = uint8(v & 15)
		hex += hex_table(val)
	}
	return hex, nil
}

func Bytes2B64(bytes []byte) (b64 string, err error) {

	var six_bit_val uint8

	for len(bytes) % 3 != 0 {
		bytes = append(bytes, byte(0))
	}

	for i := 0; i < len(bytes) / 3; i++ {

		six_bit_val = bytes[i*3] >> 2
		b64 += b64_table(six_bit_val)

		six_bit_val = ((bytes[i*3] & 3) << 4) ^ (bytes[i*3 + 1] >> 4)
		b64 += b64_table(six_bit_val)

		six_bit_val = ((bytes[i*3 + 1] & 15) << 2) ^  (bytes[i*3 + 2] >> 6)
		b64 += b64_table(six_bit_val)

		six_bit_val = bytes[i*3 + 2] & 63
		b64 += b64_table(six_bit_val)
	}

	return b64, nil

}

func hex_table(val uint8) (char string) {

	switch ; {
		case val >= 0 && val <= 9:
			return string(val + 48)

		case val >= 10:
			return string(val + 87)
	}
	return ""
}

func b64_table(val uint8) (char string) {

	switch ; {
		case val >= 0 && val <= 25:
			return string(val + 65)

		case val >= 26 && val <= 51:
			return string(val + 71)

		case val >= 52 && val <= 61:
			return string(val - 4)

		case val == 62:
			return "+"

		case val == 63:
			return "/"
	}
	return ""
}

