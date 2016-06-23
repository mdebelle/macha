package main

import (
	"database/sql"
	"net/http"
)

func postInterests(w http.ResponseWriter, r *http.Request) {

}

func getInterestId(label string, userid int) int64 {

	var id int64 = -1
	stmt, err := database.Prepare("SELECT id FROM interest WHERE label=?")
	defer stmt.Close()
	err = stmt.QueryRow(label).Scan(&id)

	if err == sql.ErrNoRows {
		smt, err := database.Prepare("INSERT interest SET label=?")
		defer smt.Close()
		checkErr(err, "getInterestId")
		res, err := smt.Exec(label)
		checkErr(err, "getInterestId")
		id, err = res.LastInsertId()

	}
	if id != -1 {
		smt, err := database.Prepare("INSERT IGNORE INTO userinterest (interestid, userid) VALUES (?, ?)")
		checkErr(err, "getInterestId")
		defer smt.Close()
		res, err := smt.Exec(id, userid)
		checkErr(err, "getInterestId")
		t, err := res.RowsAffected()
		checkErr(err, "getInterestId")
		if t == 0 {
			return -1
		}
	}
	return id
}
