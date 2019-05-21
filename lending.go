package main

import (
	"log"
	"net/http"
	"time"
)

func datetime(t time.Time) string {
	return t.Format("2006-01-02")
}

func bookingList(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	if !user.Login || user.Type != "Admin" {
		// Permission denied
		return
	}
	type displayBooking struct {
		ID      int64
		User    string
		Devices [5]int
	}
	type classBooking struct {
		Class    string
		Bookings []displayBooking
		Devices  [5]int
	}
	date, day, class := getDateClass()
	page := struct {
		User
		Date    string
		Day     string
		Classes []classBooking
	}{user, date, day, []classBooking{}}
	stmt, err := db.Prepare(`
	SELECT B.ID, U.Name, Teacher, Student, Chromebook, WAP, Projector
	FROM Bookings B, Users U
	WHERE ReturnTime > ? and LendingTime < ? and U.ID = B.User
	`)
	checkErr(err, "booking list check query prepare fatal: ")
	for c := class; c < len(className); c++ {
		thisClass := classBooking{className[c], []displayBooking{}, [5]int{}}
		rows, err := stmt.Query(date+" "+classBegin[c], date+" "+classEnd[c])
		checkErr(err, "bookings query error: ")
		defer rows.Close()
		for rows.Next() {
			b := displayBooking{}
			var t, s, c, w, p int
			err := rows.Scan(&b.ID, &b.User, &t, &s, &c, &w, &p)
			checkErr(err, "row scan fatal: ")
			b.Devices[student] = s
			b.Devices[teacher] = t
			b.Devices[chromebook] = c
			b.Devices[wap] = w
			b.Devices[projector] = p
			thisClass.Bookings = append(thisClass.Bookings, b)
			for i := range itemsName {
				thisClass.Devices[i] += b.Devices[i]
			}
		}
		page.Classes = append(page.Classes, thisClass)
	}
	checkErr(tpl.ExecuteTemplate(w, "lendingList.html", page), "Template execute fatal: ")
}

func getDateClass() (date string, day string, class int) {
	now := time.Now()
	class = -1
	for i, t := range classBegin {
		if c, err := time.Parse("2006-01-02 15:04", datetime(now)+" "+t); err != nil {
			log.Fatalln("Parse time fatal: ", datetime(now), t)
		} else if now.After(c) {
			class = i
			break
		}
	}
	if class == -1 {
		date = datetime(now.AddDate(0, 0, 1))
		day = now.AddDate(0, 0, 1).Weekday().String()
		class = 0
	} else {
		date = datetime(now)
		day = now.Weekday().String()
	}
	return
}
