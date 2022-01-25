package main

import (
	"golang_battleship/api"
	"golang_battleship/client"
	"golang_battleship/cmd"
)

const VERSION = "1.0"

func main() {
	configFlags := cmd.ParseCmdFlags()
	if configFlags.Server {
		api.Serve(configFlags.Host, configFlags.Port, configFlags.Loglevel)
	} else {
		client.Connect(configFlags.Host, configFlags.Port, configFlags.Loglevel)
	}
}
