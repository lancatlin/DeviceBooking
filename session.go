package main

import (
	"database/sql"
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func getUser(w http.ResponseWriter, r *http.Request) User {
	cookie, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		return nilUser()
	}
	row := db.QueryRow(`
	SELECT U.ID, Name, Type
	FROM Sessions S, Users U 
	WHERE S.ID = ? AND U.ID = S.User;
	`, cookie.Value)
	var uid int
	var uname string
	var utype string
	err = row.Scan(&uid, &uname, &utype)
	if err == sql.ErrNoRows {
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
		return nilUser()
	}
	go updateSession(cookie.Value)
	return newUser(uid, uname, utype)
}

func updateSession(uuid string) {
	result, err := db.Exec(`
	UPDATE Sessions 
	SET LastUsed = CURRENT_TIMESTAMP()
	WHERE ID = ? ;
	`, uuid)
	if err != nil {
		log.Fatal("Update session error: ", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rows != 1 {
		log.Println("Update session fatal affect not one row: ", rows, uuid)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if getUser(w, r).Login {
		// already login
		http.Redirect(w, r, "/", 303)
		return
	}
	type loginPage struct {
		User
		Error string
	}
	page := loginPage{nilUser(), ""}
	if r.Method != http.MethodPost {
		err := tpl.ExecuteTemplate(w, "login.html", page)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	uname := r.FormValue("email")
	password := r.FormValue("password")
	q := `
	SELECT ID, Password
	FROM Users
	WHERE Email = ? OR Name = ?
	`
	row := db.QueryRow(q, uname, uname)
	var id int
	var hashPassword []byte
	err := row.Scan(&id, &hashPassword)
	if err == sql.ErrNoRows {
		page.Error = "Fatal: Email not found"
		tpl.ExecuteTemplate(w, "login.html", page)
		return
	}
	if bcrypt.CompareHashAndPassword(hashPassword, []byte(password)) != nil {
		page.Error = "Fatal: Email or Password is uncorrect"
		tpl.ExecuteTemplate(w, "login.html", page)
		return
	}
	sessionUUID := uuid.NewV4()
	session := sessionUUID.String()
	if err != nil {
		log.Fatalln("uuid init fatal: ", err)
	}
	result, err := db.Exec(`
	INSERT INTO Sessions
	VALUES (?, ?, CURRENT_TIMESTAMP())
	`, session, id)
	if err != nil {
		log.Fatalln("db prepare fatal: ", err)
	}
	if r, err := result.RowsAffected(); err != nil || r != 1 {
		log.Fatalln("Database insert fatal: ", err, result)
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: session,
	})
	http.Redirect(w, r, "/", 303)
}

func logout(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	if !user.Login { // Didn't login
		http.Redirect(w, r, "/", 303)
	}
	cookie, _ := r.Cookie("session")
	session := cookie.Value
	stmt, err := db.Prepare(`
	DELETE FROM Sessions WHERE ID = ?
	`)
	if err != nil {
		log.Fatalln("db prepare fatal: ", err)
	}
	result, err := stmt.Exec(session)
	if err != nil {
		log.Fatal(err)
	} else {
		raws, _ := result.RowsAffected()
		if raws != 1 {
			log.Fatalln("Effect row not only one: ", result)
		}
	}
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", 303)
}
