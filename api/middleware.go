package api

import (
	"fmt"
	"golang_battleship/game"
	"golang_battleship/player"
	"net/http"

	"github.com/gorilla/mux"
)

func playerValidator(w http.ResponseWriter, r *http.Request) (*player.Player, error) {
	p, err := getPlayerFromContext(r)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func gameValidator(w http.ResponseWriter, r *http.Request) (*game.Game, error) {
	rvars := mux.Vars(r)
	gameID, ok := rvars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return nil, fmt.Errorf("no game id provided in request path")
	}
	g, err := game.GetByUUID(gameID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return nil, fmt.Errorf("no game found for id %s", gameID)
	}
	return &g, nil
}

type JWTMiddleware struct {
	jwtSigningKey []byte
	loginHandler  func(w http.ResponseWriter, r *http.Request, jwtSigningKey []byte)
}

func (jwtm JWTMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jwtm.loginHandler(w, r, jwtm.jwtSigningKey)
}

type gameValidatorHandler struct {
	gameValidator   func(w http.ResponseWriter, r *http.Request) (*game.Game, error)
	playerValidator func(w http.ResponseWriter, r *http.Request) (*player.Player, error)
	handler         func(w http.ResponseWriter, r *http.Request, p *player.Player, g *game.Game)
}

func (gv gameValidatorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p, err := gv.playerValidator(w, r)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Cannot determine player identity, %s", err.Error()))
	}
	g, err := gv.gameValidator(w, r)
	if err != nil {
		JSONErrorResponse(w, http.StatusBadRequest, err.Error())
	}
	gv.handler(w, r, p, g)
}
