package main

import (
	"database/sql"
)

type User struct {
	Id          int
	Username    string
	Firstname   string
	Lastname    string
	Email       string
	Password    []byte
	Sexe        sql.NullInt64
	Orientation sql.NullInt64
	Bio         sql.NullString
	Popularite  sql.NullInt64
	interests   []string
	images      []Image
}

type Interest struct {
	Id    int64
	Label string
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
	//  DEFAULT CHARACTER SET utf8
	//  DEFAULT COLLATE utf8_general_ci
	//
	checkErr(err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `image` (`id` int(11) NOT NULL AUTO_INCREMENT,`label` varchar(250) NOT NULL, `description` longtext NOT NULL, `userid` int(11) NOT NULL, PRIMARY KEY (`id`))")
	checkErr(err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `interest` (`id` int(11) NOT NULL AUTO_INCREMENT, `label` varchar(40) NOT NULL, PRIMARY KEY (`id`))")
	checkErr(err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `userinterest` (`id` int(11) NOT NULL AUTO_INCREMENT, `interestid` varchar(40) NOT NULL, `userid` int(11) NOT NULL, PRIMARY KEY (`id`))")
	checkErr(err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `user` ( `id` int(11) NOT NULL AUTO_INCREMENT, `username` varchar(40) NOT NULL, `firstname` varchar(40) NOT NULL, `lastname` varchar(40) NOT NULL, `email` varchar(255) NOT NULL, `password` varbinary(255) NOT NULL, `sexe` tinyint(4) DEFAULT NULL, `orientation` tinyint(4) DEFAULT NULL, `bio` longtext DEFAULT NULL, `popularite` int(11) DEFAULT NULL, PRIMARY KEY (`id`))")
	checkErr(err)

	database = db
}
