package cmd

import (
	"flag"
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

type cmdFlags struct {
	Host     string
	Port     int
	Loglevel int
	Server   bool
}

func validateLoglevel(loglevel int) error {
	if !(loglevel >= 0 && loglevel <= 3) {
		return fmt.Errorf("bad loglevel: %d", loglevel)
	}
	return nil
}

func validatePort(port int) error {
	if !(port > 0 && port < 65354) {
		return fmt.Errorf("bad port: %d", port)
	}
	return nil
}

func validateHost(host string) error {
	if ok := net.ParseIP(host); ok == nil {
		return fmt.Errorf("bad host: %s", host)
	}
	return nil
}

func setLogger(loglevel int) {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.Level(loglevel + 2)) //skip panic and fatal level, start at error
}

func ParseCmdFlags() cmdFlags {
	var host string
	var port int
	var loglevel int
	var server bool

	flag.StringVar(&host, "host", "0.0.0.0", "Server address (or interface for server mode)")
	flag.IntVar(&port, "port", 80, "Port to connect to (or to listen on for server mode)")
	flag.IntVar(&loglevel, "loglevel", 0, "Log verbosity (0 (error) - 3 (debug)")
	flag.BoolVar(&server, "server", false, "Run as server")
	flag.Parse()
	setLogger(loglevel)
	return cmdFlags{host, port, loglevel, server}
}
