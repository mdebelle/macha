package main

import (
	"database/sql"
	"errors"
	"fmt"
	"goji.io/pat"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"net/http"
)

func testEq(a, b []byte) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func checkLoginUser(username string, password []byte) (User, error) {

	var user User
	var spassword []byte

	fmt.Println("verifier le log d'un utilisateur")

	fmt.Println(username)
	err := database.QueryRow("SELECT id, Firstname, Lastname, password FROM user WHERE username=?", username).Scan(&user.Id, &user.Firstname, &user.Lastname, &spassword)
	switch {
	case err == sql.ErrNoRows:
		return user, errors.New("empty")
	case err != nil:
		return user, err
	}
	if testEq(password, spassword) {
		return user, errors.New("wrong Password")
	}
	return user, nil
}

func postUsers(w http.ResponseWriter, r *http.Request) {

	fmt.Println("ajouter un utilisateur")

	p, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	checkErr(err)

	var user = User{
		Username:  r.FormValue("username"),
		Firstname: r.FormValue("firstname"),
		Lastname:  r.FormValue("lastname"),
		Email:     r.FormValue("email"),
		Password:  p}

	smt, err := database.Prepare("INSERT user SET username=?, firstname=?, lastname=?, email=?, password=?")
	checkErr(err)
	defer smt.Close()
	_, err = smt.Exec(user.Username, user.Firstname, user.Lastname, user.Email, user.Password)
	checkErr(err)
	http.Redirect(w, r, "/", http.StatusFound)
}

type usersView struct {
	Title      string
	Stylesheet string
	Firstname  string
	Lastname   string
	Users      []User
}

func getUsers(w http.ResponseWriter, r *http.Request) {

	fmt.Println("recuperer les utilisateurs")

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	smt, err := database.Prepare("SELECT username, sexe, orientation, bio, popularite FROM user")
	checkErr(err)
	defer smt.Close()
	rows, err := smt.Query()
	checkErr(err)
	defer rows.Close()

	var users []User
	var i int
	for rows.Next() {
		users = append(users, User{})
		err := rows.Scan(&users[i].Username, &users[i].Sexe, &users[i].Orientation, &users[i].Bio, &users[i].Popularite)
		checkErr(err)
		i++
	}
	err = rows.Err()
	checkErr(err)

	renderTemplate(w, "users", &usersView{
		Title:      "All Macha Users",
		Stylesheet: "home.css",
		Users:      users})
}

// func putUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")

// }

// func deleteUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")
// }

// func putUsersSexe(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")

// }

// func putUsersOrientation(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")

// }

// func putUsersBio(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")

// }

// func putUsersFirstname(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")

// }

// func putUsersLastname(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")

// }

// func putUsersMail(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")

// }

// func putUsersPassword(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")

// }

func postUsersInterests(w http.ResponseWriter, r *http.Request) {

	fmt.Println("ajouter les interets des utiliateurs")

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	interest := r.FormValue("content")
	interestid := getInterestId(interest)

	if interestid != -1 {
		smt, err := database.Prepare("INSERT userinterest SET interestid=?, userid=?")
		checkErr(err)
		_, err = smt.Exec(interestid, session.Values["id"])
		checkErr(err)
	}
}

// func getUsersInterests(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")

// }

func deleteUsersInterests(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	fmt.Println("supprimez les interets des utiliateurs")

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	id := pat.Param(ctx, "interestid")

	smt, err := database.Prepare("DELETE * WHERE interestid=? AND userid=?")
	checkErr(err)
	_, err = smt.Exec(id, session.Values["id"])
	checkErr(err)

	http.Redirect(w, r, "/", http.StatusFound)
}

// func postUsersImages(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")

// }

// func putUsersImageProfile(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")

// }

// func deleteUsersImages(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")

// }
