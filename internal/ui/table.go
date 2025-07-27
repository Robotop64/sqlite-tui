package ui

import (
	"strings"

	"github.com/Robotop64/sqlite-tui/internal/style"
	lipgloss "github.com/charmbracelet/lipgloss"
	lgTable "github.com/charmbracelet/lipgloss/table"
)

type Column[T any] struct {
	Type  T
	Cells []T
}

type TableFocus struct {
	Column int
	Row    int
}

type FocusMode int

const (
	FocusCell FocusMode = iota
	FocusColumn
	FocusRow
)

type ColumnAlignment int

const (
	AlignLeft ColumnAlignment = iota
	AlignCenter
	AlignRight
)

type TableStyle struct {
	BaseStyle     lipgloss.Style
	ColumnWidths  []int
	MaxDimensions Dimensions
	//alignment
	ColumnAlignments []ColumnAlignment

	//display options
	Title      bool
	Scrollbars [2]bool // vertical, horizontal
	//border options
	OuterBorders      [4]bool // top, right, bottom, left
	VertColSepBorders bool
	HorHeaderBorders  [2]bool // top, bottom
}

type Viewport struct {
	Offset int
	Size   int
}

type TableWidget struct {
	Title   string
	Headers []string
	Columns []any

	FocusMode FocusMode
	Focus     TableFocus

	VertViewPort  Viewport
	HorizViewPort Viewport

	Style TableStyle
}

func generateBorders(style TableStyle) lipgloss.Border {
	base := lipgloss.NormalBorder()

	if style.Title {
		base.TopLeft = "├"
		base.TopRight = "┤"
	}
	if style.Scrollbars[1] {
		base.BottomLeft = "├"
		base.BottomRight = "┤"
	}
	// if style.Scrollbars[1] {
	// 	base.BottomLeft = "├"
	// 	base.BottomRight = "┤"
	// }

	return base
}

func generateHorScollbar(length int, fill float32, offset float32) string {
	if length <= 0 {
		return ""
	}
	base := ">" + strings.Repeat("-", length-2) + "<"
	fillLength := int(float32(length) * fill)
	if fillLength <= 0 {
		return base
	}
	fillStart := int(float32(length) * offset)
	if fillStart >= length {
		return base
	}
	if fillStart+fillLength > length {
		fillLength = length - fillStart
	}
	return base[:fillStart] + "(" + strings.Repeat("⎕", fillLength-2) + ")" + base[fillStart+fillLength:]
}

func (t *TableWidget) Render() string {

	//resolve width
	// setViewportSizes(t)
	// headers := getHeaders(t)
	// columns := getColumns(t)

	// vertScrollbar := " X\nX\nX "

	var titleSec string = ""
	if t.Style.Title {
		titleSec = style.Title.
			Width(t.Style.MaxDimensions.Width-2).
			Border(lipgloss.RoundedBorder(), true, true, false).
			AlignHorizontal(lipgloss.Center).
			SetString(t.Title).
			Render()
	}

	var horizScrollbar string = ""
	if t.Style.Scrollbars[1] {
		horizScrollbar = lipgloss.NewStyle().
			Width(t.Style.MaxDimensions.Width-2).
			Border(lipgloss.RoundedBorder(), false, true, true, true).
			AlignHorizontal(lipgloss.Center).
			SetString(
				generateHorScollbar(t.Style.MaxDimensions.Width-2, 0.5, 0.1),
			).
			Render()
	}

	var vertScrollbar string = ""
	if t.Style.Scrollbars[0] {
	}

	table := lgTable.New().
		Border(generateBorders(t.Style)).
		Width(t.Style.MaxDimensions.Width).
		Height(t.Style.MaxDimensions.Height).
		Headers(t.Headers...).
		Rows([][]string{
			{" ", " ", " "},
			{" ", " ", " "},
			{" ", " ", " "},
			{" ", " ", " "},
			{" ", " ", " "},
			{" ", " ", " "},
			{" ", " ", " "},
			{" ", " ", " "},
			{" ", " ", " "},
			{" ", " ", " "},
			{" ", " ", " "},
			{" ", " ", " "},
			{" ", " ", " "},
		}...)

	layout := [][]string{
		{titleSec},
		{table.Render(), vertScrollbar},
		{horizScrollbar},
	}

	var verticalPane []string
	for _, row := range layout {
		var horizontalPane []string
		for _, cell := range row {
			if cell != "" {
				horizontalPane = append(horizontalPane, cell)
			}
		}
		verticalPane = append(verticalPane, lipgloss.JoinHorizontal(lipgloss.Left, horizontalPane...))
	}

	view := lipgloss.JoinVertical(
		lipgloss.Top,
		verticalPane...,
	)

	return view
}

func setViewportSizes(t *TableWidget) {
	// Vertical viewport
	t.VertViewPort.Size = t.Style.MaxDimensions.Height - 3 - 1 // 3 for the header, 1 for the bottom border

	// Horizontal viewport
	allowedWidth := t.Style.MaxDimensions.Width - 2 // 2 for the borders
	allowedHeaders := t.Headers[t.VertViewPort.Offset:]

	headerWidth := 0

	allowedCols := len(t.Headers)
	for i, header := range allowedHeaders {
		newWidth := headerWidth + len(header) + 1 // +1 for right cell border //TODO replace with styled cell width
		if newWidth > allowedWidth {
			allowedCols = i
			break
		}
		headerWidth = newWidth
	}
	t.HorizViewPort.Size = allowedCols
}

func getHeaders(t *TableWidget) []string {

	leftExpander := t.HorizViewPort.Offset > 0
	rightExpander := t.HorizViewPort.Offset+t.HorizViewPort.Size < len(t.Headers)

	headers := make([]string, 0, t.HorizViewPort.Size)

	if leftExpander {
		headers = append(headers, "...")
	}

	rightExpanderSize := 0
	if rightExpander {
		rightExpanderSize = 1
	}

	headers = append(headers, t.Headers[t.HorizViewPort.Offset:t.HorizViewPort.Offset+t.HorizViewPort.Size-rightExpanderSize]...)

	if rightExpander {
		headers = append(headers, "...")
	}
	return headers
}

// func getColumns(t *TableWidget) [][]string {

// 	// underfill := t.VertViewPort.Size - len(t.Columns)

// 	columns := make([][]string, 0, t.HorizViewPort.Size)

// 	// expander := func() int {
// 	// 	if t.VertViewPort.Offset+t.VertViewPort.Size < len(t.Columns) {
// 	// 		return 1
// 	// 	}
// 	// 	return 0
// 	// }()

// 	for i := t.HorizViewPort.Offset; i < t.HorizViewPort.Offset+t.HorizViewPort.Size; i++ {
// 		column, ok := t.Columns[i].(Column[any])
// 		if !ok {
// 			continue
// 		}

// 		view := column.Cells[t.VertViewPort.Offset : t.VertViewPort.Offset+t.VertViewPort.Size]

// 		for j := range view {
// 			if str, ok := view[j].(string); ok {
// 				columns[i][j] = str
// 			} else {
// 				columns[i][j] = ""
// 			}
// 		}
// 	}

// 	// 	if t.VertViewPort.Offset > 0 {
// 	// 		columns = append([][]string{{"..."}}, columns...)
// 	// 	}
// 	// 	if expander > 0 {
// 	// 		columns = append(columns, []string{"..."})
// 	// 	}
// 	// }

// 	return columns
// }
