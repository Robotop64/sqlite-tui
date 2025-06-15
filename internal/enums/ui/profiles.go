package ui

type UiFocus int

const (
	None UiFocus = iota
	TxtInput
	TxtEdit
	ProfileList
)
