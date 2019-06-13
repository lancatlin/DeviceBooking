package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func handleInitDB(w http.ResponseWriter, r *http.Request) {
	/*
		Handle post from /init/db
	*/
	if db != nil {
		// 如果已經有資料庫
		w.WriteHeader(403)
		return
	}
	dbName := r.FormValue("db-name")
	dbUser := r.FormValue("db-user")
	dbPassword := r.FormValue("db-password")
	initDB(dbName, dbUser, dbPassword)
	initDevices()
	page := struct {
		User
		Title   string
		Content string
		Target  string
	}{
		nilUser(),
		"資料庫建立成功！",
		"資料庫已建立成功，請匯入使用者",
		`<a href="/init/user">匯入使用者</a>`,
	}
	if err := tpl.ExecuteTemplate(w, "msg.html", page); err != nil {
		log.Fatalln(err)
	}
}

func handleInitUser(w http.ResponseWriter, r *http.Request) {
	f, h, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("header: ", h)
	defer f.Close()
	initUsers(f)
	page := struct {
		Title   string
		Content string
		Target  string
	}{
		"匯入使用者資料成功！",
		"使用者資料已經匯入成功，請登入後開始使用",
		`<a href="/login">登入</a>`,
	}
	if err := tpl.ExecuteTemplate(w, "msg.html", page); err != nil {
		log.Fatalln(err)
	}
}

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

func initUsers(file io.Reader) {
	type userData struct {
		Name     string
		Email    string
		Password string
	}
	reader := csv.NewReader(file)
	var users []userData
	for {
		record, err := reader.Read()
		if err != io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		users = append(users, userData{record[0], record[1], record[2]})
		log.Println(record)
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

func initDB(dbName, dbUser, dbPassword string) {
	db, err := sql.Open("mysql", fmt.Sprintf(`%s@%s@unix(/var/run/mysqld/mysqld.sock)/?parseTime=true`, dbUser, dbPassword))
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	execSQLCommand(db, `CREATE DATABASE ?;`, dbName)
	execSQLCommand(db, `
	CREATE TABLE IF NOT EXISTS Users (
		ID INT PRIMARY KEY AUTO_INCREMENT,
		Email VARCHAR(64) UNIQUE NOT NULL,
		Name VARCHAR(32) NOT NULL,
		Type ENUM('Admin', 'Teacher') DEFAULT 'Teacher',
		Password BLOB NOT NULL
	);`)
	execSQLCommand(db, `
	CREATE TABLE IF NOT EXISTS Devices (
		ID VARCHAR(8) PRIMARY KEY,
		Type ENUM('Student-iPad', 'Teacher-iPad', 'Chromebook', 'WAP', 'WirelessProjector') NOT NULL,
		JoinDate DATETIME
	);`)
	execSQLCommand(db, `
	CREATE TABLE IF NOT EXISTS Bookings (
		ID INT PRIMARY KEY AUTO_INCREMENT,
		User INT NOT NULL,
		LendingTime DATETIME NOT NULL,
		ReturnTime DATETIME NOT NULL,
		Done BOOL DEFAULT false,
		FOREIGN KEY (User) REFERENCES Users (ID)
	);`)
	execSQLCommand(db, `
	CREATE TABLE IF NOT EXISTS BookingDevices (
		BID INT,
		Type ENUM('Student-iPad', 'Teacher-iPad', 'Chromebook', 'WAP', 'WirelessProjector'),
		Amount INT DEFAULT 0,
		FOREIGN KEY (BID) REFERENCES Bookings (ID),
		PRIMARY KEY (BID, Type)
	);`)
	execSQLCommand(db, `
	CREATE TABLE IF NOT EXISTS Records (
		Booking INT,
		Device VARCHAR(8),
		LentFrom DATETIME NOT NULL,
		LentUntil DATETIME,
		PRIMARY KEY (Booking, Device),
		FOREIGN KEY (Booking) REFERENCES Bookings (ID),
		FOREIGN KEY (Device) REFERENCES Devices (ID)
	);`)
	execSQLCommand(db, `
	CREATE TABLE IF NOT EXISTS Sessions (
		ID CHAR(36),
		User INT NOT NULL,
		LastUsed TIMESTAMP NOT NULL,
		PRIMARY KEY (ID),
		FOREIGN KEY (User) REFERENCES Users (ID)
	);`)
	execSQLCommand(db, `
	CREATE OR REPLACE VIEW UnDoneRecords AS 
	SELECT Booking, Device
	FROM Records
	WHERE LentUntil IS NULL;`)
	execSQLCommand(db, `
	CREATE OR REPLACE VIEW UnDoneBookings AS 
	SELECT ID, COUNT(Device) AS Amount
	FROM Bookings B
	LEFT JOIN UnDoneRecords R
	ON B.ID = R.Booking 
	GROUP BY ID;`)
	execSQLCommand(db, `
	CREATE OR REPLACE VIEW DevicesStatus AS 
	SELECT D.ID, COUNT(Device) AS Status, Name, D.Type
	FROM Devices D
	LEFT JOIN UnDoneRecords R
	ON D.ID = Device
	LEFT JOIN Bookings B
	ON B.ID = Booking 
	LEFT JOIN Users U
	ON U.ID = User 
	GROUP BY D.ID;`)
	execSQLCommand(db, "CREATE USER IF NOT EXISTS '?'@'localhost' IDENTIFIED BY '?';", dbUser, dbPassword)
	execSQLCommand(db, "GRANT ALL ON ?.* TO '?'@'localhost';", dbName, dbUser)
}

func execSQLCommand(db *sql.DB, command string, args ...interface{}) {
	_, err := db.Exec(command, args...)
	if err != nil {
		log.Fatalln(err)
	}
}
