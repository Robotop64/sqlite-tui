package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
)

type Fill struct{}

func (f *Fill) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	y := float32(0)
	for _, obj := range objects {
		min := obj.Size()
		obj.Resize(fyne.NewSize(size.Width, min.Height))
		obj.Move(fyne.NewPos(0, y))
		y += min.Height
	}
}

func (f *Fill) MinSize(objects []fyne.CanvasObject) fyne.Size {
	width := float32(0)
	height := float32(0)
	for _, obj := range objects {
		min := obj.MinSize()
		size := obj.Size()

		max := size.Max(min)

		if max.Width > width {
			width = max.Width
		}
		height += max.Height
	}
	return fyne.NewSize(width, height)
}

type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

func DirFromStr(dir string) Direction {
	switch dir {
	case "horizontal":
		return Horizontal
	case "vertical":
		return Vertical
	default:
		return Horizontal
	}
}

type WBox struct {
	Weights []float32 // w>0: weight, w=0: fill remaining space, w=-1: minimize
	Dir     Direction
}

func (l *WBox) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) == 0 || len(l.Weights) == 0 || len(objects) != len(l.Weights) {
		fmt.Println("Error: objects and weights length mismatch or empty")
		return
	}

	weighted_sizes := make([]float32, len(objects)) //calculated sizes of the objects along the directional axis
	total_cumulative_size := float32(0.0)
	cumulative_weight := float32(0.0)
	num_minimized := 0
	num_filling := 0

	available_size := func() float32 {
		switch l.Dir {
		case Horizontal:
			return size.Width
		case Vertical:
			return size.Height
		}
		return 0.0
	}()

	if available_size <= 0.0 {
		return
	}

	for i := 0; i < len(objects); i++ {
		weight := l.Weights[i]

		switch {
		case weight == -1.0:
			num_minimized++
			continue
		case weight == 0.0:
			num_filling++
			continue
		case weight > 0.0:
			cumulative_weight += weight
		}
	}

	if cumulative_weight > 1.0 || (cumulative_weight == 1.0 && (num_filling > 0 || num_minimized > 0)) {
		// - summed weights > 100%
		// - weighted elements fill available space (=100%) but minimized or filling are present
		fmt.Println("Error: the cumulative weight of the layout exceeds 100%")
		return
	}

	if num_minimized > 0 {
		for i := 0; i < len(objects); i++ {
			if l.Weights[i] == -1.0 {
				weighted_sizes[i] = func() float32 {
					switch l.Dir {
					case Horizontal:
						return objects[i].MinSize().Width
					case Vertical:
						return objects[i].MinSize().Height
					}
					return 0.0
				}()
				total_cumulative_size += weighted_sizes[i]
			}
		}
	}

	for i := 0; i < len(objects); i++ {
		if l.Weights[i] > 0.0 {
			weighted_sizes[i] = available_size * l.Weights[i]
			total_cumulative_size += weighted_sizes[i]
		}
	}

	if num_filling > 0 {
		for i := 0; i < len(objects); i++ {
			if l.Weights[i] == 0.0 {
				weighted_sizes[i] = (available_size - total_cumulative_size) / float32(num_filling)
			}
		}
	}

	if total_cumulative_size > available_size {
		fmt.Printf("Error: the cumulative size (%f) of the weighted objects exceeds the available space (%f)", total_cumulative_size, available_size)
		return
	}

	other_size := func() float32 {
		switch l.Dir {
		case Horizontal:
			return size.Height
		case Vertical:
			return size.Width
		}
		return 0.0
	}()

	x, y := float32(0), float32(0)
	for i := 0; i < len(objects); i++ {
		if objects[i] == nil {
			fmt.Println("Error: object at index", i, "is nil")
			continue
		}
		switch l.Dir {
		case Horizontal:
			objects[i].Move(fyne.NewPos(x, 0))
			x += weighted_sizes[i]
			objects[i].Resize(fyne.NewSize(weighted_sizes[i], other_size))
		case Vertical:
			objects[i].Move(fyne.NewPos(0, y))
			y += weighted_sizes[i]
			objects[i].Resize(fyne.NewSize(other_size, weighted_sizes[i]))
		}
	}
}

