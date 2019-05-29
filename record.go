package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (b *Booking) getRecords() (list []Record) {
	rows, err := db.Query(`
		SELECT Device, Type, LentUntil IS NOT NULL
		FROM Records R, Devices D
		WHERE Device = D.ID and Booking = ?
		ORDER BY Type;
	`, b.ID)
	if err != nil {
		log.Fatalln("Query fatal: ", err)
	}
	defer rows.Close()
	for rows.Next() {
		r := Record{}
		r.Booking = b.ID
		err := rows.Scan(&r.Device, &r.Type, &r.Done)
		if err != nil {
			log.Fatalln("Scan fatal: ", err)
		}
		list = append(list, r)
	}
	return
}

func recordList(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	if !user.Login {
		permissionDenied(w, r)
		return
	}
	bID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		// Not found
		notFound(w, r)
		return
	}
	b, err := getBooking(int64(bID))
	if err == ErrBookingNotFound {
		notFound(w, r)
		return
	}
	page := struct {
		User
		Booking
		Records []Record
	}{user, b, b.getRecords()}
	if err := tpl.ExecuteTemplate(w, "record-list.html", page); err != nil {
		log.Fatalln(err)
	}
}
