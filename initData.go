package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func randomPassword() string {
	/*
		To genrate a random password with 16 chars
	*/
	chars := make([]byte, 12)
	rand.Read(chars)
	return base64.StdEncoding.EncodeToString(chars)
}

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
	dbUser := "app"
	dbPassword := randomPassword()
	log.Println("init db")
	log.Println(dbPassword)
	err := initDB(dbName, dbUser, dbPassword)
	if err != nil {
		log.Println(err)
		tpl.ExecuteTemplate(w, "init-db.html", err.Error())
		return
	}
	log.Println("init devices")
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
		`<a href="/init/users">匯入使用者</a>`,
	}
	if err := tpl.ExecuteTemplate(w, "msg.html", page); err != nil {
		log.Fatalln(err)
	}
}

func handleInitUsers(w http.ResponseWriter, r *http.Request) {
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
		"使用者資料已經匯入成功，請重啟伺服器後開始使用",
		``,
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

func initDB(dbName, dbUser, dbPassword string) (err error) {
	cmd := exec.Command("./cmd/init-db.sh", dbName, dbUser, dbPassword)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err = cmd.Run(); err != nil {
		log.Println(err.Error())
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	file, err := os.Create(wd + "/env.json")
	if err != nil {
		return err
	}
	defer file.Close()
	data := struct {
		DBName   string
		User     string
		Password string
	}{dbName, dbUser, dbPassword}
	enc := json.NewEncoder(file)
	if err := enc.Encode(data); err != nil {
		return err
	}
	if db, err = loadDB(); err != nil {
		return err
	}
	return nil
}

func execSQLCommand(db *sql.DB, command string, args ...interface{}) {
	_, err := db.Exec(command, args...)
	if err != nil {
		log.Panic(err)
	}
}
