package main

import (
	"log"
	"net/http"
)

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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	renderTemplate(w, "chat", &chatVew{
		Header: HeadData{
			Title:      "Chat",
			Stylesheet: []string{"chat.css"}},
		Host: "localhost" + LISTEN_PORT + "/me/chat"})

}
