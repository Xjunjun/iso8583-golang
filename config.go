package iso8583

import (
	"fmt"
	"os"
	"strings"

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
	ios8583 ConfigDef
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

//NewConfig 生成8583模板
func NewConfig(input string) (iso8583Template *ConfigDef, err error) {
	var cfgfile BitConfig

	rd, err := os.OpenFile(input, os.O_RDONLY, 0600)
	if err == nil {
		defer rd.Close()
		err = yaml.NewDecoder(rd).Decode(&cfgfile)
	} else {
		//作为字符串处理
		stream := strings.NewReader(input)
		err = yaml.NewDecoder(stream).Decode(&cfgfile)
	}
	if err != nil {
		return
	}

	ios8583Tmp := ConfigDef{}
	ios8583Tmp.fieldsConfig = make(map[uint]Fielder)
	if cfgfile.BitLen != 64 && cfgfile.BitLen != 128 {
		err = fmt.Errorf("位图长度不合法[%d]", cfgfile.BitLen)
		return
	}
	ios8583Tmp.bitLen = cfgfile.BitLen >> 3

	for _, fieldCfg := range cfgfile.Fields {
		var (
			lenAttr   int
			valueAttr int
		)
		lenAttr, err = getAttrID(fieldCfg.LenAttr)
		if err != nil {
			return
		}

		valueAttr, err = getAttrID(fieldCfg.ValueAttr)
		if err != nil {
			return
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

	iso8583Template = &ios8583Tmp
	return
}

//Default 生成转换器
func Default(cfgPath string) error {

	tlp, err := NewConfig(cfgPath)
	if err != nil {
		return err
	}
	ios8583 = *tlp

	return nil

}

func init() {
	ios8583 = ConfigDef{
		bitLen: 0,
	}
	ios8583.fieldsConfig = make(map[uint]Fielder)

}
