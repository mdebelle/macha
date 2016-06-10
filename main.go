package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goji.io"
	"goji.io/pat"
	"log"
	"net/http"
)

type ResponseStatus struct {
	Status string
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

func home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home", &HomeView{
		Header: HeadData{
			Title:      "Homepage",
			Stylesheet: []string{"home.css"}}})
}

func inscription(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "inscription", &inscriptionVew{
		Header: HeadData{
			Title:      "Inscription",
			Stylesheet: []string{"inscription.css"}}})
}

func main() {

	fmt.Println("localhost:4242")

	// Connection a la base de donnee
	initdatabase()
	fmt.Println("Database Created")
	defer database.Close()

	gob.Register(UserData{})

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
