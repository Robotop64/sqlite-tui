package ui

import (
	"SQLite-GUI/internal/utils"
	"fmt"

	"fyne.io/fyne/v2"
	FContainer "fyne.io/fyne/v2/container"
	FBind "fyne.io/fyne/v2/data/binding"
	FWidget "fyne.io/fyne/v2/widget"
)

func NonValidatedEntry() *FWidget.Entry {
	entry := FWidget.NewEntry()
	entry.Validator = nil
	entry.Refresh()
	return entry
}

func NonValidatedEntryWithData(data FBind.String) *FWidget.Entry {
	entry := FWidget.NewEntryWithData(data)
	entry.Validator = nil
	entry.Refresh()
	return entry
}

// TODO make generic
func EditableTable(data *[][]string, dirtyRows *[]int) *FWidget.Table {
	type CellFocus struct {
		Row int
		Col int
	}
	var cell_edit *CellFocus
	var table *FWidget.Table

	table = FWidget.NewTable(
		func() (int, int) {
			return len(*data), len((*data)[0])
		},
		func() fyne.CanvasObject {
			lbl := FWidget.NewLabel("Placeholder")
			entry := FWidget.NewEntry()
			entry.Hide()
			return FContainer.NewStack(lbl, entry)
		},
		func(i FWidget.TableCellID, o fyne.CanvasObject) {
			cell := o.(*fyne.Container)
			val := (*data)[i.Row][i.Col]

			lbl := cell.Objects[0].(*FWidget.Label)
			entry := cell.Objects[1].(*FWidget.Entry)

			if lbl.Text != val {
				lbl.SetText(val)
			}
			if entry.Text != val {
				entry.SetText(val)
			}

			updateData := func(val string) {
				if (*data)[i.Row][i.Col] != val && !utils.Contains(*dirtyRows, i.Row) {
					*dirtyRows = append(*dirtyRows, i.Row)
					fmt.Println("Row marked as dirty:", i.Row)
				}

				//TODO Check if type casting is valid
				(*data)[i.Row][i.Col] = val
			}

			// entry.OnChanged = func(val string) {
			// 	updateData(val)
			// }

			entry.OnSubmitted = func(val string) {
				updateData(val)
				cell_edit = nil
				table.Refresh()
			}

			if cell_edit != nil && cell_edit.Row == i.Row && cell_edit.Col == i.Col {
				lbl.Hide()
				entry.Show()
			} else {
				lbl.Show()
				entry.Hide()
			}
		},
	)
	table.OnSelected = func(id FWidget.TableCellID) {
		cell_edit = &CellFocus{Row: id.Row, Col: id.Col}
		table.Refresh()
	}

	return table
}
