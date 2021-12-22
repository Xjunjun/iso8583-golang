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
	BitLen  int        `yaml:"bit_len"`  //位图长度
	MsgType FieldCfg   `yaml:"msg_type"` //报文类型
	Fields  []FieldCfg `yaml:"fields"`   //域配置
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

// 设置域配置
func newField(cfg FieldCfg) (f Fielder, err error) {
	var (
		lenAttr   int
		valueAttr int
	)
	lenAttr, err = getAttrID(cfg.LenAttr)
	if err != nil {
		return
	}

	valueAttr, err = getAttrID(cfg.ValueAttr)
	if err != nil {
		return
	}
	switch cfg.Type {
	case "number":
		f = NewNumField(cfg.FieldID, lenAttr, cfg.LenWidth, valueAttr, cfg.Max)
	case "binary":

		f = NewBinField(cfg.FieldID, lenAttr, cfg.LenWidth, valueAttr, cfg.Max)
	case "track":
		f = NewTrackField(cfg.FieldID, lenAttr, cfg.LenWidth, valueAttr, cfg.Max)
	default:
		f = NewTextField(cfg.FieldID, lenAttr, cfg.LenWidth, valueAttr, cfg.Max)
	}
	return
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

	//是否存在0域 报文类型
	if cfgfile.MsgType.Max > 0 {
		fieldCfg := cfgfile.MsgType
		if ios8583Tmp.msgTypeConfig, err = newField(fieldCfg); err != nil {
			return
		}
	}

	for _, fieldCfg := range cfgfile.Fields {
		if ios8583Tmp.fieldsConfig[uint(fieldCfg.FieldID)], err = newField(fieldCfg); err != nil {
			return
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
