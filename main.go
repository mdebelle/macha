package main

import (
	//	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html"
	"log"
	"net/http"
	// "strconv"
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
	renderTemplate(w, "home", &RequiredData{Title: "Homepage", Stylesheet: "home.css"})
}

func notfoundHandler(w http.ResponseWriter, r *http.Request) {
	field := html.EscapeString(r.URL.Path[len("/notfound/"):])
	data := RequiredData{Title: "Nothing Found", Stylesheet: "notfound.css"}
	renderTemplate(w, "notfound", &NotFoundView{Data: data, Field: field})
}

func createNewUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("creation")
	fmt.Println(r)
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}
	var user = User{Username: r.FormValue("username"),
		Firstname: r.FormValue("firstname"),
		Lastname:  r.FormValue("lastname"),
		Email:     r.FormValue("email"),
		// Sexe:        strconv.Atoi(r.FormValue("sexe")),
		// Orientation: strconv.Atoi(r.FormValue("orientation")),
		Bio: r.FormValue("bio")}

	addnewUser(&user)

	http.Redirect(w, r, "/home/", http.StatusFound)
}

func newUserHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "newUser", &RequiredData{Title: "NewUser", Stylesheet: "home.css"})
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
	data := RequiredData{Title: user.Username, Stylesheet: "profile.css"}
	renderTemplate(w, "profile", &ProfileView{Data: data, ProfileUser: user})
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
	http.HandleFunc("/newUser/", newUserHandler)
	http.HandleFunc("/notfound/", notfoundHandler)
	http.HandleFunc("/createNewUser/", createNewUserHandler)

	// stylesheet := http.FileServer(http.Dir("./css/"))

	// http.Handle("/home/css/", http.StripPrefix("/home/css/", stylesheet))
	// http.Handle("/profile/css/", http.StripPrefix("/profile/css/", stylesheet))
	// http.Handle("/notfound/css/", http.StripPrefix("/notfound/css/", stylesheet))

	log.Fatal(http.ListenAndServe(":4242", nil))
}
