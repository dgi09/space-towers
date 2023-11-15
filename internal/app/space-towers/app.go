package spacetowers

import (
	"fmt"
	"space-towers/internal/gameplay"
	"space-towers/internal/lobby"
	"space-towers/internal/pkg"
	"space-towers/internal/session"
)

type App struct {
	network    pkg.Network
	sender     *pkg.Sender
	sessionMgr *session.Mgr
	gamesMgr   *gameplay.GamesMgr
	lobby      *lobby.Lobby
}

func NewApp() *App {
	return &App{}
}

func (a *App) Init(network pkg.Network) {
	a.network = network

	a.sender = pkg.NewSender(pkg.SenderOpts{
		Network: network,
	})

	a.sessionMgr = session.NewMgr()
	a.gamesMgr = gameplay.NewGamesMgr(a.network)

	a.lobby = lobby.New(lobby.Opts{
		Network:     network,
		RoomLaucher: NewDefRoomLauncher(a.gamesMgr),
	})
}

func (a *App) Start() {
	a.lobby.Start()
}

func (a *App) OnPlayerConnected(connId string) {
}

func (a *App) OnPlayerDisconnected(connId string) {
}

func (a *App) OnPlayerMessage(connId string, data []byte) {
	p, err := pkg.ParseInPkg(data)

	if err != nil {
		fmt.Println("Error parsing package")
	}

	switch p.Type {
	case pkg.InTypeDetails:
		a.handleDetails(connId, p.Data.(pkg.InDataDetails))
	case pkg.InTypeCreateRoom:
		a.handleCreateRoom(connId, p.Data.(pkg.InDataRoom))
	case pkg.InTypeJoinRoom:
		a.handleJoinRoom(connId, p.Data.(pkg.InDataRoom))
	// case pkg.InTypeEnqeueForGame:
	// a.handleEnqueueForGame(connId)
	case pkg.InTypeGameReady:
		a.handleGameReady(connId)

	case pkg.InTypeStartRound:
		a.handleStartRound(connId)

	case pkg.InTypeRoundFinished:
		a.handleRoundFinished(connId, p.Data.(pkg.InDataRoundFinished))
	}
}
