package main

import (
	"database/sql"
	"html/template"
	"net/http"
)

type HeadData struct {
	Title      string
	Stylesheet []string
	Scripts    []string
}

type SimpleUser struct {
	Id          int64
	Bod			int
	UserName    string
}


type UserData struct {
	Id          int
	UserName    string
	FirstName   string
	LastName    string
	BirthDate   []uint8
	Email       string
	Bio         sql.NullString
	Sexe        sql.NullInt64
	Orientation sql.NullInt64
	Popularity  sql.NullInt64
	Interests   []Interest
	Matches		[]SimpleUser
}

type HomeView struct {
	Header HeadData
}

type usersView struct {
	Header HeadData
	Users  []UserData
}

type homeUserView struct {
	Header HeadData
	User   UserData
}

type inscriptionVew struct {
	Header HeadData
}

var templates = template.Must(template.ParseFiles(
	TEMPLATE_DIRECTORY+"home.html",
	TEMPLATE_DIRECTORY+"homeUser.html",
	TEMPLATE_DIRECTORY+"inscription.html",
	TEMPLATE_DIRECTORY+"users.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, v interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
