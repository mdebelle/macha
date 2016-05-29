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
	interests   []Interest
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
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `interest` (`id` int(11) NOT NULL AUTO_INCREMENT, `label` varchar(40) NOT NULL, `userid` int(11) NOT NULL, PRIMARY KEY (`id`))")
	checkErr(err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `user` ( `id` int(11) NOT NULL AUTO_INCREMENT, `username` varchar(40) NOT NULL, `firstname` varchar(40) NOT NULL, `lastname` varchar(40) NOT NULL, `email` varchar(255) NOT NULL, `sexe` tinyint(4) NOT NULL, `orientation` tinyint(4) NOT NULL, `bio` longtext NOT NULL, PRIMARY KEY (`id`))")
	checkErr(err)

	database = db
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

func getUserById(id string) (User, error) {
	rows, err := database.Query("SELECT * FROM user WHERE id=?", id)
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
