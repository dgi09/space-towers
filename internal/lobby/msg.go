package lobby

import st "space-towers/internal/spacetowers"

const (
	msgTypeCreateRoom uint8 = iota
	msgTypeJoinRoom
	msgTypeQueue
)

type msgCreateRoom struct {
	Creator    st.Player
	Name       string
	MaxPlayers uint8
}

type msgJoinRoom struct {
	Player st.Player
	Room   string
}

type msgQueue struct {
	Player st.Player
}

type playerMsg struct {
	Type uint8
	Data any
}
