package pollen

type CardType int

const (
	Bee CardType = iota
	Butterfly
	Junebug
	Wild
)

func (c CardType) String() string {
	switch c {
	case Bee:
		return "Bee"
	case Butterfly:
		return "Butterfly"
	case Junebug:
		return "Junebug"
	case Wild:
		return "Wild"
	default:
		panic("invalid card type")
	}
}
