package main

import "net/http"

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
