package main

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

func randomPassword() string {
	/*
		To genrate a random password with 16 chars
	*/
	chars := make([]byte, 20)
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
	adminPassword := r.FormValue("password")
	err := initDB(dbName, dbUser, dbPassword)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	log.Println("init devices")
	user := userData{
		"admin",
		"",
		adminPassword,
		"Admin",
	}
	if err := user.signUp(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	initDevices()
	page := struct {
		User
		Title   string
		Content string
		Target  string
	}{
		nilUser(),
		"資料庫建立成功！",
		"資料庫已建立成功，重新啟動後即可開始使用", "",
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

func initDB(dbName, dbUser, dbPassword string) (err error) {
	cmd := exec.Command("./cmd/init-db.sh", dbName, dbUser, dbPassword)
	var buf bytes.Buffer
	cmd.Stderr = &buf
	cmd.Stdout = &buf
	if err = cmd.Run(); err != nil {
		log.Println(buf.String())
		return errors.New(buf.String())
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
