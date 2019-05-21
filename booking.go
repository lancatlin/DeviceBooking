package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var classBegin = [...]string{"7:30", "8:20", "9:15", "10:10", "11:05", "12:30", "13:05", "14:00", "14:55", "15:55", "16:45"}
var classEnd = [...]string{"8:10", "9:05", "10:10", "10:55", "11:50", "13:00", "13:50", "14:45", "15:40", "16:40", "17:30"}
var className = [...]string{"早修", "C1", "C2", "C3", "C4", "午休", "C5", "C6", "C7", "C8", "C9"}

type booking struct {
	ID      int64
	Devices [5]int
	From    string
	Until   string
}

func errCheck(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

func newBooking(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	if !user.Login {
		http.Redirect(w, r, "/login", 303)
	}
	if r.Method == "POST" {
		b := &booking{}
		var err error
		b.Devices[student], _ = strconv.Atoi(r.FormValue("student"))
		b.Devices[teacher], _ = strconv.Atoi(r.FormValue("teacher"))
		b.Devices[chromebook], _ = strconv.Atoi(r.FormValue("chromebook"))
		b.Devices[wap], _ = strconv.Atoi(r.FormValue("wap"))
		b.Devices[projector], _ = strconv.Atoi(r.FormValue("wireless-projector"))
		date := r.FormValue("date")
		class, err := strconv.Atoi(r.FormValue("class"))
		errCheck(err, "")
		if class < 0 || class > 10 {
			bookingForm(w, r, []string{"請輸入正確的課堂"})
			return
		}
		if msg := b.enough(check(date, class)); len(msg) != 0 {
			bookingForm(w, r, msg)
			return
		}
		b.From = date + " " + classBegin[class] + ":00"
		b.Until = date + " " + classEnd[class] + ":00"
		log.Println(b)
		b.insertBooking(user)
		d, err := time.Parse("2006-01-02", date)
		checkErr(err, "Parse time fatal: ")
		page := struct {
			User
			booking
			Date      string
			Class     string
			Weekday   string
			ItemsName [5]string
		}{user, *b, date, className[class], d.Weekday().String(), itemsName}
		checkErr(tpl.ExecuteTemplate(w, "examination.html", page), "Execute examination fatal: ")
	} else {
		bookingForm(w, r, nil)
	}
}

func bookingForm(w http.ResponseWriter, r *http.Request, msg []string) {
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
	err := tpl.ExecuteTemplate(w, "booking.html", page)
	if err != nil {
		log.Fatal("Template execute fatal: ", err)
	}
}

func (b *booking) insertBooking(user User) {
	result, err := db.Exec(`
	INSERT INTO Bookings (User, LendingTime, ReturnTime, Student, Teacher, Chromebook, WAP, Projector)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`, user.ID, b.From, b.Until, b.Devices[student], b.Devices[teacher], b.Devices[chromebook], b.Devices[wap], b.Devices[projector])
	checkErr(err, "Insert booking value fatal: ")
	b.ID, err = result.LastInsertId()
	checkErr(err, "Get last insert ID fatal: ")
}

func (b *booking) enough(r [5]int) []string {
	msg := []string{}
	for i := 0; i < len(itemsName); i++ {
		if r[i] < b.Devices[i] {
			msg = append(msg, fmt.Sprintf("%s 只剩 %d 台", itemsName[i], r[i]))
		}
	}
	return msg
}
