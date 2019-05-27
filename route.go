package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	tpl.ExecuteTemplate(w, "index.html", user)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	/* 404 Handler */
	user := getUser(w, r)
	w.WriteHeader(http.StatusNotFound)
	checkErr(tpl.ExecuteTemplate(w, "notFound.html", user), "Render notFound fatal")
}

func permissionDenied(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	w.WriteHeader(403)
	checkErr(tpl.ExecuteTemplate(w, "permissionDenied.html", user), "Render permission denied error: ")
}

func lendForm(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	if !user.Login || user.Type != "Admin" {
		permissionDenied(w, r)
		return
	}
	var id int64
	if i, err := strconv.Atoi(mux.Vars(r)["id"]); err != nil {
		notFound(w, r)
		return
	} else {
		id = int64(i)
	}
	b, err := getBooking(id)
	if err == ErrBookingNotFound {
		notFound(w, r)
		return
	}
	if !b.ableLendout() {
		page := struct {
			User
			Title   string
			Content string
			Target  string
		}{user, "不可借出", fmt.Sprintf("%d 預約不可借出", b.ID), fmt.Sprintf(`<a href="/bookings/%d">預約項目頁面</a>`, b.ID)}
		if err := tpl.ExecuteTemplate(w, "msg.html", page); err != nil {
			log.Fatalln(err)
		}
		return
	}
	page := struct {
		User
		Booking
		ItemsName [5]string
		ItemsType [5]string
	}{user, b, itemsName, itemsType}
	log.Println(page)
	checkErr(tpl.ExecuteTemplate(w, "lending.html", page), "Execute lending form fatal: ")
}
