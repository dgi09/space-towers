package spacetowers

import (
	"space-towers/internal/gameplay"
	st "space-towers/internal/spacetowers"
)

type DefRoomLauncher struct {
	gamesMgr *gameplay.GamesMgr
}

func NewDefRoomLauncher(gamesMgr *gameplay.GamesMgr) *DefRoomLauncher {
	return &DefRoomLauncher{
		gamesMgr: gamesMgr,
	}
}

func (l *DefRoomLauncher) Launch(roomName string, players []st.Player, maxRounds uint8) {
	l.gamesMgr.CreateGame(players, maxRounds)
}
