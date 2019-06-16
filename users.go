package main

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func users(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	if user.Type != "Admin" {
		permissionDenied(w, r)
		return
	}
	rows, err := db.Query(`
	SELECT ID, Name, Email, Type
	FROM Users;
	`)
	if err != nil {
		log.Fatalln(err)
	}
	users := []User{}
	for rows.Next() {
		u := User{}
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Type); err != nil {
			log.Fatalln(err)
		}
		users = append(users, u)
	}
	page := struct {
		User
		Users []User
	}{user, users}
	if err := tpl.ExecuteTemplate(w, "users.html", page); err != nil {
		log.Fatalln(err)
	}
}

func importUsers(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	if user.Type != "Admin" {
		permissionDenied(w, r)
		return
	}
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

type userData struct {
	Name     string
	Email    string
	Password string
	Type     string
}

func initUsers(file io.Reader) {
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
		users = append(users, userData{record[0], record[1], record[2], "Teacher"})
		log.Println(record)
	}
	log.Println(users)
	for _, user := range users {
		if err := user.signUp(); err != nil {
			log.Fatalln(err)
		}
	}
}

func (user *userData) signUp() (err error) {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		INSERT INTO Users (Name, Email, Password, Type)
		VALUES (?, ?, ?, ?);
		`, user.Name, user.Email, password, user.Type)
	if err != nil {
		return err
	}
	return nil
}

func signUp(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	if user.Type != "Admin" {
		permissionDenied(w, r)
		return
	}
	u := userData{
		r.FormValue("name"),
		r.FormValue("email"),
		r.FormValue("password"),
		"Teacher",
	}
	if err := u.signUp(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	page := struct {
		User
		Title   string
		Content string
		Target  string
	}{user, "註冊成功", u.Name + "註冊成功", ""}
	if err := tpl.ExecuteTemplate(w, "msg.html", page); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
