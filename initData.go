package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func initDevices() {
	stmt, err := db.Prepare(`
	INSERT INTO Devices
	VALUES (?, ?, CURRENT_TIMESTAMP())
	`)
	checkErr(err, "Statment prepare error: ")
	insertData(stmt, makeIDSlice("ST", 70, 3), "Student-iPad")
	insertData(stmt, makeIDSlice("T", 13, 2), "Teacher-iPad")
	insertData(stmt, makeIDSlice("CB", 32, 3), "Chromebook")
	insertData(stmt, makeIDSlice("C", 90, 3), "Chromebook")
	insertData(stmt, makeIDSlice("AP", 3, 2), "WAP")
	insertData(stmt, makeIDSlice("EZ", 3, 2), "WirelessProjector")
}

func makeIDSlice(head string, number int, length int) (list []string) {
	list = make([]string, number)
	for i := 0; i < number; i++ {
		stmt := "%s%0" + strconv.Itoa(length) + "d"
		list[i] = fmt.Sprintf(stmt, head, i+1)
	}
	return list
}

func insertData(stmt *sql.Stmt, list []string, t string) {
	for _, v := range list {
		_, err := stmt.Exec(v, t)
		checkErr(err, "Insert Error: ")
	}
}

func initUsers(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("Can't open file: ", err)
	}
	defer file.Close()
	type userData struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	dec := json.NewDecoder(file)
	var users []userData
	if err := dec.Decode(&users); err != nil {
		log.Fatalln("json decode fatal: ", err)
	}
	log.Println(users)
	for _, user := range users {
		password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Password generate error: ", err)
		}
		_, err = db.Exec(`
		INSERT INTO Users (Name, Email, Password)
		VALUES (?, ?, ?);
		`, user.Name, user.Email, password)
		if err != nil {
			log.Fatalln("Insert fatal: ", err)
		}
	}
}
