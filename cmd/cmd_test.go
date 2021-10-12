package cmd

import (
	"testing"
)

func TestValidateLoglevel(t *testing.T) {
	goodLoglevels := []int{0, 3, 2}
	badLoglevels := []int{-1, 7, 44}
	for _, goodLoglevel := range goodLoglevels {
		if err := validateLoglevel(goodLoglevel); err != nil {
			t.Logf("Testing of valid loglevel %d failed", goodLoglevel)
			t.Fail()
		}
	}
	for _, badLoglevel := range badLoglevels {
		if err := validateLoglevel(badLoglevel); err == nil {
			t.Logf("Testing of invalid loglevel %d failed", badLoglevel)
			t.Fail()
		}
	}
}

func TestValidatePort(t *testing.T) {
	goodPorts := []int{1, 22, 65353, 89}
	badPorts := []int{-1, 0, 70000}
	for _, goodPort := range goodPorts {
		if err := validatePort(goodPort); err != nil {
			t.Logf("Testing of valid port %d failed", goodPort)
			t.Fail()
		}
	}
	for _, badPort := range badPorts {
		if err := validatePort(badPort); err == nil {
			t.Logf("Testing of invalid port %d failed", badPort)
			t.Fail()
		}
	}
}

func TestValidateHost(t *testing.T) {
	goodHosts := []string{"0.0.0.0", "127.0.0.1", "192.168.1.1", "123.123.123.123"}
	badHosts := []string{"256.1.1.1", "0.0.0", "a.1.1.1", "22. 3.3.3"}
	for _, goodHost := range goodHosts {
		if err := validateHost(goodHost); err != nil {
			t.Logf("Testing of valid host %s failed", goodHost)
			t.Fail()
		}
	}
	for _, badHost := range badHosts {
		if err := validateHost(badHost); err == nil {
			t.Logf("Testing of valid host %s failed", badHost)
			t.Fail()
		}
	}
}
