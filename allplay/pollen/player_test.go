package pollen

import "testing"

func TestPlayerEmptyDeckAndHand(t *testing.T) {
	p := NewPlayer("JAMES", 2, Green)

	for len(p.Hand) > 0 {
		card := p.Hand[0]
		c, err := p.PlayCard(card.ID)
		if err != nil {
			t.Error(err)
		}
		t.Error(c.Color, c.ID, c.Type, c.Value)
	}
}
