package main

import (
	"database/sql"
	"flag"
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
	var err error
	tpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalln(err)
	}
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
	filename := flag.String("init-users", "", "The json file of user's data to init database.")
	isDeviced := flag.Bool("init-deviced", false, "Whether should program init the devices data")
	flag.Parse()
	if *filename != "" {
		initUsers(*filename)
	}
	if *isDeviced {
		initDevices()
	}
	initCheck()
	log.Println("Server runs on http://localhost:8080")
	http.HandleFunc("/", index)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/booking", newBooking)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/check", checkPage)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
