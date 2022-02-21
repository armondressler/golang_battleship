package game

import (
	"fmt"
	"golang_battleship/board"
	"golang_battleship/player"
	"golang_battleship/ship"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type GameList []Game

type GameState int

type Game struct {
	Participants    []participant
	ID              uuid.UUID
	State           GameState
	Description     string
	CreationDate    time.Time
	MaxParticipants int
}

type participant struct {
	player player.Player
	ship   []ship.Ship
	board  board.Board
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

var AllGames GameList

var GameStateMap = map[GameState]string{
	0: "created",
	1: "awaiting players",
	2: "deploying ships",
	3: "running",
	4: "finished",
	5: "aborted",
}

func (p participant) String() string {
	return p.player.Name
}

func GetByUUID(uuid string) (Game, error) {
	for _, g := range AllGames {
		if g.ID.String() == uuid {
			return g, nil
		}
	}
	return Game{}, fmt.Errorf("no game found for uuid %s", uuid)
}

func (g Game) ListParticipants() []string {
	retval := []string{}
	for _, p := range g.Participants {
		retval = append(retval, p.String())
	}
	return retval
}

func (g Game) GetBoardParameters() (int, int, int) {
	if len(g.Participants) < 1 {
		return DefaultBoardsizeX, DefaultBoardsizeY, DefaultMaxships
	}
	b := g.Participants[0].board
	x, y := b.Dimensions()
	return x, y, b.MaxShips
}

func (g *Game) AddParticipant(player player.Player) error {
	if g.MaxParticipants <= len(g.Participants) {
		return fmt.Errorf("game with id %s has reached max participants (%d/%d)", g.ID, len(g.Participants), g.MaxParticipants)
	}
	boardx, boardy, maxShips := g.GetBoardParameters()
	g.Participants = append(g.Participants, participant{
		player,
		[]ship.Ship{},
		board.NewBoard(boardx, boardy, maxShips),
	})
	return nil
}

func (g *Game) RemoveParticipant(player player.Player) error {
	for i, p := range g.Participants {
		if p.player.Name == player.Name {
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
	var participants []participant
	for _, playername := range playernames {
		participants = append(participants, participant{
			*player.AllPlayersMap[playername],
			[]ship.Ship{},
			board.NewBoard(
				boardsizeX,
				boardsizeY,
				maxships)})
	}
	now := time.Now()
	g := Game{participants, gameuuid, StateCreated, description, now, maxparticipants}
	AllGames = append(AllGames, g)
	log.Info(fmt.Sprintf("Created new game %s with max. participants %d", gameuuid, maxparticipants))
	return g, nil
}
