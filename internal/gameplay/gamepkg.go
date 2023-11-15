package gameplay

const (
	gameMsgTypePlayerReady uint8 = iota
	gameMsgTypePlayerReadyToStartRound
	gameMsgTypePlayerFinishedRound
	gameMsgTypeRoundTimerFinished
)

type gameMsgPlayer struct {
	ConnectionId string
}

type gameMsgPlayerFinishedRound struct {
	ConnectionId string
	Combos       []int
}

type gameMsgTimerFinished struct {
	Round uint8
}

type gameMsg struct {
	Type uint8
	Data any
}
