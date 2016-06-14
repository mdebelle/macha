package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"goji.io/pat"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
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

func checkLoginUser(username string, password []byte) (UserData, error) {

	var user UserData
	var spassword []byte

	err := database.QueryRow("SELECT id, Firstname, Lastname, BirthDate, Email, Bio, Sexe, Orientation, Popularity, password FROM user WHERE username=?", username).Scan(
		&user.Id, &user.FirstName, &user.LastName, &user.BirthDate, &user.Email, &user.Bio, &user.Sexe, &user.Orientation, &user.Popularity, &spassword)
	switch {
	case err == sql.ErrNoRows:
		return user, errors.New("empty")
	case err != nil:
		return user, err
	}
	if bcrypt.CompareHashAndPassword(spassword, password) != nil {
		return user, errors.New("wrong Password")
	}
	user.UserName = username
	return user, nil
}

func postUsers(w http.ResponseWriter, r *http.Request) {

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

func getUsers(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	smt, err := database.Prepare("SELECT id, username, sexe, orientation, bio, popularite FROM user")
	checkErr(err)
	defer smt.Close()
	rows, err := smt.Query()
	checkErr(err)
	defer rows.Close()

	var users []UserData
	var i int
	for rows.Next() {
		users = append(users, UserData{})
		err := rows.Scan(&users[i].Id, &users[i].UserName, &users[i].Sexe, &users[i].Orientation, &users[i].Bio, &users[i].Popularity)
		checkErr(err)
		users[i].Interests = getUserInterestsList(users[i].Id)
		i++
	}
	err = rows.Err()
	checkErr(err)
	renderTemplate(w, "users", &usersView{
		Header: HeadData{
			Title:      "Profile",
			Stylesheet: []string{"users.css"}},
		Users: users})
}

func getUserInterestsList(userid int) []Interest {

	smt, err := database.Prepare("SELECT interest.id interest.label FROM interest INNER JOIN userinterest ON interest.id=userinterest.interestid WHERE userinterest.userid=?")
	checkErr(err)
	defer smt.Close()
	rows, err := smt.Query(userid)
	checkErr(err)

	var interests []Interest
	var i int
	for rows.Next() {
		interests = append(interests, Interest{})
		err := rows.Scan(&interests[i].Id, &interests[i].Label)
		checkErr(err)
		i++
	}
	return interests
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

	var interest Interest

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
	checkErr(err)
	r.Body.Close()
	err = json.Unmarshal(body, &interest)
	checkErr(err)

	interest.Id = getInterestId(interest.Label, session.Values["UserInfo"].(UserData).Id)

	if interest.Id == -1 {
		writeJson(w, ResponseStatus{Status: "ok"})
	} else {
		writeJson(w, ResponseStatus{Status: strconv.Itoa(int(interest.Id))})
	}
}

func getUsersInterests(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	smt, err := database.Prepare("SELECT interest.id, interest.label FROM interest INNER JOIN userinterest ON interest.id=userinterest.interestid WHERE userinterest.userid=?")
	checkErr(err)
	defer smt.Close()
	rows, err := smt.Query(session.Values["UserInfo"].(UserData).Id)
	checkErr(err)

	var interests []Interest
	var i int
	for rows.Next() {
		interests = append(interests, Interest{})
		err := rows.Scan(&interests[i].Id, &interests[i].Label)
		checkErr(err)
		i++
	}
	err = rows.Err()
	checkErr(err)
	writeJson(w, interests)
}

func deleteUsersInterests(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	id := pat.Param(ctx, "interestid")

	smt, err := database.Prepare("DELETE FROM userinterest WHERE interestid=? AND userid=?")
	checkErr(err)
	defer smt.Close()
	_, err = smt.Exec(id, session.Values["UserInfo"].(UserData).Id)
	checkErr(err)
	writeJson(w, ResponseStatus{Status: "ok"})

}

type PostAge struct {
	Date string
}

func postUsersAge(w http.ResponseWriter, r *http.Request) {

	var date PostAge

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
	checkErr(err)
	r.Body.Close()
	err = json.Unmarshal(body, &date)
	checkErr(err)

	smt, err := database.Prepare("UPDATE user SET user.birthdate=? WHERE id=?")
	checkErr(err)
	defer smt.Close()
	_, err = smt.Exec(date.Date, session.Values["UserInfo"].(UserData).Id)
	checkErr(err)

	// Mettre a jours la session
	var u = session.Values["UserInfo"].(UserData)

	fmt.Println(u)

	database.QueryRow("SELECT BirthDate FROM user WHERE id=?", session.Values["UserInfo"].(UserData).Id).Scan(&u.BirthDate)
	checkErr(err)
	session.Values["UserInfo"] = u
	session.Save(r, w)

	writeJson(w, ResponseStatus{Status: "ok"})

}

func getUsersAge(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	fmt.Println(session.Values["UserInfo"].(UserData))
	writeJson(w, ResponseStatus{Status: string(session.Values["UserInfo"].(UserData).BirthDate)})
}

func postUsersUsername(w http.ResponseWriter, r *http.Request) {

	var date PostAge

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
	checkErr(err)
	r.Body.Close()
	err = json.Unmarshal(body, &date)
	checkErr(err)

	smt, err := database.Prepare("UPDATE user SET user.username=? WHERE id=?")
	checkErr(err)
	defer smt.Close()
	_, err = smt.Exec(date.Date, session.Values["UserInfo"].(UserData).Id)
	checkErr(err)

	// Mettre a jours la session
	var u = session.Values["UserInfo"].(UserData)

	fmt.Println(u)

	database.QueryRow("SELECT UserName FROM user WHERE id=?", session.Values["UserInfo"].(UserData).Id).Scan(&u.UserName)
	checkErr(err)
	session.Values["UserInfo"] = u
	session.Save(r, w)

	writeJson(w, ResponseStatus{Status: "ok"})

}

func getUsersUsername(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	fmt.Println(session.Values["UserInfo"].(UserData))
	writeJson(w, ResponseStatus{Status: session.Values["UserInfo"].(UserData).UserName})
}

func postUsersFirstName(w http.ResponseWriter, r *http.Request) {

	var date PostAge

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
	checkErr(err)
	r.Body.Close()
	err = json.Unmarshal(body, &date)
	checkErr(err)

	smt, err := database.Prepare("UPDATE user SET user.FirstName=? WHERE id=?")
	checkErr(err)
	defer smt.Close()
	_, err = smt.Exec(date.Date, session.Values["UserInfo"].(UserData).Id)
	checkErr(err)

	// Mettre a jours la session
	var u = session.Values["UserInfo"].(UserData)

	fmt.Println(u)

	database.QueryRow("SELECT FirstName FROM user WHERE id=?", session.Values["UserInfo"].(UserData).Id).Scan(&u.FirstName)
	checkErr(err)
	session.Values["UserInfo"] = u
	session.Save(r, w)

	writeJson(w, ResponseStatus{Status: "ok"})

}

func getUsersFirstName(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	fmt.Println(session.Values["UserInfo"].(UserData))
	writeJson(w, ResponseStatus{Status: session.Values["UserInfo"].(UserData).FirstName})
}


func postUsersLastName(w http.ResponseWriter, r *http.Request) {

	var date PostAge

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
	checkErr(err)
	r.Body.Close()
	err = json.Unmarshal(body, &date)
	checkErr(err)

	smt, err := database.Prepare("UPDATE user SET user.LastName=? WHERE id=?")
	checkErr(err)
	defer smt.Close()
	_, err = smt.Exec(date.Date, session.Values["UserInfo"].(UserData).Id)
	checkErr(err)

	// Mettre a jours la session
	var u = session.Values["UserInfo"].(UserData)

	fmt.Println(u)

	database.QueryRow("SELECT LastName FROM user WHERE id=?", session.Values["UserInfo"].(UserData).Id).Scan(&u.LastName)
	checkErr(err)
	session.Values["UserInfo"] = u
	session.Save(r, w)

	writeJson(w, ResponseStatus{Status: "ok"})

}

func getUsersLastName(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	fmt.Println(session.Values["UserInfo"].(UserData))
	writeJson(w, ResponseStatus{Status: session.Values["UserInfo"].(UserData).LastName})
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
