package ws

import (
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

// alias websocket type Conn
type WebsocConn = websocket.Conn

// alias var websocket.JSON
var WebsocJSON = websocket.JSON

// Websocket type
type Websocket[T any] struct {
	Conss       map[*websocket.Conn]bool
	OnBroadcast func(ws *websocket.Conn, msg T)
}

func NewWS[T any](onBroadcast func(ws *websocket.Conn, msg T)) *Websocket[T] {
	return &Websocket[T]{
		Conss:       make(map[*websocket.Conn]bool),
		OnBroadcast: onBroadcast,
	}
}

func (ws Websocket[T]) HandlerWS() websocket.Handler {
	return websocket.Handler(ws.HandleWS)
}

func (ws Websocket[T]) HandleWS(websoc *websocket.Conn) {
	// logging connection
	fmt.Printf("new connection: %s\n", websoc.RemoteAddr())

	if ws.Conss == nil {
		ws.Conss = make(map[*websocket.Conn]bool)
	}

	ws.Conss[websoc] = true

	// looping
	ws.ReadLoop(websoc)

	// delete websoc key in ws.Conss
	delete(ws.Conss, websoc)
}

func (ws Websocket[T]) ReadLoop(websoc *websocket.Conn) {
	for {
		var message T
		err := websocket.JSON.Receive(websoc, &message)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("disconected %s\n", websoc.RemoteAddr())
				return
			}
			fmt.Printf("WS Error :\n%s", err.Error())
			return
		}

		ws.Broadcast(message)
	}
}

func (ws Websocket[T]) Broadcast(msg T) {
	for websoc := range ws.Conss {
		go ws.OnBroadcast(websoc, msg)
	}
}
