package iso8583

import (
	"fmt"
	"regexp"
)

//NumField 数值域
type NumField struct {
	fieldDef
}

//NewNumField 创建二进制域
func NewNumField(fieldID int, lenAttr int, lenWidth int, valueAttr int, max int) *NumField {
	nd := NumField{
		fieldDef: fieldDef{
			fieldID:   fieldID,
			lenAttr:   lenAttr,
			lenWidth:  lenWidth,
			valueAttr: valueAttr,
			max:       max,
		},
	}

	nd.SetName()

	return &nd
}

//Check 域检查
func (nd *NumField) Check(value string) error {

	//检查域长度
	if len(value) > nd.max {
		return fmt.Errorf("%s:域最大[%d]实际[%d]", nd.Name(), nd.max, len(value))
	}

	//检查域是否都是数字
	result, _ := regexp.MatchString("^\\d+$", value)
	if result == false {
		return fmt.Errorf("%s:域值非数字[%s]", nd.Name(), value)
	}

	return nil

}
