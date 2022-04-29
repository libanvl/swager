package union

type Union[T any, U any] struct {
	value_t *T
	value_u *U
}

func New[T any, U any]() Union[T, U] {
	return Union[T, U]{}
}

func (u Union[T, U]) WithT(value T) Union[T, U] {
	u.value_t = &value
	u.value_u = nil
	return u
}

func (u Union[T, U]) WithU(value U) Union[T, U] {
	u.value_t = nil
	u.value_u = &value
	return u
}

func (u Union[T, U]) Get() (*T, *U) {
	return u.value_t, u.value_u
}

func (u Union[T, U]) GetT() *T {
	return u.value_t
}

func (u Union[T, U]) GetU() *U {
	return u.value_u
}
