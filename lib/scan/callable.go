package scan

type LoxCallable interface {
	Arity() int
	Call(args []*LoxValue) (interface{}, error)
}
