package pollen

import (
	"fmt"
	"io"
	"log"
	"text/template"
)

type Board struct {
	cards  map[Position]*GardenCard
	tokens map[Position]*PollinatorToken
}

type Position struct {
	X float64
	Y float64
}

func NewBoard(token *PollinatorToken) *Board {
	b := &Board{
		cards:  make(map[Position]*GardenCard),
		tokens: make(map[Position]*PollinatorToken),
	}
	b.PlaceStarterToken(token)
	return b
}

func (b *Board) PlaceStarterToken(token *PollinatorToken) {
	b.tokens[Position{0, 0}] = token
}

// XXX finish implementing this
func (b *Board) MustPlayToken() error {
	if len(b.cards) == 0 {
		return nil
	}
	for position, _ := range b.cards {
		nw := Position{position.X + 0.5, position.Y + 0.5}
		sw := Position{position.X + 0.5, position.Y - 0.5}
		ne := Position{position.X - 0.5, position.Y + 0.5}
		se := Position{position.X - 0.5, position.Y - 0.5}
		switch {
		case b.CanPlayToken(nw) == nil:
			return fmt.Errorf("must play token at position %v", nw)
		case b.CanPlayToken(sw) == nil:
			return fmt.Errorf("must play token at position %v", sw)
		case b.CanPlayToken(ne) == nil:
			return fmt.Errorf("must play token at position %v", ne)
		case b.CanPlayToken(se) == nil:
			return fmt.Errorf("must play token at position %v", se)
		}
	}
	return nil
}

func (b *Board) PlayToken(position Position, token *PollinatorToken) error {
	err := b.CanPlayToken(position)
	if err != nil {
		return fmt.Errorf("cannot play token: %w", err)
	}
	b.tokens[position] = token
	return nil
}

func (b *Board) CanPlayToken(position Position) error {
	if _, present := b.tokens[position]; present {
		return fmt.Errorf("token already exists at position %v", position)
	}
	//Add verification that position is a whole number
	ne := Position{position.X + 0.5, position.Y + 0.5}
	se := Position{position.X + 0.5, position.Y - 0.5}
	nw := Position{position.X - 0.5, position.Y + 0.5}
	sw := Position{position.X - 0.5, position.Y - 0.5}
	_, swPresent := b.cards[sw]
	_, nwPresent := b.cards[nw]
	_, nePresent := b.cards[ne]
	_, sePresent := b.cards[se]
	switch {
	case swPresent && nwPresent:
	case sePresent && nePresent:
	case swPresent && sePresent:
	case nwPresent && nePresent:
	default:
		return fmt.Errorf("there is not two adjacent cards at position %v", position)
	}
	return nil
}

func (b *Board) PlayCard(position Position, card *GardenCard) error {
	err := b.CanPlayCard(position)
	if err != nil {
		return fmt.Errorf("cannot play card: %w", err)
	}
	b.cards[position] = card
	return nil
}

func (b *Board) CanPlayCard(position Position) error {
	_, present := b.cards[position]
	if present {
		return fmt.Errorf("card already exists at position %v", position)
	}

	//Add verification that position is a Half number
	ne := Position{position.X + 0.5, position.Y + 0.5}
	se := Position{position.X + 0.5, position.Y - 0.5}
	nw := Position{position.X - 0.5, position.Y + 0.5}
	sw := Position{position.X - 0.5, position.Y - 0.5}
	_, swPresent := b.tokens[sw]
	_, nwPresent := b.tokens[nw]
	_, nePresent := b.tokens[ne]
	_, sePresent := b.tokens[se]
	switch {
	case swPresent, nwPresent, nePresent, sePresent:
	default:
		return fmt.Errorf("position %v does not have an adjacent token", position)
	}
	return nil
}

