package main

import (
	"fmt"
	"net"
	"net/http"
)

func connectedUser(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	var v homeUserView
	if session.Values["connected"] == true {
		v.Header = HeadData{
			Title:      "Bonjour " + session.Values["UserInfo"].(UserData).FirstName + " " + session.Values["UserInfo"].(UserData).LastName,
			Stylesheet: []string{"homeUser.css"},
			Scripts:    []string{"homeUser.js"}}
		v.User, _ = session.Values["UserInfo"].(UserData)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	renderTemplate(w, "homeUser", &v)
}

func login(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	fmt.Println(session)
	if session.Values["connected"] == true {
		http.Redirect(w, r, "/me", http.StatusFound)
		return
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	userIP := net.ParseIP(ip)
	fmt.Println(userIP)

	user, err := checkLoginUser(r.FormValue("username"), []byte(r.FormValue("password")))
	checkErr(err)
	session.Values["connected"] = true
	session.Values["UserInfo"] = user
	session.Save(r, w)
	http.Redirect(w, r, "/me", http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Options.MaxAge = 0
	session.Values["connected"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
