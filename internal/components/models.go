package components

type ListModel[T any] struct {
	Items    []T
	Selected int
	Focused  int
}
