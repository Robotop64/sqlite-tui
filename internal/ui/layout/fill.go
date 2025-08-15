package layout

import "fyne.io/fyne/v2"

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
