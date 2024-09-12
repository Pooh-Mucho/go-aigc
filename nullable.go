package aigc

type Nullable[T ~int32 | ~int64 | ~float32 | ~float64 | ~bool | ~string] struct {
	Value T
	Valid bool
}

func (n *Nullable[T]) Set(value T) {
	n.Valid = true
	n.Value = value
}

func (n *Nullable[T]) Unset() {
	var defaultValue T
	n.Valid = false
	n.Value = defaultValue
}

func NewNullable[T ~int32 | ~int64 | ~float32 | ~float64 | ~bool | ~string](value T) Nullable[T] {
	return Nullable[T]{Value: value, Valid: true}
}
