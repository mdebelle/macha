package main

import (
	"log"
	"net/http"
)

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	conn := &Conn{send: make(chan []byte, 256), ws: ws}
	hub.register <- conn
	go conn.writePump()
	conn.readPump()
}

func chat(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	http.Error(w, "Not found", 404)
	// 	return
	// }
	// if r.Method != "GET" {
	// 	http.Error(w, "Method not allowed", 405)
	// 	return
	// }
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// homeTemplate.Execute(w, r.Host)

	renderTemplate(w, "chat", &chatVew{
		Header: HeadData{
			Title:      "Chat",
			Stylesheet: []string{"chat.css"}},
		Host: "localhost" + LISTEN_PORT + "/me/chat"})

}
