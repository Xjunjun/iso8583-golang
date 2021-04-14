package iso8583

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

//FieldCfg  域定义
type FieldCfg struct {
	FieldID   int    `yaml:"field_id"`   //域id
	LenAttr   string `yaml:"len_attr"`   //长度属性
	LenWidth  int    `yaml:"len_width"`  //长度域长度
	ValueAttr string `yaml:"value_attr"` //值属性
	Max       int    `yaml:"max"`        //值最大长度
	Type      string `yaml:"type"`       //域类型 根据类型选择不同对象
}

//BitConfig 域配置信息
type BitConfig struct {
	BitLen int        `yaml:"bit_len"` //位图长度
	Fields []FieldCfg `yaml:"fields"`  //域配置
}

var (
	ios8583 ios8583Def
)

func getAttrID(name string) (int, error) {
	switch name {
	case "BCDL":
		return BCDL, nil
	case "BCDR":
		return BCDR, nil
	case "NORMAL":
		return NORMAL, nil
	case "BITS":
		return BITS, nil
	}
	return 0, fmt.Errorf("未知类型[%s]", name)
}

//Default 生成转换器
func Default(cfgPath string) error {
	var cfgfile BitConfig

	rd, err := os.OpenFile(cfgPath, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}

	err = yaml.NewDecoder(rd).Decode(&cfgfile)
	if err != nil {
		return err
	}
	ios8583Tmp := ios8583Def{}
	ios8583Tmp.fieldsConfig = make(map[uint]Fielder)
	if cfgfile.BitLen != 64 && cfgfile.BitLen != 128 {
		return fmt.Errorf("位图长度不合法[%d]", cfgfile.BitLen)
	}
	ios8583Tmp.bitLen = cfgfile.BitLen >> 3

	for _, fieldCfg := range cfgfile.Fields {
		lenAttr, err := getAttrID(fieldCfg.LenAttr)
		if err != nil {
			return err
		}

		valueAttr, err := getAttrID(fieldCfg.ValueAttr)
		if err != nil {
			return err
		}
		switch fieldCfg.Type {
		case "number":
			ios8583Tmp.fieldsConfig[uint(fieldCfg.FieldID)] =
				NewNumField(fieldCfg.FieldID, lenAttr, fieldCfg.LenWidth, valueAttr, fieldCfg.Max)
		case "binary":
			ios8583Tmp.fieldsConfig[uint(fieldCfg.FieldID)] =
				NewBinField(fieldCfg.FieldID, lenAttr, fieldCfg.LenWidth, valueAttr, fieldCfg.Max)
		case "track":
			ios8583Tmp.fieldsConfig[uint(fieldCfg.FieldID)] =
				NewTrackField(fieldCfg.FieldID, lenAttr, fieldCfg.LenWidth, valueAttr, fieldCfg.Max)
		default:
			ios8583Tmp.fieldsConfig[uint(fieldCfg.FieldID)] =
				NewTextField(fieldCfg.FieldID, lenAttr, fieldCfg.LenWidth, valueAttr, fieldCfg.Max)
		}

	}

	ios8583 = ios8583Tmp

	return nil

}

func init() {
	ios8583 = ios8583Def{
		bitLen: 0,
	}
	ios8583.fieldsConfig = make(map[uint]Fielder)

}
