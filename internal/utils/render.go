package utils

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type View struct {
	Dim     Dimensions
	Content string
}

type Alignment int

const (
	Left Alignment = iota
	Center
	Right
	Top
	Bottom
)

func Overlay(back, front string, vertical, horizontal Alignment) (string, error) {
	bWidth, bHeight := lipgloss.Size(back)
	fWidth, fHeight := lipgloss.Size(front)

	if bWidth < fWidth ||
		bHeight < fHeight {
		return "", fmt.Errorf("the front view is larger than the back view")
	}

	xOffset, yOffset := getOffset(Dimensions{bWidth, bHeight}, Dimensions{fWidth, fHeight}, horizontal, vertical)

	backLines := strings.Split(back, "\n")
	frontLines := strings.Split(front, "\n")
	newLines := make([]string, len(backLines))

	for i, line := range backLines {
		if i < yOffset || i >= yOffset+fHeight {
			newLines[i] = line
			continue
		}

		newLines[i] += ansi.Cut(line, 0, xOffset)
		newLines[i] += frontLines[i-yOffset]
		newLines[i] += ansi.Cut(line, xOffset+fWidth, bWidth)
	}

	newContent := strings.Join(newLines, "\n")

	if lipgloss.Width(newContent) != bWidth ||
		lipgloss.Height(newContent) != bHeight {
		return "", fmt.Errorf("overlayed content does not match the dimensions of the back view")
	}

	return newContent, nil
}

func getOffset(backSize, frontSize Dimensions, hAlign, vAlign Alignment) (int, int) {

	var xOffset, yOffset int

	switch hAlign {
	case Left:
		xOffset = 0
	case Center:
		xOffset = (backSize.Width - frontSize.Width) / 2
	case Right:
		xOffset = backSize.Width - frontSize.Width
	}

	switch vAlign {
	case Top:
		yOffset = 0
	case Center:
		yOffset = (backSize.Height - frontSize.Height) / 2
	case Bottom:
		yOffset = backSize.Height - frontSize.Height
	}

	return xOffset, yOffset
}
