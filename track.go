package iso8583

import (
	"bytes"
	"strings"
)

//TrackField 二磁道,三磁道域
type TrackField struct {
	fieldDef
}

//NewTrackField 创建磁道域
func NewTrackField(fieldID int, lenAttr int, lenWidth int, valueAttr int, max int) *TrackField {
	td := TrackField{
		fieldDef: fieldDef{
			fieldID:   fieldID,
			lenAttr:   lenAttr,
			lenWidth:  lenWidth,
			valueAttr: valueAttr,
			max:       max,
		},
	}

	td.SetName()

	return &td
}

//Encode 组包
func (td *TrackField) Encode(bf *bytes.Buffer, value string) {
	var trackData string
	if td.valueAttr == BCDL || td.valueAttr == BCDR {
		//磁道信息bcd压缩时,需要将=转换为D
		trackData = strings.ReplaceAll(value, "=", "D")
	}

	td.fieldDef.Encode(bf, trackData)
}

/****************************
*  Decode 由于磁道信息可能
*	存在部分域加密的问题,因此
*	解包后不应该将D替换为等号,
*	应在磁道解密之后再替换,故
*	通用程序不在此重写Decode方法
******************************/
