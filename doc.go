package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gomarkdown/markdown"

	"github.com/gorilla/mux"
)

func docs(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	filename := "doc/" + mux.Vars(r)["filename"] + ".md"
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
