package main

import (
	"net/http"
)

func postInterests(w http.ResponseWriter, r *http.Request) {

}

func getInterestId(label string) int64 {

	var id int64
	stmt, err := database.Prepare("SELECT id FROM interest WHERE label=?")
	err = stmt.QueryRow(label).Scan(&id)
	defer stmt.Close()
	if err != nil {
		smt, err := database.Prepare("INSERT interests SET label=?")
		defer smt.Close()
		checkErr(err)
		res, err := smt.Exec(label)
		checkErr(err)
		id, err = res.LastInsertId()
	}
	return id
}

// func getInterests(w http.ResponseWriter, r *http.Request) {

// }

// func deleteInterests(ctx context.Context, w http.ResponseWriter, r *http.Request) {
// 	id := pat.Param(ctx, "id")
// }
