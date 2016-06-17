package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type InscriptionForm struct {
	UserName     string
	FirstName    string
	LastName     string
	Email        string
	Error        bool
	ErrorMessage string
}

func inscription(w http.ResponseWriter, r *http.Request) {

	var iform InscriptionForm
	session, err := store.Get(r, "session")
	if err == nil {
		iform = session.Values["iForm"].(InscriptionForm)
	}

	fmt.Println(iform)

	renderTemplate(w, "inscription", &inscriptionVew{
		Header: HeadData{
			Title:      "Inscription",
			Stylesheet: []string{"inscription.css"}},
		IForm: iform})
}

func checkeAlreadyExistingUsers(UserName, Email string) (string, bool) {

	var msg = ""
	var countuser, countmail int
	smt, err := database.Prepare("SELECT " +
		"SUM(CASE WHEN username=? THEN 1 ELSE 0 END) countuser, " +
		"SUM(CASE WHEN email=? THEN 1 ELSE 0 END) countmail " +
		"FROM user")

	checkErr(err)
	defer smt.Close()
	err = smt.QueryRow(UserName, Email).Scan(&countuser, &countmail)
	checkErr(err)

	if countuser > 0 {
		msg += "User already Exist "
	}
	if countmail > 0 {
		msg += "Account already Created with this email "
	}

	return msg, countuser+countmail > 0
}

func postUsers(w http.ResponseWriter, r *http.Request) {

	var iForm = InscriptionForm{
		UserName:  r.FormValue("username"),
		FirstName: r.FormValue("firstname"),
		LastName:  r.FormValue("lastname"),
		Email:     r.FormValue("email")}

	if iForm.ErrorMessage, iForm.Error = checkeAlreadyExistingUsers(iForm.UserName, iForm.Email); iForm.Error == true {
		fmt.Println("User already exists", iForm)
		session, _ := store.Get(r, "session")
		session.Values["iForm"] = iForm
		session.Save(r, w)
		http.Redirect(w, r, "/inscription", http.StatusFound)
		return
	}

	p, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	checkErr(err)

	smt, err := database.Prepare("INSERT user SET username=?, firstname=?, lastname=?, email=?, password=?")
	checkErr(err)
	defer smt.Close()
	_, err = smt.Exec(iForm.UserName, iForm.FirstName, iForm.LastName, iForm.Email, p)
	checkErr(err)
	session, _ := store.Get(r, "session")
	session.Values["iForm"] = InscriptionForm{}
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)

}
