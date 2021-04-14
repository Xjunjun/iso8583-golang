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
	trackData := strings.ReplaceAll(value, "=", "D")
	td.fieldDef.Encode(bf, trackData)
}

//Decode 解包
func (td *TrackField) Decode(br *bytes.Reader) string {
	trackData := td.fieldDef.Decode(br)
	return strings.ReplaceAll(trackData, "D", "=")
}
