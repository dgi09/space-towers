package server

import "space-towers/internal/pkg"

type Binding interface {
	Init(network pkg.Network)

	Start()

	OnPlayerConnected(connId string)

	OnPlayerDisconnected(connId string)

	OnPlayerMessage(connId string, data []byte)
}
