package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
)

func index(w http.ResponseWriter, r *http.Request) {
	if initMode {
		http.Redirect(w, r, "/init", 303)
		return
	}
	user := getUser(w, r)
	filename := "README.md"
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		notFound(w, r)
		return
	}
	html := markdown.ToHTML(file, nil, nil)
	page := struct {
		User
		Msg
	}{user, Msg{"", string(html), ""}}
	if err := tpl.ExecuteTemplate(w, "msg.html", page); err != nil {
		http.Error(w, err.Error(), 500)
	}
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
		Option  string
		Checked bool
	}{user, []Device{}, make(map[string]string), qType, qLendout}
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
		if qType != "" && qType != "all" && qType != dType {
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

func overdue(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	if !user.Login || user.Type != "Admin" {
		permissionDenied(w, r)
		return
	}
	rows, err := db.Query(`
	SELECT B.ID, U.Name, B.LendingTime, B.ReturnTime
	FROM Bookings B
	INNER JOIN UnDoneBookings UB
	ON B.ID = UB.ID
	INNER JOIN Users U
	ON B.User = U.ID
	WHERE Amount > 0 AND ReturnTime < ?;
	`, time.Now())
	if err != nil {
		log.Println(err)
	}
	page := struct {
		User
		Bookings []Booking
	}{user, []Booking{}}
	for rows.Next() {
		b := Booking{}
		if err := rows.Scan(&b.ID, &b.UName, &b.From, &b.Until); err != nil {
			log.Fatalln(err)
		}
		page.Bookings = append(page.Bookings, b)
	}
	if err := tpl.ExecuteTemplate(w, "overdue.html", page); err != nil {
		log.Fatalln(err)
	}
}

func initPage(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "init.html", nil); err != nil {
		log.Fatalln(err)
	}
}

func resetPasswordPage(w http.ResponseWriter, r *http.Request) {
	uid, err := convertID(mux.Vars(r)["id"])
	if err != nil {
		notFound(w, r)
		return
	}
	user := getUser(w, r)
	if user.Type != "Admin" && user.UID != uid {
		permissionDenied(w, r)
		return
	}
	u := User{
		UID: uid,
	}
	row := db.QueryRow(`
	SELECT Name FROM Users WHERE ID = ?;
	`, uid)
	if err := row.Scan(&u.Username); err == sql.ErrNoRows {
		notFound(w, r)
		return
	} else if err != nil {
		log.Fatalln(err)
	}
	page := struct {
		User
		U     User
		Error string
	}{
		user, u, "",
	}
	if err := tpl.ExecuteTemplate(w, "resetPassword.html", page); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
