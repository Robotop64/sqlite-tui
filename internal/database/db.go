package database

// "github.com/mattn/go-sqlite3"
// "database/sql"

type filter struct {
	Columns    []string
	MaxNumRows int
	SortCol    string //which column to sort by
	SortDir    uint8  // 0: None, 1: Ascending, 2: Descending
}

type filterMode int8

const (
	Composite filterMode = iota
	Single
)

type fieldMode int8

const (
	like fieldMode = iota
	exact
	contains
)