func (l *WBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) == 0 {
		return fyne.NewSize(0, 0)
	}

	minH := float32(0)
	minW := float32(0)
	sumW := float32(0)
	sumH := float32(0)

	for _, o := range objects {
		ms := o.MinSize()
		if ms.Height > minH {
			minH = ms.Height
		}
		if ms.Width > minW {
			minW = ms.Width
		}
		sumW += ms.Width
		sumH += ms.Height
	}

	switch l.Dir {
	case Horizontal:
		return fyne.NewSize(sumW, minH)
	case Vertical:
		return fyne.NewSize(minW, sumH)
	}

	return fyne.NewSize(0, 0)
}

//TODO: Refactor below into Boxes with Dir property

type MinVBox struct {
	MinWidth float32
} // VBox with minimum width

func (v *MinVBox) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	y := float32(0)
	for _, obj := range objects {
		min := obj.MinSize()
		obj.Resize(fyne.NewSize(size.Width, min.Height))
		obj.Move(fyne.NewPos(0, y))
		y += min.Height
	}
}

func (v *MinVBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	width := v.MinWidth
	height := float32(0)
	for _, obj := range objects {
		min := obj.MinSize()
		if min.Width > width {
			width = min.Width
		}
		height += min.Height
	}
	return fyne.NewSize(width, height)
}

type MinHBox struct {
	MinHeight float32
} // HBox with minimum height

func (h *MinHBox) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	x := float32(0)
	for _, obj := range objects {
		min := obj.MinSize()
		obj.Resize(fyne.NewSize(min.Width, size.Height))
		obj.Move(fyne.NewPos(x, 0))
		x += min.Width
	}
}

func (h *MinHBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	height := h.MinHeight
	width := float32(0)
	for _, obj := range objects {
		min := obj.MinSize()
		if min.Height > height {
			height = min.Height
		}
		width += min.Width
	}
	return fyne.NewSize(width, height)
}

type BHBox struct{} // Balanced HBox

func (b *BHBox) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	count := len(objects)
	if count == 0 {
		return
	}

	childWidth := size.Width / float32(count)
	x := float32(0)
	for _, obj := range objects {
		obj.Resize(fyne.NewSize(childWidth, size.Height))
		obj.Move(fyne.NewPos(x, 0))
		x += childWidth
	}
}

func (b *BHBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	// Minimum width = max of MinWidths * number of children
	// Minimum height = max MinHeight among all children
	count := len(objects)
	if count == 0 {
		return fyne.NewSize(0, 0)
	}

	maxMinWidth := float32(0)
	maxMinHeight := float32(0)
	for _, obj := range objects {
		min := obj.MinSize()
		if min.Width > maxMinWidth {
			maxMinWidth = min.Width
		}
		if min.Height > maxMinHeight {
			maxMinHeight = min.Height
		}
	}
	return fyne.NewSize(maxMinWidth*float32(count), maxMinHeight)
}

type BVBox struct{} //Balanced VBox

func (b *BVBox) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	count := len(objects)
	if count == 0 {
		return
	}

	childHeight := size.Height / float32(count)
	y := float32(0)
	for _, obj := range objects {
		obj.Resize(fyne.NewSize(size.Width, childHeight))
		obj.Move(fyne.NewPos(0, y))
		y += childHeight
	}
}

func (b *BVBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	count := len(objects)
	if count == 0 {
		return fyne.NewSize(0, 0)
	}

	maxMinWidth := float32(0)
	maxMinHeight := float32(0)
	for _, obj := range objects {
		min := obj.MinSize()
		if min.Width > maxMinWidth {
			maxMinWidth = min.Width
		}
		if min.Height > maxMinHeight {
			maxMinHeight = min.Height
		}
	}
	return fyne.NewSize(maxMinWidth, maxMinHeight*float32(count))
}
