package main

import (
	"database/sql"
	//	_ "github.com/go-sql-driver/mysql"
	"errors"
	"fmt"
)

type User struct {
	Id          int
	Username    string
	Firstname   string
	Lastname    string
	Email       string
	Sexe        int8
	Orientation int8
	Bio         string
	interests   []string
	images      []Image
}

type Interest struct {
	id     int
	name   string
	userid int
}

type Image struct {
	id          int
	name        string
	description string
	userid      int
}

var database *sql.DB

func initdatabase() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3307)/machadb?charset=utf8")
	checkErr(err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `image` (`id` int(11) NOT NULL AUTO_INCREMENT,`label` varchar(250) NOT NULL, `description` longtext NOT NULL, `userid` int(11) NOT NULL, PRIMARY KEY (`id`))")
	checkErr(err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `interest` (`id` int(11) NOT NULL AUTO_INCREMENT, `label` varchar(40) NOT NULL, PRIMARY KEY (`id`))")
	checkErr(err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `userinterest` (`id` int(11) NOT NULL AUTO_INCREMENT, `intersetid` varchar(40) NOT NULL, `userid` int(11) NOT NULL, PRIMARY KEY (`id`))")
	checkErr(err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `user` ( `id` int(11) NOT NULL AUTO_INCREMENT, `username` varchar(40) NOT NULL, `firstname` varchar(40) NOT NULL, `lastname` varchar(40) NOT NULL, `email` varchar(255) NOT NULL, `sexe` tinyint(4) NOT NULL, `orientation` tinyint(4) NOT NULL, `bio` longtext NOT NULL, PRIMARY KEY (`id`))")
	checkErr(err)

	database = db
}

func addnewUser(user *User) {
	smt, _ := database.Prepare("INSERT user SET username=?, firstname=?, lastname=?, email=?, sexe=?, orientation=?, bio=?")
	res, _ := smt.Exec(user.Username, user.Firstname, user.Lastname, user.Email, user.Sexe, user.Orientation, user.Bio)
	id, _ := res.LastInsertId()
	if len(user.interests) > 0 {
		addnewInterest(user.interests, id)
	}
}

func addnewInterest(interests []string, userid int64) {
	for _, interest := range interests {
		rows, err := database.Query("SELECT id FROM interest WHERE label=?", interest)
		var id int64
		if rows.Next() {
			rows.Scan(&id)
		} else {
			smt, err := database.Prepare("INSERT interest SET label=?")
			checkErr(err)
			res, err := smt.Exec(interest)
			checkErr(err)
			id, err = res.LastInsertId()
			checkErr(err)
		}
		smt, err := database.Prepare("INSERT interestuser SET interestid=?, userid=?")
		_, err = smt.Exec(id, userid)
		checkErr(err)
	}
}

func getUserByUsername(username string) (User, error) {
	rows, err := database.Query("SELECT * FROM user WHERE username=?", username)
	checkErr(err)
	defer rows.Close()
	if rows == nil {
		fmt.Println("j'y crois pas")
	}
	var user User
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.Sexe, &user.Orientation, &user.Bio)
		checkErr(err)
	} else {
		return user, errors.New("empty")
	}
	err = rows.Err()
	checkErr(err)
	return user, err
}

func getUserId(username string) (int, error) {
	rows, err := database.Query("SELECT * FROM user WHERE username=?", username)
	var user User
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.Sexe, &user.Orientation, &user.Bio)
		checkErr(err)
	} else {
		return -1, errors.New("no user found")
	}
	return user.Id, err
}

func getUserById(id string) (User, error) {
	rows, err := database.Query("SELECT * FROM user WHERE id=?", id)
	checkErr(err)
	defer rows.Close()
	var user User
	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.Sexe, &user.Orientation, &user.Bio)
		checkErr(err)
	} else {
		return user, errors.New("empty")
	}
	err = rows.Err()
	checkErr(err)
	return user, err
}
