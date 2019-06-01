package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func returnDevice(dID string) error {
	result, err := db.Exec(`
		UPDATE Records
		SET LentUntil = CURRENT_TIMESTAMP()
		WHERE LentUntil IS NULL and Device = ?;
	`, dID)
	if err != nil {
		log.Fatalln(err)
	}
	if rowNum, err := result.RowsAffected(); err != nil {
		log.Fatalln(err)
	} else if rowNum == 0 {
		return errors.New("Record Not Found")
	} else if rowNum > 1 {
		return fmt.Errorf("Return device fatal: RowsAffected is %d", rowNum)
	}
	return nil
}

func handleReturnDevice(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	if user.Type != "Admin" {
		w.WriteHeader(401)
		return
	}
	dID := r.FormValue("device")
	if err := returnDevice(dID); err != nil && err.Error() == "Record Not Found" {
		w.WriteHeader(404)
		return
	} else if err != nil {
		w.WriteHeader(403)
		return
	}
}

func getReturnList() (result []Booking) {
	/*
		回傳所有的借出中預約
	*/
	result = []Booking{}
	rows, err := db.Query(`
	SELECT ID
	FROM Bookings;
	`)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Panicln(err)
		}
		b := Booking{ID: id}
		if b.alreadyLendout() && !b.alreadyReturned() {
			b, err := getBooking(id)
			checkErr(err, "get booking fatal: ")
			result = append(result, b)
		}
	}
	return
}

func handleReturnList(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	if !user.Login || user.Type != "Admin" {
		permissionDenied(w, r)
		return
	}
	page := struct {
		User
		Bookings []Booking
	}{user, getReturnList()}
	checkErr(tpl.ExecuteTemplate(w, "returnList.html", page), "Execute fatal: ")
}
