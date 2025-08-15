package layout

import (
	"fmt"

	"fyne.io/fyne/v2"
)

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
