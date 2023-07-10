package object

import (
	"MonkeyPL/src/config"
)

type Integer struct {
	value int
}

func (i Integer) Type() config.ObjectType { return INTEGER }
func (i Integer) Value() int              { return i.value }
func NewInteger(value int) Integer {
	return Integer{value: value}
}
