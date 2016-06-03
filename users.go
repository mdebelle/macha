package main

import (
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

func checkLoginUser(username, password string) (User, error) {

	var user User
	var spassword string

	err := database.QueryRow("SELECT id, Firstname, Lastname password FROM users WHERE username=?", username).Scan(&user.Id, &user.Firstname, &user.Lastname, &spassword)
	switch {
	case err == sql.ErrNoRows:
		return user, errors.New("empty")
	case err != nil:
		return user, err
	}
	if password != spassword {
		return user, errors.New("wrong Password")
	}
	return user, nil
}

func postUsers(w http.ResponseWriter, r *http.Request) {

	fmt.Println("hello")

	p := sha1.Sum([]byte(r.FormValue("password")))

	var user = User{
		Username:  r.FormValue("username"),
		Firstname: r.FormValue("firstname"),
		Lastname:  r.FormValue("lastname"),
		Email:     r.FormValue("email"),
		Password:  string(p[:])}

	fmt.Println(user)
	smt, err := database.Prepare("INSERT user SET username=?, firstname=?, lastname=?, email=?, password=?")
	checkErr(err)
	_, err = smt.Exec(user.Username, user.Firstname, user.Lastname, user.Email, user.Password)
	checkErr(err)
	http.Redirect(w, r, "/", http.StatusFound)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Query("SELECT * FROM user")
	checkErr(err)
	defer rows.Close()
	var users []User
	var i int
	for rows.Next() {
		users = append(users, User{})
		err := rows.Scan(&users[i].Id, &users[i].Username, &users[i].Firstname, &users[i].Lastname, &users[i].Email, &users[i].Password, &users[i].Sexe, &users[i].Orientation, &users[i].Bio, &users[i].Popularite)
		checkErr(err)
		i++
	}
	err = rows.Err()
	checkErr(err)
	fmt.Println(users)
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

// func postUsersInterests(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")

// }

// func getUsersInterests(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")

// }

// func deleteUsersInterests(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")

// }

// func postUsersImages(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")

// }

// func putUsersImageProfile(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")

// }

// func deleteUsersImages(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")

// }
