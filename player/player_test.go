package player

import (
	"fmt"
	"testing"
)

func TestNewPlayer(t *testing.T) {
	goodPlayerNames := []string{"abcde", "Abcde", "ABCDE", "a", "a1bcde", "A12345_-", "A-_1234567890987654321_qwerHGF45"}
	badPlayerNames := []string{"4kjdf", "", "%çkjcd", "abcdä", "abcde1234567890987654321234567890", "sdlfkj 45"}
	duplicatePlayerNames := []string{"abcde", "abcde"}
	for _, p := range goodPlayerNames {
		fmt.Println("Validating ", p)
		if _, err := NewPlayer(p, ""); err != nil {
			t.Logf("Valid player name \"%s\" failed validation", p)
			t.Fail()
		}
	}
	for _, p := range badPlayerNames {
		if _, err := NewPlayer(p, ""); err == nil {
			t.Logf("Invalid player name \"%s\" validated successfully", p)
			t.Fail()
		}
	}
	for i, p := range duplicatePlayerNames {
		if _, err := NewPlayer(p, ""); i > 0 && err == nil {
			t.Logf("Duplicate player name \"%s\" validated successfully", p)
			t.Fail()
		}
	}
}
