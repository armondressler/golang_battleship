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
	Name             string    `json:"name"`
	PasswordHash     string    `json:"-"`
	ID               uuid.UUID `json:"id"`
	RegistrationDate time.Time `json:"-"`
	Wins             int       `json:"wins"`
	Losses           int       `json:"losses"`
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

func (l PlayerList) BestOf(ranking int) ([]*Player, error) {
	//[p1,p2,p3,p4,p5]
	if ranking < 0 {
		ranking *= -1
		retval := make([]*Player, ranking)
		if ranking >= len(l) {
			return nil, fmt.Errorf("cannot get best of %v on current amount of registered players (%v)", ranking, len(l))
		}
		copy(retval, AllPlayersList[:ranking])
		return retval, nil
	}
	retval := make([]*Player, ranking)
	k := len(AllPlayersList)
	for i := 0; i < ranking; i++ {
		retval[i] = AllPlayersList[k-i-1]
	}
	return retval, nil
}

var AllPlayersMap = make(PlayerMap)

var AllPlayersList PlayerList

func (p Player) String() string {
	return p.Name
}

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
