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
		Booking
		Status string
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
	SELECT B.ID, U.Name, B.LendingTime, B.ReturnTime
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
			err := rows.Scan(&b.ID, &b.UName, &b.From, &b.Until)
			checkErr(err, "row scan fatal: ")
			b.Devices = getBookingDevices(b.ID)
			b.Status = b.getStatus()
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
	var id int64
	if i, err := strconv.Atoi(mux.Vars(r)["id"]); err != nil {
		log.Fatalln(err)
	} else {
		id = int64(i)
	}
	log.Println("id: ", id)
	row := db.QueryRow(`
		SELECT U.Name, LendingTime, ReturnTime
		FROM Bookings B, Users U
		WHERE B.ID = ? and U.ID = B.User;
		`, id)
	page := struct {
		User
		Booking
		ItemsName   [5]string
		Status      string
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
	page.Booking.ID = id
	page.ItemsName = itemsName
	page.Devices = getBookingDevices(page.Booking.ID)
	page.Status = page.getStatus()
	page.AbleLending = page.Status == "可借出"
	checkErr(tpl.ExecuteTemplate(w, "booking.html", page), "Execute booking data page fatal: ")
}

func (b *Booking) getStatus() string {
	if b.alreadyReturned() {
		return "已歸還"
	} else if b.alreadyLendout() {
		return "已借出"
	} else if b.ableLendout() {
		return "可借出"
	} else if b.Until.Before(time.Now()) {
		return "預約過期"
	} else {
		return "尚不可借出"
	}
}

func (b *Booking) ableLendout() bool {
	/*
		Return this booking is able to lendout or not
		Check whether every booking before it after now are already lended
	*/
	if b.Until.Before(time.Now()) {
		// 已經到了歸還時間，不須借
		return false
	}
	if b.alreadyLendout() {
		// 已經全部借出，不需借
		return false
	}
	// 檢查在此事件之前的所有預約是否已經借出，如果否，則不能借出。
	rows, err := db.Query(`
			SELECT ID
			FROM Bookings
			WHERE ReturnTime > ? and ReturnTime < ?;
			`, time.Now(), b.Until)
	checkErr(err, "func ableLendout: Query fatal: ")
	for rows.Next() {
		var id int64
		checkErr(rows.Scan(&id), "func ableLendout: scan fatal: ")
		booking := Booking{
			ID:      id,
			Devices: getBookingDevices(id),
		}
		if !booking.alreadyLendout() {
			return false
		}
	}
	return true
}

func (b *Booking) alreadyLendout() bool {
	/*
		Return whether booking is all lending out
	*/
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

func (b *Booking) alreadyReturned() bool {
	if !b.alreadyLendout() {
		return false
	}
	var result bool
	row := db.QueryRow(`
		SELECT Done
		FROM Bookings
		WHERE ID = ?;
	`, b.ID)
	checkErr(row.Scan(&result), "func alreadyReturned: Scan fatal: ")
	if result {
		return true
	}
	var v int
	row = db.QueryRow(`
		SELECT COUNT(1)
		FROM Records
		WHERE Booking = ? and LentUntil is NULL;
	`, b.ID)

	if err := row.Scan(&v); err == sql.ErrNoRows {
		go func(b *Booking) {
			db.Exec(`
				UPDATE Bookings
				SET Done = true
				WHERE ID = ?;
			`, b.ID)
		}(b)
		return true
	} else if err != nil {
		log.Fatalln("func alreadyReturned: scan error: ", err)
	}
	return false
}
