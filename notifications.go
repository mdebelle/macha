package main

import (
	"fmt"
	"goji.io/pat"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

func visitedProfile(vname string, id int64) {
	t := time.Now()
	msg := vname + " a visité ton profile."
	smt, err := database.Prepare("INSERT notification SET message=?, date=?, userid=?")
	checkErr(err)
	defer smt.Close()
	_, err = smt.Exec(msg, t, id)
	checkErr(err)
}

func getNotifications(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")
	if session.Values["connected"] != true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	smt, err := database.Prepare("SELECT message, date FROM notification WHERE notification.userid=? AND notification.read=?")
	checkErr(err)
	defer smt.Close()
	rows, err := smt.Query(session.Values["UserInfo"].(UserData).Id, false)
	checkErr(err)

	var notifs []Notifications
	var i int
	for rows.Next() {
		notifs = append(notifs, Notifications{})
		err := rows.Scan(&notifs[i].Msg, &notifs[i].Date)
		checkErr(err)
		i++
	}
	err = rows.Err()
	checkErr(err)
	fmt.Println(">>>", notifs)
	writeJson(w, notifs)
}

func setReadNotifications(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	id := pat.Param(ctx, "id")
	smt, err := database.Prepare("UPDATE notification SET read=? WHERE id=?")
	checkErr(err)
	defer smt.Close()
	_, err = smt.Exec(true, id)
	checkErr(err)
	writeJson(w, ResponseStatus{Status: "ok"})
}
