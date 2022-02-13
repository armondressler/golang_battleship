package cmd

import (
	"crypto/rand"
	"flag"
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
)

type cmdFlags struct {
	Host          string
	Port          int
	Loglevel      int
	Server        bool
	JwtSigningKey []byte
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

func GetJwtSigningKeyFromEnv(envkey string) ([]byte, error) {
	e := []byte(os.Getenv(envkey))
	if len(e) == 0 {
		return []byte{}, fmt.Errorf("failed to get jwt signing key from env var %s, either empty or unset", envkey)
	}
	return e, nil
}

func GenerateJwtSigningKey(keysize int) ([]byte, error) {
	b := make([]byte, keysize)
	_, err := rand.Read(b)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func ParseCmdFlags() cmdFlags {
	var host string
	var port int
	var loglevel int
	var server bool
	var jwtSigningKey []byte

	flag.StringVar(&host, "host", "0.0.0.0", "Server address (or interface for server mode)")
	flag.IntVar(&port, "port", 80, "Port to connect to (or to listen on for server mode)")
	flag.IntVar(&loglevel, "loglevel", 0, "Log verbosity (0 (error) - 3 (debug)")
	flag.BoolVar(&server, "server", false, "Run as server")
	flag.Parse()
	setLogger(loglevel)
	jwtSigningKey, err := GetJwtSigningKeyFromEnv("BATTLESHIP_JWTSIGNINGKEY")
	if err != nil {
		log.Warn(err)
		jwtSigningKey, err = GenerateJwtSigningKey(32)
		if err != nil {
			panic(fmt.Errorf("failed to generate JWT signing key: %s", err))
		}
		log.Warn("generated JWT signing key: ", string(jwtSigningKey))
	}
	return cmdFlags{host, port, loglevel, server, jwtSigningKey}
}
