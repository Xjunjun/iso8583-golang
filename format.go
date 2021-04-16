package iso8583

import (
	"golang.org/x/text/encoding/simplifiedchinese"
)

const hextable = "0123456789ABCDEF"

//EncodeGBK  string转GBK编码
func EncodeGBK(src string) []byte {
	render, _ := simplifiedchinese.GB18030.NewEncoder().Bytes([]byte(src))
	return render
}

//DecodeGBK GBK解码string
func DecodeGBK(src []byte) string {
	render, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(src)
	return string(render)
}

//BitSet 添加位图
func BitSet(bitmap *[]byte, bit uint) {
	index := (bit - 1) >> 3
	pos := (0x08 - bit&0x07) & 0x07
	if int(bit) <= len(*bitmap)*8 {
		(*bitmap)[index] |= (1 << pos)
	}

}

//BitExist 判断位图是否存在
func BitExist(bitmap []byte, bit uint) bool {
	index := (bit - 1) >> 3
	pos := (0x08 - bit&0x07) & 0x07
	if int(bit) <= len(bitmap)*8 {
		return (bitmap)[index]&(1<<pos) != 0
	}
	return false
}

//Str2Hex ascii码拓展
func Str2Hex(str []byte) []byte {
	var i, nlen, nhigh, nlow int
	nlen = len(str)
	var deststr []byte = make([]byte, (nlen+1)/2)
	for i = 0; i < nlen; i += 2 {
		nhigh = int(str[i])
		if nhigh > 0x39 {
			nhigh -= 0x37
		} else {
			nhigh -= 0x30
		}
		if i == (nlen - 1) {
			nlow = 0x00
		} else {
			nlow = int(str[i+1])
		}
		if nlow > 0x39 {
			nlow -= 0x37
		} else {
			nlow -= 0x30
		}
		deststr[i/2] = (byte)((nhigh << 4) | (nlow & 0x0f))
	}
	return deststr
}

//Hex2Str 扩展字符串转byte数组
func Hex2Str(hexstring []byte) string {
	out := make([]byte, len(hexstring)*2)
	n := func(dst, src []byte) int {
		j := 0
		for _, v := range src {
			dst[j] = hextable[v>>4]
			dst[j+1] = hextable[v&0x0f]
			j += 2
		}
		return len(src) * 2
	}(out, hexstring)
	return string(out[:n])
}

//DecodeBCD  BCD解码
func DecodeBCD(src []byte, dType int, dataLen int) string {
	ts := Hex2Str(src)
	dst := make([]byte, dataLen)
	if len(ts) >= dataLen {
		switch dType {
		case BCDL:
			copy(dst, ts)
		case BCDR:
			copy(dst, ts[len(ts)-dataLen:])
		default:
			//默认均为右靠
			copy(dst, ts[len(ts)-dataLen:])
		}
	}
	return DecodeGBK(dst)

}

//EncodeBCD BCD编码
func EncodeBCD(data string, dType int, dataLen int) []byte {

	res := make([]byte, dataLen*2)

	tmp := EncodeGBK(data)

	if len(tmp) <= len(res) {

		switch dType {
		case BCDL:
			copy(res, tmp)
		case BCDR:
			copy(res[dataLen*2-len(tmp):], tmp)
		default:
			//默认均为右靠
			copy(res[dataLen*2-len(tmp):], tmp)
		}
	}

	return Str2Hex(res)
}

//Pack 8583组包
func Pack(data map[int]string, msgType []byte) (res []byte, err error) {
	res, err = ios8583.Pack(data, msgType)
	return
}

//Unpack 8583解包,从位图开始
func Unpack(msg []byte) (res map[int]string, err error) {
	res, err = ios8583.Unpack(msg)
	return
}
