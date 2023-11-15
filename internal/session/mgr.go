package session

import (
	"errors"
	st "space-towers/internal/spacetowers"
	"time"
)

type Session struct {
	Begin  time.Time
	Player st.Player
}

type Mgr struct {
	sessions map[string]Session
}

func NewMgr() *Mgr {
	return &Mgr{
		sessions: make(map[string]Session),
	}
}

func (m *Mgr) StartSession(p st.Player) {
	s := Session{
		Begin:  time.Now(),
		Player: p,
	}

	m.sessions[p.ConnectionId] = s
}

func (m *Mgr) Exists(connId string) bool {
	_, ok := m.sessions[connId]
	return ok
}

func (m *Mgr) Get(connId string) (*Session, error) {
	s, ok := m.sessions[connId]

	if !ok {
		return nil, errors.New("session not found")
	}

	return &s, nil
}
