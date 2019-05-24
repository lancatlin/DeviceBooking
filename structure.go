package main

import "time"

const (
	student = iota
	teacher
	chromebook
	wap
	projector
)

var itemsName = [5]string{"學生機", "教師機", "Chromebook", "無線基地台", "無線投影機"}
var itemsType = [5]string{"Student-iPad", "Teacher-iPad", "Chromebook", "WAP", "WirelessProjector"}
var classBegin = [...]string{"7:30", "8:20", "9:15", "10:10", "11:05", "12:30", "13:05", "14:00", "14:55", "15:55", "16:45"}
var classEnd = [...]string{"8:10", "9:05", "10:10", "10:55", "11:50", "13:00", "13:50", "14:45", "15:40", "16:40", "17:30"}
var className = [...]string{"早修", "C1", "C2", "C3", "C4", "午休", "C5", "C6", "C7", "C8", "C9"}
var typeToIndex = map[string]int{"Student-iPad": 0, "Teacher-iPad": 1, "Chromebook": 2, "WAP": 3, "WirelessProjector": 4}

// User is the struct for template executing
type User struct {
	ID       int
	Username string
	Type     string
	Login    bool
}

func newUser(id int, name, utype string) User {
	return User{id, name, utype, true}
}

func nilUser() User {
	var user User
	user.Login = false
	return user
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
