package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

var iPadCount, chromebookCount int

func initCheck() {
	row := db.QueryRow(`
	SELECT COUNT(1) FROM Devices WHERE Type = 'Student-iPad';
	`)
	checkErr(row.Scan(&iPadCount), "Query count of student iPads")
	row = db.QueryRow(`
	SELECT COUNT(1) FROM Devices WHERE Type = 'Chromebook';
	`)
	checkErr(row.Scan(&chromebookCount), "Query count of student chromebooks")
}

func checkPage(w http.ResponseWriter, r *http.Request) {
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
		IPads       [11][6]int
		Chromebooks [11][6]int
	}{}
	p, err := strconv.Atoi(r.FormValue("p"))
	if err != nil {
		p = 0
	}
	page.Dates, page.IPads, page.Chromebooks = check(p)
	page.Classes = className
	page.Times = classBegin
	page.User = user
	checkErr(tpl.ExecuteTemplate(w, "check.html", page), "Execute Fatal: ")
}

func check(p int) (dates [6]string, iPads, chromebooks [11][6]int) {
	now := time.Now()
	begin := now.AddDate(0, 0, -1*int(now.Weekday())+7*p+1)
	stmt, err := db.Prepare(`SELECT SUM(Student), SUM(Chromebook) FROM Bookings WHERE ReturnTime > ? AND LendingTime < ?;`)
	checkErr(err, "Prepare statement fatal: ")
	for d := 0; d < 6; d++ {
		dates[d] = begin.AddDate(0, 0, d).Format("01-02")
		date := begin.AddDate(0, 0, d).Format("2006-01-02")
		for c := 0; c < 11; c++ {
			from, end := date+" "+classBegin[c]+":00", date+" "+classEnd[c]+":00"
			log.Println(from, end)
			var iPad, chrome int
			row := stmt.QueryRow(from, end)
			err := row.Scan(&iPad, &chrome)
			if err != nil {
				iPad, chrome = 0, 0
			}
			log.Println(iPad, chrome)
			iPads[c][d] = iPadCount - int(iPad)
			chromebooks[c][d] = chromebookCount - chrome
		}
	}
	log.Println(iPads, chromebooks)
	return
}
