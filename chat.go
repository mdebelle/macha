package main

import (
	"goji.io/pat"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	listenAddr = "localhost" + LISTEN_PORT // server address
)

var (
	Message       = websocket.Message        // codec for string, []byte
	ActiveClients = make(map[ClientConn]int) // map containing clients
)

// Client connection consists of the websocket and the client ip
type ClientConn struct {
	websocket      *websocket.Conn
	clientIP       string
	clientUsername string
	chatId         int
}

// WebSocket server to handle chat between clients
func SockServer(ws *websocket.Conn) {
	var err error
	var clientMessage string

	// cleanup on server side
	defer ws.Close()

	r := ws.Request()

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		return
	}

	var u = session.Values["UserInfo"].(UserData)

	log.Println("Client connected:", r.RemoteAddr, u.UserName)
	sockCli := ClientConn{ws, r.RemoteAddr, u.UserName, u.ChatId}
	ActiveClients[sockCli] = 0
	log.Println("Number of clients connected ...", len(ActiveClients))

	for {
		if err = Message.Receive(ws, &clientMessage); err != nil {
			log.Println("Websocket Disconnected waiting", err.Error())
			delete(ActiveClients, sockCli)
			log.Println("Number of clients still connected ...", len(ActiveClients))
			return
		}

		clientMessage = sockCli.clientUsername + " Said: " + clientMessage
		log.Println("clientMessage")
		for cs, _ := range ActiveClients {
			if u.ChatId == cs.chatId {
				if err = Message.Send(cs.websocket, clientMessage); err != nil {
					// we could not send the message to a peer
					log.Println("Could not send message to ", cs.clientUsername, err.Error())
				}
			}
		}
	}
}

// RootHandler renders the template for the root page
func RootHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log.Println("Roothandler")

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusNetworkAuthenticationRequired)
		return
	}
	var u = session.Values["UserInfo"].(UserData)

	log.Println(u)
	chatname := pat.Param(ctx, "chatname")
	chatname = strings.Replace(chatname, "me", strconv.FormatInt(int64(u.Id), 10), 1)

	log.Println(chatname)
	smt, err := database.Prepare("SELECT id FROM chatroom WHERE chatname_one=? OR chatname_two=?")
	checkErr(err, "RootHandler")
	var id int64
	smt.QueryRow(chatname, chatname).Scan(&id)
	if id == 0 {
		http.Redirect(w, r, "/me", http.StatusFound)
		return
	}
	u.ChatId = int(id)
	session.Values["UserInfo"] = u
	session.Save(r, w)
	log.Println("conneciton to chatroom", id)
	renderTemplate(w, "chat", &chatVew{
		Header: HeadData{
			Title:      "ChatRoom",
			Stylesheet: []string{"chat.css"},
			Scripts:    []string{"chat.js"}},
		Host: listenAddr})
}
