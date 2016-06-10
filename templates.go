package main

import (
	"html/template"
	"net/http"
)

type HeadData struct {
	Title      string
	Stylesheet []string
	Scripts    []string
}

type UserData struct {
	Id          int
	UserName    string
	FirstName   string
	LastName    string
	Email       string
	Bio         string
	Sexe        int
	Orientation int
	Popularity  int
	Interests   []Interest
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
