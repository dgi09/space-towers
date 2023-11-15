package pkg

import (
	"bytes"
	"encoding/json"
	st "space-towers/internal/spacetowers"
)

const (
	outTypeReady uint8 = iota
	outTypePrepareGame
	outTypeAwaitRoundStart
	outTypeStartRound
	outTypeEndGame
	outTypeForceRoundFinish
)

type Protocol struct {
	n          Network
	recipients []string
}

func (p *Protocol) sendPkg(t uint8, data any) {
	var dataBytes []byte = nil

	if data != nil {
		dataRes, err := json.Marshal(data)
		if err == nil {
			dataBytes = dataRes
		}
	}

	totalLen := 1 + len(dataBytes)

	buf := bytes.Buffer{}
	buf.Grow(totalLen)
	buf.WriteByte(byte(t))
	buf.Write(dataBytes)

	stream := buf.Bytes()

	for _, recipient := range p.recipients {
		p.n.SendPkg(recipient, stream)
	}
}

func (p *Protocol) Ready() {
	p.sendPkg(outTypeReady, nil)
}

func (p *Protocol) PrepareGame(players []string, maxRounds uint8, roundDuration uint16) {
	type prepareGameMsg struct {
		Players       []string `json:"players"`
		MaxRounds     uint8    `json:"maxRounds"`
		RoundDuration uint16   `json:"roundDuration"`
	}

	p.sendPkg(outTypePrepareGame, prepareGameMsg{
		Players:       players,
		MaxRounds:     maxRounds,
		RoundDuration: roundDuration,
	})
}

func (p *Protocol) AwaitRoundStart(score st.GameScore) {
	type awaitRoundStartMsg struct {
		Score st.GameScore `json:"score"`
	}

	p.sendPkg(outTypeAwaitRoundStart, awaitRoundStartMsg{
		Score: score,
	})
}

func (p *Protocol) StartRound(round uint8, duration uint16, deck st.RoundDeck) {
	type startRoundMsg struct {
		Round    uint8        `json:"round"`
		Duration uint16       `json:"duration"`
		Deck     st.RoundDeck `json:"deck"`
	}

	p.sendPkg(outTypeStartRound, startRoundMsg{
		Round:    round,
		Duration: duration,
		Deck:     deck,
	})
}

func (p *Protocol) ForceRoundFinish() {
	p.sendPkg(outTypeForceRoundFinish, nil)
}

func (p *Protocol) EndGame(score st.GameScore) {
	type endGameMsg struct {
		Score st.GameScore `json:"score"`
	}

	p.sendPkg(outTypeEndGame, endGameMsg{
		Score: score,
	})
}
