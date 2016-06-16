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
	"time"
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

func transformAge(dob []uint8) int {
	const layouttime = "2006-01-02"
	tnow := time.Now()
	t, _ := time.Parse(layouttime, string(dob))
	return tnow.Year() - t.Year()
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

	// Creation d'un compte
	mux.HandleFunc(pat.Get("/inscription"), inscription)
	mux.HandleFunc(pat.Post("/users"), postUsers)

	// Connexion Deconnexion utilisateur
	mux.HandleFunc(pat.Post("/login"), login)
	mux.HandleFunc(pat.Get("/logout"), logout)
	mux.HandleFunc(pat.Get("/me"), connectedUser)

	// User's interests Road
	mux.HandleFunc(pat.Post("/users/me/interests/"), postUsersInterests)
	mux.HandleFunc(pat.Get("/users/me/interests/"), getUsersInterests)
	mux.HandleFuncC(pat.Delete("/users/me/interests/:interestid"), deleteUsersInterests)
	// User's age Road
	mux.HandleFunc(pat.Put("/users/me/age/"), postUsersAge)
	mux.HandleFunc(pat.Get("/users/me/age/"), getUsersAge)
	// User's username Road
	mux.HandleFunc(pat.Put("/users/me/username/"), postUsersUsername)
	mux.HandleFunc(pat.Get("/users/me/username/"), getUsersUsername)
	// User's firstname Road
	mux.HandleFunc(pat.Put("/users/me/firstname/"), postUsersFirstName)
	mux.HandleFunc(pat.Get("/users/me/firstname/"), getUsersFirstName)
	// User's lastname Road
	mux.HandleFunc(pat.Put("/users/me/lastname/"), postUsersLastName)
	mux.HandleFunc(pat.Get("/users/me/lastname/"), getUsersLastName)

	// Public Profile
	mux.HandleFuncC(pat.Get("/users/:id"), publicProfile)
	mux.HandleFuncC(pat.Put("/users/:id/like/"), likeAnUser)
	mux.HandleFuncC(pat.Put("/users/:id/unlike/"), unlikeAnUser)
	// Notifications
	mux.HandleFuncC(pat.Put("/notifications/:id"), setReadNotifications)
	mux.HandleFunc(pat.Get("/notifications"), getNotifications)

	//Matches
	mux.HandleFunc(pat.Get("/users/me/matches/"), getUsersMatches)

	// Interests
	mux.HandleFunc(pat.Post("/interests/"), postInterests)

	// Static Files
	mux.HandleFunc(pat.Get("/css/*"), staticCssFiles)
	mux.HandleFunc(pat.Get("/js/*"), staticJsFiles)

	log.Fatal(http.ListenAndServe(":4242", mux))

}
