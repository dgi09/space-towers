package spacetowers

import "space-towers/internal/server"

func Cmd() {
	app := NewApp()

	server := server.New(server.Opts{
		Port:    8080,
		Binding: app,
	})

	server.Start()
}
