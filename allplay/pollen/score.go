package pollen

import (
	"github.com/gofrs/uuid"
	"github.com/hibooboo2/ggames/allplay/logger"
	"github.com/hibooboo2/ggames/allplay/pollen/cards"
	"github.com/hibooboo2/ggames/allplay/pollen/colors"
	"github.com/hibooboo2/ggames/allplay/pollen/position"
	"github.com/hibooboo2/ggames/allplay/pollen/token"
)

type GameScore struct {
	BeeMeeplesLeft       int
	JunebugMeeplesLeft   int
	ButterflyMeeplesLeft int
	Scores               map[colors.Color]*PlayerScore
}

type PlayerScore struct {
	c                colors.Color
	BeeMeeples       int
	JunebugMeeples   int
	ButterflyMeeples int
}

func NewGameScore(players int) *GameScore {
	gameScore := &GameScore{
		BeeMeeplesLeft:       10,
		JunebugMeeplesLeft:   10,
		ButterflyMeeplesLeft: 10,
		Scores: map[colors.Color]*PlayerScore{
			colors.Purple: {c: colors.Purple},
			colors.Green:  {c: colors.Green},
			colors.Pink:   {c: colors.Pink},
			colors.Orange: {c: colors.Orange},
		},
	}

	if players > 2 {
		gameScore.BeeMeeplesLeft = 16
		gameScore.JunebugMeeplesLeft = 16
		gameScore.ButterflyMeeplesLeft = 16
	}
	return gameScore
}

func (gs *GameScore) OutOfMeeples() bool {
	return gs.BeeMeeplesLeft == 0 || gs.JunebugMeeplesLeft == 0 || gs.ButterflyMeeplesLeft == 0
}

func (b *Board) UpdateScores() {
	logger.AtLevelf(logger.LBoard|logger.LToken|logger.LPosition|logger.LScore, "Calculating score for tokens\n\n")
	if b.Scores == nil {
		b.Scores = NewGameScore(b.Players)
	}

	for position, token := range b.tokens {
		if token.IsSurrounded() {
			logger.AtLevelf(logger.LBoard|logger.LToken|logger.LPosition|logger.LScore, "Token %v is surrounded\n\n", position)
			continue
		}

		if b.Scores.OutOfMeeples() {
			logger.AtLevelf(logger.LBoard|logger.LToken|logger.LPosition|logger.LScore, "Out of meeples\n\n")
			return
		}

		if b.TokenIsSurrounded(token.ID) {
			beeMeepleColor, juneBugMeepleColor, butterflyMeepleColor := b.ScoreToken(token.ID)
			playerScore, ok := b.Scores.Scores[beeMeepleColor]
			if ok {
				b.Scores.BeeMeeplesLeft--
				playerScore.BeeMeeples++
			}
			playerScore, ok = b.Scores.Scores[juneBugMeepleColor]
			if ok {
				b.Scores.JunebugMeeplesLeft--
				playerScore.JunebugMeeples++
			}
			playerScore, ok = b.Scores.Scores[butterflyMeepleColor]
			if ok {
				b.Scores.ButterflyMeeplesLeft--
				playerScore.ButterflyMeeples++
			}

			token.SetSurrounded(true)
			logger.AtLevelf(logger.LBoard|logger.LToken|logger.LPosition|logger.LScore, "Token %v is surrounded position %v", token.ID, position)
		}
	}
}

func (b *Board) ScoreToken(tokenID uuid.UUID) (colors.Color, colors.Color, colors.Color) {
	var p *position.Position
	var tokenType token.TokenType
	for p2, token := range b.tokens {
		if token.ID == tokenID {
			p = &p2
			tokenType = token.Type
			break
		}
	}

	if p == nil {
		return 0, 0, 0
	}

	logger.AtLevel(logger.LScore, "Calculating token meeple winners: ", p)

	//Add verification that position is a whole number
	ne := position.Position{p.X + 0.5, p.Y + 0.5}
	se := position.Position{p.X + 0.5, p.Y - 0.5}
	nw := position.Position{p.X - 0.5, p.Y + 0.5}
	sw := position.Position{p.X - 0.5, p.Y - 0.5}
	swCard, swPresent := b.cards[sw]
	nwCard, nwPresent := b.cards[nw]
	neCard, nePresent := b.cards[ne]
	seCard, sePresent := b.cards[se]
	if !(swPresent && nwPresent && nePresent && sePresent) {
		return 0, 0, 0
	}

	return scoreForColors([]*cards.GardenCard{swCard, nwCard, neCard, seCard}, tokenType)
}

func scoreForColors(cardsToScore []*cards.GardenCard, ct token.TokenType) (colors.Color, colors.Color, colors.Color) {
	scores := make(map[MeepleType]map[colors.Color]int)
	scores[BeeMeeple] = make(map[colors.Color]int)
	scores[JunebugMeeple] = make(map[colors.Color]int)
	scores[ButterflyMeeple] = make(map[colors.Color]int)

	for _, card := range cardsToScore {
		switch ct {
		case token.BeeToken:
			if card.Type == cards.Bee || card.Type == cards.Wild {
				scores[BeeMeeple][card.C] += card.Value
			}
		case token.JunebugToken:
			if card.Type == cards.Junebug || card.Type == cards.Wild {
				scores[JunebugMeeple][card.C] += card.Value
			}
		case token.ButterflyToken:
			if card.Type == cards.Butterfly || card.Type == cards.Wild {
				scores[ButterflyMeeple][card.C] += card.Value
			}
		case token.BeeJunebugToken:
			if card.Type == cards.Bee || card.Type == cards.Wild {
				scores[BeeMeeple][card.C] += card.Value
			}
			if card.Type == cards.Junebug || card.Type == cards.Wild {
				scores[JunebugMeeple][card.C] += card.Value
			}
		case token.BeeButterflyToken:
			if card.Type == cards.Bee || card.Type == cards.Wild {
				scores[BeeMeeple][card.C] += card.Value
			}
			if card.Type == cards.Butterfly || card.Type == cards.Wild {
				scores[ButterflyMeeple][card.C] += card.Value
			}
		case token.JunebugButterFlyToken:
			if card.Type == cards.Junebug || card.Type == cards.Wild {
				scores[JunebugMeeple][card.C] += card.Value
			}
			if card.Type == cards.Butterfly || card.Type == cards.Wild {
				scores[ButterflyMeeple][card.C] += card.Value
			}
		case token.BeeJunebugButterFlyToken:
			if card.Type == cards.Bee || card.Type == cards.Wild {
				scores[BeeMeeple][card.C] += card.Value
			}
			if card.Type == cards.Junebug || card.Type == cards.Wild {
				scores[JunebugMeeple][card.C] += card.Value
			}
			if card.Type == cards.Butterfly || card.Type == cards.Wild {
				scores[ButterflyMeeple][card.C] += card.Value
			}
		}
	}

	return getMaxColor(scores[BeeMeeple]), getMaxColor(scores[JunebugMeeple]), getMaxColor(scores[ButterflyMeeple])
}

func getMaxColor(colorScores map[colors.Color]int) colors.Color {
	var maxColor colors.Color
	var maxScore int
	for color, score := range colorScores {
		switch {
		case score == maxScore:
			maxColor = 0
		case score > maxScore:
			maxColor = color
			maxScore = score
		}
	}
	return maxColor
}
