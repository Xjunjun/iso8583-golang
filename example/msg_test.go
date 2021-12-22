package test

import (
	"testing"

	"github.com/Xjunjun/iso8583-golang"
)

type LogUser struct {
	t *testing.T
}

func (log *LogUser) Info(args ...interface{}) {
	log.t.Log(args...)
}

func (log *LogUser) Infof(template string, args ...interface{}) {
	log.t.Logf(template, args...)
}

func Test8583PosMsg(t *testing.T) {
	var err error
	err = iso8583.Default("8583-pos.yml")
	if err != nil {
		t.Log(err)
		return
	}

	iso8583.SetLogger(&LogUser{t})

	iso8583Data := make(map[int]string)
	iso8583Data[0] = "0200"
	iso8583Data[2] = "6212142400000000105"
	iso8583Data[3] = "000000"
	iso8583Data[4] = "000000510000"
	iso8583Data[11] = "000484"
	iso8583Data[14] = "2907"
	iso8583Data[22] = "052"
	iso8583Data[23] = "001"
	iso8583Data[25] = "00"
	iso8583Data[26] = "12"
	iso8583Data[35] = "6212142400000000105=04141234"
	iso8583Data[41] = "10673470"
	iso8583Data[42] = "986474810800165"
	iso8583Data[46] = "290000000000000000000000000000000000000000218.204.252.14100620000000"
	iso8583Data[49] = "156"
	iso8583Data[53] = "0600000000000000"
	iso8583Data[55] = "9F26088560B4F34F53F7E49F2701809F101307010103A02012010A0100000500000725A5649F3704D159248F9F36021ACB950500000008009A032104079C01009F02060000005100005F2A02015682027C009F1A0201569F03060000000000009F330390C8C09F34033F00009F3501229F1E0831323334353637388408A0000003330101019F090200309F410400000002"
	iso8583Data[60] = "220000020006"
	iso8583Data[64] = "3541314239313946"

	out2, err := iso8583.Pack(iso8583Data)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("报文:", iso8583.Hex2Str(out2))

	unpack2, err := iso8583.Unpack(out2)
	if err != nil {
		t.Log(err)
		return
	}

	for i := 0; i < 128; i++ {
		if value, ok := unpack2[i]; ok {
			t.Log("域:", i, "值:", value)
		}
	}
}

func Test8583InsMsg(t *testing.T) {
	var err error
	config, err := iso8583.NewConfig("8583-ins.yml")
	if err != nil {
		t.Log(err)
		return
	}

	iso8583.SetLogger(&LogUser{t})

	iso8583Data := make(map[int]string)
	iso8583Data[0] = "0420"
	iso8583Data[2] = "6212142400000000105"
	iso8583Data[3] = "000000"
	iso8583Data[4] = "000000510000"
	iso8583Data[7] = "0415162807"
	iso8583Data[11] = "000484"
	iso8583Data[14] = "2907"
	iso8583Data[18] = "5411"
	iso8583Data[22] = "052"
	iso8583Data[23] = "001"
	iso8583Data[25] = "00"
	iso8583Data[26] = "12"
	iso8583Data[32] = "09860000"
	iso8583Data[33] = "48210000"
	iso8583Data[35] = "6212142400000000105=04141234"
	iso8583Data[37] = "110516036413"
	iso8583Data[41] = "10673470"
	iso8583Data[42] = "986474810800165"
	iso8583Data[43] = "这是一个测试商户哦                      "
	iso8583Data[49] = "156"
	iso8583Data[53] = "0600000000000000"
	iso8583Data[55] = "9F26088560B4F34F53F7E49F2701809F101307010103A02012010A0100000500000725A5649F3704D159248F9F36021ACB950500000008009A032104079C01009F02060000005100005F2A02015682027C009F1A0201569F03060000000000009F330390C8C09F34033F00009F3501229F1E0831323334353637388408A0000003330101019F090200309F410400000002"
	iso8583Data[60] = "000005"
	iso8583Data[128] = "3541314239313946"

	out2, err := config.Pack(iso8583Data)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("报文:", iso8583.Hex2Str(out2))

	unpack2, err := config.Unpack(out2)
	if err != nil {
		t.Log(err)
		return
	}

	for i := 0; i <= 128; i++ {
		if value, ok := unpack2[i]; ok {
			t.Log("域:", i, "值:", value)
		}
	}
}
