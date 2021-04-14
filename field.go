package iso8583

import (
	"bytes"
	"fmt"
	"strconv"
)

const (
	//BCDL 左靠bcd
	BCDL = iota
	//BCDR 右靠bcd
	BCDR
	//NORMAL 正常
	NORMAL
	//BITS 二进制
	BITS
)

//fieldDef  域定义
type fieldDef struct {
	fieldID   int    //域id
	lenAttr   int    //长度属性
	lenWidth  int    //长度域长度
	valueAttr int    //值属性
	max       int    //值最大长度
	name      string //域名称
}

func (fd *fieldDef) SetName() {
	fd.name = fmt.Sprintf("BID[%03d]", fd.fieldID)
}

//Name 域名称
func (fd *fieldDef) Name() string {
	return fd.name
}

//Check 域检查
func (fd *fieldDef) Check(value string) error {
	//检查域长度
	if len(value) > fd.max {
		return fmt.Errorf("%s:域最大[%d]实际[%d]", fd.Name(), fd.max, len(value))
	}

	return nil
}

//Encode 组包
func (fd *fieldDef) Encode(bf *bytes.Buffer, value string) {
	var (
		dataLen int
		sLen    string
		data    []byte
	)
	dataLen = len(value)
	switch fd.valueAttr {
	case NORMAL:
		data = EncodeGBK(value)
	case BCDL, BCDR:
		data = EncodeBCD(value, fd.valueAttr, (dataLen+1)/2)
	case BITS:
		dataLen = len(value) / 2
		data = Str2Hex([]byte(value))
	default:
		data = EncodeGBK(value)
	}

	if fd.lenWidth > 0 {
		sLen = fmt.Sprintf("%0*d", fd.lenWidth, dataLen)
		switch fd.lenAttr {
		case NORMAL:
			bf.WriteString(sLen)
		case BCDL, BCDR:
			bf.Write(EncodeBCD(sLen, fd.lenAttr, (fd.lenWidth+1)/2))
		}

	}
	bf.Write(data)

}

//Decode 解包
func (fd *fieldDef) Decode(br *bytes.Reader) string {
	var (
		realLen int //真实长度
	)
	if fd.lenWidth > 0 {
		switch fd.lenAttr {
		case NORMAL:
			tmp := make([]byte, fd.lenWidth)
			br.Read(tmp)
			realLen, _ = strconv.Atoi(string(tmp))
		case BCDL, BCDR:
			tmp := make([]byte, (fd.lenWidth+1)/2)
			br.Read(tmp)
			realLen, _ = strconv.Atoi(DecodeBCD(tmp, fd.lenAttr, fd.lenWidth))
		}
	} else {
		realLen = fd.max
	}

	//读取值域
	switch fd.valueAttr {
	case NORMAL:
		tmp := make([]byte, realLen)
		br.Read(tmp)
		return DecodeGBK(tmp)
	case BCDL, BCDR:
		tmp := make([]byte, (realLen+1)/2)
		br.Read(tmp)
		return DecodeBCD(tmp, fd.valueAttr, realLen)
	default:
		tmp := make([]byte, realLen)
		br.Read(tmp)
		return Hex2Str(tmp)
	}

}

//Print 打印域信息
func (fd *fieldDef) Print(value string) {
	if fd.lenWidth > 0 {
		sLen := fmt.Sprintf("%0*d", fd.lenWidth, len(value))
		log.Infof("%s: Len[%s] Val[%s]", fd.Name(), sLen, value)
	} else {
		log.Infof("%s: Val[%s]", fd.Name(), value)
	}
}

//Fielder 域属性
type Fielder interface {
	Check(value string) error              //域检查
	Encode(bf *bytes.Buffer, value string) //域组包
	Decode(br *bytes.Reader) string        //域解包
	Name() string                          //域名称
	Print(value string)                    //打印
}

//ios8583Def 8583报文结构定义
type ios8583Def struct {
	bitLen       int
	fieldsConfig map[uint]Fielder
}
