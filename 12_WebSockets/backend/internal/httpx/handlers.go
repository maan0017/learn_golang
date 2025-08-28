package httpx

import (
	"go/web-sockets/internal/chat"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

func WebSocketHandler(hub *chat.Hub, allowedOriginsCSV string) http.Handler {
	allowed := parseCSV(allowedOriginsCSV)

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return true // allow non-browser clients
			}
			for _, o := range allowed {
				if o == origin {
					return true
				}
			}
			return false
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "upgrade failed", http.StatusBadRequest)
			return
		}

		client := chat.NewClient(hub, conn)
		hub.Register <- client

		// One goroutine reads, one writes
		go client.WritePump()
		go client.ReadPump()
	})
}

func parseCSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
