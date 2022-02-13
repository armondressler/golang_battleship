package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"golang_battleship/game"
	"golang_battleship/player"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	ws "github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const VERSION = "1.0"
const JWT_COOKIE_NAME = "battleship_jwt"

type RegisterPlayerBody struct {
	Playername     string `json:"name"`
	PasswordBCrypt string `json:"passwordbcrypt"`
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
	if _, err := player.NewPlayer(b.Playername, b.PasswordBCrypt); err != nil {
		resp, _ := json.Marshal(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}
	log.Info(fmt.Printf("registered new player %s", b.Playername))
	w.WriteHeader(http.StatusOK)
}

func JoinGame(w http.ResponseWriter, r *http.Request) {
	//gonna need some session handling ...
}

func LeaveGame(w http.ResponseWriter, r *http.Request) {

}

func CheckJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(JWT_COOKIE_NAME)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		t, err := jwt.ParseWithClaims(c.Value, jwt.StandardClaims{}, nil)
		if err != nil || !t.Valid {
			log.Warn("Malformed or invalid token, ", err)
			w.WriteHeader(http.StatusUnauthorized)
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		claims, ok := t.Claims.(jwt.MapClaims)
		if !ok {
			log.Warn("Malformed claims")
			w.WriteHeader(http.StatusUnauthorized)
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		log.Warn("SUBJECT IS: ", claims["Subject"])
		h.ServeHTTP(w, r)
	})
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
	games := make(map[string]map[string]string)
	for _, g := range game.AllGames {
		gamestats := map[string]string{
			"participants": strings.Join(g.ListParticipants(), ","),
			"state":        game.GameStateMap[g.State],
		}
		games[g.ID.String()] = gamestats
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

type JWTSingingKeyHandler struct {
	JWTSingingKey []byte
	Handler       func(w http.ResponseWriter, r *http.Request, jwtSigningKey []byte)
}

func (h JWTSingingKeyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Handler(w, r, h.JWTSingingKey)
}

func Serve(addr string, port int, jwtSigningKey []byte) {
	r := mux.NewRouter()
	needsAuth := mux.NewRouter()
	needsAuth.Use(CheckJWT)

	r.Path("/players").Methods("POST").HandlerFunc(RegisterPlayer)
	r.Path("/players").Methods("GET").HandlerFunc(Scoreboard)
	r.Path("/login").Methods("POST").Handler(JWTSingingKeyHandler{JWTSingingKey: jwtSigningKey, Handler: Login})
	needsAuth.Path("/games").Methods("GET").HandlerFunc(ListGames)
	log.Fatal(http.ListenAndServe(addr+":"+fmt.Sprint(port), r))
}
