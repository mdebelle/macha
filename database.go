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
	db, err := sql.Open("mysql", DB_CONNECTION)
	checkErr(err, "database")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `image` ( " +
		"`id` int(11) NOT NULL AUTO_INCREMENT, " +
		"`label` varchar(250) NOT NULL, " +
		"`description` longtext NOT NULL, " +
		"`userid` int(11) NOT NULL, " +
		"PRIMARY KEY (`id`))")
	checkErr(err, "database")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `interest` ( " +
		"`id` int(11) NOT NULL AUTO_INCREMENT, " +
		"`label` varchar(40) NOT NULL, " +
		"PRIMARY KEY (`id`))")
	checkErr(err, "database")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `userinterest` ( " +
		"`id` int(11) NOT NULL AUTO_INCREMENT, " +
		"`interestid` int(11) NOT NULL, " +
		"`userid` int(11) NOT NULL, " +
		"PRIMARY KEY (`id`), UNIQUE KEY `idunik` (`interestid`,`userid`))")
	checkErr(err, "database")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `user` ( " +
		"`id` int(11) NOT NULL AUTO_INCREMENT, " +
		"`username` varchar(40) NOT NULL, " +
		"`firstname` varchar(40) NOT NULL, " +
		"`lastname` varchar(40) NOT NULL, " +
		"`birthdate` date DEFAULT NULL, " +
		"`email` varchar(255) NOT NULL, " +
		"`password` varbinary(255) NOT NULL, " +
		"`sexe` tinyint(4) DEFAULT NULL, " +
		"`orientation` tinyint(4) DEFAULT NULL, " +
		"`bio` longtext DEFAULT NULL, " +
		"`popularity` int(11) DEFAULT NULL, " +
		"PRIMARY KEY (`id`))")
	checkErr(err, "database")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `notification` ( " +
		"`id` int(11) NOT NULL AUTO_INCREMENT, " +
		"`message` text, " +
		"`date` date NOT NULL, " +
		"`read` tinyint(1) NOT NULL DEFAULT '0', " +
		"`userid` int(11) NOT NULL, " +
		"PRIMARY KEY (`id`))")
	checkErr(err, "database")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `last_connexion` ( " +
		"`id` int(11) NOT NULL AUTO_INCREMENT, " +
		"`date` datetime DEFAULT NULL, " +
		"`userid` int(11) NOT NULL, " +
		"PRIMARY KEY (`id`))")
	checkErr(err, "database")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `chatroom` ( " +
		"`id` int(11) NOT NULL AUTO_INCREMENT, " +
		"`chatname_one` varchar(80) NOT NULL, " +
		"`chatname_two` varchar(80) NOT NULL, " +
		"PRIMARY KEY (`id`))")
	checkErr(err, "database")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `message` ( " +
		"`id` int(11) NOT NULL AUTO_INCREMENT, " +
		"`chatid` varchar(80) NOT NULL, " +
		"`msg` longtext NOT NULL, " +
		"`date` varchar(80) NOT NULL, " +
		"PRIMARY KEY (`id`))")
	checkErr(err, "database")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `closedroom` ( " +
		"`id` int(11) NOT NULL AUTO_INCREMENT, " +
		"`chatid` varchar(80) NOT NULL, " +
		"`userid` varchar(80) NOT NULL, " +
		"PRIMARY KEY (`id`))")
	checkErr(err, "database")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `blockeduser` ( " +
		"`id` int(11) NOT NULL AUTO_INCREMENT, " +
		"`baduserid` varchar(80) NOT NULL, " +
		"`userid` varchar(80) NOT NULL, " +
		"PRIMARY KEY (`id`))")
	checkErr(err, "database")
	database = db
}
