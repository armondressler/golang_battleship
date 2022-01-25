package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var apiHandlerMap = map[string]http.HandlerFunc{
	"/echo":  echo,
	"/hello": hello,
}

func hello(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
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
	upgrader := websocket.Upgrader{}
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

func Serve(addr string, port int, loglevel int) {
	for urlPath, handlerFunction := range apiHandlerMap {
		log.Info(fmt.Sprintf("Launching handler for %s", urlPath))
		http.HandleFunc(urlPath, handlerFunction)
	}
	log.Fatal(http.ListenAndServe(addr+":"+fmt.Sprint(port), nil))
}
