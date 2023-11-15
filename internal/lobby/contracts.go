package lobby

import st "space-towers/internal/spacetowers"

type RoomLauncher interface {
	Launch(roomName string, players []st.Player, maxRounds uint8)
}
