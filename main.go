package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

var db *sql.DB
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
	var err error
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connection := fmt.Sprintf("%s:%s@unix(/var/run/mysqld/mysqld.sock)/%s", user, password, dbName)
	db, err = sql.Open("mysql", connection)
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	log.Println("Success init")
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/booking", booking)
	http.ListenAndServe(":8080", nil)
}
