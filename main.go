package main

import (
	"html"
	"log"
	"net/http"

	"github.com/dleung/gotail"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/v1/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)

		// Listens to messages from the client and closes the connection when necessary
		go func(conn *websocket.Conn) {
			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					conn.Close()
				}
			}
		}(conn)

		// Sends data to the client
		go func(conn *websocket.Conn) {
			// ch := time.Tick(5 * time.Second)
			// for range ch {
			tail, err := gotail.NewTail("test.txt", gotail.Config{Timeout: 10})
			check(err)
			for line := range tail.Lines {
				conn.WriteJSON(Log{
					LogLine: html.EscapeString(line),
				})
			}
			// }
		}(conn)
	})

	http.ListenAndServe(":3000", nil)
}

// Log is the struct for the logfile you want to tail
// it doesnt need to be complicated, a single property per line is fine
type Log struct {
	LogLine string `json:"logline"`
}

// Just a simple error checking wrapper that logs errors to console
func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
