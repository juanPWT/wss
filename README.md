# WSS

go package websocket based by "golang.org/x/net/websocket", This package focuses on broadcasting only and can only receive and send JSON, maybe it will be developed further to suit the use of websocket, the reason I created this package was only if I wanted to use websocket to only send and receive JSON, I just had to use this package hahahah.

### INSTALL

```
go get github.com/juanPWT/wss
```

### USAGE

```go
var port = "127.0.0.1:8000"

type Message struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

func handleBroadcast(w *ws.WebsocConn, msg Message) {
	err := ws.WebsocJSON.Send(w, msg)
	if err != nil {
		fmt.Println(err)
		w.Close()
	}
}

func main() {

	websocket := ws.NewWS(handleBroadcast)

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

```
