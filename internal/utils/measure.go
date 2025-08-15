package utils

import "fyne.io/fyne/v2"

func Max[T ~int | ~float32 | ~float64](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T ~int | ~float32 | ~float64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

type Dimensions[T ~int | ~float32 | ~float64] struct {
	Width  T
	Height T
}

func (d Dimensions[T]) ToFyneSize() fyne.Size {
	return fyne.NewSize(float32(d.Width), float32(d.Height))
}
func FyneToDimensions(size fyne.Size) Dimensions[float32] {
	return Dimensions[float32]{
		Width:  float32(size.Width),
		Height: float32(size.Height),
	}
}

func FitSize[T ~int | ~float32 | ~float64](item, container Dimensions[T]) Dimensions[T] {
	return Dimensions[T]{
		Width:  Min(item.Width, container.Width),
		Height: Min(item.Height, container.Height),
	}
}
