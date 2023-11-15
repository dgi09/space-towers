package gameplay

import (
	"space-towers/internal/pkg"
	st "space-towers/internal/spacetowers"

	"github.com/google/uuid"
)

type GamesMgr struct {
	network     pkg.Network
	games       map[string]*Game
	playerGames map[string]string
}

func NewGamesMgr(network pkg.Network) *GamesMgr {
	return &GamesMgr{
		network:     network,
		games:       make(map[string]*Game),
		playerGames: make(map[string]string),
	}
}

func (g *GamesMgr) CreateGame(players []st.Player, maxRounds uint8) {
	id := uuid.New().String()

	for _, player := range players {
		g.playerGames[player.ConnectionId] = id
	}

	game := NewGame(GameOpts{
		Id:        id,
		Network:   g.network,
		Players:   players,
		MaxRounds: maxRounds,
	})

	g.games[id] = game

	game.Start()
}

func (g *GamesMgr) RemoveGame(id string) {
	game, ok := g.games[id]
	if !ok {
		return
	}

	for _, player := range game.GetPlayers() {
		delete(g.playerGames, player.ConnectionId)
	}

	delete(g.games, id)
}

func (g *GamesMgr) GetGameByPlayer(connectionId string) *Game {
	gameId, ok := g.playerGames[connectionId]
	if !ok {
		return nil
	}

	game, ok := g.games[gameId]
	if !ok {
		return nil
	}

	return game
}
