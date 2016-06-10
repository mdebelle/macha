package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goji.io"
	"goji.io/pat"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// type MinimumInfo struct {
// 	Title      string
// 	Stylesheet string
// 	Firstname  string
// 	Lastname   string
// }

type ResponseStatus struct {
	Status string
}

func home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home", &HomeView{
		Header: HeadData{
			Title:      "Homepage",
			Stylesheet: []string{"home.css"}}})
}

func connectedUser(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	var v homeUserView
	if session.Values["connected"] == true {
		v.Header = HeadData{
			Title:      "Bonjour " + session.Values["Firstname"].(string) + " " + session.Values["Lastname"].(string),
			Stylesheet: []string{"homeUser.css"},
			Scripts:    []string{}}
		v.User = session.Values["UserInfo"].(UserData)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	renderTemplate(w, "homeUser", &v)
}

func inscription(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "inscription", &inscriptionVew{
		Header: HeadData{
			Title:      "Inscription",
			Stylesheet: []string{"inscription.css"}}})
}

func login(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	if session.Values["connected"] == true {
		http.Redirect(w, r, "/me", http.StatusFound)
		return
	}
	username := r.FormValue("username")
	password, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	checkErr(err)

	user, err := checkLoginUser(username, password)
	checkErr(err)
	session.Values["connected"] = true
	session.Values["UserInfo"] = user
	session.Save(r, w)
	http.Redirect(w, r, "/me", http.StatusFound)

}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Values["connected"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func writeJson(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		panic(err)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	fmt.Println("localhost:4242")

	// Connection a la base de donnee
	initdatabase()
	fmt.Println("Database Created")
	defer database.Close()

	mux := goji.NewMux()

	mux.HandleFunc(pat.Get("/"), home)

	mux.HandleFunc(pat.Get("/inscription"), inscription)
	mux.HandleFunc(pat.Post("/login"), login)
	mux.HandleFunc(pat.Get("/logout"), logout)
	mux.HandleFunc(pat.Get("/me"), connectedUser)

	mux.HandleFunc(pat.Post("/users"), postUsers)
	mux.HandleFunc(pat.Get("/users"), getUsers)

	// mux.HandleFuncC(pat.Put("/users/:id"), putUsers)
	// mux.HandleFuncC(pat.Delete("/users/:id"), deleteUsers)

	// mux.HandleFuncC(pat.Put("/users/me/sexe/"), putUsersSexe)
	// mux.HandleFuncC(pat.Put("/users/me/orientation/"), putUsersOrientation)
	// mux.HandleFuncC(pat.Put("/users/me/bio"), putUsersBio)
	// mux.HandleFuncC(pat.Put("/users/me/firstname"), putUsersFirstname)
	// mux.HandleFuncC(pat.Put("/users/me/lastname"), putUsersLastname)
	// mux.HandleFuncC(pat.Put("/users/me/mail"), putUsersMail)
	// mux.HandleFuncC(pat.Put("/users/me/password"), putUsersPassword)

	// User's interests Road
	mux.HandleFunc(pat.Post("/users/me/interests/"), postUsersInterests)
	mux.HandleFunc(pat.Get("/users/me/interests/"), getUsersInterests)
	mux.HandleFuncC(pat.Delete("/users/me/interests/:interestid"), deleteUsersInterests)

	// // User's images Road
	// mux.HandleFuncC(pat.Post("/users/:id/images/"), postUsersImages)
	// mux.HandleFuncC(pat.Put("/users/:id/images/:idimage"), putUsersImageProfile)
	// mux.HandleFuncC(pat.Delete("/users/:id/images/:idimage"), deleteUsersImages)

	// // Interests
	mux.HandleFunc(pat.Post("/interests/"), postInterests)
	// mux.HandleFunc(pat.Get("/interests/"), getInterests)
	// mux.HandleFuncC(pat.Delete("/interests/:id"), deleteInterests)

	// // Search
	// mux.HandleFuncC(pat.Get("search/interests/:id/Users"), getInterestsUsers)
	// mux.HandleFuncC(pat.Get("search/sex/:id/Users"), getSexUsers)
	// mux.HandleFuncC(pat.Get("search/orientation/:id/Users"), getOrientationUsers)
	// mux.HandleFuncC(pat.Get("search/sexe/:sexeid/orientation/:orientationid/Users"), getSexeOrientationUsers)

	mux.HandleFunc(pat.Get("/css/*"), staticCssFiles)
	mux.HandleFunc(pat.Get("/js/*"), staticJsFiles)

	log.Fatal(http.ListenAndServe(":4242", mux))

}
