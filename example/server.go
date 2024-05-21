package main

import (
	"fmt"
	"net/http"

	"github.com/juanPWT/wss/ws"
)

var port = "127.0.0.1:8000"

type Message struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

func handleBroadcast(w *ws.WebsocConn, msg interface{}) {
	err := ws.WebsocJSON.Send(w, msg)
	if err != nil {
		fmt.Printf("WS Error: %s\n", err)
		w.Close()
	}
}

func main() {

	websocket := ws.NewWS(Message{}, handleBroadcast)

	// mux
	mux := http.NewServeMux()
	mux.Handle("/ws", websocket.HandlerWS())

	handler := MiddlewareAllowCors(mux)

	fmt.Println("server starting in ", port)
	err := http.ListenAndServe(port, handler)
	if err != nil {
		fmt.Println("failed start server")
	}
}

func MiddlewareAllowCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
