package spacetowers

type CardSuit uint8

const (
	Clubs CardSuit = iota
	Diamonds
	Hearts
	Spades
)

type CardValue uint8

const (
	Ace CardValue = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

type Card struct {
	Suit CardSuit
	Rank CardValue
}

type Deck []Card

type Player struct {
	ConnectionId string
	Name         string
}

type RoundDeck struct {
	Extra Deck
	Row1  Deck
	Row2  Deck
	Row3  Deck
	Row4  Deck
}

type GameScoreEntry struct {
	Player string
	Score  uint
}

type GameScore []GameScoreEntry
