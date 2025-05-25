package database

type filter struct {
	Columns    []string
	MaxNumRows int
	SortCol    string //which column to sort by
	SortDir    uint8  // 0: None, 1: Ascending, 2: Descending
}
