package main

import (
	"encoding/json"
	"fmt"
	"goji.io/pat"
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
	database.QueryRow("SELECT LastName FROM user WHERE id=?", session.Values["UserInfo"].(UserData).Id).Scan(&u.LastName)
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

func uptdateUsersBio(w http.ResponseWriter, r *http.Request) {

	var data PostAge

	fmt.Println("lololo")

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
	checkErr(err)
	r.Body.Close()
	err = json.Unmarshal(body, &data)
	checkErr(err)

	smt, err := database.Prepare("UPDATE user SET user.Bio=? WHERE id=?")
	checkErr(err)
	defer smt.Close()
	fmt.Println(data.Date)
	_, err = smt.Exec(data.Date, session.Values["UserInfo"].(UserData).Id)
	checkErr(err)

	var u = session.Values["UserInfo"].(UserData)
	u.Bio.String = data.Date
	session.Values["UserInfo"] = u
	session.Save(r, w)

	writeJson(w, ResponseStatus{Status: "ok"})
}

func getUsersMatches(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	smt, err := database.Prepare("SELECT id, Username, Birthdate FROM user WHERE id!=?")
	checkErr(err)
	defer smt.Close()
	rows, err := smt.Query(session.Values["UserInfo"].(UserData).Id)
	checkErr(err)

	var users []SimpleUser
	var i int
	for rows.Next() {
		users = append(users, SimpleUser{})
		var dob []uint8
		err := rows.Scan(&users[i].Id, &users[i].UserName, &dob)
		if dob != nil {
			users[i].Bod = transformAge(dob)
			checkErr(err)
		}
		checkErr(err)
		i++
	}

	writeJson(w, users)
}

func publicProfile(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	var user SimpleUser
	var c bool
	id, _ := strconv.ParseInt(pat.Param(ctx, "id"), 10, 64)

	session, _ := store.Get(r, "session")
	if session.Values["connected"] == true {
		c = true
	}

	smt, err := database.Prepare("SELECT user.username, user.birthdate FROM user WHERE id=?")
	checkErr(err)
	var dob []uint8
	smt.QueryRow(id).Scan(&user.UserName, &dob)
	user.Id = id
	user.Bod = transformAge(dob)
	if c == false {
		renderTemplate(w, "publicProfile", &publicProfileView{
			Header: HeadData{
				Title:      "Profile",
				Stylesheet: []string{"publicProfile.css"}},
			Connection: false,
			Profile:    user})
		visitedProfile("unknown", id)
	} else {
		renderTemplate(w, "publicProfile", &publicProfileView{
			Header: HeadData{
				Title:      "Profile",
				Stylesheet: []string{"publicProfile.css"}},
			Connection: true,
			User:       session.Values["UserInfo"].(UserData),
			Profile:    user})
		visitedProfile(session.Values["UserInfo"].(UserData).UserName, id)
	}

}

func likeAnUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	profileid := pat.Param(ctx, "id")
	fmt.Println(profileid)

}

func unlikeAnUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	profileid := pat.Param(ctx, "id")
	fmt.Println(profileid)

}
