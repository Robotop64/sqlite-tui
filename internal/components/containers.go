package components

type Table struct {
	Headers   []string
	Rows      [][]string
	ActiveRow int
	ActiveCol int
}

type List struct {
	Items      []string
	ActiveItem int
}
