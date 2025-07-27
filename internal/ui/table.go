package ui

import (
	"fmt"
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
	FocusTable
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
	Headers    bool
	Scrollbars [2]bool // vertical, horizontal
	//border options
	OuterBorders      [4]bool // top, right, bottom, left
	VertColSepBorders bool
	HorHeaderBorders  [2]bool // top, bottom
}

type ViewPortMode int

const (
	ViewportNormal ViewPortMode = iota
	ViewportFill
)

type Viewport struct {
	Offset int
	Size   int
	Mode   ViewPortMode
}

type TableWidget struct {
	Title   string
	Headers []string
	columns []any

	FocusMode FocusMode
	Focus     TableFocus

	VertViewPort  Viewport
	HorizViewPort Viewport

	Style TableStyle
}

func (t *TableWidget) SetHeaders(headers ...string) {
	t.Headers = headers
}
func (t *TableWidget) SetColumns(columns ...any) error {
	rows := len(columns[0].([]any))
	for _, col := range columns {
		if len(col.([]any)) != rows {
			return fmt.Errorf("all columns must have the same number of rows")
		}
	}
	t.columns = columns
	return nil
}

func (t *TableWidget) Render() string {

	tableDims := t.Style.MaxDimensions
	if t.Style.Title {
		tableDims.Height -= 2 // 1 for the title, 1 for the top border
	}
	if t.Style.Scrollbars[0] {
		tableDims.Width -= 2 // 2 for the vertical scrollbar
	}
	if t.Style.Scrollbars[1] {
		tableDims.Height -= 2 // 2 for the horizontal scrollbar
	}

	//resolve width
	setViewportSizes(t)
	// headers := getHeaders(t)
	// columns := getColumns(t)

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
		scrollWidth := tableDims.Width - 2
		fill := float32(len(t.Headers)) / float32(t.HorizViewPort.Size)
		offset := float32(len(t.Headers)) / float32(t.HorizViewPort.Offset)
		horizScrollbar = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), false, true, true, true).
			SetString(
				generateHorScollbar(scrollWidth, fill, offset),
			).
			Render()
	}

	var vertScrollbar string = ""
	if t.Style.Scrollbars[0] {
		scrollHeight := tableDims.Height - 2
		if t.Style.HorHeaderBorders[0] {
			scrollHeight -= 1 // 1 for the upper header border
		}
		if t.Style.HorHeaderBorders[1] {
			scrollHeight -= 1 // 1 for the bottom header border
		}
		border := lipgloss.RoundedBorder()
		if t.Style.Title {
			border.TopRight = "┤"
		}
		fill := float32(len(t.Headers)) / float32(t.VertViewPort.Size)
		offset := float32(len(t.Headers)) / float32(t.VertViewPort.Offset)
		vertScrollbar = lipgloss.NewStyle().
			Border(border, true, true, true, false).
			SetString(
				generateVertScrollbar(scrollHeight, fill, offset),
			).
			Render()
	}

	table := lgTable.New().
		Border(generateBorders(t.Style)).
		Width(tableDims.Width).
		Height(tableDims.Height).
		Headers(getHeaders(t)...)

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
		if len(horizontalPane) > 0 {
			verticalPane = append(verticalPane, lipgloss.JoinHorizontal(lipgloss.Left, horizontalPane...))
		}
	}

	view := lipgloss.JoinVertical(
		lipgloss.Top,
		verticalPane...,
	)

	return view
}

// helper functions ======================

func generateBorders(style TableStyle) lipgloss.Border {
	base := lipgloss.RoundedBorder()

	if style.Title {
		base.TopLeft = "├"
		base.TopRight = "┤"
	}
	if style.Scrollbars[0] {
		base.TopRight = "┬"
		base.BottomRight = "┴"
	}
	if style.Scrollbars[1] {
		base.BottomLeft = "├"
		base.BottomRight = "┤"
	}
	if style.Scrollbars[0] && style.Scrollbars[1] {
		base.BottomRight = "┼"
	}

	return base
}

func generateHorScollbar(length int, fill float32, offset float32) string {
	if length <= 0 {
		return ""
	}
	base := "|" + strings.Repeat("-", length-2) + "|"
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
	return base[:fillStart] + strings.Repeat("█", fillLength) + base[fillStart+fillLength:]
}

func generateVertScrollbar(length int, fill float32, offset float32) string {
	if length <= 0 {
		return ""
	}

	base := make([]string, 0, length)
	base = append(base, "┬")
	base = append(base[:1], strings.Split(strings.Repeat("│", length-2), "")...)
	base = append(base, "┴")

	fillLength := int(float32(length) * fill)
	if fillLength <= 0 {
		return strings.Join(base, "\n")
	}
	fillStart := int(float32(length) * offset)
	if fillStart >= length {
		return strings.Join(base, "\n")
	}
	if fillStart+fillLength > length {
		fillLength = length - fillStart
	}
	base = append(base[:fillStart], append(strings.Split(strings.Repeat("█", fillLength), ""), base[fillStart+fillLength:]...)...)

	return strings.Join(base, "\n")
}

func setViewportSizes(t *TableWidget) {
	// Vertical viewport
	t.VertViewPort.Size = t.Style.MaxDimensions.Height - 1 // 1 for the bottom border
	if t.Style.Title {
		t.VertViewPort.Size -= 2 // 1 for the title label, 1 for the bottom border
	}
	if t.Style.Headers {
		t.VertViewPort.Size -= 1 // 1 for the header label row
		if t.Style.HorHeaderBorders[0] {
			t.VertViewPort.Size -= 1 // 1 for the upper header border
		}
		if t.Style.HorHeaderBorders[1] {
			t.VertViewPort.Size -= 1 // 1 for the bottom header border
		}
	}
	if t.Style.Scrollbars[1] {
		t.VertViewPort.Size -= 2 // 2 for the horizontal scrollbar
	}

	// Horizontal viewport
	if t.HorizViewPort.Mode == ViewportFill {
		allowedWidth := t.Style.MaxDimensions.Width - 2 // 2 for the borders
		allowedHeaders := t.Headers[t.HorizViewPort.Offset:]

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
}

func getHeaders(t *TableWidget) []string {

	return t.Headers[t.HorizViewPort.Offset : t.HorizViewPort.Offset+t.HorizViewPort.Size]
}

func getColumns(t *TableWidget) [][]string {

	// underfill := t.VertViewPort.Size - len(t.Columns)

	columns := make([][]string, 0, t.HorizViewPort.Size)
	for i := 0; i < t.HorizViewPort.Size; i++ {
		columns[i] = make([]string, t.VertViewPort.Size)
	}

	// expander := func() int {
	// 	if t.VertViewPort.Offset+t.VertViewPort.Size < len(t.Columns) {
	// 		return 1
	// 	}
	// 	return 0
	// }()

	// for i := t.HorizViewPort.Offset; i < t.HorizViewPort.Offset+t.HorizViewPort.Size; i++ {
	// 	column, ok := t.Columns[i].(Column[any])
	// 	if !ok {
	// 		continue
	// 	}

	// 	view := column.Cells[t.VertViewPort.Offset : t.VertViewPort.Offset+t.VertViewPort.Size]

	// 	for j := range view {
	// 		if str, ok := view[j].(string); ok {
	// 			columns[i][j] = str
	// 		} else {
	// 			columns[i][j] = ""
	// 		}
	// 	}
	// }

	// 	if t.VertViewPort.Offset > 0 {
	// 		columns = append([][]string{{"..."}}, columns...)
	// 	}
	// 	if expander > 0 {
	// 		columns = append(columns, []string{"..."})
	// 	}
	// }

	return columns
}
