package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"goji.io"
	"goji.io/pat"
	// "html"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	tmplDir = "tmpl/"
	cssDir  = "./css/"
	dataDir = "data/"
)

func home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home", &RequiredData{Title: "Homepage", Stylesheet: "home.css"})
}

func inscription(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "inscription", &RequiredData{Title: "Inscription", Stylesheet: "inscription.css"})
}

func staticfiles(w http.ResponseWriter, r *http.Request) {
	static_file := r.URL.Path[len("/css/"):]
	if len(static_file) != 0 {
		f, err := http.Dir(cssDir).Open(static_file)
		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, r, static_file, time.Now(), content)
			return
		}
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("hello")
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

	mux.HandleFunc(pat.Post("/users"), postUsers)
	mux.HandleFunc(pat.Get("/users"), getUsers)

	// mux.HandleFuncC(pat.Put("/users/:id"), putUsers)
	// mux.HandleFuncC(pat.Delete("/users/:id"), deleteUsers)

	// mux.HandleFuncC(pat.Put("/users/:id/sexe/:genre"), putUsersSexe)
	// mux.HandleFuncC(pat.Put("/users/:id/orientation/:orientation"), putUsersOrientation)
	// mux.HandleFuncC(pat.Put("/users/:id/bio"), putUsersBio)
	// mux.HandleFuncC(pat.Put("/users/:id/firstname"), putUsersFirstname)
	// mux.HandleFuncC(pat.Put("/users/:id/lastname"), putUsersLastname)
	// mux.HandleFuncC(pat.Put("/users/:id/mail"), putUsersMail)
	// mux.HandleFuncC(pat.Put("/users/:id/password"), putUsersPassword)

	// // User's interests Road
	// mux.HandleFuncC(pat.Post("/users/:id/interests/:interestid"), postUsersInterests)
	// mux.HandleFuncC(pat.Get("/users/:id/interests"), getUsersInterests)
	// mux.HandleFuncC(pat.Delete("/users/:id/interests/:interestid"), deleteUsersInterests)

	// // User's images Road
	// mux.HandleFuncC(pat.Post("/users/:id/images/"), postUsersImages)
	// mux.HandleFuncC(pat.Put("/users/:id/images/:idimage"), putUsersImageProfile)
	// mux.HandleFuncC(pat.Delete("/users/:id/images/:idimage"), deleteUsersImages)

	// // Interests
	// mux.HandleFunc(pat.Post("/interests/"), postInterests)
	// mux.HandleFunc(pat.Get("/interests/"), getInterests)
	// mux.HandleFuncC(pat.Delete("/interests/:id"), deleteInterests)

	// // Search
	// mux.HandleFuncC(pat.Get("search/interests/:id/Users"), getInterestsUsers)
	// mux.HandleFuncC(pat.Get("search/sex/:id/Users"), getSexUsers)
	// mux.HandleFuncC(pat.Get("search/orientation/:id/Users"), getOrientationUsers)
	// mux.HandleFuncC(pat.Get("search/sexe/:sexeid/orientation/:orientationid/Users"), getSexeOrientationUsers)

	mux.HandleFunc(pat.Get("/css/*"), staticfiles)

	log.Fatal(http.ListenAndServe(":4242", mux))

}
