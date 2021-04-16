package iso8583

//TextField 文本域
type TextField struct {
	fieldDef
}

//NewTextField 创建文本域
func NewTextField(fieldID int, lenAttr int, lenWidth int, valueAttr int, max int) *TextField {
	td := TextField{
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
