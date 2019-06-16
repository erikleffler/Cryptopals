package libcrypto

import (
	"fmt"
)

func Hex2Bytes(hex string) (bytes []byte, err error) {

	if len(hex) % 2 != 0 {
		return nil, fmt.Errorf("hex2Bytes recieved string with an odd number of letters")
	}

	bytes = make([]byte, len(hex) / 2)
	var byte_val uint8
	var val uint8

	for i, v := range []rune(hex) {

		val, err = hex_backward_table(v)
		if err != nil {
			return nil, err
		}

		if i % 2 == 0 {
			byte_val = 16 * val
		} else {
			byte_val += val
			bytes[i / 2] = byte(byte_val)
		}
	}

	return bytes, nil
}

func B642Bytes(b64 string) (bytes []byte, err error) {

	if len(b64) % 4 != 0 {
		return nil, fmt.Errorf("B642Bytes recieved string with an insufficient padding")
	}

	bytes = make([]byte, (len(b64) / 4) * 3)
	var val uint8

	for i := 0; i < len(b64) / 4; i++ {

		val, err = b64_backward_table(rune(b64[i*4]))
		if err != nil {
			return nil, err
		}
		bytes[i*3] = (val << 2)

		val, err = b64_backward_table(rune(b64[i*4+1]))
		if err != nil {
			return nil, err
		}
		bytes[i*3] ^= val >> 4
		bytes[i*3+1] =  (val & 15) << 4

		val,err = b64_backward_table(rune(b64[i*4+2]))
		if err != nil {
			return nil, err
		}
		bytes[i*3+1] ^= val >> 2
		bytes[i*3+2] = (val & 3) << 6

		val,err = b64_backward_table(rune(b64[i*4+3]))
		if err != nil {
			return nil, err
		}
		bytes[i*3+2] ^= val

	}

	return bytes, nil
}

func Bytes2Hex(bytes []byte) (hex string, err error) {

	var val uint8
	for _, v := range bytes {

		val = uint8(v >> 4)
		hex += hex_forward_table(val)

		val = uint8(v & 15)
		hex += hex_forward_table(val)
	}
	return hex, nil
}

func Bytes2B64(bytes []byte) (b64 string, err error) {

	var six_bit_val uint8

	padding_indicator := ""

	for len(bytes) % 3 != 0 {
		bytes = append(bytes, byte(0))
		padding_indicator += "="
	}

	for i := 0; i < len(bytes) / 3; i++ {

		six_bit_val = bytes[i*3] >> 2
		b64 += b64_forward_table(six_bit_val)

		six_bit_val = ((bytes[i*3] & 3) << 4) ^ (bytes[i*3 + 1] >> 4)
		b64 += b64_forward_table(six_bit_val)

		six_bit_val = ((bytes[i*3 + 1] & 15) << 2) ^  (bytes[i*3 + 2] >> 6)
		b64 += b64_forward_table(six_bit_val)

		six_bit_val = bytes[i*3 + 2] & 63
		b64 += b64_forward_table(six_bit_val)
	}


	b64 = b64[:len(b64)-len(padding_indicator)] + padding_indicator

	return b64, nil

}

func hex_forward_table(val uint8) (char string) {

	switch ; {
		case val >= 0 && val <= 9:
			return string(val + 48)

		case val >= 10:
			return string(val + 87)
	}
	return ""
}

func hex_backward_table(char rune) (val uint8, err error) {

	val = uint8(char)

	switch ; {
		case (val >= 48) && (val <= 57):
			val -= 48

		case (val >= 97) && (val <= 122):
			val -= 87

		case (val >= 65) && (val <= 90):
			val -= 55

		default:
			return 0, fmt.Errorf("hex_backward_table recieved string with a non hex character: %c", char)
	}
	return val, nil
}

func b64_forward_table(val uint8) (char string) {

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

func b64_backward_table(char rune) (val uint8, err error) {

	val = uint8(char)

	switch ; {
		case (val >= 65) && (val <= 90):
			val -= 65

		case (val >= 97) && (val <= 122):
			val -= 71

		case (val >= 48) && (val <= 57):
			val += 4

		case val == 43:
			val = 62

		case val == 47:
			val = 63

		case val == 61:
			val = 0

		default:
			return 0, fmt.Errorf("b64_backward_table recieved string with a non hex character: %c", char)
	}
	return val, nil
}
