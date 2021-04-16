package iso8583

import "fmt"

//Logger 日志接口
type Logger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
}

type logDefault struct {
}

func (std *logDefault) Info(args ...interface{}) {
	fmt.Println(args...)
}

func (std *logDefault) Infof(template string, args ...interface{}) {
	fmt.Printf(template, args...)
	fmt.Println()
}

var log Logger

func init() {
	log = &logDefault{}
}

//SetLogger 设置自定义日志
func SetLogger(logusr Logger) {
	log = logusr
}
