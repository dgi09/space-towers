package lobby

import (
	"space-towers/internal/pkg"
	st "space-towers/internal/spacetowers"
)

type room struct {
	maxPlayers uint8
	name       string
	players    map[string]st.Player
	maxRounds  uint8
}

type Opts struct {
	Network     pkg.Network
	RoomLaucher RoomLauncher
}

type Lobby struct {
	rooms map[string]room

	//channels
	termChan     chan struct{}
	waitTermChan chan struct{}
	plMsgChan    chan playerMsg

	//deps
	pkgSender *pkg.Sender
	launcher  RoomLauncher
}

func New(opts Opts) *Lobby {
	return &Lobby{
		rooms: make(map[string]room),

		termChan:     make(chan struct{}),
		waitTermChan: make(chan struct{}),
		plMsgChan:    make(chan playerMsg),

		pkgSender: pkg.NewSender(pkg.SenderOpts{
			Network: opts.Network,
		}),
		launcher: opts.RoomLaucher,
	}
}

func (l *Lobby) Start() {
	go l.work()
}

func (l *Lobby) Term() {
	close(l.termChan)

	<-l.waitTermChan
}

func (l *Lobby) CreateRoom(name string, maxPlayers uint8, creator st.Player) {
	msg := playerMsg{
		Type: msgTypeCreateRoom,
		Data: msgCreateRoom{
			Name:       name,
			MaxPlayers: maxPlayers,
			Creator:    creator,
		},
	}

	l.pushPlayerMsg(msg)
}

func (l *Lobby) JoinRoom(room string, player st.Player) {
	msg := playerMsg{
		Type: msgTypeJoinRoom,
		Data: msgJoinRoom{
			Room:   room,
			Player: player,
		},
	}

	l.pushPlayerMsg(msg)
}

func (l *Lobby) Queue(player st.Player) {
	msg := playerMsg{
		Type: msgTypeQueue,
		Data: msgQueue{
			Player: player,
		},
	}

	l.pushPlayerMsg(msg)
}

func (l *Lobby) pushPlayerMsg(msg playerMsg) {
	l.plMsgChan <- msg
}

func (l *Lobby) getFirstRoom() *room {
	for _, room := range l.rooms {
		return &room
	}

	return nil
}

func (l *Lobby) work() {
	for {
		select {
		case msg := <-l.plMsgChan:
			l.onPlayerMsg(msg)
		case <-l.termChan:
			close(l.waitTermChan)
			return
		}
	}
}

func (l *Lobby) onPlayerMsg(msg playerMsg) {
	switch msg.Type {

	case msgTypeCreateRoom:
		l.workCreateRoom(msg.Data.(msgCreateRoom))

	case msgTypeJoinRoom:
		l.workJoinRoom(msg.Data.(msgJoinRoom))

	case msgTypeQueue:
		l.workQueue(msg.Data.(msgQueue))
	}
}

func (l *Lobby) evalRoom(room room) {
	if room.maxPlayers == uint8(len(room.players)) {
		l.workStartRoom(room)
	}
}

func (l *Lobby) workStartRoom(room room) {
	delete(l.rooms, room.name)

	players := make([]st.Player, 0, len(room.players))
	for _, player := range room.players {
		players = append(players, player)
	}

	l.launcher.Launch(room.name, players, room.maxRounds)
}

func (l *Lobby) workCreateRoom(msg msgCreateRoom) {
	_, ok := l.rooms[msg.Name]

	if ok {
		return
	}

	room := room{
		maxPlayers: msg.MaxPlayers,
		maxRounds:  5,
		name:       msg.Name,
		players:    make(map[string]st.Player),
	}

	room.players[msg.Creator.ConnectionId] = msg.Creator
	l.rooms[room.name] = room

	l.evalRoom(room)
}

func (l *Lobby) workJoinRoom(msg msgJoinRoom) {
	room, ok := l.rooms[msg.Room]

	if !ok {
		return
	}

	room.players[msg.Player.ConnectionId] = msg.Player

	l.evalRoom(room)
}

func (l *Lobby) workQueue(msg msgQueue) {
	if len(l.rooms) == 0 {
		r := room{
			maxPlayers: 2,
			name:       "room",
			players:    make(map[string]st.Player),
		}

		r.players[msg.Player.ConnectionId] = msg.Player
		l.rooms[r.name] = r

		return
	}

	room := l.getFirstRoom()

	l.workJoinRoom(msgJoinRoom{
		Room: room.name,

		Player: msg.Player,
	})
}
