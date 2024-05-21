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

type Websocket struct {
	Conss       map[*websocket.Conn]bool
	TypeMessage interface{}
	OnBroadcast func(ws *websocket.Conn, msg interface{})
}

func NewWS(typeMessage interface{}, onBroadcast func(ws *websocket.Conn, msg interface{})) *Websocket {
	return &Websocket{
		Conss:       make(map[*websocket.Conn]bool),
		TypeMessage: typeMessage,
		OnBroadcast: onBroadcast,
	}
}

func (ws Websocket) HandlerWS() websocket.Handler {
	return websocket.Handler(ws.HandleWS)
}

func (ws Websocket) HandleWS(websoc *websocket.Conn) {
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

func (ws Websocket) ReadLoop(websoc *websocket.Conn) {
	for {
		message := &ws.TypeMessage
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

func (ws Websocket) Broadcast(msg interface{}) {
	for websoc := range ws.Conss {
		go ws.OnBroadcast(websoc, msg)
	}
}
