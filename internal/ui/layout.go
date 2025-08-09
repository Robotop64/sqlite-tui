package ui

import "fyne.io/fyne/v2"

type MinVBox struct {
	MinWidth float32
}

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
}

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

type BHBox struct{}

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

type BVBox struct{}

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
