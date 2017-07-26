// CrossGate加解密库，用于游戏加解密
// C语言翻译为GO语言
// version 1.0 beta
// by koangel
// email: jackliu100@gmail.com
// 2017/7/23
package CGCrypt

import (
	"math"
)

/*
* used by bitstream routines
 */
var bitstream_maxbyte int
var bitstream_bitaddr int
var bitstream_buf []byte

/* initialize bitstream for output */
func initOutputBitStream(buf []byte, buflen int) {
	bitstream_bitaddr = 0
	bitstream_maxbyte = buflen
	bitstream_buf = buf
}

/* initialize bitstream for input */
func initInputBitStream(buf []byte, buflen int) {
	bitstream_bitaddr = 0
	bitstream_maxbyte = buflen
	bitstream_buf = buf
}

/*
* read from bit stream. used only from 1 bit to 8 bits
* this is a base routine
 */
func readInputBitStreamBody(bwidth int) uint32 {
	mod := bitstream_bitaddr % 8
	byteaddr := bitstream_bitaddr / 8
	/* return if excess */
	if byteaddr >= bitstream_maxbyte {
		return 0
	}

	if bwidth >= 1 && bwidth <= 8 {
		b1 := uint((uint(bitstream_buf[byteaddr]) & uint(saacproto_modifymask_first[mod][bwidth])) >> uint(mod))
		b2 := uint((uint(bitstream_buf[byteaddr+1]) & uint(saacproto_modifymask_second[mod][bwidth])) << uint(8-mod))
		bitstream_bitaddr += bwidth
		return uint32(b1 | b2)
	} else {
		return 0
	}
}

/*
*  read from bit stream. used from 1 bit to 32 bits
*
 */
func readInputBitStream(bwidth int) uint32 {
	if bwidth <= 0 {
		return 0
	} else if bwidth >= 1 && bwidth <= 8 {
		return readInputBitStreamBody(bwidth)
	} else if bwidth >= 9 && bwidth <= 16 {
		first := readInputBitStreamBody(8)
		second := readInputBitStreamBody(bwidth - 8)
		return first + (second << 8)
	} else if bwidth >= 17 && bwidth <= 24 {
		first := readInputBitStreamBody(8)
		second := readInputBitStreamBody(8)
		third := readInputBitStreamBody(bwidth - 8)
		return first + (second << 8) + (third << 16)
	} else if bwidth >= 25 && bwidth <= 32 {
		first := readInputBitStreamBody(8)
		second := readInputBitStreamBody(8)
		third := readInputBitStreamBody(8)
		forth := readInputBitStreamBody(bwidth - 8)
		return first + (second << 8) + (third << 16) + (forth << 24)
	}
	return 0
}

/*
* write to a bitstream. only used from 1 bit to 8 bits
* this is a base routine.
 */
func writeOutputBitStreamBody(bwidth int, b byte) int {
	mod := bitstream_bitaddr % 8
	byteaddr := bitstream_bitaddr / 8
	/* return error if excess */
	if bitstream_maxbyte <= (byteaddr + 1) {
		return -1
	}
	bitstream_buf[byteaddr] &= byte(saacproto_modifymask_first[mod][bwidth])
	bitstream_buf[byteaddr] |= byte((int(b) << uint(mod)) & saacproto_modifymask_first[mod][bwidth])
	bitstream_buf[byteaddr+1] &= byte(saacproto_modifymask_second[mod][bwidth])
	bitstream_buf[byteaddr+1] |= byte((int(b) >> uint(8-mod)) & saacproto_modifymask_second[mod][bwidth])
	bitstream_bitaddr += bwidth
	return byteaddr + 1
}

/*
* write to a bitstream. used from 1 bits to 32 bits
* returns -1 if error or buffer excession
 */
func writeOutputBitStream(bwidth int, dat uint) int {
	var ret int = 0
	if bwidth <= 0 {
		return -1
	} else if bwidth >= 1 && bwidth <= 8 {
		if writeOutputBitStreamBody(bwidth, byte(dat)) < 0 {
			return -1
		}
	} else if bwidth > 8 && bwidth <= 16 {
		if writeOutputBitStreamBody(8, byte(dat&0xff)) < 0 {
			return -1
		}
		if writeOutputBitStreamBody(bwidth-8, byte((dat>>8)&0xff)) < 0 {
			return -1
		}
	} else if bwidth > 16 && bwidth <= 24 {
		if writeOutputBitStreamBody(8, byte(dat&0xff)) < 0 {
			return -1
		}
		if writeOutputBitStreamBody(8, byte((dat>>8)&0xff)) < 0 {
			return -1
		}
		if writeOutputBitStreamBody(bwidth-16, byte((dat>>16)&0xff)) < 0 {
			return -1
		}
	} else if bwidth > 24 && bwidth <= 32 {
		if writeOutputBitStreamBody(8, byte(dat&0xff)) < 0 {
			return -1
		}
		if writeOutputBitStreamBody(8, byte((dat>>8)&0xff)) < 0 {
			return -1
		}
		if writeOutputBitStreamBody(8, byte((dat>>16)&0xff)) < 0 {
			return -1
		}
		if writeOutputBitStreamBody(bwidth-24, byte((dat>>24)&0xff)) < 0 {
			return -1
		}
	} else {
		return -1
	}
	return ret
}

