package interpret

import (
	"fmt"
	"github.com/hrumst/gox-lox/lib/scan"
)

type LoxClassInstance struct {
	class  *LoxClass
	fields map[string]*scan.LoxValue
}

func NewLoxClassInstance(class *LoxClass) *LoxClassInstance {
	return &LoxClassInstance{
		fields: make(map[string]*scan.LoxValue),
		class:  class,
	}
}

func (li *LoxClassInstance) String() string {
	return fmt.Sprintf("[class instance] %s", li.class.declaration.Name.Lexeme)
}

func (li *LoxClassInstance) Get(name scan.Token) (*scan.LoxValue, error) {
	if value, ok := li.fields[name.Lexeme]; ok {
		return value, nil
	}
	if method := li.class.findMethod(name.Lexeme); method != nil {
		funcMethod := method.(*LoxFunction)
		return scan.NewCallableLoxValue(funcMethod.bind(li)), nil
	}
	return nil, NewRuntimeError(fmt.Sprintf("undefined property '%s'", name.Lexeme), &name)
}

func (li *LoxClassInstance) Set(name scan.Token, value *scan.LoxValue) {
	li.fields[name.Lexeme] = value
}
