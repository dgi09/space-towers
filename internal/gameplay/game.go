package gameplay

import (
	"sort"
	"space-towers/internal/pkg"
	st "space-towers/internal/spacetowers"
	"time"
)

type gameState uint8

const (
	awaitingPlayersReady gameState = iota
	awaitingRoundStart
	roundInProgress
)

type GamePlayer struct {
	ConnectionId     string
	Name             string
	Ready            bool
	ReadyToPlayRound bool
	InRound          bool
	Score            uint
}

type GameOpts struct {
	Id        string
	Network   pkg.Network
	Players   []st.Player
	MaxRounds uint8
}

type Game struct {
	id            string
	maxRounds     uint8
	roundDuration uint16
	currentRound  uint8
	state         gameState
	playersList   []st.Player
	players       map[string]*GamePlayer

	//channels
	termChan      chan struct{}
	termComplChan chan struct{}
	playerMsgChan chan gameMsg

	//deps
	pkgSender *pkg.Sender
}

func NewGame(opts GameOpts) *Game {
	players := make(map[string]*GamePlayer)
	recipients := make([]string, len(opts.Players))
	for _, player := range opts.Players {
		players[player.ConnectionId] = &GamePlayer{
			ConnectionId:     player.ConnectionId,
			Name:             player.Name,
			Score:            0,
			Ready:            false,
			ReadyToPlayRound: false,
			InRound:          false,
		}

		recipients = append(recipients, player.ConnectionId)
	}

	return &Game{
		id: opts.Id,

		maxRounds:     opts.MaxRounds,
		roundDuration: 60,
		currentRound:  0,
		state:         awaitingPlayersReady,
		playersList:   opts.Players,
		players:       players,

		// channels
		termChan:      make(chan struct{}),
		termComplChan: make(chan struct{}),

		playerMsgChan: make(chan gameMsg),

		// deps
		pkgSender: pkg.NewSender(pkg.SenderOpts{
			Network:    opts.Network,
			Recipients: recipients,
		}),
	}
}

func (g *Game) GetPlayers() []st.Player {
	return g.playersList
}

func (g *Game) Start() error {
	go g.gameplay()

	return nil
}

func (g *Game) Stop() error {
	close(g.termChan)
	<-g.termComplChan

	return nil
}

func (g *Game) PlayerReady(connectionId string) {
	msg := gameMsg{
		Type: gameMsgTypePlayerReady,
		Data: gameMsgPlayer{
			ConnectionId: connectionId,
		},
	}

	g.pushMsg(msg)
}

func (g *Game) PlayerReadyToStartRound(connectionId string) {
	msg := gameMsg{
		Type: gameMsgTypePlayerReadyToStartRound,
		Data: gameMsgPlayer{
			ConnectionId: connectionId,
		},
	}

	g.pushMsg(msg)
}

func (g *Game) PlayerRoundFinished(connectionId string, combos []int) {
	msg := gameMsg{
		Type: gameMsgTypePlayerFinishedRound,
		Data: gameMsgPlayerFinishedRound{
			ConnectionId: connectionId,
			Combos:       combos,
		},
	}

	g.pushMsg(msg)
}

func (g *Game) RoundTimerFinished(round uint8) {
	msg := gameMsg{
		Type: gameMsgTypeRoundTimerFinished,
		Data: gameMsgTimerFinished{
			Round: round,
		},
	}

	g.pushMsg(msg)
}

func (g *Game) pushMsg(msg gameMsg) error {
	g.playerMsgChan <- msg

	return nil
}

func (g *Game) gameplay() {
	g.gpSendPlayersPrepareGame()

	for {
		select {

		case pmsg := <-g.playerMsgChan:
			g.gpHandleMsg(pmsg)

		case <-g.termChan:
			close(g.termComplChan)
			return
		}
	}
}

func (g *Game) gpGetScore() st.GameScore {
	score := make(st.GameScore, 0, len(g.playersList))

	for _, player := range g.playersList {
		score = append(score, st.GameScoreEntry{
			Player: player.Name,
			Score:  g.players[player.ConnectionId].Score,
		})
	}

	sort.SliceStable(score, func(i, j int) bool {
		return score[i].Score > score[j].Score
	})

	return score
}

