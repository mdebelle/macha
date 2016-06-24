package main

import (
	// "error"
	"goji.io/pat"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"strconv"
)

var (
	ActiveClients = make(map[ClientConn]int)
	Message       = websocket.Message
)

type ClientConn struct {
	websocket *websocket.Conn
	clientIP  string
	userId    int
}

func serveWs(ws *websocket.Conn) {

	log.Println("hey")
	var err error
	var msg string

	defer ws.Close()

	r := ws.Request()
	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		log.Println("je close session")
		return
	}

	u := session.Values["UserInfo"].(UserData)

	sockCli := ClientConn{ws, r.RemoteAddr, u.Id}

	for {
		if err = Message.Receive(ws, &msg); err != nil {
			delete(ActiveClients, sockCli)
			log.Println("je close message")
			return
		}

		msg = session.Values["UserInfo"].(UserData).UserName + " : " + msg
		for cs, _ := range ActiveClients {
			// if cs.userId == u.ChatId || cs.userId == u.Id {
			if err = Message.Send(cs.websocket, msg); err != nil {
				log.Println("Could not send message to ", cs.clientIP, err.Error())
			}
			// }
		}
	}

}

func chat(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	// log.Println(session)
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	u := session.Values["UserInfo"].(UserData)
	id := pat.Param(ctx, "id")
	u.ChatId, _ = strconv.Atoi(id)
	session.Values["UserInfo"] = u
	session.Save(r, w)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	renderTemplate(w, "chat", &chatVew{
		Header: HeadData{
			Title:      "Chat",
			Stylesheet: []string{"chat.css"}},
		Host: "localhost" + LISTEN_PORT})

}
