package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func bookingList(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	if !user.Login {
		permissionDenied(w, r)
		return
	}
	type displayBooking struct {
		ID      int
		User    string
		Devices [5]int
	}
	type classBooking struct {
		Class    string
		Bookings []displayBooking
		Devices  [5]int
	}
	date, day := getDateClass()
	page := struct {
		User
		Date    string
		Day     string
		Classes []classBooking
	}{user, date, day, []classBooking{}}
	stmt, err := db.Prepare(`
	SELECT B.ID, U.Name
	FROM Bookings B, Users U
	WHERE ReturnTime > ? and LendingTime < ? and U.ID = B.User
	`)
	checkErr(err, "booking list check query prepare fatal: ")
	for c := 0; c < len(className); c++ {
		thisClass := classBooking{className[c], []displayBooking{}, [5]int{}}
		rows, err := stmt.Query(date+" "+classBegin[c], date+" "+classEnd[c])
		checkErr(err, "bookings query error: ")
		defer rows.Close()
		for rows.Next() {
			b := displayBooking{}
			err := rows.Scan(&b.ID, &b.User)
			checkErr(err, "row scan fatal: ")
			b.Devices = getBookingDevices(b.ID)
			thisClass.Bookings = append(thisClass.Bookings, b)
			for i := range itemsName {
				thisClass.Devices[i] += b.Devices[i]
			}
		}
		page.Classes = append(page.Classes, thisClass)
	}
	checkErr(tpl.ExecuteTemplate(w, "bookingList.html", page), "Template execute fatal: ")
}

func getDateClass() (date string, day string) {
	now := time.Now()
	date = now.Format("2006-01-02")
	day = now.Weekday().String()
	return
}

func bookingPage(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	if !user.Login {
		permissionDenied(w, r)
		return
	}
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	log.Println("id: ", id)
	row := db.QueryRow(`
		SELECT U.Name, LendingTime, ReturnTime
		FROM Bookings B, Users U
		WHERE B.ID = ? and U.ID = B.User;
		`, id)
	page := struct {
		User
		BID         int
		UName       string
		From        time.Time
		Until       time.Time
		Devices     [5]int
		ItemsName   [5]string
		AbleLending bool
	}{}
	page.User = user
	if err := row.Scan(&page.UName, &page.From, &page.Until); err == sql.ErrNoRows {
		// 404 no this booking
		notFound(w, r)
		return
	} else if err != nil {
		log.Fatalln("Query booking error: ", err)
		return
	}
	page.BID = id
	page.ItemsName = itemsName
	page.Devices = getBookingDevices(page.BID)
	now := time.Now()
	if now.Before(page.Until) {

	}

	checkErr(tpl.ExecuteTemplate(w, "booking.html", page), "Execute booking data page fatal: ")
}

func ableLendout(b Booking) bool {
	/*
		Return this booking is able to lendout or not
		Check whether every booking before it after now are already lended
	*/
	/*
		rows := db.QueryRow(`
			SELECT ID
			FROM Bookings
			WHERE ReturnDate
			`)
	*/
	return false
}

func alreadyLendout(b Booking) bool {
	stmt, err := db.Prepare(`
	SELECT COUNT(1)
	FROM Records R, Devices D
	WHERE R.Booking = ? and D.ID = R.Device and D.Type = ?;
	`)
	checkErr(err, "Prepare count booking record fatal: ")
	for i, t := range itemsType {
		row := stmt.QueryRow(b.ID, t)
		var count int
		if err := row.Scan(&count); err == sql.ErrNoRows {
			return false
		} else if err != nil {
			log.Fatalln(err)
		}
		if count < b.Devices[i] {
			return false
		}
	}
	return true
}

func getBookingDevices(id int) (devices [5]int) {
	rows, err := db.Query(`
	SELECT Type, Amount
	FROM BookingDevices
	WHERE BID = ?
	ORDER BY Type;
	`, id)
	checkErr(err, "Query booking devices fatal: ")
	i := 0
	for rows.Next() {
		var t string
		var amount int
		checkErr(rows.Scan(&t, &amount), "getBookingDevices Scan rows fatal: ")
		for itemsType[i] != t && i < 5 {
			i++
		}
		devices[i] = amount
	}
	return
}
