package main

import "net/http"

func booking(w http.ResponseWriter, r *http.Request) {
	type BookingPage struct {
		User
		Days    [5]string
		Classes [11]string
	}
	if r.Method == http.MethodPost {

	} else {
		tpl.ExecuteTemplate(w, "booking.html", nil)
	}
}
