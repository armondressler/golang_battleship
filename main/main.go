package main

import (
	"golang_battleship/api"
	"golang_battleship/client"
	"golang_battleship/cmd"
)

func main() {
	configFlags := cmd.ParseCmdFlags()
	if configFlags.Server {
		api.Serve(configFlags.Host, configFlags.Port, configFlags.JwtSigningKey)
	} else {
		client.Connect(configFlags.Host, configFlags.Port)
	}
}
