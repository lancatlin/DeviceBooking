package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// 紀錄預約資料
type Booking struct {
	ID      int64
	Devices [5]int
	From    time.Time
	Until   time.Time
	UName   string
}

func newBooking(w http.ResponseWriter, r *http.Request) {
	// Handle /booking
	user := getUser(w, r)
	if !user.Login {
		permissionDenied(w, r)
		return
	}
	if r.Method == "POST" {
		b := Booking{}
		var err error
		b.Devices[student], _ = strconv.Atoi(r.FormValue("student"))
		b.Devices[teacher], _ = strconv.Atoi(r.FormValue("teacher"))
		b.Devices[chromebook], _ = strconv.Atoi(r.FormValue("chromebook"))
		b.Devices[wap], _ = strconv.Atoi(r.FormValue("wap"))
		b.Devices[projector], _ = strconv.Atoi(r.FormValue("wireless-projector"))
		date := r.FormValue("date")
		class, err := strconv.Atoi(r.FormValue("class"))
		checkErr(err, "")
		if class < 0 || class > 10 {
			bookingForm(w, r, []string{"請輸入正確的課堂"})
			return
		}
		if msg := b.enough(check(date, class)); len(msg) != 0 {
			bookingForm(w, r, msg)
			return
		}
		b.From = parseSQLTime(date + " " + classBegin[class] + ":00")
		b.Until = parseSQLTime(date + " " + classEnd[class] + ":00")
		b.insertBooking(user)
		log.Println("b.ID: ", b.ID)
		checkErr(err, "Parse time fatal: ")
		page := struct {
			User
			Booking
			Date      string
			Class     string
			Weekday   string
			ItemsName [5]string
		}{user, b, date, className[class], b.From.Weekday().String(), itemsName}
		log.Println("page.ID", page.Booking.ID)
		checkErr(tpl.ExecuteTemplate(w, "examination.html", page), "Execute examination fatal: ")
	} else {
		bookingForm(w, r, nil)
	}
}

func bookingForm(w http.ResponseWriter, r *http.Request, msg []string) {
	// Render booking form and error messages
	type bookingPage struct {
		User
		Classes [11]string
		Min     string
		Max     string
		Msg     []string
	}
	var page bookingPage
	page.User = getUser(w, r)
	page.Msg = msg
	now := time.Now()
	max := now.AddDate(0, 1, 0)
	page.Min = now.Format("2006-01-02")
	page.Max = max.Format("2006-01-02")
	page.Classes = className
	err := tpl.ExecuteTemplate(w, "newBooking.html", page)
	if err != nil {
		log.Fatal("Template execute fatal: ", err)
	}
}

func (b *Booking) insertBooking(user User) {
	// insert booking data into database
	result, err := db.Exec(`
	INSERT INTO Bookings (User, LendingTime, ReturnTime)
	VALUES (?, ?, ?);
	`, user.ID, b.From, b.Until)
	checkErr(err, "Insert booking value fatal: ")
	b.ID, err = result.LastInsertId()
	checkErr(err, "Get last insert ID fatal: ")

	stmt, err := db.Prepare(`
	INSERT INTO BookingDevices
	VALUES (?, ?, ?);
	`)
	checkErr(err, "Insert query prepare fatal: ")
	for i, v := range b.Devices {
		if v > 0 {
			_, err := stmt.Exec(b.ID, itemsType[i], v)
			checkErr(err, "Insert booking devices fatal: ")
		}
	}
}

func (b *Booking) enough(r [5]int) []string {
	/*
		Check whether devices are enough for this booking.
		Return slice of error messages. len == 0 if all pass.
	*/
	msg := []string{}
	for i := 0; i < len(itemsName); i++ {
		if r[i] < b.Devices[i] {
			msg = append(msg, fmt.Sprintf("%s 只剩 %d 台", itemsName[i], r[i]))
		}
	}
	return msg
}
