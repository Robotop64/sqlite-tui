package layout

import "fyne.io/fyne/v2"

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
