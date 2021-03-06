package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

const (
	student = iota
	teacher
	chromebook
	wap
	projector
)

const (
	StatusLending        = "借出中"
	StatusReturned       = "已還入"
	StatusAbleLendout    = "可借出"
	StatusOverdue        = "預約過期"
	StatusNotAbleLendout = "不可借出"
)

var (
	// ErrBookingNotFound is returned by getBooking(id int64)
	ErrBookingNotFound = errors.New("Cannot found current booking by ID")
)

var itemsName = [5]string{"學生機", "教師機", "Chromebook", "無線基地台", "無線投影機"}
var itemsType = [5]string{"Student-iPad", "Teacher-iPad", "Chromebook", "WAP", "WirelessProjector"}

// var classBegin = [...]string{"7:30", "8:20", "9:15", "10:10", "11:05", "12:30", "13:05", "14:00", "14:55", "15:55", "16:45"}
var classBegin = [11]time.Duration{}
var classEnd = [11]time.Duration{}
var className = [...]string{"早修", "C1", "C2", "C3", "C4", "午休", "C5", "C6", "C7", "C8", "C9"}
var typeToIndex = map[string]int{"Student-iPad": 0, "Teacher-iPad": 1, "Chromebook": 2, "WAP": 3, "WirelessProjector": 4}

// User is the structure for template executing
type User struct {
	UID      int64
	Username string
	Type     string
	Login    bool
	Email    string
}

// Booking 紀錄預約資料
type Booking struct {
	ID      int64
	Devices [5]int
	From    time.Time
	Until   time.Time
	UName   string
	Status  string
}

// Record structure
type Record struct {
	Device  string
	Type    string
	Booking int64
	Done    bool
}

// Msg is the structure of the data msg.html need
type Msg struct {
	Title   string
	Content string
	Target  string
}

func init() {
	begin := [...]string{"7h30m", "8h20m", "9h15m", "10h10m", "11h05m", "12h30m", "13h05m", "14h", "14h55m", "15h55m", "16h45m"}
	end := [11]string{"8h10m", "9h05m", "10h10m", "10h55m", "11h50m", "13h", "13h50m", "14h45m", "15h40m", "16h40m", "17h30m"}
	for i := 0; i < 11; i++ {
		var err error
		classBegin[i], err = time.ParseDuration(begin[i])
		if err != nil {
			log.Fatalln(err)
		}
		classEnd[i], err = time.ParseDuration(end[i])
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func newUser(id int64, name, utype string) User {
	return User{id, name, utype, true, ""}
}

func nilUser() User {
	var user User
	user.Login = false
	return user
}

func loadUser(uid int64) (user User, err error) {
	row := db.QueryRow(`
	SELECT Name, Type, Email
	FROM Users WHERE ID = ?;
	`, uid)
	user = User{}
	user.UID = uid
	if err = row.Scan(&user.Username, &user.Type, &user.Email); err == sql.ErrNoRows {
		return user, err
	} else if err != nil {
		log.Fatal(err)
	}
	return user, nil
}

func datetime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func parseSQLTime(s string) time.Time {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	checkErr(err, "Parse time fatal: ")
	return t
}

func parseClass(date, class string) time.Time {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", date+" "+class+":00", time.Local)
	checkErr(err, "Parse class fatal: ")
	return t
}

func zeroClock(t time.Time) time.Time {
	h, m, s := t.Clock()
	return t.Add(-1 * (time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s)*time.Second))
}

func formatDuration(d time.Duration) string {
	h, m := int(d.Hours())%24, int(d.Minutes())%60
	return fmt.Sprintf("%d:%d", h, m)
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Panicln(msg, err)
	}
}
