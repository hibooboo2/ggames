package token

import (
	"sync"

	"github.com/gofrs/uuid"
	"github.com/hibooboo2/ggames/allplay/logger"
	"github.com/hibooboo2/ggames/allplay/pollen/position"
)

type PollinatorToken struct {
	ID           uuid.UUID
	Position     *position.Position
	Type         TokenType
	isSurrounded bool
	l            sync.Mutex
	b            *TokenBag
	isUsed       bool
}

func (p *PollinatorToken) Play() {
	p.b.ConsumeToken(p.ID)
	p.isUsed = true
	logger.Tokenf("Token %s played", p.ID)
}

func (p *PollinatorToken) IsUsed() bool {
	return p.isUsed
}

func (p *PollinatorToken) IsSurrounded() bool {
	p.l.Lock()
	defer p.l.Unlock()
	pos := position.Position{}
	if p.Position != nil {
		pos = *p.Position
	}
	logger.AtLevelf(logger.LToken|logger.LScore, "%s %v is surrounded: %v", p.ID, pos, p.isSurrounded)
	return p.isSurrounded
}

func (p *PollinatorToken) SetSurrounded(surrounded bool) {
	p.l.Lock()
	p.isSurrounded = surrounded
	p.l.Unlock()
	p.IsSurrounded()
}

type TokenType int

const (
	BeeToken TokenType = 1 << iota
	JunebugToken
	ButterflyToken

	BeeJunebugToken       TokenType = BeeToken | JunebugToken
	BeeButterflyToken     TokenType = BeeToken | ButterflyToken
	JunebugButterFlyToken TokenType = JunebugToken | ButterflyToken

	BeeJunebugButterFlyToken TokenType = BeeToken | JunebugToken | ButterflyToken
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

func (t TokenType) Image() string {
	switch t {
	case BeeToken:
		return "Token_Bee.jpg"
	case JunebugToken:
		return "Token_Junebug.jpg"
	case ButterflyToken:
		return "Token_Butterfly.jpg"
	case BeeJunebugToken:
		return "Token_BeeJunebug.jpg"
	case BeeButterflyToken:
		return "Token_BeeButterfly.jpg"
	case JunebugButterFlyToken:
		return "Token_JunebugButterFly.jpg"
	case BeeJunebugButterFlyToken:
		return "Token_BeeJunebugButterFly.jpg"
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
	tokens map[uuid.UUID]*PollinatorToken
}

func NewTokenBag(shuffle bool) *TokenBag {
	tb := &TokenBag{map[uuid.UUID]*PollinatorToken{}}
	tks := createPollinatorTokens(tb)
	for _, tk := range tks {
		_, duplicateID := tb.tokens[tk.ID]
		if duplicateID {
			panic("duplicate token ID")
		}
		tb.tokens[tk.ID] = tk
	}

	return tb
}

func (tb *TokenBag) OutOfTokens() bool {
	return len(tb.tokens) == 0
}

func (tb *TokenBag) TakeNextToken() *PollinatorToken {
	tks := tb.GetTokens(1)
	if len(tks) == 0 {
		return nil
	}
	return tks[0]
}

func (tb *TokenBag) GetTokens(n int) []*PollinatorToken {
	if n == 0 {
		return nil
	}

	if len(tb.tokens) == 0 {
		return nil
	}

	if len(tb.tokens) < n {
		n = len(tb.tokens)
	}
	tks := make([]*PollinatorToken, 0, n)
	i := 0
	for _, tk := range tb.tokens {
		if tk.IsUsed() {
			panic("token is used")
		}
		tks = append(tks, tk)
		i++
		if i == n {
			break
		}
	}
	return tks
}

func (tb *TokenBag) GetToken(tokenID uuid.UUID) *PollinatorToken {
	for _, t := range tb.tokens {
		if t.ID == tokenID {
			return t
		}
	}
	return nil
}

func (tb *TokenBag) ConsumeToken(tokenID uuid.UUID) {
	logger.AtLevelf(logger.LBoard|logger.LPosition|logger.LToken, "%s trying to consume", tokenID)
	_, ok := tb.tokens[tokenID]
	if !ok {
		panic("token already consumed")
	}
	delete(tb.tokens, tokenID)
}

func (tb *TokenBag) HasToken(tokenID uuid.UUID) bool {
	_, has := tb.tokens[tokenID]
	return has
}

func createPollinatorTokens(b *TokenBag) []*PollinatorToken {
	singleTokenCreators := []func(b *TokenBag) *PollinatorToken{NewBeeToken, NewJunebugToken, NewButterflyToken}
	tokens := []*PollinatorToken{}
	for _, creator := range singleTokenCreators {
		tokens = append(tokens, NewTokenGroup(creator, b, 5)...)
	}
	doubleTokenCreators := []func(b *TokenBag) *PollinatorToken{NewBeeJunebugToken, NewBeeButterflyToken, NewJunebugButterFlyToken}
	for _, creator := range doubleTokenCreators {
		tokens = append(tokens, NewTokenGroup(creator, b, 9)...)
	}
	tokens = append(tokens, NewTokenGroup(NewBeeJunebugButterFlyToken, b, 2)...)
	return tokens
}

func NewTokenGroup(tokenCreator func(b *TokenBag) *PollinatorToken, b *TokenBag, n int) []*PollinatorToken {
	tokens := make([]*PollinatorToken, n)
	for i := 0; i < n; i++ {
		tokens[i] = tokenCreator(b)
	}
	return tokens
}

func NewToken(t TokenType, b *TokenBag) *PollinatorToken {
	return &PollinatorToken{
		ID:   uuid.Must(uuid.NewV4()),
		Type: t,
		b:    b,
	}
}

func NewBeeToken(b *TokenBag) *PollinatorToken {
	return NewToken(BeeToken, b)
}

func NewJunebugToken(b *TokenBag) *PollinatorToken {
	return NewToken(JunebugToken, b)
}

func NewButterflyToken(b *TokenBag) *PollinatorToken {
	return NewToken(ButterflyToken, b)
}

func NewBeeJunebugToken(b *TokenBag) *PollinatorToken {
	return NewToken(BeeJunebugToken, b)
}

func NewBeeButterflyToken(b *TokenBag) *PollinatorToken {
	return NewToken(BeeButterflyToken, b)
}

func NewJunebugButterFlyToken(b *TokenBag) *PollinatorToken {
	return NewToken(JunebugButterFlyToken, b)
}

func NewBeeJunebugButterFlyToken(b *TokenBag) *PollinatorToken {
	return NewToken(BeeJunebugButterFlyToken, b)
}
