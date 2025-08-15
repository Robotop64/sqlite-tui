package layout

type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

func DirFromStr(dir string) Direction {
	switch dir {
	case "horizontal":
		return Horizontal
	case "vertical":
		return Vertical
	default:
		return Horizontal
	}
}
