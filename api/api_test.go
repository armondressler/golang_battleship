package api

import (
	"golang_battleship/game"
	"golang_battleship/player"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/steinfletcher/apitest"
)

func TestRegister(t *testing.T) {
	finish := make(chan struct{})
	r := mux.NewRouter()
	r.HandleFunc("/home", Scoreboard)
	r.HandleFunc("/register", RegisterPlayer)

	go func() {
		if err := http.ListenAndServe("127.0.0.1:8080", r); err != nil {
			panic(err)
		}
	}()

	go func() {
		cli := &http.Client{
			Timeout: time.Second * 1,
		}

		apitest.New().
			EnableNetworking(cli).
			Post("http://127.0.0.1:8080/register").
			JSON(`{"name": "Rudolf"}`).
			Expect(t).
			Status(http.StatusOK).
			End()

		apitest.New().
			EnableNetworking(cli).
			Post("http://127.0.0.1:8080/register").
			JSON(`{"name": "Dagobert"}`).
			Expect(t).
			Status(http.StatusOK).
			End()

		apitest.New().
			EnableNetworking(cli).
			Post("http://127.0.0.1:8080/register").
			JSON(`{"name": "%Dagobert"}`).
			Expect(t).
			Status(http.StatusBadRequest).
			End()

		apitest.New().
			EnableNetworking(cli).
			Post("http://127.0.0.1:8080/register").
			JSON(`{"name": "wwwbbbbbbbbbbiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiicdssdffsdf"}`).
			Expect(t).
			Status(http.StatusBadRequest).
			End()

		apitest.New().
			EnableNetworking(cli).
			Get("http://127.0.0.1:8080/home").
			Expect(t).
			//Body(`["Rudolf","Dagobert"]`).
			Status(http.StatusOK).
			End()
		finish <- struct{}{}
	}()
	<-finish
}

func TestListGames(t *testing.T) {
	finish := make(chan struct{})
	r := mux.NewRouter()
	r.HandleFunc("/games", ListGames)
	player.NewPlayer("Rudolf", "")
	player.NewPlayer("Dagobert", "")
	game.NewGame(12, 12, 6, "New Game", 2, "Rudolf", "Dagobert")

	go func() {
		if err := http.ListenAndServe("127.0.0.1:8080", r); err != nil {
			panic(err)
		}
	}()

	go func() {
		cli := &http.Client{
			Timeout: time.Second * 1,
		}

		apitest.New().
			EnableNetworking(cli).
			Get("http://127.0.0.1:8080/games").
			Expect(t).
			Status(http.StatusOK).
			End()

		finish <- struct{}{}
	}()
	<-finish
}
