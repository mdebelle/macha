package main

import (
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var errorWrong = errors.New("Wrong Password")

func connectedUser(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "session")
	checkErr(err)
	var v homeUserView

	if session.Values["connected"] == true {
		v.Header = HeadData{

			Title: "Bonjour " + session.Values["UserInfo"].(UserData).FirstName + " " + session.Values["UserInfo"].(UserData).LastName,

			Stylesheet: []string{"homeUser.css"},
			Scripts:    []string{"homeUser.js"}}
		v.User, _ = session.Values["UserInfo"].(UserData)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	renderTemplate(w, "homeUser", &v)
}

func checkLoginUser(username string, password []byte) (UserData, error) {

	var user UserData
	var spassword []byte

	err := database.QueryRow("SELECT id, Firstname, Lastname, BirthDate, Email, Bio, Sexe, Orientation, Popularity, password FROM user WHERE username=?", username).Scan(
		&user.Id, &user.FirstName, &user.LastName, &user.BirthDate, &user.Email, &user.Bio, &user.Sexe, &user.Orientation, &user.Popularity, &spassword)
	switch {
	case err == sql.ErrNoRows:
		return user, sql.ErrNoRows
	case err != nil:
		return user, err
	}
	if bcrypt.CompareHashAndPassword(spassword, password) != nil {
		return user, errorWrong
	}
	user.UserName = username
	return user, nil
}

func login(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	fmt.Println(session)
	if session.Values["connected"] == true {
		fmt.Println("dejaconnecte")
		http.Redirect(w, r, "/me", http.StatusFound)
		return
	}

	user, err := checkLoginUser(r.FormValue("username"), []byte(r.FormValue("password")))
	fmt.Println(err)
	switch {
	case err == sql.ErrNoRows:
		http.Redirect(w, r, "/", http.StatusFound)
		return
	case err == errorWrong:
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
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
