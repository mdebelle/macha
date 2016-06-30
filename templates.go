package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

type HeadData struct {
	Title      string
	Stylesheet []string
	Scripts    []string
}

func (h HeadData) String() string {
	return fmt.Sprintf("Title: %s,\nStylsheet: %v\nScripts: %v\n",
		h.Title, h.Stylesheet, h.Scripts)
}

type SimpleUser struct {
	Id       int64
	Bod      int
	UserName string
	Bio      string
}

func (s SimpleUser) String() string {
	return fmt.Sprintf("Id: %d %s\nAge: %d\nBio: %s\n",
		s.Id, s.UserName, s.Bod, s.Bio)
}

type Notifications struct {
	Id     int64
	Date   []uint8
	Read   bool
	UserId int64
	Msg    string
}

func (n Notifications) String() string {
	return fmt.Sprintf("Id: %d \tDate: %s \tState: %t\nfor: %d\nContent: %s",
		n.Id, string(n.Date), n.Read, n.UserId, n.Msg)
}

type UserData struct {
	Id            int
	UserName      string
	FirstName     string
	LastName      string
	BirthDate     []uint8
	Email         string
	Bio           sql.NullString
	Sexe          sql.NullInt64
	Orientation   sql.NullInt64
	Popularity    sql.NullInt64
	Interests     []Interest
	Matches       []SimpleUser
	Notifs        []Notifications
	LastConnexion string
	ChatId        int
}

func (u UserData) String() string {
	return fmt.Sprintf("Id: %d \tUsername: %d\nFirstname: %d \tLastname: %s\nDate de Naissance: %s\nEmail: %s\nBio: %s\n"+
		"Sexe: %d \tOrientation: %d\n", u.Id, u.UserName, u.FirstName, u.LastName, string(u.BirthDate), u.Email, u.Bio.String,
		u.Sexe.Int64, u.Orientation.Int64)
}

type HomeView struct {
	Header HeadData
}

type ErrorView struct {
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

type publicProfileView struct {
	Header     HeadData
	Connection bool
	User       UserData
	Profile    SimpleUser
}

type inscriptionVew struct {
	Header HeadData
	IForm  InscriptionForm
}

type chatVew struct {
	Header HeadData
	Host   string
}

var templates = template.Must(template.ParseFiles(
	TEMPLATE_DIRECTORY+"home.html",
	TEMPLATE_DIRECTORY+"homeUser.html",
	TEMPLATE_DIRECTORY+"publicProfile.html",
	TEMPLATE_DIRECTORY+"inscription.html",
	TEMPLATE_DIRECTORY+"chat.html",
	TEMPLATE_DIRECTORY+"users.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, v interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
