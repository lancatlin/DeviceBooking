package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func booking(w http.ResponseWriter, r *http.Request) {
	type BookingPage struct {
		User
		Days    [6]string
		Classes [11]string
		Times   [11]string
	}
	p, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		p = 0
	}
	var page BookingPage
	page.User = getUser(w, r)
	now := time.Now()
	begin := now.AddDate(0, 0, -1*int(now.Weekday())+7*p+1)
	page.Days = [6]string{}
	for i := 0; i < 6; i++ {
		day := begin.AddDate(0, 0, i)
		page.Days[i] = fmt.Sprintf("%d/%d", day.Month(), day.Day())
	}
	page.Classes = [...]string{"早修", "E1", "E2", "E3", "E4", "午休", "E5", "E6", "E7", "E8", "E9"}
	page.Times = [...]string{"7:30", "8:20", "9:15", "10:10", "11:05", "12:30", "13:05", "14:00", "14:55", "15:55", "16:45"}
	if r.Method == http.MethodPost {

	} else {
		err := tpl.ExecuteTemplate(w, "booking.html", page)
		if err != nil {
			log.Fatal("Template execute fatal: ", err)
		}
	}
}
