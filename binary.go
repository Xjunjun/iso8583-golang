package iso8583

import (
	"fmt"
)

//BinField 二进制域
type BinField struct {
	fieldDef
}

//NewBinField 创建二进制域
func NewBinField(fieldID int, lenAttr int, lenWidth int, valueAttr int, max int) *BinField {
	bd := BinField{
		fieldDef: fieldDef{
			fieldID:   fieldID,
			lenAttr:   lenAttr,
			lenWidth:  lenWidth,
			valueAttr: valueAttr,
			max:       max,
		},
	}

	bd.SetName()

	return &bd
}

//Check 域校验
func (bd *BinField) Check(value string) error {
	if len(value)/2 > bd.max {
		return fmt.Errorf("%s:域最大[%d]实际[%d]", bd.Name(), bd.max, len(value)/2)
	}
	return nil
}
