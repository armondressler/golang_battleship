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
	Participants    []participant         `json:"participants"`
	ID              uuid.UUID             `json:"-"`
	State           GameState             `json:"state"`
	Description     string                `json:"description"`
	CreationDate    time.Time             `json:"creation_date"`
	MaxParticipants int                   `json:"max_participants"`
	BoardParameters board.BoardParameters `json:"board_parameters"`
}

type participant struct {
	Player player.Player `json:"player"`
	board  board.Board   `json:"-"`
}

const (
	StateCreated GameState = iota
	StateAwaitingPlayers
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

const ValidGameIDRegex = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"

var AllGames []Game

var GameStateMap = map[GameState]string{
	0: "created",
	1: "awaiting players",
	2: "deploying ships",
	3: "running",
	4: "finished",
	5: "aborted",
}

func (p participant) String() string {
	return p.Player.Name
}

func (p participant) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func GetByUUID(uuid string) (Game, error) {
	for _, g := range AllGames {
		if g.ID.String() == uuid {
			return g, nil
		}
	}
	return Game{}, fmt.Errorf("no game found for uuid %s", uuid)
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

func (g GameState) String() string {
	return GameStateMap[g]
}

func (g GameState) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}

func (g Game) String() string {
	json_map := map[string]interface{}{
		"participants":     g.ListParticipants(),
		"max_participants": g.MaxParticipants,
		"state":            GameStateMap[g.State],
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
	g.Participants = append(g.Participants, participant{
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

func NewGame(boardsizeX, boardsizeY, maxships int, description string, maxparticipants int, playernames ...string) (Game, error) {
	if boardsizeX < 10 || boardsizeY < 10 {
		return Game{}, fmt.Errorf("boardsize (%d * %d) too small", boardsizeX, boardsizeY)
	}
	if maxships < 1 {
		return Game{}, fmt.Errorf("maximum ship capacity (%d) too small", maxships)
	}
	if description == "" {
		description = DefaultDescription
	}
	if maxparticipants == 0 {
		maxparticipants = DefaultMaxParticipants
	}
	gameuuid := uuid.New()
	g := Game{[]participant{}, gameuuid, StateCreated, description, time.Now(), maxparticipants, board.BoardParameters{
		SizeX: boardsizeX, SizeY: boardsizeY, MaxShips: maxships},
	}

	for _, playername := range playernames {
		p, _ := player.GetByName(playername)
		g.AddParticipant(p)
	}

	AllGames = append(AllGames, g)
	log.Info(fmt.Sprintf("Created new game %s with max. participants %d", gameuuid, maxparticipants))
	return g, nil
}
