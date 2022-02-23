package api

import (
	"encoding/json"
	"fmt"

	"golang_battleship/game"
	"golang_battleship/player"
	"net/http"

	"github.com/gorilla/mux"
	ws "github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	VERSION               = "1.0"
	JWT_COOKIE_NAME       = "battleship_jwt"
	PASSWORD_REHASH_COUNT = 10
)

type Playername string

type RegisterPlayerBody struct {
	Playername string `json:"name"`
	Password   string `json:"password"`
}

type CreateGameBody struct {
	BoardsizeX  int    `json:"boardsizeX,omitempty"`
	BoardsizeY  int    `json:"boardsizeY,omitempty"`
	Maxships    int    `json:"maxships,omitempty"`
	Maxplayers  int    `json:"maxplayers,omitempty"`
	Description string `json:"description,omitempty"`
}

type ScoreboardEntry struct {
	Name   string `json:"name"`
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
}

type ScoreboardBody []ScoreboardEntry

func Version(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	retval, _ := json.Marshal(VERSION)
	w.Write(retval)
}

func RegisterPlayer(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var b RegisterPlayerBody
	if err := decoder.Decode(&b); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := player.NewPlayer(b.Playername, b.Password); err != nil {
		resp, _ := json.Marshal(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	log.Info(fmt.Printf("registered new player %s", b.Playername))
	w.WriteHeader(http.StatusOK)
}

func JoinGame(w http.ResponseWriter, r *http.Request) {
	p, err := getPlayerFromContext(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rvars := mux.Vars(r)
	gameID, ok := rvars["id"]
	if !ok {
		w.Write([]byte("no game id provided in request path"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	g, err := game.GetByUUID(gameID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("no game found for id %s", gameID)))
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = g.AddParticipant(p)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func LeaveGame(w http.ResponseWriter, r *http.Request) {

}

func CreateGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var c CreateGameBody
	if err := decoder.Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	g, err := game.NewGame(c.BoardsizeX, c.BoardsizeY, c.Maxships, c.Description, c.Maxplayers)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	body := make(map[string]string)
	body["id"] = g.ID.String()
	retval, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	w.Write(retval)
}

func DeleteGame(w http.ResponseWriter, r *http.Request) {

}

func ListGames(w http.ResponseWriter, r *http.Request) {
	games := make(map[string]game.Game)
	for _, g := range game.AllGames {
		games[g.ID.String()] = g
	}
	if retval, err := json.Marshal(games); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(retval)
	}
}

func Scoreboard(w http.ResponseWriter, r *http.Request) {
	scoreboard := ScoreboardBody{}
	for _, p := range player.AllPlayersList {
		entry := ScoreboardEntry{Name: p.Name, Wins: p.Wins, Losses: p.Losses}
		scoreboard = append(scoreboard, entry)
	}
	log.Info(scoreboard)
	retval, err := json.Marshal(scoreboard)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(retval)
}

func echo(w http.ResponseWriter, r *http.Request) {
	upgrader := ws.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

type JWTMiddleware struct {
	jwtSigningKey []byte
	loginHandler  func(w http.ResponseWriter, r *http.Request, jwtSigningKey []byte)
}

func (jwtm JWTMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jwtm.loginHandler(w, r, jwtm.jwtSigningKey)
}

func Serve(addr string, port int, jwtSigningKey []byte) {
	jwtm := JWTMiddleware{jwtSigningKey: jwtSigningKey, loginHandler: Login}
	defaultRouter := mux.NewRouter()
	needsAuthRouter := defaultRouter.NewRoute().Subrouter()
	needsAuthRouter.Use(jwtm.CheckJWT)
	pw, _ := hashPassword("armon", PASSWORD_REHASH_COUNT)
	player.NewPlayer("armon", pw)
	defaultRouter.Path("/login").Methods("POST").Handler(jwtm)

	needsAuthRouter.Path("/players").Methods("POST").HandlerFunc(RegisterPlayer)
	needsAuthRouter.Path("/players").Methods("GET").HandlerFunc(Scoreboard)
	needsAuthRouter.Path("/games").Methods("GET").HandlerFunc(ListGames)
	needsAuthRouter.Path(fmt.Sprintf("/games/{id:%s}/join", game.ValidGameIDRegex)).HandlerFunc(JoinGame)
	log.Fatal(http.ListenAndServe(addr+":"+fmt.Sprint(port), defaultRouter))
}
