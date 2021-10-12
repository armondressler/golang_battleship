package cmd

import (
	"flag"
	"fmt"
	"net"
)

type cmdFlags struct {
	host     string
	port     int
	loglevel int
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

func ParseCmdFlags() cmdFlags {
	var host string
	var port int
	var loglevel int

	flag.StringVar(&host, "host", "0.0.0.0", "Interface to listen on.")
	flag.IntVar(&port, "port", 80, "Port to listen on.")
	flag.IntVar(&loglevel, "loglevel", 0, "Log verbosity (0 (fatal) - 3 (debug)")
	flag.Parse()
	return cmdFlags{host, port, loglevel}
}
