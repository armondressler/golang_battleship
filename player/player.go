package player

import (
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/google/uuid"
)

type PlayerMap map[string]*Player

type PlayerList []*Player

type Player struct {
	Name             string
	PasswordHash     string
	ID               uuid.UUID
	RegistrationDate time.Time
	Wins, Losses     int
}

func (l PlayerList) Len() int {
	return len(l)
}

func (l PlayerList) Less(i, j int) bool {
	return l[i].Wins < l[j].Wins || (l[i].Wins == l[j].Wins && l[i].Losses > l[j].Losses) || (l[i].Wins == l[j].Wins && l[i].Losses == l[j].Losses && sort.StringsAreSorted([]string{l[i].Name, l[j].Name}))
}

func (l PlayerList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

var AllPlayersMap = make(PlayerMap)

var AllPlayersList PlayerList

func (p *Player) ScoreWin() {
	p.Wins += 1
	sort.Sort(AllPlayersList)
}

func (p *Player) ScoreLoss() {
	p.Losses += 1
	sort.Sort(AllPlayersList)
}

func GetByName(playername string) (Player, error) {
	for _, existingPlayer := range AllPlayersMap {
		if playername == existingPlayer.Name {
			return *existingPlayer, nil
		}
	}
	return Player{}, fmt.Errorf("player name %s doesnt exist", playername)
}

func NewPlayer(name string, passwordHash string) (*Player, error) {
	r := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-_]{0,31}$`)
	if !r.MatchString(name) {
		return &Player{}, fmt.Errorf("player name %s doesn't meet requirements (starts with a letter, "+
			"only letters or numbers allowed, max size 32 characters)", name)
	}
	for _, existingPlayer := range AllPlayersMap {
		if name == existingPlayer.Name {
			return &Player{}, fmt.Errorf("player name %s is already taken", name)
		}
	}
	id := uuid.New()
	now := time.Now().UTC()
	p := Player{name, passwordHash, id, now, 0, 0}
	AllPlayersMap[name] = &p
	AllPlayersList = append(AllPlayersList, &p)
	sort.Sort(AllPlayersList)
	return &p, nil
}

func DeletePlayer(name string) (Player, error) {
	p, ok := AllPlayersMap[name]
	if !ok {
		return Player{}, fmt.Errorf("no player with name \"%s\" found", name)
	}
	delete(AllPlayersMap, name)
	for i, p := range AllPlayersList {
		if p.Name == name {
			AllPlayersList = append(AllPlayersList[:i], AllPlayersList[i+1:]...)
			sort.Sort(AllPlayersList)
		}
	}
	return *p, nil
}
