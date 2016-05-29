package main

import (
	//	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html"
	"log"
	"net/http"
)

const (
	tmplDir = "tmpl/"
	dataDir = "data/"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {

	if len(r.URL.Path[len("/home/"):]) > 0 {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "home", &DefaultView{Title: "Homepage", Stylesheet: "home.css"})
}

func notfoundHandler(w http.ResponseWriter, r *http.Request) {
	field := html.EscapeString(r.URL.Path[len("/notfound/"):])
	renderTemplate(w, "notfound", &NotFoundView{Title: "Nothing Found", Stylesheet: "notfound.css", Field: field})
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}
	fmt.Println(r.Method)
	username := html.EscapeString(r.URL.Path[len("/profile/"):])
	user, err := getUserByUsername(username)
	if err != nil {
		http.Redirect(w, r, "/notfound/"+username, http.StatusFound)
		return
	}
	renderTemplate(w, "profile", &ProfileView{Title: user.Username, Stylesheet: "profile.css", ProfileUser: user})
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("hello")
		panic(err)
	}
}

func main() {

	fmt.Println("localhost:4242")

	// static file

	// Connection a la base de donnee
	initdatabase()
	fmt.Println("Database Created")
	defer database.Close()

	// roads
	http.HandleFunc("/home/", homeHandler)
	http.HandleFunc("/profile/", profileHandler)
	http.HandleFunc("/notfound/", notfoundHandler)

	stylesheet := http.FileServer(http.Dir("./css/"))

	http.Handle("/home/css/", http.StripPrefix("/home/css/", stylesheet))
	http.Handle("/profile/css/", http.StripPrefix("/profile/css/", stylesheet))
	http.Handle("/notfound/css/", http.StripPrefix("/notfound/css/", stylesheet))

	log.Fatal(http.ListenAndServe(":4242", nil))
}
