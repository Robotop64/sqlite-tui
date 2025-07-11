package utils

type FocusElement struct {
	Left  *FocusElement
	Right *FocusElement
	Up    *FocusElement
	Down  *FocusElement
}

func (m *FocusElement) Move(direction string) *FocusElement {
	var next *FocusElement = nil
	switch direction {
	case "left":
		if m.Left != nil {
			next = m.Left
		}
	case "right":
		if m.Right != nil {
			next = m.Right
		}
	case "up":
		if m.Up != nil {
			next = m.Up
		}
	case "down":
		if m.Down != nil {
			next = m.Down
		}
	}
	if next == nil {
		return m
	} else {
		return next
	}
}

func FocusChain(elements []*FocusElement, direction string) {
	switch direction {
	case "horizontal":
		for i := 0; i < len(elements)-1; i++ {
			elements[i].Right = elements[i+1]
			elements[i+1].Left = elements[i]
		}
	case "vertical":
		for i := 0; i < len(elements)-1; i++ {
			elements[i].Down = elements[i+1]
			elements[i+1].Up = elements[i]
		}
	}
}

func FocusGrid(elements [][]*FocusElement) {
	for row := 0; row < len(elements); row++ {
		for col := 0; col < len(elements[row]); col++ {
			sel := elements[row][col]

			// Reset links to avoid overwriting
			sel.Left, sel.Right, sel.Up, sel.Down = nil, nil, nil, nil

			// Assign neighbors only if they exist
			if col-1 >= 0 && elements[row][col-1] != sel {
				sel.Left = elements[row][col-1]
			}
			if col+1 < len(elements[row]) && elements[row][col+1] != sel {
				sel.Right = elements[row][col+1]
			}
			if row-1 >= 0 && elements[row-1][col] != sel {
				sel.Up = elements[row-1][col]
			}
			if row+1 < len(elements) && elements[row+1][col] != sel {
				sel.Down = elements[row+1][col]
			}
		}
	}
}
