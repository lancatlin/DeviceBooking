package main

import (
	"database/sql"
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
	a, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		notFound(w, r)
		return
	}
	id = int64(a)
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
	checkErr(tpl.ExecuteTemplate(w, "lending.html", page), "Execute lending form fatal: ")
}

func devices(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	qType := r.FormValue("type")
	qLendout := r.FormValue("Lendout") == "true"
	if !user.Login {
		permissionDenied(w, r)
	}
	rows, err := db.Query(`
	SELECT * FROM DevicesStatus;
	`)
	if err != nil {
		log.Fatalln(err)
	}
	type Device struct {
		ID     string
		Status bool
		Uname  string
	}
	page := struct {
		User
		Devices []Device
		Types   map[string]string
	}{user, []Device{}, make(map[string]string)}
	for rows.Next() {
		var device string
		var status bool
		var uname sql.NullString
		var name string
		var dType string
		if err := rows.Scan(&device, &status, &uname, &dType); err != nil {
			log.Fatal(err)
		}
		if !uname.Valid {
			name = ""
		} else {
			name = uname.String
		}
		if qType != "" && qType != dType {
			continue
		}
		if qLendout && !status {
			continue
		}
		page.Devices = append(page.Devices, Device{device, status, name})
	}
	for i := range itemsType {
		page.Types[itemsType[i]] = itemsName[i]
	}
	if err := tpl.ExecuteTemplate(w, "devices.html", page); err != nil {
		log.Fatal(err)
	}
}
