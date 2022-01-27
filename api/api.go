package api

import (
	"encoding/json"
	"fmt"

	"golang_battleship/player"
	"net/http"

	ws "github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const VERSION = "1.0"

var apiHandlerMap = map[string]http.HandlerFunc{
	"/echo":     echo,
	"/hello":    hello,
	"/home":     home,
	"/register": register,
	"/version":  version,
}

type RegisterPlayer struct {
	Name string `json:"name"`
}

func version(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	retval, _ := json.Marshal(VERSION)
	w.Write(retval)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var p RegisterPlayer
	err := decoder.Decode(&p)
	player.NewPlayer(p.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var playerlist []string
	for _, p := range player.Players {
		playerlist = append(playerlist, p.String())
	}
	retval, err := json.Marshal(playerlist)
	log.Info("PLAYERLIST: ", string(retval))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(retval)
}

func hello(w http.ResponseWriter, r *http.Request) {
	upgrader := ws.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()

	err = c.WriteMessage(0, []byte("hello"))
	if err != nil {
		log.Println("write:", err)
	}
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

func Serve(addr string, port int) {
	for urlPath, handlerFunction := range apiHandlerMap {
		log.Info(fmt.Sprintf("Launching handler for %s", urlPath))
		http.HandleFunc(urlPath, handlerFunction)
	}
	log.Fatal(http.ListenAndServe(addr+":"+fmt.Sprint(port), nil))
}
