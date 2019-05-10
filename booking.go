package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

var classBegin = [...]string{"7:30", "8:20", "9:15", "10:10", "11:05", "12:30", "13:05", "14:00", "14:55", "15:55", "16:45"}
var classEnd = [...]string{"8:10", "9:05", "10:10", "10:55", "11:50", "13:00", "13:50", "14:45", "15:40", "16:40", "17:30"}

func errCheck(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

func booking(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		type booking struct {
			Student    int
			Teacher    int
			Chromebook int
			WAP        int
			Projector  int
			From       string
			Until      string
		}
		b := &booking{}
		var err error
		b.Student, err = strconv.Atoi(r.FormValue("student"))
		errCheck(err, "")
		b.Teacher, err = strconv.Atoi(r.FormValue("teacher"))
		errCheck(err, "")
		b.Chromebook, err = strconv.Atoi(r.FormValue("chromebook"))
		errCheck(err, "")
		b.WAP, err = strconv.Atoi(r.FormValue("wap"))
		errCheck(err, "")
		b.Projector, err = strconv.Atoi(r.FormValue("wireless-projector"))
		errCheck(err, "")
		date := r.FormValue("date")
		class, err := strconv.Atoi(r.FormValue("class"))
		errCheck(err, "")
		if class < 0 || class > 10 {
			log.Fatalln("error class time: ", class)
		}
		b.From = date + " " + classBegin[class] + ":00"
		b.Until = date + " " + classEnd[class] + ":00"
		log.Println(b)
	} else {
		type bookingPage struct {
			User
			Classes [11]string
			Min     string
			Max     string
		}
		var page bookingPage
		page.User = getUser(w, r)
		now := time.Now()
		max := now.AddDate(0, 1, 0)
		page.Min = now.Format("2006-01-02")
		page.Max = max.Format("2006-01-02")
		page.Classes = [...]string{"早修", "E1", "E2", "E3", "E4", "午休", "E5", "E6", "E7", "E8", "E9"}
		err := tpl.ExecuteTemplate(w, "booking.html", page)
		if err != nil {
			log.Fatal("Template execute fatal: ", err)
		}
	}
}