func Encode64(in []byte) []byte {
	var i int
	var use_bytes int
	var address int = 0
	len := len(in)
	out := make([]byte, len)
	out[0] = 0
	for i = 0; ; i += 3 {
		var in1 byte
		var in2 byte
		var in3 byte
		var out1 byte
		var out2 byte
		var out3 byte
		var out4 byte
		if i >= len {
			break
		}

		if i >= (len - 1) { /* the last letter ( to be thrown away ) */
			in1 = in[i] & 0xff
			in2 = 0
			in3 = 0
			use_bytes = 2
		} else if i >= (len - 2) { /* the last 2 letters ( process only 1 byte)*/
			in1 = in[i] & 0xff
			in2 = in[i+1] & 0xff
			in3 = 0
			use_bytes = 3
		} else { /* there are more or equal than 3 letters */
			in1 = in[i] & 0xff
			in2 = in[i+1] & 0xff
			in3 = in[i+2] & 0xff
			use_bytes = 4
		}
		out1 = ((in1 & 0xfc) >> 2) & 0x3f
		out2 = ((in1 & 0x03) << 4) | (((in2 & 0xf0) >> 4) & 0x0f)
		out3 = ((in2 & 0x0f) << 2) | (((in3 & 0xc0) >> 6) & 0x03)
		out4 = (in3 & 0x3f)
		if use_bytes >= 2 {
			out[address] = base64_charset[out1]
			address++
			out[address] = base64_charset[out2]
			address++
			out[address] = 0
		}
		if use_bytes >= 3 {
			out[address] = base64_charset[out3]
			address++
			out[address] = 0
		}
		if use_bytes >= 4 {
			out[address] = base64_charset[out4]
			address++
			out[address] = 0
		}
	}

	return out
}

func Decode64(in []byte) []byte {
	var in1 byte
	var in2 byte
	var in3 byte
	var in4 byte
	var out1 byte
	var out2 byte
	var out3 byte
	var use_bytes int
	var address int = 0

	len := len(in)
	out := make([]byte, len)
	var i int
	for i = 0; ; i += 4 {
		if in[i] == 0 {
			break
		} else if in[i+1] == 0 { /* the last letter */
			break
		} else if in[i+2] == 0 { /* the last 2 letters */
			in1 = base64_reversecharset[in[i]]
			in2 = base64_reversecharset[in[i+1]]
			in3 = 0
			in4 = 0
			use_bytes = 1
		} else if in[i+3] == 0 { /* the last  3 letters */
			in1 = base64_reversecharset[in[i]]
			in2 = base64_reversecharset[in[i+1]]
			in3 = base64_reversecharset[in[i+2]]
			in4 = 0
			use_bytes = 2
		} else { /* process 4 letters */
			in1 = base64_reversecharset[in[i]]
			in2 = base64_reversecharset[in[i+1]]
			in3 = base64_reversecharset[in[i+2]]
			in4 = base64_reversecharset[in[i+3]]
			use_bytes = 3
		}
		out1 = (in1 << 2) | (((in2 & 0x30) >> 4) & 0x0f)
		out2 = ((in2 & 0x0f) << 4) | (((in3 & 0x3c) >> 2) & 0x0f)
		out3 = ((in3 & 0x03) << 6) | (in4 & 0x3f)
		if use_bytes >= 1 {
			out[address] = out1
			address++
		}
		if use_bytes >= 2 {
			out[address] = out2
			address++
		}
		if use_bytes >= 3 {
			out[address] = out3
			address++
		}
		if use_bytes != 3 {
			break
		}
	}

	return out
}

func abs(val int) int {
	return int(math.Abs(float64(val)))
}

func jDecode(src []byte, key int) (decoded []byte, decodedlen int) {
	decoded = make([]byte, len(src))
	var sum byte = 0
	var i int
	var srclen = len(src)
	decodedlen = srclen - 1
	if decodedlen == 0 {
		return /* return error if length is 0 */
	}
	sum = src[abs(key%(decodedlen))]
	for i = 0; i < srclen; i++ {
		if abs((key % (decodedlen))) > i {
			decoded[i] = src[i] - byte(int(sum)*int((i*i)%3))
		}
		if abs((key % (decodedlen))) < i {
			decoded[i-1] = src[i] - byte(int(sum)*int((i*i)%7))
		}
	}
	for i = 0; i < decodedlen; i++ {
		if ((key % 7) == (i % 5)) || ((key % 2) == (i % 2)) {
			decoded[i] = ^decoded[i]
		}
	}

	return
}

func jEncode(src []byte, key int, maxencodedlen int) (encoded []byte, encodedlen int) {
	var sum byte = 0
	var i int
	srclen := len(src)
	encoded = make([]byte, srclen)
	encodedlen = srclen
	if srclen+1 > maxencodedlen {
		encodedlen = maxencodedlen
		for i = 0; i < encodedlen; i++ {
			encoded[i] = src[i]
		}
	}
	if srclen+1 <= maxencodedlen {
		encodedlen = srclen + 1
		for i = 0; i < srclen; i++ {
			sum = sum + src[i]
			if ((key % 7) == (i % 5)) || ((key % 2) == (i % 2)) {
				src[i] = ^src[i]
			}
		}
		for i = 0; i < encodedlen; i++ {
			if abs((key % srclen)) > i {
				encoded[i] = src[i] + byte(int(sum)*int((i*i)%3))
			}

			if abs((key % srclen)) == i {
				encoded[i] = sum
			}
			if abs((key % srclen)) < i {
				encoded[i] = src[i-1] + byte(int(sum)*int((i*i)%7))
			}
		}
	}

	return
}
