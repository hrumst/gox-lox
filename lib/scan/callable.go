package scan

type LoxCallable interface {
	String() string
	Arity() int
	Call(args []*LoxValue) (*LoxValue, error)
}

type LoxClass LoxCallable

type LoxClassInstance interface {
	String() string
	Get(name Token) (*LoxValue, error)
	Set(name Token, value *LoxValue)
}
