package main

import (
	"database/sql"
	"fmt"
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
	d := r.FormValue("date")
	var date time.Time
	var err error
	if d == "" {
		date = zeroClock(time.Now())
	} else {
		date, err = time.ParseInLocation("2006-01-02", d, time.Local)
		checkErr(err, "Parse date fatal: ")
	}
	day := date.Weekday().String()
	page := struct {
		User
		Date      time.Time
		Day       string
		Classes   []classBooking
		Yesterday time.Time
		Tomorrow  time.Time
	}{user, date, day, []classBooking{}, date.AddDate(0, 0, -1), date.AddDate(0, 0, 1)}
	stmt, err := db.Prepare(`
	SELECT B.ID, U.Name, B.LendingTime, B.ReturnTime
	FROM Bookings B, Users U
	WHERE ReturnTime > ? and LendingTime < ? and U.ID = B.User
	`)
	checkErr(err, "booking list check query prepare fatal: ")
	for c := 0; c < len(className); c++ {
		thisClass := classBooking{className[c], []displayBooking{}, [5]int{}}
		rows, err := stmt.Query(date.Add(classBegin[c]), date.Add(classEnd[c]))
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

func newRecord(w http.ResponseWriter, r *http.Request) {
	var bID int64
	a, err := strconv.Atoi(r.FormValue("bid"))
	if err != nil {
		w.WriteHeader(404)
		return
	}
	bID = int64(a)
	dID := r.FormValue("device")
	// 檢查設備是否已借出
	row := db.QueryRow(`
		SELECT COUNT(1)
		FROM Records R, Devices D
		WHERE R.Device = D.ID and D.ID = ? and LentUntil IS NULL;
	`, dID)
	var v int
	if err := row.Scan(&v); v != 0 {
		// 已被借出
		w.WriteHeader(403)
		return
	} else if err != nil {
		log.Fatalln(err)
	}
	var dType string
	row = db.QueryRow(`
	SELECT Type FROM Devices WHERE ID = ?;
	`, dID)
	if err := row.Scan(&dType); err == sql.ErrNoRows {
		w.WriteHeader(404)
		return
	} else if err != nil {
		log.Fatalln(err)
	}
	// 檢查預約是否已滿
	b, err := getBooking(bID)
	if err == ErrBookingNotFound {
		w.WriteHeader(401)
		return
	}
	if !b.ableLendout() {
		w.WriteHeader(401)
		return
	}
	row = db.QueryRow(`
	SELECT COUNT(1)
	FROM Records R, Devices D
	WHERE Booking = ? and R.Device = D.ID and Type = ?;
	`, b.ID, dType)
	var amount int
	if err := row.Scan(&amount); err == sql.ErrNoRows {
		amount = 0
	} else if err != nil {
		log.Fatalln(err)
	}
	i := typeToIndex[dType]
	if amount == b.Devices[i] {
		w.WriteHeader(405)
		return
	} else if amount > b.Devices[i] {
		log.Fatalln("Amount more than booking !", amount, b.Devices[i])
	}
	// 借出設備
	result, err := db.Exec(`
	INSERT INTO Records
	VALUES (?, ?, ?, NULL);
	`, b.ID, dID, time.Now())
	checkErr(err, "Insert fatal: ")
	rID, err := result.LastInsertId()
	checkErr(err, "")
	_, err = fmt.Fprintf(w, `
	{
		"type": "%s",
		"amount": %d,
		"recordID": "%d",
		"done": %t
	}`, dType, amount+1, rID, b.alreadyLendout())
	checkErr(err, "")
}
