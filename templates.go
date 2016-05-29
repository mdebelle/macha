package main

import (
	"html/template"
	"net/http"
)

type DefaultView struct {
	Title      string
	Stylesheet string
}

type NotFoundView struct {
	Title      string
	Stylesheet string
	Field      string
}

var templates = template.Must(template.ParseFiles(
	tmplDir+"home.html",
	tmplDir+"profile.html",
	tmplDir+"notfound.html",
	tmplDir+"search.html",
	tmplDir+"tchat.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, v interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
