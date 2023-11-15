package gameplay

import (
	"math/rand"
	st "space-towers/internal/spacetowers"
	"time"
)

var StartDeck []st.Card = []st.Card{
	{Suit: st.Clubs, Rank: st.Ace},
	{Suit: st.Clubs, Rank: st.Two},
	{Suit: st.Clubs, Rank: st.Three},
	{Suit: st.Clubs, Rank: st.Four},
	{Suit: st.Clubs, Rank: st.Five},
	{Suit: st.Clubs, Rank: st.Six},
	{Suit: st.Clubs, Rank: st.Seven},
	{Suit: st.Clubs, Rank: st.Eight},
	{Suit: st.Clubs, Rank: st.Nine},
	{Suit: st.Clubs, Rank: st.Ten},
	{Suit: st.Clubs, Rank: st.Jack},
	{Suit: st.Clubs, Rank: st.Queen},
	{Suit: st.Clubs, Rank: st.King},

	{Suit: st.Diamonds, Rank: st.Ace},
	{Suit: st.Diamonds, Rank: st.Two},
	{Suit: st.Diamonds, Rank: st.Three},
	{Suit: st.Diamonds, Rank: st.Four},
	{Suit: st.Diamonds, Rank: st.Five},
	{Suit: st.Diamonds, Rank: st.Six},
	{Suit: st.Diamonds, Rank: st.Seven},
	{Suit: st.Diamonds, Rank: st.Eight},
	{Suit: st.Diamonds, Rank: st.Nine},
	{Suit: st.Diamonds, Rank: st.Ten},
	{Suit: st.Diamonds, Rank: st.Jack},
	{Suit: st.Diamonds, Rank: st.Queen},
	{Suit: st.Diamonds, Rank: st.King},

	{Suit: st.Hearts, Rank: st.Ace},
	{Suit: st.Hearts, Rank: st.Two},
	{Suit: st.Hearts, Rank: st.Three},
	{Suit: st.Hearts, Rank: st.Four},
	{Suit: st.Hearts, Rank: st.Five},
	{Suit: st.Hearts, Rank: st.Six},
	{Suit: st.Hearts, Rank: st.Seven},
	{Suit: st.Hearts, Rank: st.Eight},
	{Suit: st.Hearts, Rank: st.Nine},
	{Suit: st.Hearts, Rank: st.Ten},
	{Suit: st.Hearts, Rank: st.Jack},
	{Suit: st.Hearts, Rank: st.Queen},
	{Suit: st.Hearts, Rank: st.King},

	{Suit: st.Spades, Rank: st.Ace},
	{Suit: st.Spades, Rank: st.Two},
	{Suit: st.Spades, Rank: st.Three},
	{Suit: st.Spades, Rank: st.Four},
	{Suit: st.Spades, Rank: st.Five},
	{Suit: st.Spades, Rank: st.Six},
	{Suit: st.Spades, Rank: st.Seven},
	{Suit: st.Spades, Rank: st.Eight},
	{Suit: st.Spades, Rank: st.Nine},
	{Suit: st.Spades, Rank: st.Ten},
	{Suit: st.Spades, Rank: st.Jack},
	{Suit: st.Spades, Rank: st.Queen},
	{Suit: st.Spades, Rank: st.King},
}

func GenRandomRoundDeck(reduce uint8) st.RoundDeck {
	res := st.RoundDeck{
		Row1: make(st.Deck, 0, 10),
		Row2: make(st.Deck, 0, 9),
		Row3: make(st.Deck, 0, 6),
		Row4: make(st.Deck, 0, 3),
	}

	randS := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(randS)

	currentCard := 0
	deck := make(st.Deck, len(StartDeck))
	copy(deck, StartDeck)

	deckLen := len(deck)

	for i := 0; i < 10; i++ {
		randCardIndex := rand.Intn(deckLen - currentCard)
		res.Row1 = append(res.Row1, deck[randCardIndex])
		deck[randCardIndex] = deck[deckLen-currentCard-1]
		currentCard++
	}

	for i := 0; i < 9; i++ {
		randCardIndex := rand.Intn(deckLen - currentCard)
		res.Row2 = append(res.Row2, deck[randCardIndex])
		deck[randCardIndex] = deck[deckLen-currentCard-1]
		currentCard++
	}

	for i := 0; i < 6; i++ {
		randCardIndex := rand.Intn(deckLen - currentCard)
		res.Row3 = append(res.Row3, deck[randCardIndex])
		deck[randCardIndex] = deck[deckLen-currentCard-1]
		currentCard++
	}

	for i := 0; i < 3; i++ {
		randCardIndex := rand.Intn(deckLen - currentCard)
		res.Row4 = append(res.Row4, deck[randCardIndex])
		deck[randCardIndex] = deck[deckLen-currentCard-1]
		currentCard++
	}

	extraLen := deckLen - currentCard - int(reduce)
	res.Extra = make(st.Deck, 0, extraLen)

	for i := 0; i < extraLen; i++ {
		randCardIndex := rand.Intn(deckLen - currentCard)
		res.Extra = append(res.Extra, deck[randCardIndex])
		deck[randCardIndex] = deck[deckLen-currentCard-1]
		currentCard++
	}

	return res
}
