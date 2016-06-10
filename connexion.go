package main

import (
	// "golang.org/x/crypto/bcrypt"
	"net/http"
)

func connectedUser(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	var v homeUserView
	if session.Values["connected"] == true {
		v.Header = HeadData{
			Title:      "Bonjour " + session.Values["Firstname"].(string) + " " + session.Values["Lastname"].(string),
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
	if session.Values["connected"] == true {
		http.Redirect(w, r, "/me", http.StatusFound)
		return
	}
	user, err := checkLoginUser(r.FormValue("username"), []byte(r.FormValue("password")))
	checkErr(err)
	session.Values["connected"] = true
	session.Values["UserInfo"] = user
	session.Save(r, w)
	http.Redirect(w, r, "/me", http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Values["connected"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
