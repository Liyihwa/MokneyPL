package object

import "MonkeyPL/src/config"

/*
这里定义的是寄宿语言(MonkeyPL)中的变量用到的接口
*/
type ObjectValue interface{}

type Object interface {
	Type() config.ObjectType
}
