package position

import (
	"fmt"
	"strconv"
	"strings"
)

type Position struct {
	X float64
	Y float64
}

func (p *Position) Enc() string {
	return fmt.Sprintf("%f:%f", p.X, p.Y)
}

func ParsePosition(s string) (*Position, error) {
	splitVar := strings.Split(s, ":")
	if len(splitVar) != 2 {
		return nil, fmt.Errorf("invalid position: %q", s)
	}
	x, err := strconv.ParseFloat(splitVar[0], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid float x: %q", splitVar[0])
	}
	y, err := strconv.ParseFloat(splitVar[1], 64)
	if err != nil {
		return nil, fmt.Errorf("invalid float y: %q", splitVar[1])
	}

	return &Position{
		X: x,
		Y: y,
	}, nil
}
