package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

func postInterests(w http.ResponseWriter, r *http.Request) {

}

func getInterestId(label string, userid int) int64 {

	fmt.Println("je cherche l'id d'un interet ou je l'ajoute")

	var id int64 = -1
	stmt, err := database.Prepare("SELECT id FROM interest WHERE label=?")
	defer stmt.Close()
	err = stmt.QueryRow(label).Scan(&id)

	if err == sql.ErrNoRows {
		smt, err := database.Prepare("INSERT interest SET label=?")
		defer smt.Close()
		checkErr(err)
		res, err := smt.Exec(label)
		checkErr(err)
		id, err = res.LastInsertId()

	}
	if id != -1 {
		smt, err := database.Prepare("INSERT IGNORE INTO userinterest (interestid, userid) VALUES (?, ?)")
		checkErr(err)
		defer smt.Close()
		res, err := smt.Exec(id, userid)
		checkErr(err)
		t, err := res.RowsAffected()
		checkErr(err)
		fmt.Println("rows", t)
		if t == 0 {
			return -1
		}
	}
	return id
}

// func getInterests(w http.ResponseWriter, r *http.Request) {

// }

// func deleteInterests(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")
// }
