package api

import (
	"encoding/json"
	"golang_battleship/game"
	"golang_battleship/player"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/steinfletcher/apitest"
)

func TestRegisterPlayer(t *testing.T) {
	finish := make(chan struct{})
	r := mux.NewRouter()
	r.HandleFunc("/players", Scoreboard)
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
			Get("http://127.0.0.1:8080/players").
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
	g, _ := game.NewGame(12, 12, 6, "New Game", 2, "Rudolf", "Dagobert")

	go func() {
		if err := http.ListenAndServe("127.0.0.1:8080", r); err != nil {
			panic(err)
		}
	}()

	go func() {
		cli := &http.Client{
			Timeout: time.Second * 1,
		}

		r := apitest.New().
			EnableNetworking(cli).
			Get("http://127.0.0.1:8080/games").
			Expect(t).
			Status(http.StatusOK).
			End()

		body, _ := ioutil.ReadAll(r.Response.Body)
		var json_body map[string]json.RawMessage
		var json_body_game map[string]json.RawMessage
		err := json.Unmarshal(body, &json_body)
		if err != nil {
			t.Errorf("failed to find games in json response dict")
			t.Fail()
		}
		err = json.Unmarshal(json_body[g.ID.String()], &json_body_game)
		if err != nil {
			t.Errorf("failed to find game id %s in json response dict", g.ID.String())
			t.Fail()
		}
		var json_body_game_participants []string
		err = json.Unmarshal(json_body_game["participants"], &json_body_game_participants)
		if err != nil {
			t.Errorf("failed to find key participants in game with id %s of json response dict", g.ID.String())
			t.Fail()
		}
		for _, p := range json_body_game_participants {
			found := 0
			for _, pn := range g.Participants {
				if pn.String() == p {
					found = 1
					break
				}
			}
			if found == 0 {
				t.Errorf("got unknown participant %s", p)
				t.Fail()
			}
		}
		finish <- struct{}{}
	}()
	<-finish
}

func TestCreateGame(t *testing.T) {
	finish := make(chan struct{})
	r := mux.NewRouter()
	r.HandleFunc("/games", CreateGame)

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
			Post("http://127.0.0.1:8080/games").
			Body(`{}`).
			Expect(t).
			Status(http.StatusOK).
			End()

		finish <- struct{}{}
	}()
	<-finish
}
