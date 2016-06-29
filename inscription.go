package main

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
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
		if session.Values["iForm"] != nil {
			iform = session.Values["iForm"].(InscriptionForm)
		}
	}

	renderTemplate(w, "inscription", &inscriptionVew{
		Header: HeadData{
			Title:      "Inscription",
			Stylesheet: []string{"inscription.css"}},
		IForm: iform})
}

func checkeAlreadyExistingUsers(UserName, Email string) (string, bool) {

	var msg = ""
	var countuser, countmail sql.NullInt64
	smt, err := database.Prepare("SELECT " +
		"SUM(CASE WHEN username=? THEN 1 ELSE 0 END) countuser, " +
		"SUM(CASE WHEN email=? THEN 1 ELSE 0 END) countmail " +
		"FROM user")

	checkErr(err, "checkeAlreadyExistingUsers")
	defer smt.Close()
	err = smt.QueryRow(UserName, Email).Scan(&countuser, &countmail)
	checkErr(err, "checkeAlreadyExistingUsers")

	if countuser.Valid && countuser.Int64 > 0 {
		msg += "User already Exist "
	}
	if countmail.Valid && countmail.Int64 > 0 {
		msg += "Account already Created with this email "
	}

	return msg, countuser.Int64+countuser.Int64 > 0
}

func postUsers(w http.ResponseWriter, r *http.Request) {

	var iForm = InscriptionForm{
		UserName:  r.FormValue("username"),
		FirstName: r.FormValue("firstname"),
		LastName:  r.FormValue("lastname"),
		Email:     r.FormValue("email")}

	if iForm.ErrorMessage, iForm.Error = checkeAlreadyExistingUsers(iForm.UserName, iForm.Email); iForm.Error == true {
		session, _ := store.Get(r, "session")
		session.Values["iForm"] = iForm
		session.Save(r, w)
		http.Redirect(w, r, "/inscription", http.StatusFound)
		return
	}

	p, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	checkErr(err, "postUsers")
	tok, err := bcrypt.GenerateFromPassword([]byte(time.Now().String()), bcrypt.DefaultCost)

	smt, err := database.Prepare("INSERT user SET username=?, firstname=?, lastname=?, email=?, token=?")
	checkErr(err, "postUsers")
	defer smt.Close()
	res, err := smt.Exec(iForm.UserName, iForm.FirstName, iForm.LastName, iForm.Email, tok)
	checkErr(err, "postUsers")
	id, err := res.LastInsertId()

	smt, err = database.Prepare("INSERT pw SET userid=?, password=?")
	checkErr(err, "postUsers")
	_, err = smt.Exec(id, p)
	checkErr(err, "postUsers")

	session, _ := store.Get(r, "session")
	session.Values["iForm"] = InscriptionForm{}
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)

}
