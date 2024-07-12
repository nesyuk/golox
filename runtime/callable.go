package runtime

type LoxCallable interface {
	call(Iner) interface{}
}
