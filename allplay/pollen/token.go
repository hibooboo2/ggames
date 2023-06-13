package pollen

import (
	"math/rand"

	"github.com/gofrs/uuid"
	"github.com/hibooboo2/ggames/allplay/logger"
)

type PollinatorToken struct {
	ID   uuid.UUID
	Type TokenType
}

type TokenType int

const (
	BeeToken TokenType = iota
	JunebugToken
	ButterflyToken
	BeeJunebugToken
	BeeButterflyToken
	JunebugButterFlyToken
	BeeJunebugButterFlyToken
)

func (t TokenType) String() string {
	switch t {
	case BeeToken:
		return "BeeToken"
	case JunebugToken:
		return "JunebugToken"
	case ButterflyToken:
		return "ButterflyToken"
	case BeeJunebugToken:
		return "BeeJunebugToken"
	case BeeButterflyToken:
		return "BeeButterflyToken"
	case JunebugButterFlyToken:
		return "JunebugButterFlyToken"
	case BeeJunebugButterFlyToken:
		return "BeeJunebugButterFlyToken"
	default:
		panic("unknown token type")
	}
}

func (t TokenType) Color() string {
	switch t {
	case BeeToken:
		return "#ff0000"
	case JunebugToken:
		return "#00ff00"
	case ButterflyToken:
		return "#0000ff"
	case BeeJunebugToken:
		return "#ff00ff"
	case BeeButterflyToken:
		return "#00ffff"
	case JunebugButterFlyToken:
		return "#ffff00"
	case BeeJunebugButterFlyToken:
		return "#9900ff"
	default:
		panic("unknown token type")
	}
}

type TokenBag struct {
	tokens []*PollinatorToken
}

func NewTokenBag() *TokenBag {
	tb := &TokenBag{}
	tb.tokens = createPollinatorTokens()
	rand.Shuffle(len(tb.tokens), func(i, j int) {
		tb.tokens[i], tb.tokens[j] = tb.tokens[j], tb.tokens[i]
	})
	return tb
}

func (tb *TokenBag) GetTokens(n int) []*PollinatorToken {
	if len(tb.tokens) == 0 {
		return nil
	}

	if len(tb.tokens) < n {
		n = len(tb.tokens)
	}

	return tb.tokens[:n]
}

func (tb *TokenBag) GetToken(tokenID uuid.UUID) *PollinatorToken {
	for _, t := range tb.tokens {
		if t.ID == tokenID {
			return t
		}
	}
	return nil
}

func (tb *TokenBag) ConsumeTokens(tokens []*PollinatorToken) {
	for _, token := range tokens {
		tb.ConsumeToken(token.ID)
	}
}

func (tb *TokenBag) ConsumeToken(tokenID uuid.UUID) {
	for i := range tb.tokens {
		if tb.tokens[i].ID == tokenID {
			tb.tokens = append(tb.tokens[:i], tb.tokens[i+1:]...)
			logger.Tokenln("There are now ", len(tb.tokens), " tokens left")
			return
		}
	}
}

func (tb *TokenBag) HasToken(tokenID uuid.UUID) bool {
	for i := range tb.tokens {
		if tb.tokens[i].ID == tokenID {
			return true
		}
	}
	return false
}

func createPollinatorTokens() []*PollinatorToken {
	singleTokenCreators := []func() *PollinatorToken{NewBeeToken, NewJunebugToken, NewButterflyToken}
	tokens := []*PollinatorToken{}
	for _, creator := range singleTokenCreators {
		tokens = append(tokens, NewTokenGroup(creator, 5)...)
	}
	doubleTokenCreators := []func() *PollinatorToken{NewBeeJunebugToken, NewBeeButterflyToken, NewJunebugButterFlyToken}
	for _, creator := range doubleTokenCreators {
		tokens = append(tokens, NewTokenGroup(creator, 9)...)
	}
	tokens = append(tokens, NewTokenGroup(NewBeeJunebugButterFlyToken, 2)...)
	return tokens
}

func NewTokenGroup(tokenCreator func() *PollinatorToken, n int) []*PollinatorToken {
	tokens := make([]*PollinatorToken, n)
	for i := 0; i < n; i++ {
		tokens[i] = tokenCreator()
	}
	return tokens
}

func NewBeeToken() *PollinatorToken {
	return &PollinatorToken{
		ID:   uuid.Must(uuid.NewV4()),
		Type: BeeToken,
	}
}

func NewJunebugToken() *PollinatorToken {
	return &PollinatorToken{
		ID:   uuid.Must(uuid.NewV4()),
		Type: JunebugToken,
	}
}

func NewButterflyToken() *PollinatorToken {
	return &PollinatorToken{
		ID:   uuid.Must(uuid.NewV4()),
		Type: ButterflyToken,
	}
}

func NewBeeJunebugToken() *PollinatorToken {
	return &PollinatorToken{
		ID:   uuid.Must(uuid.NewV4()),
		Type: BeeJunebugToken,
	}
}

func NewBeeButterflyToken() *PollinatorToken {
	return &PollinatorToken{
		ID:   uuid.Must(uuid.NewV4()),
		Type: BeeButterflyToken,
	}
}

func NewJunebugButterFlyToken() *PollinatorToken {
	return &PollinatorToken{
		ID:   uuid.Must(uuid.NewV4()),
		Type: JunebugButterFlyToken,
	}
}

func NewBeeJunebugButterFlyToken() *PollinatorToken {
	return &PollinatorToken{
		ID:   uuid.Must(uuid.NewV4()),
		Type: BeeJunebugButterFlyToken,
	}
}
