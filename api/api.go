package api

import (
	"encoding/json"
	"fmt"

	"golang_battleship/game"
	"golang_battleship/player"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	ws "github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	VERSION               = "1.0"
	JWT_COOKIE_NAME       = "battleship_jwt"
	PASSWORD_REHASH_COUNT = 10
)

func Version(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, http.StatusOK, VersionResponseBody{Version: VERSION})
}

func RegisterPlayer(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var b RegisterPlayerBody
	if err := decoder.Decode(&b); err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, "Failed to decode JSON body")
		return
	}
	p, err := player.NewPlayer(b.Playername, b.Password)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed to create new player, %s", err.Error()))
		return
	}
	log.Info(fmt.Printf("registered new player %s", b.Playername))
	JSONResponse(w, http.StatusOK, RegisterPlayerResponseBody{ID: p.ID.String()})
}

func JoinGame(w http.ResponseWriter, r *http.Request, p *player.Player, g *game.Game) {
	err := g.AddParticipant(*p)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed to join game with id %s, %s", g.ID, err))
		return
	}
	JSONResponse(w, http.StatusOK, JoinGameResponseBody{ID: g.ID.String()})
}

func LeaveGame(w http.ResponseWriter, r *http.Request, p *player.Player, g *game.Game) {
	err := g.RemoveParticipant(*p)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed to leave game with id %s, %s", g.ID, err))
		return
	}
	JSONResponse(w, http.StatusOK, LeaveGameResponseBody{ID: g.ID.String()})
}

func CreateGame(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	c := CreateGameBody{
		BoardsizeX:  game.DefaultBoardsizeX,
		BoardsizeY:  game.DefaultBoardsizeY,
		Maxships:    game.DefaultMaxships,
		Maxplayers:  game.DefaultMaxParticipants,
		Description: game.DefaultDescription,
	}
	if err := decoder.Decode(&c); err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, "Failed to parse request")
		return
	}
	g, err := game.NewGame(c.BoardsizeX, c.BoardsizeY, c.Maxships, c.Description, c.Maxplayers)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Failed to create new game, %s", err))
		return
	}
	JSONResponse(w, http.StatusOK, CreateGameResponseBody{ID: g.ID.String()})
}

func ListGames(w http.ResponseWriter, r *http.Request) {
	games := make(map[string]game.Game)
	for _, g := range game.AllGames {
		games[g.ID.String()] = g
	}
	JSONResponse(w, http.StatusOK, games)
}

func Scoreboard(w http.ResponseWriter, r *http.Request) {
	scoreboard := ScoreboardResponseBody{}
	for _, p := range player.AllPlayersList {
		entry := ScoreboardEntry{Name: p.Name, Wins: p.Wins, Losses: p.Losses}
		scoreboard = append(scoreboard, entry)
	}
	JSONResponse(w, http.StatusOK, scoreboard)
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

func DeleteGame(w http.ResponseWriter, r *http.Request, p *player.Player, g *game.Game) {
	err := game.DeleteByUUID(g.ID.String())
	if err != nil {
		JSONErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete game with id %s, %s", g.ID, err))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func JSONErrorResponse(w http.ResponseWriter, httpStatus int, message string) {
	e := ErrorResponseBody{
		Message: message,
	}
	body, _ := json.Marshal(e)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	if len(message) > 0 {
		w.Write(body)
	}
}

func JSONResponse(w http.ResponseWriter, httpStatus int, content interface{}) {
	body, err := json.Marshal(content)
	if err != nil {
		JSONErrorResponse(w, http.StatusInternalServerError, "")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(body)
}

func Serve(addr string, port int, jwtSigningKey []byte, csrfAuthKey []byte) {
	jwtm := JWTMiddleware{jwtSigningKey: jwtSigningKey, loginHandler: Login}
	csrfm := csrf.Protect(csrfAuthKey)
	defaultRouter := mux.NewRouter()

	needsAuthRouter := defaultRouter.NewRoute().Subrouter()
	needsAuthRouter.Use(jwtm.CheckJWT)
	needsAuthRouter.Use(csrfm)

	pw, _ := hashPassword("armon", PASSWORD_REHASH_COUNT)
	player.NewPlayer("armon", pw)

	defaultRouter.Path("/").Methods("GET").Handler(http.RedirectHandler("/static/html/login.html", http.StatusPermanentRedirect))
	defaultRouter.Path("/login").Methods("POST").Handler(jwtm)
	defaultRouter.Path("/version").Methods("GET").HandlerFunc(Version)
	defaultRouter.PathPrefix("/static").Methods("GET").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	needsAuthRouter.Path("/players").Methods("POST").HandlerFunc(RegisterPlayer)
	needsAuthRouter.Path("/players").Methods("GET").HandlerFunc(Scoreboard)
	needsAuthRouter.Path("/games").Methods("GET").HandlerFunc(ListGames)
	needsAuthRouter.Path("/games").Methods("POST").HandlerFunc(CreateGame)

	needsAuthRouter.Path(fmt.Sprintf("/games/{id:%s}", game.ValidGameIDRegex)).Methods("DELETE").Handler(
		gameValidatorHandler{
			gameValidator:   gameValidator,
			playerValidator: playerValidator,
			handler:         DeleteGame,
		})
	needsAuthRouter.Path(fmt.Sprintf("/games/{id:%s}/join", game.ValidGameIDRegex)).Methods("GET").Handler(
		gameValidatorHandler{
			gameValidator:   gameValidator,
			playerValidator: playerValidator,
			handler:         JoinGame,
		})
	needsAuthRouter.Path(fmt.Sprintf("/games/{id:%s}/leave", game.ValidGameIDRegex)).Methods("GET").Handler(
		gameValidatorHandler{
			gameValidator:   gameValidator,
			playerValidator: playerValidator,
			handler:         LeaveGame,
		})

	log.Fatal(http.ListenAndServe(addr+":"+fmt.Sprint(port), defaultRouter))
}
