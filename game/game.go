package game

import (
	"encoding/json"
	"fmt"
	"golang_battleship/board"
	"golang_battleship/player"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type GameState int

type Game struct {
	Participants    []Participant         `json:"participants"`
	ID              uuid.UUID             `json:"id"`
	State           GameState             `json:"state"`
	Description     string                `json:"description"`
	CreationDate    time.Time             `json:"creation_date"`
	MaxParticipants int                   `json:"max_participants"`
	BoardParameters board.BoardParameters `json:"board_parameters"`
}

type Participant struct {
	Player player.Player `json:"player"`
	board  board.Board   `json:"-"`
}

const (
	StateOpen GameState = iota
	StateDeployingShips
	StateRunning
	StateFinished
	StateAborted
)

const (
	DefaultBoardsizeX      = 12
	DefaultBoardsizeY      = 12
	DefaultMaxships        = 5
	DefaultMaxParticipants = 2
	DefaultDescription     = "Join Me"
)

const ValidGameIDRegex = "[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"

var AllGames []*Game

var GameStateMap = map[string]GameState{
	"open":            0,
	"deploying ships": 1,
	"running":         2,
	"finished":        3,
	"aborted":         4,
}

func (p Participant) String() string {
	return p.Player.Name
}

func (p Participant) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func GetByUUID(uuid string) (*Game, error) {
	for _, g := range AllGames {
		if g.ID.String() == uuid {
			return g, nil
		}
	}
	return nil, fmt.Errorf("no game found for uuid %s", uuid)
}

func DeleteByUUID(uuid string) error {
	index := -1
	for i, g := range AllGames {
		if g.ID.String() == uuid {
			index = i
		}
	}
	if index == -1 {
		return fmt.Errorf("no game found for uuid %s", uuid)
	}
	AllGames = append(AllGames[:index], AllGames[index+1:]...)
	return nil
}

func (g Game) String() string {
	json_map := map[string]interface{}{
		"id":               g.ID.String(),
		"participants":     g.ListParticipants(),
		"max_participants": g.MaxParticipants,
		"state":            g.State,
		"description":      g.Description,
		"creation_date":    g.CreationDate,
		"board_parameters": g.BoardParameters,
	}
	s, _ := json.Marshal(json_map)
	return string(s[:])
}

func (g Game) ListParticipants() []string {
	retval := []string{}
	for _, p := range g.Participants {
		retval = append(retval, p.String())
	}
	return retval
}

func (g *Game) AddParticipant(player player.Player) error {
	if g.MaxParticipants <= len(g.Participants) {
		return fmt.Errorf("game with id %s has reached max participants (%d/%d)", g.ID, len(g.Participants), g.MaxParticipants)
	}
	for _, p := range g.ListParticipants() {
		if p == player.Name {
			return fmt.Errorf("player %s is already participant of game with id %s", p, g.ID)
		}
	}
	g.Participants = append(g.Participants, Participant{
		player,
		board.NewBoard(g.BoardParameters),
	})
	return nil
}

func (g *Game) RemoveParticipant(player player.Player) error {
	for i, p := range g.Participants {
		if p.Player.Name == player.Name {
			g.Participants = append(g.Participants[:i], g.Participants[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("no participant with name %s found for game with id %s", player.Name, g.ID)
}

func NewGame(boardsizeX, boardsizeY, maxships int, description string, maxparticipants int, playernames ...string) (*Game, error) {
	if boardsizeX < 10 || boardsizeY < 10 {
		return &Game{}, fmt.Errorf("boardsize (%d * %d) too small", boardsizeX, boardsizeY)
	}
	if maxships < 1 {
		return &Game{}, fmt.Errorf("maximum ship capacity (%d) too small", maxships)
	}
	if len(description) == 0 {
		description = DefaultDescription
	}
	if maxparticipants == 0 {
		maxparticipants = DefaultMaxParticipants
	}
	gameuuid := uuid.New()
	g := Game{[]Participant{}, gameuuid, StateOpen, description, time.Now(), maxparticipants, board.BoardParameters{
		SizeX: boardsizeX, SizeY: boardsizeY, MaxShips: maxships},
	}

	for _, playername := range playernames {
		p, _ := player.GetByName(playername)
		g.AddParticipant(p)
	}

	AllGames = append(AllGames, &g)
	log.Info(fmt.Sprintf("Created new game %s with max. participants %d", gameuuid, maxparticipants))
	return &g, nil
}
