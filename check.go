package main

import (
	"net/http"
	"strconv"
	"time"
)

var count = [5]int{}

func initCheck() {
	/*
		Init total devices number
	*/
	stmt, err := db.Prepare(`
	SELECT COUNT(1) FROM Devices WHERE Type = ?;
	`)
	checkErr(err, "func initCheck: prepare statment fatal: ")
	for i, v := range itemsType {
		row := stmt.QueryRow(v)
		var value int
		checkErr(row.Scan(&value), "total "+v+" query fatal: ")
		count[i] = value
	}
}

func checkPage(w http.ResponseWriter, r *http.Request) {
	/*
		Handle /check
		Return remaining devices of each lesson.
		Choose week by "page"
	*/
	user := getUser(w, r)
	if !user.Login {
		http.Redirect(w, r, "/login", 303)
		return
	}
	page := struct {
		User
		Dates       [6]string
		Classes     [11]string
		Times       [11]string
		IPads       [6][11]int
		Chromebooks [6][11]int
	}{}
	p, err := strconv.Atoi(r.FormValue("p"))
	if err != nil {
		p = 0
	}
	page.Dates, page.IPads, page.Chromebooks = checkWeek(p)
	page.Classes = className
	page.Times = classBegin
	page.User = user
	checkErr(tpl.ExecuteTemplate(w, "check.html", page), "Execute Fatal: ")
}

func checkWeek(p int) (dates [6]string, iPads, chromebooks [6][11]int) {
	now := time.Now()
	begin := now.AddDate(0, 0, -1*int(now.Weekday())+7*p+1)
	for d := 0; d < 6; d++ {
		dates[d] = begin.AddDate(0, 0, d).Format("01-02")
		date := begin.AddDate(0, 0, d).Format("2006-01-02")
		for c := 0; c < 11; c++ {
			result := check(date, c)
			iPads[d][c] = result[student]
			chromebooks[d][c] = result[chromebook]
		}
	}
	return
}

func check(date string, class int) (result [5]int) {
	/*
		Check how many remaining devices in a lesson
	*/
	result = [5]int{}
	from, end := date+" "+classBegin[class]+":00", date+" "+classEnd[class]+":00"
	stmt, err := db.Prepare(`
	SELECT SUM(Amount)
	FROM Bookings B, BookingDevices BD
	WHERE ReturnTime > ? AND LendingTime < ? and BID = ID and Type = ?;
	`)
	checkErr(err, "prepare check statment fatal: ")
	for i, v := range itemsType {
		row := stmt.QueryRow(from, end, v)
		var value int
		if err := row.Scan(&value); err != nil {
			value = 0
		}
		result[i] = count[i] - value
	}
	return
}