// XXX this is not working missing some playable locations
func (b *Board) CardLocationsPlayable() map[Position]struct{} {
	positions := map[Position]struct{}{}
	for position := range b.tokens {
		nw := Position{position.X + 0.5, position.Y + 0.5}
		sw := Position{position.X + 0.5, position.Y - 0.5}
		ne := Position{position.X - 0.5, position.Y + 0.5}
		se := Position{position.X - 0.5, position.Y - 0.5}

		if b.CanPlayCard(nw) == nil {
			positions[nw] = struct{}{}
		}
		if b.CanPlayCard(sw) == nil {
			positions[sw] = struct{}{}
		}
		if b.CanPlayCard(ne) == nil {
			positions[ne] = struct{}{}
		}
		if b.CanPlayCard(se) == nil {
			positions[se] = struct{}{}
		}
	}
	return positions
}

func (b *Board) Render(w io.Writer, p *Player) error {
	boardTmpl := template.New("board")
	boardTmpl = boardTmpl.Funcs(template.FuncMap{
		"tokenStyle": func(token *PollinatorToken, position Position) string {
			tokenStyle := fmt.Sprintf(`background-color: %s; left: %dpx; bottom: %dpx;`,
				token.Type.Color(), int((position.X*2*25)+15), int((position.Y*2*25)+15))
			log.Println(tokenStyle)
			return tokenStyle
		},
		"cardStyle": func(card *GardenCard, position Position) string {
			tokenStyle := fmt.Sprintf(`background-color: %s; left: %dpx; bottom: %dpx;`,
				card.Color, int(position.X*2*25), int(position.Y*2*25))
			log.Println(tokenStyle)
			return tokenStyle
		},
		"playableStyle": func(position Position, o interface{}) string {
			tokenStyle := fmt.Sprintf(`background-color: orange; opacity: 0.5; left: %dpx; bottom: %dpx;`,
				int(position.X*2*25), int(position.Y*2*25))
			log.Println(tokenStyle)
			return tokenStyle
		},
	})

	boardTmpl = template.Must(boardTmpl.Parse(`
	{{ $debug :=.Debug }}
	{{ $player :=.Player}}
	<!DOCTYPE html>
		<head>
			<style>
                .center {
                    background-color: red;
					position: relative;
					left: 50%;
					top: 50%;
                }
				.token {
					width: 20px;
					height: 20px;
					border-radius: 100%;
					position: absolute;
					font-size: 4px;
				}
				.card {
					width: 50px;
					height: 50px;
					position: absolute;
				}
				.playableCard {
					width: 50px;
					height: 50px;
					position: absolute;
				}
				.board {
					position:absolute;
					top:0px;
					right:0px;
					bottom:0px;
					left:0px;
				}
				.centered {
					position: absolute;
					top: 50%;
					left: 50%;
					transform: translate(-50%, -50%);
					color: #00FF00;
					font-size: 5px;
				}
			</style>
		</head>
		<body>
            <div class="board">
				<div class="center">
					{{range $position, $card :=.Cards}}
						<div class="card" style="{{ cardStyle $card $position }}">
                            <div>
								<img class="card" src="/images/{{ $card.Name }}.png" title="{{ $card.Name }}">
								{{ if $debug }}
									<div class="centered"> Position {{ $position }}</div>
								{{end}}
								</img>
							</div>
						</div>
					{{end}}
					{{range $position, $token :=.Tokens}}
						<div class="token" style="{{ tokenStyle $token $position }}">
							{{$token.Type}}
						</div>
					{{end}}
					{{range $position, $empty := .PlayableCards}}
						<div class="playableCard" style="{{ playableStyle $position 0 }}">
                            <div>
								<img class="card" src="/images/Back_{{$player.Color}}.png">
									{{ if $debug }}
										<div class="centered"> Position {{ $position }}</div>
									{{end}}
								</img>
							</div>
						</div>
					{{end}}
				</div>
			</div>
		<body>
	<html>
	`))

	return boardTmpl.Execute(w, struct {
		Cards         map[Position]*GardenCard
		Tokens        map[Position]*PollinatorToken
		PlayableCards map[Position]struct{}
		Debug         bool
		Player        *Player
	}{b.cards, b.tokens, b.CardLocationsPlayable(), false, p})
}
