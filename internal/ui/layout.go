package ui

type Element interface{} // Either a Widget or Style Element (Separator, Spacer, etc.)
type Position struct {
	X int
	Y int
}
type Dimensions struct {
	Width  int
	Height int
}

type Layout struct {
	Widgets    []Element
	Positions  []Position
	Dimensions []Dimensions
}

func (l *Layout) Render() string {
	return l.Widgets[0].(Widget).Render()
}
