package api

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
	BoardsizeX  int    `json:"boardsizeX,omitempty"`
	BoardsizeY  int    `json:"boardsizeY,omitempty"`
	Maxships    int    `json:"maxships,omitempty"`
	Maxplayers  int    `json:"maxplayers,omitempty"`
	Description string `json:"description,omitempty"`
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
