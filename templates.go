package main

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles(
	tmplDir+"home.html",
	tmplDir+"profile.html",
	tmplDir+"search.html",
	tmplDir+"tchat.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
