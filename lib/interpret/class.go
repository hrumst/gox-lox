package interpret

import (
	"fmt"
	"github.com/hrumst/gox-lox/lib/parse"
	"github.com/hrumst/gox-lox/lib/scan"
)

type LoxClass struct {
	interpreter *Interpreter
	declaration *parse.StmtClass
	methods     map[string]scan.LoxCallable
	superClass  *LoxClass
}

func NewLoxClass(
	interpreter *Interpreter,
	declaration *parse.StmtClass,
	superClass *LoxClass,
	methods map[string]scan.LoxCallable,
) *LoxClass {
	return &LoxClass{
		interpreter: interpreter,
		declaration: declaration,
		methods:     methods,
		superClass:  superClass,
	}
}

func (l *LoxClass) String() string {
	return fmt.Sprintf("[class] %s", l.declaration.Name.Lexeme)
}

func (l *LoxClass) Call(args []*scan.LoxValue) (*scan.LoxValue, error) {
	instance := NewLoxClassInstance(l)
	initializer := l.findMethod("init")
	if initializer != nil {
		initFunc := initializer.(*LoxFunction)
		if _, err := initFunc.bind(instance).Call(args); err != nil {
			return nil, err
		}
	}
	return scan.NewClassInstanceLoxValue(instance), nil
}

func (l *LoxClass) Arity() int {
	initializer := l.findMethod("init")
	if initializer != nil {
		return initializer.Arity()
	}
	return 0
}

func (l *LoxClass) findMethod(name string) scan.LoxCallable {
	if method, ok := l.methods[name]; ok {
		return method
	}
	if l.superClass != nil {
		return l.superClass.findMethod(name)
	}
	return nil
}
