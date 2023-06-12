package pollen

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"strings"
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
	positions := b.GetTokensMustPlay()
	if len(positions) == 0 {
		return nil
	}
	return fmt.Errorf("must play tokens at the following positions: %v", positions)
}

func (b *Board) GetTokensMustPlay() []Position {
	if len(b.cards) == 0 {
		return nil
	}
	tokensMustPlay := map[Position]struct{}{}
	for position, _ := range b.cards {
		nw := Position{position.X + 0.5, position.Y + 0.5}
		sw := Position{position.X + 0.5, position.Y - 0.5}
		ne := Position{position.X - 0.5, position.Y + 0.5}
		se := Position{position.X - 0.5, position.Y - 0.5}

		if b.CanPlayToken(nw) == nil {
			tokensMustPlay[nw] = struct{}{}
		}
		if b.CanPlayToken(sw) == nil {
			tokensMustPlay[sw] = struct{}{}
		}
		if b.CanPlayToken(ne) == nil {
			tokensMustPlay[ne] = struct{}{}
		}
		if b.CanPlayToken(se) == nil {
			tokensMustPlay[se] = struct{}{}
		}
	}
	var positions []Position
	for position := range tokensMustPlay {
		positions = append(positions, position)
	}
	return positions
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

var boardTmpl = template.Must(template.New("board").Funcs(template.FuncMap{
	"tokenStyle": func(token *PollinatorToken, position Position) string {
		tokenStyle := fmt.Sprintf(`background-color: %s; left: %dpx; bottom: %dpx;`,
			token.Type.Color(), int((position.X*2*25)+15), int((position.Y*2*25)+15))
		return tokenStyle
	},
	"cardStyle": func(card *GardenCard, position Position) string {
		tokenStyle := fmt.Sprintf(`background-color: %s; left: %dpx; bottom: %dpx;`,
			card.Color, int(position.X*2*25), int(position.Y*2*25))
		return tokenStyle
	},
	"playableStyle": func(position Position, o interface{}) string {
		tokenStyle := fmt.Sprintf(`background-color: orange; opacity: 0.5; left: %dpx; bottom: %dpx;`,
			int(position.X*2*25), int(position.Y*2*25))
		return tokenStyle
	},
	"handStyle": func(card GardenCard, p int) string {
		tokenStyle := fmt.Sprintf(`left: %dpx; bottom: 0px;`, p*75)
		return ""
		return tokenStyle
	},
}).Parse(`
	{{ $debug :=.Debug }}
	{{ $player :=.Player}}
	{{ $gameID :=.GameID}}
	{{ $isPlayerTurn :=.IsPlayerTurn}}
	<div class="board">
		<div class="center">
			{{range $position, $card :=.Cards}}
				<div class="card" style="{{ cardStyle $card $position }}">
					<div>
						<img class="card" src="/static/images/{{ $card.Name }}.png" title="{{ $card.Name }}">
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
			{{if $isPlayerTurn}}
				{{range $position, $empty := .PlayableCards}}
					<div class="playableCard" style="{{ playableStyle $position 0 }}" onclick='playCard("{{$gameID}}",getCardToPlay(),{{$position.X}},{{$position.Y}})'>
						<div>
							<img class="card" src="/static/images/Back_{{$player.Color}}.png">
								{{ if $debug }}
									<div class="centered"> Position {{ $position }}</div>
								{{end}}
							</img>
						</div>
					</div>
				{{ end }}
			{{ end}}
			<div class="bottom hand">
				<div class="cardHolder">
					{{range $cardNum, $card := .Hand}}
						<div class="handCard" id="{{ $card.ID }}" style="{{ handStyle $card $cardNum }}" onclick='setCardToPlay("{{ $card.ID }}")'>
							<div>
								<img class="handCard" id="{{ $card.ID }}_img" src="/static/images/{{ $card.Name }}.png">
									{{ if $debug }}
										<div class="centered"> Position {{ $cardNum }}</div>
									{{end}}
								</img>
							</div>
						</div>
					{{end}}
				</div>
			</div>
		</div>
	</div>
	`))

func (b *Board) Render(w io.Writer, p *Player, g *Game) error {
	buff := bytes.NewBuffer(nil)

	err := boardTmpl.Execute(buff, struct {
		Cards         map[Position]*GardenCard
		Tokens        map[Position]*PollinatorToken
		PlayableCards map[Position]struct{}
		Debug         bool
		Player        *Player
		GameID        string
		Hand          []GardenCard
		IsPlayerTurn  bool
	}{
		b.cards,
		b.tokens,
		b.CardLocationsPlayable(),
		false,
		p,
		g.id.String(),
		p.Hand,
		g.activePlayer().Username == p.Username,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "data: %s\n\n", base64.StdEncoding.EncodeToString(buff.Bytes()))
	return nil
}
