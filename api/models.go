package api

import (
	"golang_battleship/board"
	"golang_battleship/game"
	"time"
)

type ErrorResponseBody struct {
	Message string `json:"message"`
}

type Playername string

type VersionResponseBody struct {
	Version string `json:"version"`
}

type RegisterPlayerBody struct {
	Playername string `json:"name"`
	Password   string `json:"password"`
}

type RegisterPlayerResponseBody struct {
	ID string `json:"id"`
}

type CreateGameBody struct {
	BoardParameters board.BoardParameters `json:"board_parameters"`
	MaxPlayers      int                   `json:"max_players,omitempty"`
	Description     string                `json:"description,omitempty"`
}

type GetGameResponseBody struct {
	ID           string             `json:"id"`
	State        game.GameState     `json:"state"`
	CreationDate time.Time          `json:"creation_date"`
	Participants []game.Participant `json:"participants"`
	CreateGameBody
}

type CreateGameResponseBody struct {
	ID string `json:"id"`
}

type LeaveGameResponseBody struct {
	ID string `json:"id"`
}

type JoinGameResponseBody struct {
	ID string `json:"id"`
}

type ScoreboardEntry struct {
	Name   string `json:"name"`
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
}

type ScoreboardResponseBody []ScoreboardEntry
