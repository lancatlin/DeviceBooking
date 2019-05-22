package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// 紀錄預約資料
type booking struct {
	ID      int64
	Devices [5]int
	From    string
	Until   string
}

func newBooking(w http.ResponseWriter, r *http.Request) {
	// Handle /booking
	user := getUser(w, r)
	if !user.Login {
		permissionDenied(w, r)
		return
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
		checkErr(err, "")
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

func (b *booking) insertBooking(user User) {
	// insert booking data into database
	result, err := db.Exec(`
	INSERT INTO Bookings (User, LendingTime, ReturnTime, Student, Teacher, Chromebook, WAP, Projector)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`, user.ID, b.From, b.Until, b.Devices[student], b.Devices[teacher], b.Devices[chromebook], b.Devices[wap], b.Devices[projector])
	checkErr(err, "Insert booking value fatal: ")
	b.ID, err = result.LastInsertId()
	checkErr(err, "Get last insert ID fatal: ")
}

func (b *booking) enough(r [5]int) []string {
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

func getBooking(w http.ResponseWriter, r *http.Request) {
	/*
		Handle booking data
		Transfer request to "/booking/new" and "/booking/list"
	*/
	// 是否要檢查登入？
	switch url := r.URL.Path; url {
	case "list":
		bookingList(w, r)
	case "new":
		newBooking(w, r)
	default:
		user := getUser(w, r)
		if !user.Login {
			permissionDenied(w, r)
			return
		}
		id, err := strconv.Atoi(url)
		if err != nil {
			notFound(w, r)
			return
		}
		log.Println("id: ", id)
		row := db.QueryRow(`
		SELECT U.Name, LendingTime, ReturnTime, Student, Teacher, Chromebook, WAP, Projector
		FROM Bookings B, Users U
		WHERE B.ID = ? and U.ID = B.User;
		`, id)
		page := struct {
			User
			ID        int
			UName     string
			From      string
			Until     string
			Devices   [5]int
			ItemsName [5]string
		}{}
		page.User = user
		var s, t, c, wa, p int
		if err := row.Scan(&page.UName, &page.From, &page.Until, &s, &t, &c, &wa, &p); err == sql.ErrNoRows {
			// 404 no this booking
			notFound(w, r)
			return
		} else if err != nil {
			log.Fatalln("Query booking error: ", err)
			return
		}
		page.ID = id
		page.Devices[student] = s
		page.Devices[teacher] = t
		page.Devices[chromebook] = c
		page.Devices[wap] = wa
		page.Devices[projector] = p
		page.ItemsName = itemsName
		checkErr(tpl.ExecuteTemplate(w, "booking.html", page), "Execute booking data page fatal: ")
	}
}
