package main

import "net/http"

func booking(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

	} else {
		bookingData := struct {
		}{}
		tpl.ExecuteTemplate(w, "booking.html", nil)
	}
}
