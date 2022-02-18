package iso8583

import (
	"bytes"
	"errors"
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
	v := EncodeGBK(value)
	//检查域长度
	if len(v) > fd.max {
		return fmt.Errorf("%s:域最大[%d]实际[%d]", fd.Name(), fd.max, len(v))
	}

	if fd.lenWidth == 0 && len(v) != fd.max {
		return fmt.Errorf("%s:域数据未满足定长要求", fd.Name())
	}

	return nil
}

//Encode 组包
func (fd *fieldDef) Encode(bf *bytes.Buffer, value string) int {
	start := bf.Len()
	var (
		dataLen int
		sLen    string
		data    []byte
	)

	gbkValue := EncodeGBK(value)
	if fd.lenWidth > 0 {
		dataLen = len(gbkValue)
	} else {
		dataLen = fd.max
	}

	switch fd.valueAttr {
	case NORMAL:
		data = gbkValue
	case BCDL, BCDR:
		data = EncodeBCD(value, fd.valueAttr, (dataLen+1)/2)
	case BITS:
		dataLen = dataLen / 2
		data = Str2Hex([]byte(value))
	default:
		data = gbkValue
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
	return bf.Len() - start
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
		vl := len(value)
		if fd.valueAttr == BITS {
			vl /= 2
		}
		sLen := fmt.Sprintf("%0*d", fd.lenWidth, vl)
		log.Infof("%s: Len[%s] Val[%s]", fd.Name(), sLen, value)
	} else {
		log.Infof("%s: Val[%s]", fd.Name(), value)
	}
}

//Fielder 域属性
type Fielder interface {
	Check(value string) error                  //域检查
	Encode(bf *bytes.Buffer, value string) int //域组包
	Decode(br *bytes.Reader) string            //域解包
	Name() string                              //域名称
	Print(value string)                        //打印
}

//ConfigDef 8583报文结构定义
type ConfigDef struct {
	bitLen        int
	msgTypeConfig Fielder
	fieldsConfig  map[uint]Fielder
}

//Pack 8583组包
func (iso *ConfigDef) Pack(data map[int]string) (res []byte, err error) {

	var (
		buffer     bytes.Buffer
		msgTypeLen int
	)

	if iso.msgTypeConfig != nil {
		// 存在报文头 取0域进行赋值
		if value, ok := data[0]; ok {
			iso.msgTypeConfig.Print(value)
			err = iso.msgTypeConfig.Check(value)
			if err != nil {
				return
			}
			msgTypeLen = iso.msgTypeConfig.Encode(&buffer, value)
		} else {
			err = errors.New("缺失报文类型")
			log.Info(err)
			return
		}
	}

	bitMap := make([]byte, iso.bitLen)

	buffer.Write(bitMap)

	filedNum := uint(iso.bitLen * 8)

	for i := uint(2); i <= filedNum; i++ {
		if value, ok := data[int(i)]; ok {
			if fieldConfig, ok := iso.fieldsConfig[i]; ok {
				//域配置存在
				BitSet(&bitMap, i)
				if i > 64 {
					//第二位图存在
					BitSet(&bitMap, 1)
				}
				fieldConfig.Print(value)
				err = fieldConfig.Check(value)
				if err != nil {
					return
				}
				fieldConfig.Encode(&buffer, value)
			}
		}

	}
	res = buffer.Bytes()
	copy(res[msgTypeLen:msgTypeLen+len(bitMap)], bitMap)

	return

}

//Unpack 8583解包,从报文头开始
func (iso *ConfigDef) Unpack(msg []byte) (res map[int]string, msgLen int, err error) {
	//获取位图
	res = make(map[int]string)
	stream := bytes.NewReader(msg)

	if iso.msgTypeConfig != nil {
		//存在报文类型
		value := iso.msgTypeConfig.Decode(stream)
		iso.msgTypeConfig.Print(value)
		if err != nil {
			log.Info("解析报文类型失败", err)
			return
		}
		res[0] = value
	}

	bitMap := make([]byte, iso.bitLen)

	stream.Read(bitMap)
	filedNum := uint(iso.bitLen * 8)

	for i := uint(2); i <= filedNum; i++ {
		if BitExist(bitMap, i) {
			if fieldConfig, ok := iso.fieldsConfig[i]; ok {
				value := fieldConfig.Decode(stream)
				fieldConfig.Print(value)
				err = fieldConfig.Check(value)
				if err != nil {
					return
				}
				res[int(i)] = value
			} else {
				return nil, 0, fmt.Errorf("域[%d]配置不存在", i)
			}
		}
	}

	msgLen = int(stream.Size()) - stream.Len() //8583报文体长度

	return
}
