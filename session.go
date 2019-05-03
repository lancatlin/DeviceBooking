package main

import (
	"database/sql"
	"log"
	"net/http"
)

func getUser(w http.ResponseWriter, r *http.Request) User {
	cookie, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		return nilUser()
	}
	row := db.QueryRow(`
	SELECT U.ID, Name, Type
	FROM Sessions S, Users U 
	WHERE S.ID = ? AND U.ID = S.User
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
	UPDATE Sessions S
	SET LastUsed = TIMESTAMP()
	WHERE ID = ?
	`, uuid)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rows != 1 {
		log.Fatalln("Update session fatal affect not one row: ", rows)
	}
}
