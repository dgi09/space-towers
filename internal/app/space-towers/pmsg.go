package spacetowers

import (
	"space-towers/internal/pkg"
	"space-towers/internal/spacetowers"
)

func (a *App) handleEnqueueForGame(connId string) {
	s, err := a.sessionMgr.Get(connId)

	if err != nil {
		return
	}

	a.lobby.Queue(s.Player)
}

func (a *App) handleDetails(connId string, data pkg.InDataDetails) {
	if a.sessionMgr.Exists(connId) {
		return
	}

	a.sessionMgr.StartSession(spacetowers.Player{
		ConnectionId: connId,
		Name:         data.Name,
	})

	p := a.sender.To(connId)
	p.Ready()
}

func (a *App) handleCreateRoom(connId string, data pkg.InDataRoom) {
	s, err := a.sessionMgr.Get(connId)

	if err != nil {
		return
	}

	a.lobby.CreateRoom(data.Room, 2, s.Player)
}

func (a *App) handleJoinRoom(connId string, data pkg.InDataRoom) {
	s, err := a.sessionMgr.Get(connId)

	if err != nil {
		return
	}

	a.lobby.JoinRoom(data.Room, s.Player)
}

func (a *App) handleGameReady(connId string) {
	g := a.gamesMgr.GetGameByPlayer(connId)

	if g == nil {
		return
	}

	g.PlayerReady(connId)
}

func (a *App) handleStartRound(connId string) {
	g := a.gamesMgr.GetGameByPlayer(connId)

	if g == nil {
		return
	}

	g.PlayerReadyToStartRound(connId)
}

func (a *App) handleRoundFinished(connId string, data pkg.InDataRoundFinished) {
	g := a.gamesMgr.GetGameByPlayer(connId)

	if g == nil {
		return
	}

	g.PlayerRoundFinished(connId, data.Combos)
}
