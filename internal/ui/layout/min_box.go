package layout

import "fyne.io/fyne/v2"

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
