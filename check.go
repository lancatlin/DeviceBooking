package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
)

var count = [5]int{}
var checkStmt *sql.Stmt

func initCheck() {
	/*
		Init total devices number
	*/
	var s, t, c, w, p int
	stmt, err := db.Prepare(`
	SELECT COUNT(1) FROM Devices WHERE Type = ?;
	`)
	checkErr(err, "Statement prepare fatal: ")
	row := stmt.QueryRow("Student-iPad")
	checkErr(row.Scan(&s), "Query count of student iPads")
	row = stmt.QueryRow("Teacher-iPad")
	checkErr(row.Scan(&t), "Query count of student Teacher iPads")
	row = stmt.QueryRow("Chromebook")
	checkErr(row.Scan(&c), "Query count of student chromebooks")
	row = stmt.QueryRow("WAP")
	checkErr(row.Scan(&w), "Query count of student WAPs")
	row = stmt.QueryRow("WirelessProjector")
	checkErr(row.Scan(&t), "Query count of student Projectors")
	count[student] = s
	count[teacher] = t
	count[chromebook] = c
	count[wap] = w
	count[projector] = p
	checkStmt, err = db.Prepare(`SELECT SUM(Student), SUM(Teacher), SUM(Chromebook), SUM(WAP), SUM(Projector) FROM Bookings WHERE ReturnTime > ? AND LendingTime < ?;`)
	checkErr(err, "Prepare checking Stmt fatal: ")
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
	var s, t, c, w, p int
	row := checkStmt.QueryRow(from, end)
	err := row.Scan(&s, &t, &c, &w, &p)
	if err != nil {
		s, t, c, w, p = 0, 0, 0, 0, 0
	}
	result[student] = count[student] - s
	result[teacher] = count[teacher] - t
	result[chromebook] = count[chromebook] - c
	result[wap] = count[wap] - w
	result[projector] = count[projector] - p
	return
}