func (g *Game) gpNextRound() {
	g.state = roundInProgress
	g.currentRound++

	for _, player := range g.players {
		player.InRound = true
	}

	deck := GenRandomRoundDeck(g.currentRound)

	go func(game *Game, round uint8) {
		timer := time.NewTimer(time.Duration(g.roundDuration) * time.Second)
		<-timer.C

		game.RoundTimerFinished(round)
	}(g, g.currentRound)

	g.pkgSender.ToAll().StartRound(g.currentRound, g.roundDuration, deck)
}

func (g *Game) gpAwaitNextRound() {
	g.state = awaitingRoundStart

	for _, player := range g.players {
		player.ReadyToPlayRound = false
		player.InRound = false
	}

	score := g.gpGetScore()

	g.pkgSender.ToAll().AwaitRoundStart(score)
}

func (g *Game) gpSendPlayersPrepareGame() {
	playersNames := make([]string, len(g.playersList))

	for _, player := range g.playersList {
		playersNames = append(playersNames, player.Name)
	}

	g.pkgSender.ToAll().PrepareGame(playersNames, g.maxRounds, g.roundDuration)
}

func (g *Game) gpHandleMsg(pmsg gameMsg) {
	switch pmsg.Type {
	case gameMsgTypePlayerReady:
		g.gpHandlePlayerReady(pmsg.Data.(gameMsgPlayer))

	case gameMsgTypePlayerReadyToStartRound:
		g.gpHandlePlayerReadyToStartRound(pmsg.Data.(gameMsgPlayer))

	case gameMsgTypePlayerFinishedRound:
		g.gpHandlePlayerFinishedRound(pmsg.Data.(gameMsgPlayerFinishedRound))

	case gameMsgTypeRoundTimerFinished:
		g.handleRoundTimerFinished(pmsg.Data.(gameMsgTimerFinished))
	}
}

func (g *Game) gpHandlePlayerReady(pmsg gameMsgPlayer) {
	player, ok := g.players[pmsg.ConnectionId]
	if !ok {
		return
	}

	player.Ready = true

	readyPlayers := 0
	for _, player := range g.players {
		if player.Ready {
			readyPlayers++
		}
	}

	if readyPlayers == len(g.players) {
		g.gpAwaitNextRound()
	}
}

func (g *Game) gpHandlePlayerReadyToStartRound(pmsg gameMsgPlayer) {
	player, ok := g.players[pmsg.ConnectionId]
	if !ok {
		return
	}

	player.ReadyToPlayRound = true

	readyPlayers := 0
	for _, player := range g.players {
		if player.ReadyToPlayRound {
			readyPlayers++
		}
	}

	if readyPlayers == len(g.players) {
		g.gpNextRound()
	}
}

func (g *Game) gpHandlePlayerFinishedRound(pmsg gameMsgPlayerFinishedRound) {
	if g.state != roundInProgress {
		return
	}

	player, ok := g.players[pmsg.ConnectionId]
	if !ok {
		return
	}

	player.InRound = false

	for _, combo := range pmsg.Combos {
		points := 1
		if combo > 1 && combo <= len(ComboPoints) {
			points = ComboPoints[combo-1]
		}
		player.Score += uint(points)
	}

	roundFinished := true
	for _, player := range g.players {
		if player.InRound {
			roundFinished = false
			break
		}
	}

	if roundFinished {
		if g.currentRound == g.maxRounds {
			g.pkgSender.ToAll().EndGame(g.gpGetScore())

		} else {
			g.gpAwaitNextRound()
		}
	}
}

func (g *Game) handleRoundTimerFinished(pmsg gameMsgTimerFinished) {
	if g.state != roundInProgress {
		return
	}

	if pmsg.Round != g.currentRound {
		return
	}

	for _, player := range g.players {
		if player.InRound {
			p := g.pkgSender.To(player.ConnectionId)
			p.ForceRoundFinish()
		}
	}
}
