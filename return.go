package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func returnDevice(dID string) (count int, err error) {
	row := db.QueryRow(`
	SELECT COUNT(1) FROM UnDoneRecords
	WHERE Booking = (SELECT Booking FROM UnDoneRecords WHERE Device = ?);
	`, dID)
	if err = row.Scan(&count); err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, errors.New("record not found")
	}
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
		return 0, errors.New("Record Not Found")
	} else if rowNum > 1 {
		return 0, fmt.Errorf("Return device fatal: RowsAffected is %d", rowNum)
	}
	return count - 1, nil
}

func handleReturnDevice(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	if user.Type != "Admin" {
		w.WriteHeader(401)
		return
	}
	dID := r.FormValue("device")
	count, err := returnDevice(dID)
	if err != nil && err.Error() == "Record Not Found" {
		w.WriteHeader(404)
		return
	} else if err != nil {
		w.WriteHeader(403)
		return
	}
	if count == 0 {
		w.WriteHeader(201)
	}
}

func getLendingList() (result []Booking, err error) {
	/*
		回傳所有的借出中預約
	*/
	result = []Booking{}
	rows, err := db.Query(`
	SELECT UB.ID
	FROM UnDoneBookings UB, Bookings B
	WHERE UB.ID = B.ID AND Amount != 0
	ORDER BY LendingTime;
	`)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			log.Panicln(err)
		}
		b, _ := getBooking(id)
		b.Status = StatusLending
		result = append(result, b)
	}
	return
}

func handleLendingList(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	if !user.Login || user.Type != "Admin" {
		permissionDenied(w, r)
		return
	}
	bookings, err := getLendingList()
	if err != nil {
		log.Panicln(err)
	}
	page := struct {
		User
		Bookings []Booking
	}{user, bookings}
	checkErr(tpl.ExecuteTemplate(w, "returnList.html", page), "Execute fatal: ")
}
