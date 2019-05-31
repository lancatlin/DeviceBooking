package main

import (
	"database/sql"
	"log"
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
		Dates    [7]time.Time
		Classes  [11]string
		Times    [11]time.Duration
		Devices  [7][11][5]int
		Previous int
		Next     int
	}{}
	p, err := strconv.Atoi(r.FormValue("p"))
	if err != nil {
		p = 0
	}
	page.Dates, page.Devices = checkWeek(p)
	page.Classes = className
	page.Times = classBegin
	page.User = user
	page.Previous = p - 1
	page.Next = p + 1
	checkErr(tpl.ExecuteTemplate(w, "check.html", page), "Execute Fatal: ")
}

func check(begin, end time.Time) (amount [5]int) {
	stmt, err := db.Prepare(`
	SELECT SUM(Amount)
	FROM BookingDevices BD, Bookings B
	WHERE Type = ? AND BID = ID 
	AND ReturnTime > ? AND LendingTime < ?;
	`)
	if err != nil {
		log.Fatalln(err)
	}
	for i, t := range itemsType {
		row := stmt.QueryRow(t, begin, end)
		var value sql.NullInt64
		if err := row.Scan(&value); err != nil {
			log.Panic(err)
		}
		if value.Valid {
			amount[i] = count[i] - int(value.Int64)
		} else {
			amount[i] = count[i]
		}
	}
	return
}

func checkWeek(w int) (dates [7]time.Time, amounts [7][11][5]int) {
	now := time.Now()
	begin := zeroClock(now.AddDate(0, 0, -1*int(now.Weekday())+1+7*w))
	for i := 0; i < 7; i++ {
		date := begin.AddDate(0, 0, i)
		dates[i] = date
		for j := 0; j < 11; j++ {
			amounts[i][j] = check(date.Add(classBegin[j]), date.Add(classEnd[j]))
		}
	}
	return
}
