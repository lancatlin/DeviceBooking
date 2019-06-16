package main

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"strings"

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
	f, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	msg := initUsers(f)
	if len(msg) != 0 {
		page := struct {
			User
			Title   string
			Content string
			Target  string
		}{
			user,
			"匯入過程發生錯誤",
			strings.Join(msg, "<br>"),
			"<a href='/users'>前一頁</a>",
		}
		if err := tpl.ExecuteTemplate(w, "msg.html", page); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
	page := struct {
		User
		Title   string
		Content string
		Target  string
	}{
		user,
		"匯入使用者資料成功！",
		"使用者資料已經匯入成功",
		`<a href='/users'>回到使用者管理</a>`,
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

func initUsers(file io.Reader) (msg []string) {
	reader := csv.NewReader(file)
	var users []userData
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		users = append(users, userData{record[0], record[1], record[2], "Teacher"})
	}
	msg = []string{}
	for _, user := range users[1:] {
		if err := user.signUp(); err != nil {
			msg = append(msg, err.Error())
		}
	}
	return
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
