package pollen

type Color int

const (
	Purple Color = iota
	Green
	Pink
	Orange
)

func (c Color) String() string {
	switch c {
	case Purple:
		return "Purple"
	case Green:
		return "Green"
	case Pink:
		return "Pink"
	case Orange:
		return "Orange"
	default:
		panic("unknown color")
	}
}
